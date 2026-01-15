-- Migration v24: Medical Image Analysis (Expanded for Public Health)
-- Adiciona tabelas para análise de imagens médicas via MedGemma
-- Suporta: Receitas, Feridas, Malária, TB, Testes Rápidos, Lesões Cutâneas, Úlceras, Pé Diabético

-- Tabela de análises de imagens médicas
CREATE TABLE IF NOT EXISTS medical_image_analysis (
    id BIGSERIAL PRIMARY KEY,
    idoso_id BIGINT NOT NULL REFERENCES idosos(id) ON DELETE CASCADE,
    image_type VARCHAR(50) NOT NULL CHECK (image_type IN (
        'prescription',           -- Receita médica
        'wound',                  -- Ferida genérica
        'lab_result',            -- Resultado de exame
        'medication_photo',      -- Foto de medicamento
        'malaria_smear',         -- Esfregaço de malária
        'chest_xray',            -- Raio-X de tórax (TB)
        'rapid_test',            -- Teste rápido (COVID, HIV, Dengue)
        'skin_lesion',           -- Lesão cutânea (Mpox, melanoma)
        'pressure_ulcer',        -- Úlcera de pressão (escara)
        'diabetic_foot',         -- Pé diabético
        'other'
    )),
    image_url TEXT,
    geolocation POINT,                    -- Coordenadas GPS para epidemiologia
    test_metadata JSONB,                  -- Metadados do teste (fabricante, lote, etc)
    analysis_result JSONB NOT NULL,
    severity VARCHAR(20),
    requires_medical_attention BOOLEAN DEFAULT false,
    caregiver_notified BOOLEAN DEFAULT false,
    notified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Índices para performance
CREATE INDEX idx_medical_image_idoso ON medical_image_analysis(idoso_id);
CREATE INDEX idx_medical_image_type ON medical_image_analysis(image_type);
CREATE INDEX idx_medical_image_severity ON medical_image_analysis(severity);
CREATE INDEX idx_medical_image_date ON medical_image_analysis(created_at DESC);

-- View para receitas analisadas
CREATE OR REPLACE VIEW v_analyzed_prescriptions AS
SELECT 
    mia.id,
    mia.idoso_id,
    i.nome as idoso_nome,
    mia.analysis_result->>'doctor_name' as medico,
    mia.analysis_result->>'doctor_crm' as crm,
    mia.analysis_result->>'prescription_date' as data_receita,
    jsonb_array_length(mia.analysis_result->'medications') as total_medicamentos,
    mia.created_at as data_analise
FROM medical_image_analysis mia
JOIN idosos i ON i.id = mia.idoso_id
WHERE mia.image_type = 'prescription'
ORDER BY mia.created_at DESC;

-- View para feridas analisadas
CREATE OR REPLACE VIEW v_analyzed_wounds AS
SELECT 
    mia.id,
    mia.idoso_id,
    i.nome as idoso_nome,
    mia.analysis_result->>'type' as tipo_lesao,
    mia.analysis_result->>'size' as tamanho,
    mia.severity as gravidade,
    mia.requires_medical_attention as requer_atencao,
    mia.caregiver_notified as cuidador_notificado,
    mia.created_at
FROM medical_image_analysis mia
JOIN idosos i ON i.id = mia.idoso_id
WHERE mia.image_type = 'wound'
ORDER BY mia.created_at DESC;

-- View para alertas pendentes
CREATE OR REPLACE VIEW v_medical_image_alerts AS
SELECT 
    mia.id,
    mia.idoso_id,
    i.nome as idoso_nome,
    mia.image_type,
    mia.severity,
    mia.created_at,
    EXTRACT(EPOCH FROM (NOW() - mia.created_at))/60 as minutes_since_analysis
FROM medical_image_analysis mia
JOIN idosos i ON i.id = mia.idoso_id
WHERE mia.requires_medical_attention = true
  AND mia.caregiver_notified = false
  AND mia.created_at >= NOW() - INTERVAL '24 hours'
ORDER BY mia.severity DESC, mia.created_at DESC;

-- Função para marcar notificação enviada
CREATE OR REPLACE FUNCTION mark_medical_image_notified(analysis_id BIGINT)
RETURNS VOID AS $$
BEGIN
    UPDATE medical_image_analysis
    SET caregiver_notified = true,
        notified_at = NOW()
    WHERE id = analysis_id;
END;
$$ LANGUAGE plpgsql;

-- 🌍 VIEWS PARA EPIDEMIOLOGIA E SAÚDE PÚBLICA

-- View para casos de malária
CREATE OR REPLACE VIEW v_malaria_cases AS
SELECT 
    mia.id,
    mia.idoso_id,
    i.nome as paciente_nome,
    mia.analysis_result->>'parasitemia' as parasitemia,
    mia.analysis_result->>'species' as especie,
    mia.severity as gravidade,
    mia.geolocation,
    mia.created_at as data_diagnostico
FROM medical_image_analysis mia
JOIN idosos i ON i.id = mia.idoso_id
WHERE mia.image_type = 'malaria_smear'
  AND mia.analysis_result->>'result' = 'POSITIVE'
ORDER BY mia.created_at DESC;

-- View para triagem de tuberculose
CREATE OR REPLACE VIEW v_tb_screening AS
SELECT 
    mia.id,
    mia.idoso_id,
    i.nome as paciente_nome,
    mia.analysis_result->>'tb_probability' as probabilidade_tb,
    mia.analysis_result->>'findings' as achados,
    mia.requires_medical_attention as requer_confirmacao,
    mia.geolocation,
    mia.created_at
FROM medical_image_analysis mia
JOIN idosos i ON i.id = mia.idoso_id
WHERE mia.image_type = 'chest_xray'
ORDER BY mia.created_at DESC;

-- View para testes rápidos (COVID, HIV, Dengue)
CREATE OR REPLACE VIEW v_rapid_tests AS
SELECT 
    mia.id,
    mia.idoso_id,
    i.nome as paciente_nome,
    mia.test_metadata->>'test_type' as tipo_teste,
    mia.test_metadata->>'manufacturer' as fabricante,
    mia.analysis_result->>'result' as resultado,
    mia.analysis_result->>'confidence' as confianca,
    mia.geolocation,
    mia.created_at
FROM medical_image_analysis mia
JOIN idosos i ON i.id = mia.idoso_id
WHERE mia.image_type = 'rapid_test'
ORDER BY mia.created_at DESC;

-- View para lesões cutâneas (Mpox, melanoma)
CREATE OR REPLACE VIEW v_skin_lesions AS
SELECT 
    mia.id,
    mia.idoso_id,
    i.nome as paciente_nome,
    mia.analysis_result->>'lesion_type' as tipo_lesao,
    mia.analysis_result->>'mpox_probability' as prob_mpox,
    mia.analysis_result->>'melanoma_risk' as risco_melanoma,
    mia.severity as gravidade,
    mia.requires_medical_attention as requer_atencao,
    mia.created_at
FROM medical_image_analysis mia
JOIN idosos i ON i.id = mia.idoso_id
WHERE mia.image_type = 'skin_lesion'
ORDER BY mia.created_at DESC;

-- View para mapa epidemiológico (heatmap)
CREATE OR REPLACE VIEW v_epidemiological_heatmap AS
SELECT 
    image_type as tipo_caso,
    COUNT(*) as total_casos,
    COUNT(*) FILTER (WHERE severity IN ('ALTO', 'CRÍTICO')) as casos_graves,
    geolocation::TEXT as coordenadas,
    DATE_TRUNC('day', created_at) as data
FROM medical_image_analysis
WHERE geolocation IS NOT NULL
  AND created_at >= NOW() - INTERVAL '30 days'
  AND image_type IN ('malaria_smear', 'chest_xray', 'rapid_test', 'skin_lesion')
GROUP BY image_type, geolocation::TEXT, DATE_TRUNC('day', created_at)
ORDER BY data DESC, total_casos DESC;

-- Índice geoespacial para queries de proximidade
CREATE INDEX idx_medical_image_geolocation ON medical_image_analysis USING GIST(geolocation)
WHERE geolocation IS NOT NULL;

-- Índice para test_metadata
CREATE INDEX idx_medical_image_test_metadata ON medical_image_analysis USING GIN(test_metadata)
WHERE test_metadata IS NOT NULL;

-- Comentários para documentação
COMMENT ON TABLE medical_image_analysis IS 'Análises de imagens médicas realizadas pelo MedGemma - Suporta saúde pública e epidemiologia';
COMMENT ON COLUMN medical_image_analysis.image_type IS 'Tipo de imagem: prescription, wound, malaria_smear, chest_xray, rapid_test, skin_lesion, pressure_ulcer, diabetic_foot, other';
COMMENT ON COLUMN medical_image_analysis.geolocation IS 'Coordenadas GPS para rastreamento epidemiológico (Point: latitude, longitude)';
COMMENT ON COLUMN medical_image_analysis.test_metadata IS 'Metadados do teste: fabricante, lote, tipo de teste, etc';
COMMENT ON COLUMN medical_image_analysis.analysis_result IS 'Resultado completo da análise em formato JSON';
COMMENT ON VIEW v_analyzed_prescriptions IS 'Receitas médicas analisadas com extração de medicamentos';
COMMENT ON VIEW v_analyzed_wounds IS 'Feridas e lesões analisadas com avaliação de gravidade';
COMMENT ON VIEW v_medical_image_alerts IS 'Alertas de imagens médicas que requerem atenção e ainda não foram notificados';
COMMENT ON VIEW v_malaria_cases IS 'Casos positivos de malária detectados via análise de esfregaço';
COMMENT ON VIEW v_tb_screening IS 'Triagens de tuberculose via raio-X de tórax';
COMMENT ON VIEW v_rapid_tests IS 'Resultados de testes rápidos (COVID, HIV, Dengue, etc)';
COMMENT ON VIEW v_skin_lesions IS 'Lesões cutâneas analisadas (Mpox, melanoma, etc)';
COMMENT ON VIEW v_epidemiological_heatmap IS 'Dados agregados para mapa epidemiológico (heatmap) - últimos 30 dias';
