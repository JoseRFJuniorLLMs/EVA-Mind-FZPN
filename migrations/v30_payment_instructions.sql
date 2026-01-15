-- Migration v30: Payment Instructions Table
-- Instruções bancárias para pagamentos internacionais (Wise, Nomad)

-- Tabela de instruções de pagamento
CREATE TABLE IF NOT EXISTS payment_instructions (
    id SERIAL PRIMARY KEY,
    
    -- Provider (Wise ou Nomad)
    provider VARCHAR(20) NOT NULL CHECK (provider IN ('wise', 'nomad')),
    
    -- Moeda
    currency VARCHAR(3) NOT NULL,
    
    -- Detalhes bancários (JSON flexível)
    details_json JSONB NOT NULL,
    /* Exemplo de estrutura:
    {
        "account_holder": "EVA Payments Ltd",
        "iban": "DE89370400440532013000",
        "swift": "COBADEFF",
        "bank_name": "Commerzbank",
        "routing_number": "026009593",  // Para USD/ACH
        "account_number": "1234567890",
        "reference_template": "EVA-{user_id}-{timestamp}"
    }
    */
    
    -- Status
    active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Índices
CREATE INDEX idx_payment_instructions_provider ON payment_instructions(provider);
CREATE INDEX idx_payment_instructions_currency ON payment_instructions(currency);
CREATE INDEX idx_payment_instructions_active ON payment_instructions(active);

-- Índice composto para busca rápida
CREATE INDEX idx_payment_instructions_provider_currency 
ON payment_instructions(provider, currency, active);

-- View: Instruções ativas
CREATE OR REPLACE VIEW v_active_payment_instructions AS
SELECT 
    id,
    provider,
    currency,
    details_json,
    created_at,
    updated_at
FROM payment_instructions
WHERE active = TRUE
ORDER BY provider, currency;

-- Função: Obter instruções de pagamento
CREATE OR REPLACE FUNCTION get_payment_instructions(
    p_provider VARCHAR(20),
    p_currency VARCHAR(3)
) RETURNS TABLE (
    id INT,
    provider VARCHAR(20),
    currency VARCHAR(3),
    details_json JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        pi.id,
        pi.provider,
        pi.currency,
        pi.details_json
    FROM payment_instructions pi
    WHERE pi.provider = p_provider
      AND pi.currency = p_currency
      AND pi.active = TRUE
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- Função: Gerar referência de pagamento única
CREATE OR REPLACE FUNCTION generate_payment_reference(
    p_user_id BIGINT,
    p_provider VARCHAR(20)
) RETURNS VARCHAR(100) AS $$
DECLARE
    v_timestamp BIGINT;
    v_reference VARCHAR(100);
BEGIN
    v_timestamp := EXTRACT(EPOCH FROM NOW())::BIGINT;
    v_reference := UPPER(p_provider) || '-' || p_user_id || '-' || v_timestamp;
    
    RETURN v_reference;
END;
$$ LANGUAGE plpgsql;

-- Trigger: Atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_payment_instructions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_payment_instructions_updated_at
    BEFORE UPDATE ON payment_instructions
    FOR EACH ROW
    EXECUTE FUNCTION update_payment_instructions_updated_at();

-- Dados iniciais: Wise (EUR)
INSERT INTO payment_instructions (provider, currency, details_json, active)
VALUES (
    'wise',
    'EUR',
    '{
        "account_holder": "EVA Payments Ltd",
        "iban": "DE89370400440532013000",
        "swift": "COBADEFF",
        "bic": "COBADEFF",
        "bank_name": "Commerzbank AG",
        "bank_address": "Frankfurt, Germany",
        "reference_template": "EVA-{user_id}-{timestamp}",
        "instructions": [
            "Use o código de referência fornecido",
            "Transferências levam 1-3 dias úteis",
            "Envie comprovante após transferência"
        ]
    }'::JSONB,
    TRUE
) ON CONFLICT DO NOTHING;

-- Dados iniciais: Wise (USD)
INSERT INTO payment_instructions (provider, currency, details_json, active)
VALUES (
    'wise',
    'USD',
    '{
        "account_holder": "EVA Payments LLC",
        "routing_number": "026009593",
        "account_number": "1234567890",
        "account_type": "Checking",
        "bank_name": "Bank of America",
        "bank_address": "New York, NY, USA",
        "swift": "BOFAUS3N",
        "reference_template": "EVA-{user_id}-{timestamp}",
        "instructions": [
            "ACH transfers: 2-3 business days",
            "Wire transfers: Same day",
            "Include reference code in transfer description"
        ]
    }'::JSONB,
    TRUE
) ON CONFLICT DO NOTHING;

-- Dados iniciais: Wise (GBP)
INSERT INTO payment_instructions (provider, currency, details_json, active)
VALUES (
    'wise',
    'GBP',
    '{
        "account_holder": "EVA Payments UK Ltd",
        "sort_code": "23-14-70",
        "account_number": "12345678",
        "iban": "GB29NWBK60161331926819",
        "swift": "NWBKGB2L",
        "bank_name": "NatWest",
        "bank_address": "London, UK",
        "reference_template": "EVA-{user_id}-{timestamp}",
        "instructions": [
            "UK bank transfers: Same day",
            "International transfers: 1-2 days",
            "Use reference code provided"
        ]
    }'::JSONB,
    TRUE
) ON CONFLICT DO NOTHING;

-- Dados iniciais: Nomad (USD)
INSERT INTO payment_instructions (provider, currency, details_json, active)
VALUES (
    'nomad',
    'USD',
    '{
        "account_holder": "EVA Payments",
        "routing_number": "084009519",
        "account_number": "9876543210",
        "account_type": "Checking",
        "bank_name": "Nomad Global",
        "bank_address": "New York, NY, USA",
        "reference_template": "EVA-NOMAD-{user_id}-{timestamp}",
        "instructions": [
            "ACH transfers only",
            "Processing time: 2-5 business days",
            "Include full reference code"
        ]
    }'::JSONB,
    TRUE
) ON CONFLICT DO NOTHING;

-- Dados iniciais: Nomad (EUR)
INSERT INTO payment_instructions (provider, currency, details_json, active)
VALUES (
    'nomad',
    'EUR',
    '{
        "account_holder": "EVA Payments Europe",
        "iban": "LT123456789012345678",
        "swift": "REVOLT21",
        "bic": "REVOLT21",
        "bank_name": "Nomad Europe",
        "bank_address": "Vilnius, Lithuania",
        "reference_template": "EVA-NOMAD-{user_id}-{timestamp}",
        "instructions": [
            "SEPA transfers: 1-2 business days",
            "SWIFT transfers: 2-3 business days",
            "Always include reference code"
        ]
    }'::JSONB,
    TRUE
) ON CONFLICT DO NOTHING;

-- Comentários para documentação
COMMENT ON TABLE payment_instructions IS 'Instruções bancárias para pagamentos internacionais';
COMMENT ON COLUMN payment_instructions.provider IS 'Provider: wise ou nomad';
COMMENT ON COLUMN payment_instructions.currency IS 'Moeda: EUR, USD, GBP, etc';
COMMENT ON COLUMN payment_instructions.details_json IS 'Detalhes bancários completos em JSON';
COMMENT ON COLUMN payment_instructions.active IS 'Se as instruções estão ativas';

COMMENT ON FUNCTION get_payment_instructions IS 'Retorna instruções de pagamento para provider e moeda';
COMMENT ON FUNCTION generate_payment_reference IS 'Gera código de referência único para pagamento';
