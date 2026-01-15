-- Migration v27: Family Voices
-- Tabela para gerenciar vozes clonadas dos familiares

CREATE TABLE IF NOT EXISTS family_voices (
    id BIGSERIAL PRIMARY KEY,
    idoso_id BIGINT NOT NULL REFERENCES idosos(id) ON DELETE CASCADE,
    familiar_nome VARCHAR(100) NOT NULL,
    familiar_tipo VARCHAR(50) NOT NULL, -- 'filha', 'filho', 'neta', 'neto', 'esposa', 'esposo'
    emotion_type VARCHAR(20) NOT NULL CHECK (emotion_type IN ('serious', 'loving', 'happy')),
    voice_id VARCHAR(200) NOT NULL UNIQUE, -- ID retornado pelo EVA-Voice
    file_size_kb DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Índices
CREATE INDEX idx_family_voices_idoso ON family_voices(idoso_id);
CREATE INDEX idx_family_voices_tipo ON family_voices(familiar_tipo);
CREATE INDEX idx_family_voices_emotion ON family_voices(emotion_type);

-- View: Vozes disponíveis por idoso
CREATE OR REPLACE VIEW v_available_voices AS
SELECT 
    i.id as idoso_id,
    i.nome as idoso_nome,
    fv.familiar_nome,
    fv.familiar_tipo,
    fv.emotion_type,
    fv.voice_id,
    fv.created_at
FROM idosos i
LEFT JOIN family_voices fv ON fv.idoso_id = i.id
ORDER BY i.id, fv.familiar_tipo, fv.emotion_type;

-- Comentários
COMMENT ON TABLE family_voices IS 'Vozes clonadas dos familiares para cada idoso';
COMMENT ON COLUMN family_voices.emotion_type IS 'Tipo de emoção: serious (séria), loving (carinhosa), happy (animada)';
COMMENT ON COLUMN family_voices.voice_id IS 'ID da voz no microserviço EVA-Voice';
