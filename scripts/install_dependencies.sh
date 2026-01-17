#!/bin/bash
# Install Python dependencies for EVA-Mind-FZPN scripts

echo "🔧 Instalando dependências Python para EVA-Mind-FZPN..."
echo ""

# Atualizar pip
echo "📦 Atualizando pip..."
python3 -m pip install --upgrade pip --break-system-packages

# Instalar dependências
echo ""
echo "📥 Instalando bibliotecas..."

pip3 install qdrant-client --break-system-packages
pip3 install requests --break-system-packages

echo ""
echo "✅ Dependências instaladas com sucesso!"
echo ""
echo "Bibliotecas instaladas:"
echo "  • qdrant-client (para Qdrant vector DB)"
echo "  • requests (para HTTP requests)"
echo ""
