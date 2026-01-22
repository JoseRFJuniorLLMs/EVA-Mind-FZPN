-- migrations/004_fix_embedding_dimension.sql
-- ============================================
-- Correção crítica: Ajustar dimensão de 1536 para 768
-- Modelo: text-embedding-004 (Gemini)
-- Data: 2026-01-20
-- ============================================

BEGIN;

-- Passo 1: Criar nova coluna com dimensão correta
ALTER TABLE episodic_memories 
ADD COLUMN embedding_new vector(768);

COMMENT ON COLUMN episodic_memories.embedding_new IS 
'Embeddings do Gemini text-embedding-004 (768 dimensões)';

-- Passo 2: Marcar dados antigos para re-processamento
-- (dados com dimensão errada são inválidos)
UPDATE episodic_memories 
SET embedding_new = NULL;

-- Passo 3: Remover coluna antiga
ALTER TABLE episodic_memories 
DROP COLUMN embedding;

-- Passo 4: Renomear nova coluna
ALTER TABLE episodic_memories 
RENAME COLUMN embedding_new TO embedding;

-- Passo 5: Atualizar função de busca
DROP FUNCTION IF EXISTS search_similar_memories(BIGINT, vector, INT, FLOAT);

CREATE OR REPLACE FUNCTION search_similar_memories(
    p_idoso_id BIGINT,
    p_query_embedding vector(768),  -- ✅ Dimensão corrigida
    p_limit INT DEFAULT 10,
    p_min_similarity FLOAT DEFAULT 0.5
)
RETURNS TABLE (
    id BIGINT,
    content TEXT,
    speaker TEXT,
    timestamp TIMESTAMPTZ,
    emotion TEXT,
    importance FLOAT,
    topics TEXT[], -- ✅ Alterado para TEXT[] para consistência com schema
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        em.id,
        em.content,
        em.speaker::TEXT,
        em.timestamp,
        em.emotion::TEXT,
        em.importance,
        em.topics,
        1 - (em.embedding <=> p_query_embedding) AS similarity
    FROM episodic_memories em
    WHERE em.idoso_id = p_idoso_id
      AND em.embedding IS NOT NULL
      AND (1 - (em.embedding <=> p_query_embedding)) >= p_min_similarity
    ORDER BY em.embedding <=> p_query_embedding
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;

COMMIT;

-- Verificação
SELECT 
    COUNT(*) as total_memories,
    COUNT(embedding) as with_embedding,
    COUNT(*) - COUNT(embedding) as need_reembedding
FROM episodic_memories;
