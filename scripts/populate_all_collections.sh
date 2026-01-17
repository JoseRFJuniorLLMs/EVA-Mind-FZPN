#!/bin/bash
# Script Master para Popular Todas as Collections do EVA-Mind-FZPN
# Executa: Nasrudin + Esopo + Zen (Koans + Somático)

set -e  # Para em caso de erro

echo "======================================================================="
echo "🚀 EVA-MIND-FZPN - POPULAÇÃO COMPLETA DO QDRANT"
echo "======================================================================="
echo ""

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Verificar se está no diretório correto
if [ ! -f "scripts/populate_nasrudin_with_lacan.py" ]; then
    echo -e "${RED}❌ Erro: Execute este script do diretório raiz EVA-Mind-FZPN${NC}"
    exit 1
fi

# Verificar se Qdrant está rodando
echo -e "${BLUE}🔍 Verificando Qdrant...${NC}"
if ! curl -s http://localhost:6333/collections > /dev/null; then
    echo -e "${RED}❌ Qdrant não está rodando em localhost:6333${NC}"
    echo "   Inicie o Qdrant primeiro: sudo systemctl start qdrant"
    exit 1
fi
echo -e "${GREEN}✅ Qdrant está online${NC}"
echo ""

# Verificar se Ollama está rodando
echo -e "${BLUE}🔍 Verificando Ollama...${NC}"
if ! curl -s http://localhost:11434/api/tags > /dev/null; then
    echo -e "${RED}❌ Ollama não está rodando em localhost:11434${NC}"
    echo "   Inicie o Ollama primeiro: ollama serve"
    exit 1
fi
echo -e "${GREEN}✅ Ollama está online${NC}"
echo ""

# Verificar modelo de embedding
echo -e "${BLUE}🔍 Verificando modelo nomic-embed-text...${NC}"
if ! ollama list | grep -q "nomic-embed-text"; then
    echo -e "${YELLOW}⚠️  Modelo nomic-embed-text não encontrado${NC}"
    echo "   Baixando modelo..."
    ollama pull nomic-embed-text
fi
echo -e "${GREEN}✅ Modelo nomic-embed-text disponível${NC}"
echo ""

# Função para executar script Python com tratamento de erro
run_script() {
    local script_name=$1
    local description=$2
    
    echo "======================================================================="
    echo -e "${BLUE}📥 $description${NC}"
    echo "======================================================================="
    
    if python3 "scripts/$script_name"; then
        echo -e "${GREEN}✅ $description - CONCLUÍDO${NC}"
        echo ""
        return 0
    else
        echo -e "${RED}❌ $description - FALHOU${NC}"
        echo "   Verifique os logs acima para detalhes"
        return 1
    fi
}

# Contador de sucesso
SUCCESS_COUNT=0
TOTAL_COUNT=3

# 1. Popular Nasrudin Stories
if run_script "populate_nasrudin_with_lacan.py" "Populando Histórias de Nasrudin (270 histórias)"; then
    ((SUCCESS_COUNT++))
fi

# 2. Popular Aesop Fables
if run_script "populate_aesop_fables.py" "Populando Fábulas de Esopo (~300 fábulas)"; then
    ((SUCCESS_COUNT++))
fi

# 3. Popular Zen Content
if run_script "populate_zen_content.py" "Populando Conteúdo Zen (Koans + Exercícios Somáticos)"; then
    ((SUCCESS_COUNT++))
fi

# Resumo Final
echo "======================================================================="
echo -e "${BLUE}📊 RESUMO FINAL${NC}"
echo "======================================================================="
echo ""

# Consultar Qdrant para estatísticas
echo "Consultando collections no Qdrant..."
echo ""

for collection in "nasrudin_stories" "aesop_fables" "zen_koans" "somatic_exercises"; do
    if curl -s "http://localhost:6333/collections/$collection" | grep -q "points_count"; then
        points=$(curl -s "http://localhost:6333/collections/$collection" | grep -o '"points_count":[0-9]*' | grep -o '[0-9]*')
        echo -e "${GREEN}✅ $collection: $points pontos${NC}"
    else
        echo -e "${YELLOW}⚠️  $collection: não encontrada${NC}"
    fi
done

echo ""
echo "======================================================================="

if [ $SUCCESS_COUNT -eq $TOTAL_COUNT ]; then
    echo -e "${GREEN}🎉 SUCESSO TOTAL! Todas as collections foram populadas!${NC}"
    echo ""
    echo "Sistema EVA-Mind-FZPN está completo com:"
    echo "  🎭 Nasrudin (Paradoxo/Humor) - Tipos Emocionais"
    echo "  📚 Esopo (Moral/Lógica) - Tipos Racionais"
    echo "  🧘 Zen Koans (Insight/Silêncio) - Overthinking"
    echo "  🫁 Exercícios Somáticos (Aterramento) - Crises"
else
    echo -e "${YELLOW}⚠️  PARCIALMENTE CONCLUÍDO: $SUCCESS_COUNT de $TOTAL_COUNT scripts executados com sucesso${NC}"
    echo ""
    echo "Verifique os erros acima e execute novamente os scripts que falharam."
fi

echo "======================================================================="
echo ""

# Mostrar próximos passos
echo -e "${BLUE}📋 PRÓXIMOS PASSOS:${NC}"
echo ""
echo "1. Verificar collections:"
echo "   curl http://localhost:6333/collections | jq"
echo ""
echo "2. Testar busca em uma collection:"
echo "   curl -X POST http://localhost:6333/collections/nasrudin_stories/points/scroll \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"limit\": 5, \"with_payload\": true}'"
echo ""
echo "3. Implementar backend Go (pkg/nasrudin/, pkg/aesop/, pkg/zen/)"
echo ""
echo "4. Testar integração completa"
echo ""
echo "======================================================================="
