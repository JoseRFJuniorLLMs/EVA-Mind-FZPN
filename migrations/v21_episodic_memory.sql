-- ================================================
-- Migração v21: Memória Episódica (RAG)
-- Data: 2026-01-15
-- Objetivo: Adicionar sistema de memória semântica
-- ================================================

-- 1. Instalar extensão pgvector
CREATE EXTENSION IF NOT EXISTS vector;

-- 2. Criar tabela de memórias episódicas
CREATE TABLE episodic_memories (
    id SERIAL PRIMARY KEY,
    idoso_id INT NOT NULL REFERENCES idosos(id) ON DELETE CASCADE,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Conteúdo da memória
    speaker VARCHAR(20) NOT NULL CHECK (speaker IN ('user', 'assistant')),
    content TEXT NOT NULL,
    
    -- Embedding vetorial (1536 dimensões - Gemini text-embedding-004)
    embedding vector(1536),
    
    -- Metadados semânticos
    emotion VARCHAR(50), -- 'feliz', 'triste', 'neutro', 'ansioso', 'confuso'
    importance FLOAT DEFAULT 0.5 CHECK (importance BETWEEN 0.0 AND 1.0),
    topics TEXT[], -- ['saúde', 'família', 'medicamento', 'lazer']
    
    -- Contexto da conversa
    session_id VARCHAR(100),
    call_history_id INT REFERENCES historico_ligacoes(id) ON DELETE SET NULL,
    
    -- Auditoria
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Índices para busca semântica (IVFFlat)
-- IVFFlat é otimizado para datasets grandes (>10k vetores)
-- lists = sqrt(total_rows) é uma boa heurística inicial
CREATE INDEX idx_episodic_embedding ON episodic_memories 
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

-- 4. Índices convencionais para queries frequentes
CREATE INDEX idx_episodic_idoso_time ON episodic_memories(idoso_id, timestamp DESC);
CREATE INDEX idx_episodic_importance ON episodic_memories(importance DESC) WHERE importance > 0.7;
CREATE INDEX idx_episodic_session ON episodic_memories(session_id) WHERE session_id IS NOT NULL;

-- 5. Trigger para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_episodic_memory_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_episodic_memory_update
    BEFORE UPDATE ON episodic_memories
    FOR EACH ROW
    EXECUTE FUNCTION update_episodic_memory_timestamp();

-- 6. View para memórias recentes (útil para debug)
CREATE OR REPLACE VIEW recent_memories AS
SELECT 
    em.id,
    i.nome as idoso_nome,
    em.speaker,
    LEFT(em.content, 100) as content_preview,
    em.emotion,
    em.importance,
    em.topics,
    em.timestamp
FROM episodic_memories em
JOIN idosos i ON em.idoso_id = i.id
ORDER BY em.timestamp DESC
LIMIT 100;

-- 7. Função auxiliar para busca semântica
-- Retorna as K memórias mais similares a uma query
CREATE OR REPLACE FUNCTION search_similar_memories(
    p_idoso_id INT,
    p_query_embedding vector(1536),
    p_limit INT DEFAULT 5,
    p_min_similarity FLOAT DEFAULT 0.5
)
RETURNS TABLE (
    memory_id INT,
    content TEXT,
    speaker VARCHAR(20),
    memory_timestamp TIMESTAMPTZ,
    emotion VARCHAR(50),
    importance FLOAT,
    topics TEXT[],
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        em.id,
        em.content,
        em.speaker,
        em.timestamp,
        em.emotion,
        em.importance,
        em.topics,
        1 - (em.embedding <=> p_query_embedding) as sim
    FROM episodic_memories em
    WHERE em.idoso_id = p_idoso_id
      AND (1 - (em.embedding <=> p_query_embedding)) >= p_min_similarity
    ORDER BY em.embedding <=> p_query_embedding
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;

-- 8. Comentários para documentação
COMMENT ON TABLE episodic_memories IS 'Armazena memórias conversacionais para sistema RAG (Retrieval-Augmented Generation)';
COMMENT ON COLUMN episodic_memories.embedding IS 'Vetor de 1536 dimensões gerado por Gemini text-embedding-004';
COMMENT ON COLUMN episodic_memories.importance IS 'Score de importância: 0.0 (trivial) a 1.0 (crítico)';
COMMENT ON COLUMN episodic_memories.topics IS 'Tópicos extraídos via LLM (ex: saúde, família, medicação)';

-- 9. Dados de exemplo (para testes)
-- Comentar em produção
/*
INSERT INTO episodic_memories (idoso_id, speaker, content, emotion, importance, topics)
VALUES 
    (1, 'user', 'Meu neto ganhou uma bicicleta vermelha no aniversário', 'feliz', 0.8, ARRAY['família', 'celebração']),
    (1, 'assistant', 'Que legal! Ele deve ter adorado. Qual a idade dele?', 'neutro', 0.5, ARRAY['família']),
    (1, 'user', 'Ele fez 8 anos', 'feliz', 0.7, ARRAY['família', 'aniversário']);
*/

COMMIT;
