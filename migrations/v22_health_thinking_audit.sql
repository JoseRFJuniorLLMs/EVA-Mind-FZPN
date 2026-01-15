-- Migration v22: Health Thinking Audit
-- Adiciona tabelas para auditoria de análises de saúde via Gemini Thinking Mode

-- Tabela de auditoria de análises de saúde
CREATE TABLE IF NOT EXISTS health_thinking_audit (
    id BIGSERIAL PRIMARY KEY,
    idoso_id BIGINT NOT NULL REFERENCES idosos(id) ON DELETE CASCADE,
    concern TEXT NOT NULL,
    thought_process JSONB,
    risk_level VARCHAR(20) CHECK (risk_level IN ('BAIXO', 'MÉDIO', 'ALTO', 'CRÍTICO')),
    recommended_actions JSONB,
    seek_medical_care BOOLEAN DEFAULT false,
    urgency_level VARCHAR(20) CHECK (urgency_level IN ('immediate', 'within_24h', 'within_week', 'routine')),
    caregiver_notified BOOLEAN DEFAULT false,
    notified_at TIMESTAMP,
    final_answer TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Índices para performance
CREATE INDEX idx_health_audit_idoso ON health_thinking_audit(idoso_id);
CREATE INDEX idx_health_audit_risk ON health_thinking_audit(risk_level);
CREATE INDEX idx_health_audit_date ON health_thinking_audit(created_at DESC);
CREATE INDEX idx_health_audit_urgency ON health_thinking_audit(urgency_level);
CREATE INDEX idx_health_audit_notified ON health_thinking_audit(caregiver_notified);

-- View para dashboard de saúde
CREATE OR REPLACE VIEW v_health_concerns_summary AS
SELECT 
    i.id as idoso_id,
    i.nome as idoso_nome,
    COUNT(*) as total_concerns,
    COUNT(*) FILTER (WHERE risk_level = 'CRÍTICO') as critical_count,
    COUNT(*) FILTER (WHERE risk_level = 'ALTO') as high_count,
    COUNT(*) FILTER (WHERE risk_level = 'MÉDIO') as medium_count,
    COUNT(*) FILTER (WHERE risk_level = 'BAIXO') as low_count,
    COUNT(*) FILTER (WHERE caregiver_notified = true) as notified_count,
    MAX(created_at) as last_concern_date
FROM idosos i
LEFT JOIN health_thinking_audit hta ON hta.idoso_id = i.id
WHERE hta.created_at >= NOW() - INTERVAL '30 days'
GROUP BY i.id, i.nome;

-- View para alertas críticos não notificados
CREATE OR REPLACE VIEW v_critical_alerts_pending AS
SELECT 
    hta.id,
    hta.idoso_id,
    i.nome as idoso_nome,
    hta.concern,
    hta.risk_level,
    hta.urgency_level,
    hta.created_at,
    EXTRACT(EPOCH FROM (NOW() - hta.created_at))/60 as minutes_since_alert
FROM health_thinking_audit hta
JOIN idosos i ON i.id = hta.idoso_id
WHERE hta.risk_level IN ('CRÍTICO', 'ALTO')
  AND hta.caregiver_notified = false
  AND hta.created_at >= NOW() - INTERVAL '24 hours'
ORDER BY hta.created_at DESC;

-- Função para marcar notificação enviada
CREATE OR REPLACE FUNCTION mark_caregiver_notified(audit_id BIGINT)
RETURNS VOID AS $$
BEGIN
    UPDATE health_thinking_audit
    SET caregiver_notified = true,
        notified_at = NOW()
    WHERE id = audit_id;
END;
$$ LANGUAGE plpgsql;

-- Comentários para documentação
COMMENT ON TABLE health_thinking_audit IS 'Auditoria de análises de saúde realizadas pelo Gemini Thinking Mode';
COMMENT ON COLUMN health_thinking_audit.thought_process IS 'Array JSON com os passos de raciocínio do modelo';
COMMENT ON COLUMN health_thinking_audit.risk_level IS 'Nível de risco identificado: BAIXO, MÉDIO, ALTO, CRÍTICO';
COMMENT ON COLUMN health_thinking_audit.urgency_level IS 'Urgência para buscar atendimento médico';
COMMENT ON VIEW v_health_concerns_summary IS 'Resumo de preocupações de saúde por idoso nos últimos 30 dias';
COMMENT ON VIEW v_critical_alerts_pending IS 'Alertas críticos que ainda não foram notificados aos cuidadores';
