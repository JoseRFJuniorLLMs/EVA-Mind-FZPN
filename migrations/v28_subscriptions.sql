-- Migration v28: Subscriptions Table
-- Sistema de assinaturas para planos Basic, Gold e Diamond

-- Tabela principal de assinaturas
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    
    -- Status da assinatura
    status VARCHAR(20) DEFAULT 'trialing' CHECK (status IN (
        'active',      -- Assinatura ativa e paga
        'past_due',    -- Pagamento atrasado (grace period)
        'canceled',    -- Cancelada pelo usuário ou sistema
        'trialing'     -- Período de teste
    )),
    
    -- Plano contratado
    plan_tier VARCHAR(20) NOT NULL CHECK (plan_tier IN (
        'basic',       -- Gratuito/básico
        'gold',        -- Premium intermediário
        'diamond'      -- Premium completo
    )),
    
    -- Período de cobrança
    current_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    current_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- Método de pagamento padrão
    payment_method_default VARCHAR(50), -- 'stripe_card', 'asaas_pix', 'bitcoin', etc
    
    -- ID externo no gateway (Stripe, Asaas, etc)
    external_subscription_id VARCHAR(255) UNIQUE,
    
    -- Metadata adicional (frequência, cupons, etc)
    metadata_json JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT check_period_dates CHECK (current_period_end > current_period_start)
);

-- Índices para performance
CREATE INDEX idx_subscriptions_user ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_end_date ON subscriptions(current_period_end);
CREATE INDEX idx_subscriptions_tier ON subscriptions(plan_tier);

-- Índice único: um usuário só pode ter uma assinatura ativa por tier
CREATE UNIQUE INDEX idx_subscriptions_user_tier 
ON subscriptions(user_id, plan_tier) 
WHERE status = 'active';

-- Índice composto para queries de expiração
CREATE INDEX idx_subscriptions_status_end 
ON subscriptions(status, current_period_end) 
WHERE status = 'active';

-- View: Assinaturas ativas
CREATE OR REPLACE VIEW v_active_subscriptions AS
SELECT 
    s.id,
    s.user_id,
    u.email,
    u.nome,
    s.plan_tier,
    s.status,
    s.current_period_start,
    s.current_period_end,
    s.payment_method_default,
    EXTRACT(EPOCH FROM (s.current_period_end - NOW())) / 86400 as days_remaining,
    CASE 
        WHEN s.current_period_end < NOW() THEN 'expired'
        WHEN s.current_period_end < NOW() + INTERVAL '7 days' THEN 'expiring_soon'
        ELSE 'active'
    END as period_status
FROM subscriptions s
JOIN usuarios u ON u.id = s.user_id
WHERE s.status IN ('active', 'trialing')
ORDER BY s.current_period_end ASC;

-- View: Assinaturas por plano
CREATE OR REPLACE VIEW v_subscriptions_by_plan AS
SELECT 
    plan_tier,
    status,
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE current_period_end > NOW()) as active_count,
    COUNT(*) FILTER (WHERE current_period_end < NOW()) as expired_count
FROM subscriptions
GROUP BY plan_tier, status
ORDER BY plan_tier, status;

-- Função: Estender período de assinatura
CREATE OR REPLACE FUNCTION extend_subscription_period(
    p_subscription_id BIGINT,
    p_days INT DEFAULT 30
) RETURNS VOID AS $$
BEGIN
    UPDATE subscriptions
    SET current_period_end = current_period_end + (p_days || ' days')::INTERVAL,
        status = 'active',
        updated_at = NOW()
    WHERE id = p_subscription_id;
    
    RAISE NOTICE 'Subscription % extended by % days', p_subscription_id, p_days;
END;
$$ LANGUAGE plpgsql;

-- Função: Verificar acesso de assinatura
CREATE OR REPLACE FUNCTION check_subscription_access(p_user_id BIGINT)
RETURNS BOOLEAN AS $$
DECLARE
    has_access BOOLEAN;
BEGIN
    SELECT EXISTS(
        SELECT 1
        FROM subscriptions
        WHERE user_id = p_user_id
          AND status = 'active'
          AND current_period_end > NOW()
    ) INTO has_access;
    
    RETURN has_access;
END;
$$ LANGUAGE plpgsql;

-- Função: Obter tier da assinatura ativa
CREATE OR REPLACE FUNCTION get_user_subscription_tier(p_user_id BIGINT)
RETURNS VARCHAR(20) AS $$
DECLARE
    v_tier VARCHAR(20);
BEGIN
    SELECT plan_tier INTO v_tier
    FROM subscriptions
    WHERE user_id = p_user_id
      AND status = 'active'
      AND current_period_end > NOW()
    ORDER BY 
        CASE plan_tier
            WHEN 'diamond' THEN 1
            WHEN 'gold' THEN 2
            WHEN 'basic' THEN 3
        END
    LIMIT 1;
    
    RETURN COALESCE(v_tier, 'basic');
END;
$$ LANGUAGE plpgsql;

-- Trigger: Atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_subscriptions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_subscriptions_updated_at();

-- Tabela de histórico de mudanças de status
CREATE TABLE IF NOT EXISTS subscription_status_history (
    id BIGSERIAL PRIMARY KEY,
    subscription_id BIGINT NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    old_status VARCHAR(20),
    new_status VARCHAR(20),
    reason TEXT,
    changed_by BIGINT REFERENCES usuarios(id),
    changed_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_status_history_subscription ON subscription_status_history(subscription_id);
CREATE INDEX idx_status_history_date ON subscription_status_history(changed_at DESC);

-- Trigger: Registrar mudanças de status
CREATE OR REPLACE FUNCTION log_subscription_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        INSERT INTO subscription_status_history (
            subscription_id,
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

CREATE TRIGGER trigger_log_subscription_status
    AFTER UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION log_subscription_status_change();

-- Comentários para documentação
COMMENT ON TABLE subscriptions IS 'Assinaturas de usuários (Basic, Gold, Diamond)';
COMMENT ON COLUMN subscriptions.status IS 'Status: active, past_due, canceled, trialing';
COMMENT ON COLUMN subscriptions.plan_tier IS 'Plano: basic (gratuito), gold, diamond';
COMMENT ON COLUMN subscriptions.current_period_end IS 'Data de expiração do período atual';
COMMENT ON COLUMN subscriptions.external_subscription_id IS 'ID no gateway de pagamento (Stripe, etc)';
COMMENT ON COLUMN subscriptions.metadata_json IS 'Dados adicionais: {frequency: monthly/yearly, coupon: xxx}';

COMMENT ON FUNCTION extend_subscription_period IS 'Estende período de assinatura por N dias';
COMMENT ON FUNCTION check_subscription_access IS 'Verifica se usuário tem assinatura ativa';
COMMENT ON FUNCTION get_user_subscription_tier IS 'Retorna tier da assinatura ativa (ou basic se nenhuma)';
