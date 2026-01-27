#!/bin/bash
# Neo4j Audit Script - EVA Mind
# Executa queries de diagnóstico no Neo4j

NEO4J_USER="neo4j"
NEO4J_PASS="Debian23"
NEO4J_URI="bolt://localhost:7687"

echo "=================================================="
echo "🔍 AUDITORIA NEO4J - EVA MIND"
echo "Data: $(date)"
echo "=================================================="

# Função para executar query
run_query() {
    local title="$1"
    local query="$2"
    echo ""
    echo "📊 $title"
    echo "--------------------------------------------------"
    cypher-shell -u "$NEO4J_USER" -p "$NEO4J_PASS" -a "$NEO4J_URI" "$query" 2>/dev/null || echo "❌ Erro na query"
}

# 1. Contagem total de nós
run_query "1. CONTAGEM DE NÓS POR TIPO" \
"MATCH (n) RETURN labels(n) AS tipo, count(n) AS quantidade ORDER BY quantidade DESC;"

# 2. Contagem de relacionamentos
run_query "2. CONTAGEM DE RELACIONAMENTOS" \
"MATCH ()-[r]->() RETURN type(r) AS tipo, count(r) AS quantidade ORDER BY quantidade DESC;"

# 3. Labels existentes
run_query "3. LABELS EXISTENTES" \
"CALL db.labels() YIELD label RETURN label ORDER BY label;"

# 4. Relationship types existentes
run_query "4. TIPOS DE RELACIONAMENTO" \
"CALL db.relationshipTypes() YIELD relationshipType RETURN relationshipType ORDER BY relationshipType;"

# 5. Todas as pessoas/idosos
run_query "5. TODAS AS PESSOAS (IDOSOS)" \
"MATCH (p:Person) RETURN p.id AS idoso_id, p.created AS criado LIMIT 20;"

# 6. Todos os eventos/conversas (últimos 50)
run_query "6. ÚLTIMAS CONVERSAS (50)" \
"MATCH (p:Person)-[:EXPERIENCED]->(e:Event)
RETURN p.id AS idoso_id, e.speaker AS quem, substring(e.content, 0, 80) AS mensagem, e.timestamp AS quando
ORDER BY e.timestamp DESC LIMIT 50;"

# 7. Significantes mais frequentes
run_query "7. SIGNIFICANTES MAIS FREQUENTES" \
"MATCH (s:Significante)
RETURN s.idoso_id AS idoso, s.word AS palavra, s.frequency AS freq, s.emotional_valence AS valencia
ORDER BY s.frequency DESC LIMIT 20;"

# 8. Tópicos mais mencionados
run_query "8. TÓPICOS MAIS MENCIONADOS" \
"MATCH (p:Person)-[r:MENTIONED]->(t:Topic)
RETURN p.id AS idoso, t.name AS topico, r.count AS vezes
ORDER BY r.count DESC LIMIT 20;"

# 9. Emoções registradas
run_query "9. EMOÇÕES REGISTRADAS" \
"MATCH (p:Person)-[r:FEELS]->(em:Emotion)
RETURN p.id AS idoso, em.name AS emocao, r.count AS vezes
ORDER BY r.count DESC LIMIT 20;"

# 10. Demandas registradas
run_query "10. DEMANDAS (FDPN)" \
"MATCH (p:Person)-[:DEMANDS]->(d:Demand)
RETURN p.id AS idoso, d.type AS tipo, substring(d.text, 0, 60) AS texto, d.urgency AS urgencia
ORDER BY d.timestamp DESC LIMIT 20;"

# 11. Padrões comportamentais
run_query "11. PADRÕES COMPORTAMENTAIS" \
"MATCH (p:Person)-[:HAS_PATTERN]->(pat:Pattern)
RETURN p.id AS idoso, pat.name AS padrao, pat.type AS tipo, pat.occurrences AS ocorrencias, pat.confidence AS confianca
ORDER BY pat.occurrences DESC LIMIT 20;"

# 12. Dados específicos do usuário 1121
run_query "12. DADOS DO USUÁRIO 1121" \
"MATCH (p:Person {id: 1121})-[r]->(n)
RETURN type(r) AS relacao, labels(n) AS tipo_destino, count(*) AS quantidade;"

# 13. Conversas do usuário 1121
run_query "13. CONVERSAS DO USUÁRIO 1121" \
"MATCH (p:Person {id: 1121})-[:EXPERIENCED]->(e:Event)
RETURN e.speaker AS quem, substring(e.content, 0, 100) AS mensagem, e.emotion AS emocao, e.timestamp AS quando
ORDER BY e.timestamp DESC LIMIT 30;"

# 14. Verificar índices
run_query "14. ÍNDICES EXISTENTES" \
"SHOW INDEXES;"

# 15. Estatísticas do banco
run_query "15. ESTATÍSTICAS GERAIS" \
"MATCH (n) WITH count(n) AS nodes
MATCH ()-[r]->() WITH nodes, count(r) AS rels
RETURN nodes AS total_nos, rels AS total_relacionamentos;"

echo ""
echo "=================================================="
echo "✅ AUDITORIA COMPLETA"
echo "=================================================="
