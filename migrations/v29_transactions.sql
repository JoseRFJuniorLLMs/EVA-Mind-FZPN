-- Migration v29: Transactions Table
-- Sistema de transações financeiras (pagamentos e extrato)

-- Tabela principal de transações
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    
    -- Relacionamentos
    subscription_id BIGINT REFERENCES subscriptions(id) ON DELETE SET NULL,
    user_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    
    -- Valores
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'BRL',
    
    -- Gateway de pagamento
    provider VARCHAR(20) NOT NULL CHECK (provider IN (
        'stripe',      -- Cartão de crédito internacional
        'asaas',       -- Pix brasileiro
        'opennode',    -- Bitcoin Lightning
        'wise',        -- Transferência internacional
        'nomad'        -- Conta global
    )),
    
    -- Status da transação
    status VARCHAR(30) DEFAULT 'pending' CHECK (status IN (
        'pending',           -- Aguardando pagamento
        'paid',              -- Pago com sucesso
        'failed',            -- Falhou
        'waiting_approval',  -- Aguardando aprovação manual (Wise/Nomad)
        'refunded'           -- Estornado
    )),
    
    -- Referência externa (ID no gateway)
    external_ref VARCHAR(255) UNIQUE,
    
    -- Comprovante (para pagamentos manuais)
    proof_url VARCHAR(512), -- URL no Google Cloud Storage
    
    -- Motivo de falha (se aplicável)
    failure_reason TEXT,
    
    -- Metadata adicional
    metadata_json JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT check_positive_amount CHECK (amount > 0),
    CONSTRAINT check_external_ref_not_empty CHECK (
        external_ref IS NULL OR length(external_ref) > 0
    )
);

-- Índices para performance
CREATE INDEX idx_transactions_user ON transactions(user_id);
CREATE INDEX idx_transactions_subscription ON transactions(subscription_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_provider ON transactions(provider);
CREATE INDEX idx_transactions_external_ref ON transactions(external_ref);
CREATE INDEX idx_transactions_created ON transactions(created_at DESC);

-- Índices compostos para queries frequentes
CREATE INDEX idx_transactions_user_created 
ON transactions(user_id, created_at DESC);

CREATE INDEX idx_transactions_status_provider 
ON transactions(status, provider) 
WHERE status = 'waiting_approval';

CREATE INDEX idx_transactions_user_status 
ON transactions(user_id, status);

-- View: Histórico de transações com detalhes
CREATE OR REPLACE VIEW v_transaction_history AS
SELECT 
    t.id,
    t.user_id,
    u.email,
    u.nome as user_name,
    t.subscription_id,
    s.plan_tier,
    t.amount,
    t.currency,
    t.provider,
    t.status,
    t.external_ref,
    t.proof_url,
    t.failure_reason,
    t.created_at,
    t.updated_at,
    EXTRACT(EPOCH FROM (t.updated_at - t.created_at)) / 60 as processing_time_minutes
FROM transactions t
JOIN usuarios u ON u.id = t.user_id
LEFT JOIN subscriptions s ON s.id = t.subscription_id
ORDER BY t.created_at DESC;

-- View: Transações pendentes de aprovação (pagamentos)
CREATE OR REPLACE VIEW v_pending_payment_approvals AS
SELECT 
    t.id as transaction_id,
    t.user_id,
    u.email,
    u.nome as user_name,
    t.amount,
    t.currency,
    t.provider,
    t.proof_url,
    t.created_at,
    EXTRACT(EPOCH FROM (NOW() - t.created_at)) / 3600 as hours_waiting
FROM transactions t
JOIN usuarios u ON u.id = t.user_id
WHERE t.status = 'waiting_approval'
ORDER BY t.created_at ASC;

-- View: Estatísticas de transações
CREATE OR REPLACE VIEW v_transaction_stats AS
SELECT 
    provider,
    status,
    currency,
    COUNT(*) as total_transactions,
    SUM(amount) as total_amount,
    AVG(amount) as avg_amount,
    MIN(amount) as min_amount,
    MAX(amount) as max_amount,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '24 hours') as last_24h,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') as last_7d,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') as last_30d
FROM transactions
GROUP BY provider, status, currency
ORDER BY provider, status;

-- View: Taxa de sucesso por gateway
CREATE OR REPLACE VIEW v_payment_success_rate AS
SELECT 
    provider,
    COUNT(*) as total_attempts,
    COUNT(*) FILTER (WHERE status = 'paid') as successful,
    COUNT(*) FILTER (WHERE status = 'failed') as failed,
    COUNT(*) FILTER (WHERE status = 'pending') as pending,
    ROUND(
        100.0 * COUNT(*) FILTER (WHERE status = 'paid') / NULLIF(COUNT(*), 0),
        2
    ) as success_rate_percent
