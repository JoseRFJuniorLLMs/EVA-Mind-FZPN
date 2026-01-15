-- Migration v25: A/B Testing Framework
-- Sistema de testes A/B para comparar Thinking Mode vs Normal Mode

-- Tabela de configuração de testes A/B
CREATE TABLE IF NOT EXISTS ab_test_config (
    id SERIAL PRIMARY KEY,
    test_name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    enabled BOOLEAN DEFAULT true,
    percentage_group_a INT DEFAULT 50 CHECK (percentage_group_a >= 0 AND percentage_group_a <= 100),
    group_a_name VARCHAR(50) DEFAULT 'thinking_mode',
    group_b_name VARCHAR(50) DEFAULT 'normal_mode',
    started_at TIMESTAMP DEFAULT NOW(),
    ended_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Tabela de atribuição de usuários a grupos
CREATE TABLE IF NOT EXISTS ab_test_assignments (
    id BIGSERIAL PRIMARY KEY,
    test_name VARCHAR(100) NOT NULL REFERENCES ab_test_config(test_name),
    idoso_id BIGINT NOT NULL REFERENCES idosos(id) ON DELETE CASCADE,
    assigned_group VARCHAR(50) NOT NULL,
    assigned_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(test_name, idoso_id)
);

-- Tabela de métricas de A/B testing
CREATE TABLE IF NOT EXISTS ab_test_metrics (
    id BIGSERIAL PRIMARY KEY,
    test_name VARCHAR(100) NOT NULL REFERENCES ab_test_config(test_name),
    idoso_id BIGINT NOT NULL REFERENCES idosos(id) ON DELETE CASCADE,
    assigned_group VARCHAR(50) NOT NULL,
    metric_type VARCHAR(100) NOT NULL,
    metric_value DECIMAL(10,4),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Índices
CREATE INDEX idx_ab_assignments_test ON ab_test_assignments(test_name);
CREATE INDEX idx_ab_assignments_idoso ON ab_test_assignments(idoso_id);
CREATE INDEX idx_ab_metrics_test ON ab_test_metrics(test_name);
CREATE INDEX idx_ab_metrics_group ON ab_test_metrics(assigned_group);
CREATE INDEX idx_ab_metrics_type ON ab_test_metrics(metric_type);
CREATE INDEX idx_ab_metrics_date ON ab_test_metrics(created_at DESC);

-- View: Resumo de performance por grupo
CREATE OR REPLACE VIEW v_ab_test_performance AS
SELECT 
    m.test_name,
    m.assigned_group as grupo,
    m.metric_type as metrica,
    COUNT(*) as total_amostras,
    AVG(m.metric_value) as media,
    STDDEV(m.metric_value) as desvio_padrao,
    MIN(m.metric_value) as minimo,
    MAX(m.metric_value) as maximo,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY m.metric_value) as mediana
FROM ab_test_metrics m
GROUP BY m.test_name, m.assigned_group, m.metric_type
ORDER BY m.test_name, m.metric_type, m.assigned_group;

-- View: Comparação direta entre grupos
CREATE OR REPLACE VIEW v_ab_test_comparison AS
SELECT 
    a.test_name,
    a.metrica as metric_type,
    a.grupo as grupo_a,
    a.media as media_a,
    a.total_amostras as amostras_a,
    b.grupo as grupo_b,
    b.media as media_b,
    b.total_amostras as amostras_b,
    ROUND(((b.media - a.media) / NULLIF(a.media, 0) * 100)::numeric, 2) as diferenca_percentual,
    CASE 
        WHEN b.media > a.media THEN 'Grupo B Melhor'
        WHEN b.media < a.media THEN 'Grupo A Melhor'
        ELSE 'Empate'
    END as vencedor
FROM v_ab_test_performance a
JOIN v_ab_test_performance b 
    ON a.test_name = b.test_name 
    AND a.metrica = b.metrica 
    AND a.grupo < b.grupo
ORDER BY a.test_name, a.metrica;

-- View: Distribuição de usuários por grupo
CREATE OR REPLACE VIEW v_ab_test_distribution AS
SELECT 
    test_name,
    assigned_group as grupo,
    COUNT(*) as total_usuarios,
    ROUND(100.0 * COUNT(*) / SUM(COUNT(*)) OVER (PARTITION BY test_name), 2) as percentual
FROM ab_test_assignments
GROUP BY test_name, assigned_group
ORDER BY test_name, assigned_group;

-- Função: Atribuir usuário a grupo de teste
CREATE OR REPLACE FUNCTION assign_ab_test_group(
    p_test_name VARCHAR(100),
    p_idoso_id BIGINT
) RETURNS VARCHAR(50) AS $$
DECLARE
    v_group VARCHAR(50);
    v_percentage INT;
    v_group_a_name VARCHAR(50);
    v_group_b_name VARCHAR(50);
    v_hash_value INT;
BEGIN
    -- Verificar se já existe atribuição
    SELECT assigned_group INTO v_group
    FROM ab_test_assignments
    WHERE test_name = p_test_name AND idoso_id = p_idoso_id;
    
    IF FOUND THEN
        RETURN v_group;
    END IF;
    
    -- Buscar configuração do teste
    SELECT percentage_group_a, group_a_name, group_b_name
    INTO v_percentage, v_group_a_name, v_group_b_name
    FROM ab_test_config
    WHERE test_name = p_test_name AND enabled = true;
    
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Teste A/B % não encontrado ou desabilitado', p_test_name;
    END IF;
    
    -- Hash determinístico baseado no ID
    v_hash_value := (hashtext(p_test_name || '_' || p_idoso_id::TEXT)::BIGINT % 100);
    
    -- Atribuir grupo
    IF v_hash_value < v_percentage THEN
        v_group := v_group_a_name;
    ELSE
        v_group := v_group_b_name;
    END IF;
    
    -- Salvar atribuição
    INSERT INTO ab_test_assignments (test_name, idoso_id, assigned_group)
    VALUES (p_test_name, p_idoso_id, v_group)
    ON CONFLICT (test_name, idoso_id) DO NOTHING;
    
    RETURN v_group;
END;
$$ LANGUAGE plpgsql;

-- Função: Registrar métrica de A/B test
CREATE OR REPLACE FUNCTION log_ab_test_metric(
    p_test_name VARCHAR(100),
    p_idoso_id BIGINT,
    p_metric_type VARCHAR(100),
    p_metric_value DECIMAL(10,4),
    p_metadata JSONB DEFAULT NULL
) RETURNS VOID AS $$
DECLARE
    v_group VARCHAR(50);
BEGIN
    -- Obter grupo do usuário
    v_group := assign_ab_test_group(p_test_name, p_idoso_id);
    
    -- Registrar métrica
    INSERT INTO ab_test_metrics (
        test_name,
        idoso_id,
        assigned_group,
        metric_type,
        metric_value,
        metadata
    ) VALUES (
        p_test_name,
        p_idoso_id,
        v_group,
        p_metric_type,
        p_metric_value,
        p_metadata
    );
END;
$$ LANGUAGE plpgsql;

-- Inserir teste padrão: Thinking Mode vs Normal Mode
INSERT INTO ab_test_config (
    test_name,
    description,
    percentage_group_a,
    group_a_name,
    group_b_name
) VALUES (
    'health_triage_mode',
    'Comparação entre Thinking Mode e Normal Mode para triagem de saúde',
    50,
    'thinking_mode',
    'normal_mode'
) ON CONFLICT (test_name) DO NOTHING;

-- Comentários
COMMENT ON TABLE ab_test_config IS 'Configuração de testes A/B';
COMMENT ON TABLE ab_test_assignments IS 'Atribuição de usuários a grupos de teste';
COMMENT ON TABLE ab_test_metrics IS 'Métricas coletadas durante testes A/B';
COMMENT ON FUNCTION assign_ab_test_group IS 'Atribui usuário a grupo de teste de forma determinística';
COMMENT ON FUNCTION log_ab_test_metric IS 'Registra métrica de A/B test';
COMMENT ON VIEW v_ab_test_performance IS 'Performance agregada por grupo de teste';
COMMENT ON VIEW v_ab_test_comparison IS 'Comparação direta entre grupos A e B';
COMMENT ON VIEW v_ab_test_distribution IS 'Distribuição de usuários por grupo';
