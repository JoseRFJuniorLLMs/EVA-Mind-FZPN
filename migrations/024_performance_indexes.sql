-- ============================================================================
-- PERFORMANCE MIGRATION: Indices para otimizacao de queries
-- Issue: Full table scans em queries frequentes
-- Fix: Criar indices compostos para colunas mais usadas
-- Impacto esperado: 3x melhoria em tempo de query
-- ============================================================================

-- ============================================================================
-- TABELA: idosos
-- Queries frequentes: by-cpf, by-id
-- ============================================================================

-- Indice para busca por CPF (usado em login, validacao)
CREATE INDEX IF NOT EXISTS idx_idosos_cpf
ON idosos(cpf);

-- Indice para busca por telefone (usado em chamadas)
CREATE INDEX IF NOT EXISTS idx_idosos_telefone
ON idosos(telefone);

-- ============================================================================
-- TABELA: memories (memorias do paciente)
-- Queries frequentes: GetRecent, por idoso_id + timestamp
-- ============================================================================

-- Indice composto para busca de memorias recentes por paciente
CREATE INDEX IF NOT EXISTS idx_memories_idoso_timestamp
ON memories(idoso_id, created_at DESC);

-- Indice para busca por tipo de memoria
CREATE INDEX IF NOT EXISTS idx_memories_type
ON memories(memory_type);

-- ============================================================================
-- TABELA: transcriptions (transcricoes de audio)
-- Queries frequentes: por sessao, por idoso, por data
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_transcriptions_idoso_timestamp
ON transcriptions(idoso_id, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_transcriptions_session
ON transcriptions(session_id);

-- ============================================================================
-- TABELA: interaction_cognitive_load (carga cognitiva)
-- Queries frequentes: ultimas 24h por paciente
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_cognitive_load_patient_time
ON interaction_cognitive_load(patient_id, timestamp DESC);

-- ============================================================================
-- TABELA: ethical_boundary_events (eventos eticos)
-- Queries frequentes: por paciente, por severidade
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_ethical_events_patient_severity
ON ethical_boundary_events(patient_id, severity, timestamp DESC);

-- ============================================================================
-- TABELA: clinical_decision_explanations (explicacoes clinicas)
-- Queries frequentes: por paciente, tipo de decisao
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_clinical_decisions_patient
ON clinical_decision_explanations(patient_id, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_clinical_decisions_type
ON clinical_decision_explanations(decision_type);

-- ============================================================================
-- TABELA: trajectory_simulations (simulacoes de trajetoria)
-- Queries frequentes: por paciente, data de simulacao
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_trajectory_patient_date
ON trajectory_simulations(patient_id, simulation_date DESC);

-- ============================================================================
-- TABELA: persona_sessions (sessoes de persona)
-- Queries frequentes: por paciente, tipo de persona
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_persona_sessions_patient
ON persona_sessions(patient_id, started_at DESC);

-- ============================================================================
-- TABELA: exit_protocols (protocolos de saida)
-- Queries frequentes: por paciente, fase atual
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_exit_protocols_patient
ON exit_protocols(patient_id, current_phase);

-- ============================================================================
-- TABELA: deep_memory_events (eventos de memoria profunda)
-- Queries frequentes: por paciente, tipo de evento
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_deep_memory_patient_type
ON deep_memory_events(patient_id, event_type, detected_at DESC);

-- ============================================================================
-- TABELA: device_tokens (tokens FCM)
-- Queries frequentes: por idoso_id
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_device_tokens_idoso
ON device_tokens(idoso_id);

-- ============================================================================
-- TABELA: enneagram_results (resultados Eneagrama)
-- Queries frequentes: por idoso_id
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_enneagram_idoso
ON enneagram_results(idoso_id);

-- ============================================================================
-- PARTIAL INDEXES (para queries com filtros especificos)
-- ============================================================================

-- Indice parcial para eventos eticos de alta severidade (mais rapido para alertas)
CREATE INDEX IF NOT EXISTS idx_ethical_high_severity
ON ethical_boundary_events(patient_id, timestamp DESC)
WHERE severity IN ('high', 'critical');

-- Indice parcial para memorias importantes (starred/flagged)
CREATE INDEX IF NOT EXISTS idx_memories_important
ON memories(idoso_id, created_at DESC)
WHERE is_important = true;

-- ============================================================================
-- ANALYZE para atualizar estatisticas do planner
-- ============================================================================

ANALYZE idosos;
ANALYZE memories;
ANALYZE transcriptions;
ANALYZE interaction_cognitive_load;
ANALYZE ethical_boundary_events;
ANALYZE clinical_decision_explanations;
ANALYZE trajectory_simulations;
ANALYZE persona_sessions;

-- ============================================================================
-- COMENTARIO FINAL
-- Executar: psql -U postgres -d eva -f 024_performance_indexes.sql
-- Ou via make migrate
-- ============================================================================