FROM transactions
WHERE created_at >= NOW() - INTERVAL '30 days'
GROUP BY provider
ORDER BY success_rate_percent DESC;

-- Função: Atualizar status de transação
CREATE OR REPLACE FUNCTION update_transaction_status(
    p_transaction_id BIGINT,
    p_new_status VARCHAR(30),
    p_failure_reason TEXT DEFAULT NULL
) RETURNS VOID AS $$
BEGIN
    UPDATE transactions
    SET status = p_new_status,
        failure_reason = p_failure_reason,
        updated_at = NOW()
    WHERE id = p_transaction_id;
    
    RAISE NOTICE 'Transaction % updated to status: %', p_transaction_id, p_new_status;
END;
$$ LANGUAGE plpgsql;

-- Função: Buscar transação por referência externa
CREATE OR REPLACE FUNCTION get_transaction_by_external_ref(p_external_ref VARCHAR(255))
RETURNS TABLE (
    id BIGINT,
    user_id BIGINT,
    subscription_id BIGINT,
    amount DECIMAL(10,2),
    status VARCHAR(30),
    provider VARCHAR(20)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        t.id,
        t.user_id,
        t.subscription_id,
        t.amount,
        t.status,
        t.provider
    FROM transactions t
    WHERE t.external_ref = p_external_ref;
END;
$$ LANGUAGE plpgsql;

-- Função: Processar pagamento confirmado
CREATE OR REPLACE FUNCTION process_payment_confirmation(
    p_transaction_id BIGINT,
    p_external_ref VARCHAR(255) DEFAULT NULL
) RETURNS VOID AS $$
DECLARE
    v_subscription_id BIGINT;
    v_amount DECIMAL(10,2);
BEGIN
    -- Atualizar status da transação
    UPDATE transactions
    SET status = 'paid',
        external_ref = COALESCE(p_external_ref, external_ref),
        updated_at = NOW()
    WHERE id = p_transaction_id
    RETURNING subscription_id, amount INTO v_subscription_id, v_amount;
    
    -- Se houver subscription associada, estender período
    IF v_subscription_id IS NOT NULL THEN
        -- Determinar dias baseado no valor (simplificado)
        PERFORM extend_subscription_period(
            v_subscription_id,
            CASE 
                WHEN v_amount >= 500 THEN 365  -- Anual
                ELSE 30                         -- Mensal
            END
        );
    END IF;
    
    RAISE NOTICE 'Payment confirmed for transaction %', p_transaction_id;
END;
$$ LANGUAGE plpgsql;

-- Trigger: Atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_transactions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_transactions_updated_at();

-- Tabela de log de mudanças de status de transação
CREATE TABLE IF NOT EXISTS transaction_status_log (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    old_status VARCHAR(30),
    new_status VARCHAR(30),
    changed_at TIMESTAMP DEFAULT NOW(),
    metadata_json JSONB DEFAULT '{}'
);

CREATE INDEX idx_transaction_log_transaction ON transaction_status_log(transaction_id);
CREATE INDEX idx_transaction_log_date ON transaction_status_log(changed_at DESC);

-- Trigger: Registrar mudanças de status de transação
CREATE OR REPLACE FUNCTION log_transaction_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        INSERT INTO transaction_status_log (
            transaction_id,
            old_status,
            new_status,
            changed_at
        ) VALUES (
            NEW.id,
            OLD.status,
            NEW.status,
            NOW()
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_log_transaction_status
    AFTER UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION log_transaction_status_change();

-- Comentários para documentação
COMMENT ON TABLE transactions IS 'Transações financeiras (pagamentos e extrato)';
COMMENT ON COLUMN transactions.provider IS 'Gateway: stripe, asaas, opennode, wise, nomad';
COMMENT ON COLUMN transactions.status IS 'Status: pending, paid, failed, waiting_approval, refunded';
COMMENT ON COLUMN transactions.external_ref IS 'ID único no gateway de pagamento';
COMMENT ON COLUMN transactions.proof_url IS 'URL do comprovante no Google Cloud Storage';
COMMENT ON COLUMN transactions.metadata_json IS 'Dados adicionais: {invoice_id, session_id, etc}';

COMMENT ON FUNCTION update_transaction_status IS 'Atualiza status de uma transação';
COMMENT ON FUNCTION get_transaction_by_external_ref IS 'Busca transação por referência externa';
COMMENT ON FUNCTION process_payment_confirmation IS 'Processa confirmação de pagamento e estende assinatura';
