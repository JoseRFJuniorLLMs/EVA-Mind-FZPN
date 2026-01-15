-- ================================================
-- Migração v23: Catálogo de Vozes e Troca Dinâmica
-- Data: 2026-01-15
-- Objetivo: Permitir troca de voz em tempo real
-- ================================================

-- Tabela de vozes disponíveis (catálogo)
CREATE TABLE IF NOT EXISTS eva_voices (
    id SERIAL PRIMARY KEY,
    voice_name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL, -- Nome amigável
    gender VARCHAR(20), -- 'feminina', 'masculina', 'neutra'
    tone VARCHAR(50), -- 'calorosa', 'suave', 'enérgica', 'grave'
    language VARCHAR(10) DEFAULT 'pt-BR',
    provider VARCHAR(50) DEFAULT 'gemini',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Popular com vozes do Gemini (21 vozes disponíveis)
INSERT INTO eva_voices (voice_name, display_name, gender, tone) VALUES
-- As 5 Clássicas (Padrão)
('Aoede', 'Aoede - Calorosa', 'feminina', 'calorosa'),
('Charon', 'Charon - Grave', 'masculina', 'grave'),
('Fenrir', 'Fenrir - Enérgica', 'masculina', 'enérgica'),
('Kore', 'Kore - Suave', 'feminina', 'suave'),
('Puck', 'Puck - Animada', 'neutra', 'animada'),

-- As Novas (Gemini 2.5 / Live API) - Femininas
('Zephyr', 'Zephyr - Brilhante', 'feminina', 'brilhante'),
('Leda', 'Leda - Jovem', 'feminina', 'energética'),
('Callirrhoe', 'Callirrhoe - Descontraída', 'feminina', 'descontraída'),
('Autonoe', 'Autonoe - Acolhedora', 'feminina', 'acolhedora'),
('Despina', 'Despina - Aveludada', 'feminina', 'aveludada'),
('Erinome', 'Erinome - Professoral', 'feminina', 'séria'),
('Laomedeia', 'Laomedeia - Otimista', 'feminina', 'otimista'),

-- As Novas (Gemini 2.5 / Live API) - Masculinas
('Orus', 'Orus - Seguro', 'masculina', 'firme'),
('Enceladus', 'Enceladus - Sussurrada', 'masculina', 'suave'),
('Iapetus', 'Iapetus - Direto', 'masculina', 'claro'),
('Umbriel', 'Umbriel - Calmo', 'masculina', 'calmo'),
('Algieba', 'Algieba - Macia', 'masculina', 'suave'),
('Algenib', 'Algenib - Profunda', 'masculina', 'áspera'),
('Rasalgethi', 'Rasalgethi - Informativa', 'masculina', 'formal'),
('Alnilam', 'Alnilam - Firme', 'masculina', 'séria');

-- Adicionar campo de preferência de voz na tabela idosos
ALTER TABLE idosos ADD COLUMN IF NOT EXISTS preferred_voice VARCHAR(50) DEFAULT 'Aoede';
ALTER TABLE idosos ADD CONSTRAINT fk_idosos_voice FOREIGN KEY (preferred_voice) REFERENCES eva_voices(voice_name);

-- Tabela de histórico de trocas de voz (analytics)
CREATE TABLE IF NOT EXISTS voice_change_history (
    id SERIAL PRIMARY KEY,
    idoso_id INT NOT NULL REFERENCES idosos(id) ON DELETE CASCADE,
    old_voice VARCHAR(50),
    new_voice VARCHAR(50) NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT NOW(),
    change_method VARCHAR(50) -- 'voice_command', 'app_settings', 'tool_call'
);

CREATE INDEX idx_voice_history_idoso ON voice_change_history(idoso_id);

-- View para estatísticas de vozes mais populares
CREATE OR REPLACE VIEW voice_popularity AS
SELECT 
    v.voice_name,
    v.display_name,
    COUNT(DISTINCT i.id) as users_count,
    v.gender,
    v.tone
FROM eva_voices v
LEFT JOIN idosos i ON i.preferred_voice = v.voice_name
WHERE v.is_active = true
GROUP BY v.voice_name, v.display_name, v.gender, v.tone
ORDER BY users_count DESC;

-- Função para obter voz aleatória (para demonstração)
CREATE OR REPLACE FUNCTION get_random_voice()
RETURNS VARCHAR(50) AS $$
DECLARE
    random_voice VARCHAR(50);
BEGIN
    SELECT voice_name INTO random_voice
    FROM eva_voices
    WHERE is_active = true
    ORDER BY RANDOM()
    LIMIT 1;
    
    RETURN random_voice;
END;
$$ LANGUAGE plpgsql;

-- Comentários
COMMENT ON TABLE eva_voices IS 'Catálogo de vozes disponíveis para EVA';
COMMENT ON TABLE voice_change_history IS 'Histórico de trocas de voz por idoso';
COMMENT ON COLUMN idosos.preferred_voice IS 'Voz preferida do idoso (pode ser alterada em tempo real)';

COMMIT;
