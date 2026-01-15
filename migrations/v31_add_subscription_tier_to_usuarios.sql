-- Migration v31: Add subscription_tier to usuarios table
-- Adiciona campo de tier de assinatura na tabela de usuários

-- Adicionar coluna subscription_tier
ALTER TABLE usuarios
ADD COLUMN IF NOT EXISTS subscription_tier VARCHAR(20) DEFAULT 'basic' 
CHECK (subscription_tier IN ('basic', 'gold', 'diamond'));

-- Criar índice para queries de tier
CREATE INDEX IF NOT EXISTS idx_usuarios_subscription_tier 
ON usuarios(subscription_tier);

-- Atualizar usuários existentes para tier basic
UPDATE usuarios
SET subscription_tier = 'basic'
WHERE subscription_tier IS NULL;

-- Comentário
COMMENT ON COLUMN usuarios.subscription_tier IS 'Tier da assinatura: basic (gratuito), gold, diamond';
