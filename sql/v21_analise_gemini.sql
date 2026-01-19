CREATE TABLE IF NOT EXISTS analise_gemini (
    id SERIAL PRIMARY KEY,
    idoso_id INTEGER NOT NULL, -- FK to usuarios/idosos handled optionally
    tipo VARCHAR(50) NOT NULL, -- 'AUDIO', 'GRAPH', 'TEXT'
    conteudo JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_analise_gemini_idoso_created ON analise_gemini(idoso_id, created_at DESC);
