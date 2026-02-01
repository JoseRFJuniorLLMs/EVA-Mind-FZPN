-- ============================================================================
-- MIGRATION: Adiciona suporte a idioma por idoso
-- Sistema internacional - EVA detecta idioma do usuário
-- ============================================================================

-- Adicionar coluna idioma na tabela idosos
ALTER TABLE idosos ADD COLUMN IF NOT EXISTS idioma VARCHAR(10) DEFAULT 'pt-BR';

-- Criar índice para consultas rápidas
CREATE INDEX IF NOT EXISTS idx_idosos_idioma ON idosos(idioma);

-- Comentário
COMMENT ON COLUMN idosos.idioma IS 'Idioma preferido do idoso (pt-BR, en-US, es-ES, fr-FR, etc.)';

-- ============================================================================
-- IDIOMAS SUPORTADOS:
-- pt-BR = Português Brasileiro (default)
-- pt-PT = Português de Portugal
-- en-US = English (US)
-- en-GB = English (UK)
-- es-ES = Español
-- es-MX = Español (México)
-- fr-FR = Français
-- de-DE = Deutsch
-- it-IT = Italiano
-- ja-JP = 日本語
-- zh-CN = 中文
-- ============================================================================

SELECT 'Migration 025_idioma_idoso complete!' AS status;
