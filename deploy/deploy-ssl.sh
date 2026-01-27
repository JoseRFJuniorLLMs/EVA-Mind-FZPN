#!/bin/bash
# ============================================
# EVA - Script de Deploy com SSL
# Data: 2026-01-27
# ============================================
#
# COMO USAR:
# 1. No servidor: cd /var/www/eva-mind
# 2. git pull
# 3. chmod +x deploy/deploy-ssl.sh
# 4. ./deploy/deploy-ssl.sh
# ============================================

set -e

echo "============================================"
echo "🔒 EVA - Deploy com SSL"
echo "============================================"

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 1. Parar servicos
echo -e "${YELLOW}⏹️  Parando servicos...${NC}"
pkill -f "python.*main.py" 2>/dev/null || true
pkill -f "python.*scheduler_api" 2>/dev/null || true
pkill -f "eva-mind-fzpn" 2>/dev/null || true
sleep 2

# 2. Configurar nginx
echo -e "${YELLOW}🔧 Configurando nginx...${NC}"
cp deploy/nginx-ssl.conf /etc/nginx/sites-available/eva-ssl
ln -sf /etc/nginx/sites-available/eva-ssl /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Testar e recarregar nginx
nginx -t
if [ $? -eq 0 ]; then
    systemctl reload nginx
    echo -e "${GREEN}✅ Nginx configurado!${NC}"
else
    echo -e "${RED}❌ Erro na configuracao do nginx!${NC}"
    exit 1
fi

# 3. Build do Go (EVA-Mind)
echo -e "${YELLOW}🔨 Compilando EVA-Mind...${NC}"
cd /var/www/eva-mind
go build -o eva-mind-fzpn .
echo -e "${GREEN}✅ Build concluido!${NC}"

# 4. Iniciar EVA-Mind na porta interna
echo -e "${YELLOW}🚀 Iniciando EVA-Mind na porta 8091...${NC}"
export PORT=8091
nohup ./eva-mind-fzpn > /var/log/eva-mind.log 2>&1 &
sleep 2

# 5. Iniciar Python Backend (se existir)
if [ -d "/root/EVA-back" ]; then
    echo -e "${YELLOW}🚀 Iniciando Python Backend na porta 8001...${NC}"
    cd /root/EVA-back/eva-enterprise
    export PORT=8001
    nohup python main.py > /var/log/eva-back.log 2>&1 &
    sleep 2
fi

# 6. Verificar
echo ""
echo -e "${YELLOW}🔍 Verificando servicos...${NC}"
echo ""

echo "Portas internas (backends):"
ss -tlnp | grep -E "8001|8091" || echo "  Nenhum backend ativo"

echo ""
echo "Portas publicas (nginx com SSL):"
ss -tlnp | grep -E "443|8000|8090" | grep nginx || echo "  Nginx nao esta nas portas SSL"

echo ""
echo "============================================"
echo -e "${GREEN}✅ Deploy concluido!${NC}"
echo ""
echo "URLs disponiveis:"
echo "  - https://eva-ia.org (Frontend)"
echo "  - https://eva-ia.org:8000/api/v1 (API Python)"
echo "  - wss://eva-ia.org:8090/ws/video (WebSocket Video)"
echo "  - wss://eva-ia.org:8090/ws/pcm (WebSocket Audio)"
echo ""
echo "Logs:"
echo "  - tail -f /var/log/eva-mind.log"
echo "  - tail -f /var/log/eva-back.log"
echo "============================================"
