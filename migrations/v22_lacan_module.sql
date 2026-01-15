-- ================================================
-- Migração v22: Módulo Lacaniano
-- Data: 2026-01-15
-- Objetivo: Suporte a análise psicanalítica lacaniana
-- ================================================

-- Tabela de marcadores de transferência
CREATE TABLE IF NOT EXISTS transferencia_markers (
    id SERIAL PRIMARY KEY,
    idoso_id INT NOT NULL REFERENCES idosos(id) ON DELETE CASCADE,
    marker_type VARCHAR(50) NOT NULL, -- 'filial', 'materna', 'paterna', 'conjugal', 'fraternal'
    phrase TEXT NOT NULL, -- Frase que indicou transferência
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_transferencia_idoso ON transferencia_markers(idoso_id);
CREATE INDEX idx_transferencia_type ON transferencia_markers(marker_type);

-- Tabela de significantes recorrentes
CREATE TABLE IF NOT EXISTS significantes_recorrentes (
    id SERIAL PRIMARY KEY,
    idoso_id INT NOT NULL REFERENCES idosos(id) ON DELETE CASCADE,
    palavra VARCHAR(100) NOT NULL,
    frequencia INT DEFAULT 1,
    contextos TEXT[], -- Contextos onde a palavra apareceu
    primeira_ocorrencia TIMESTAMPTZ DEFAULT NOW(),
    ultima_ocorrencia TIMESTAMPTZ DEFAULT NOW(),
    ultima_interpelacao TIMESTAMPTZ, -- Quando EVA interpelou sobre esse significante
    
    UNIQUE(idoso_id, palavra)
);

CREATE INDEX idx_signif_idoso_freq ON significantes_recorrentes(idoso_id, frequencia DESC);
CREATE INDEX idx_signif_palavra ON significantes_recorrentes(palavra);

-- View para análise de transferência
CREATE OR REPLACE VIEW transferencia_summary AS
SELECT 
    idoso_id,
    marker_type,
    COUNT(*) as occurrences,
    MAX(timestamp) as last_occurrence
FROM transferencia_markers
WHERE timestamp > NOW() - INTERVAL '30 days'
GROUP BY idoso_id, marker_type
ORDER BY idoso_id, occurrences DESC;

-- View para significantes mais frequentes
CREATE OR REPLACE VIEW top_significantes AS
SELECT 
    idoso_id,
    palavra,
    frequencia,
    ultima_ocorrencia,
    ultima_interpelacao,
    CASE 
        WHEN ultima_interpelacao IS NULL THEN true
        WHEN AGE(NOW(), ultima_interpelacao) > INTERVAL '7 days' THEN true
        ELSE false
    END as deve_interpelar
FROM significantes_recorrentes
WHERE frequencia >= 3
ORDER BY frequencia DESC;

-- Comentários
COMMENT ON TABLE transferencia_markers IS 'Rastreia padrões de transferência psicanalítica';
COMMENT ON TABLE significantes_recorrentes IS 'Rastreia significantes (palavras-chave) recorrentes na fala do idoso';
COMMENT ON COLUMN significantes_recorrentes.ultima_interpelacao IS 'Quando EVA perguntou sobre o significado dessa palavra';

COMMIT;
