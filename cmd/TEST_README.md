# FZPN Validation Test Suite

## Objetivo

Este script valida automaticamente os 3 pilares da arquitetura FZPN:

1. **FDPN (Priming Engine):** Latência de priming < 10ms
2. **Zeta (Personality Router):** Mudanças corretas de personalidade
3. **Lacan (Signifier Service):** Rastreamento de significantes emocionais
4. **Co-Intelligence:** Anti-Sycophancy (Mollick)

## Pré-requisitos

### 1. Neo4j Rodando
```bash
# Verificar se Neo4j está acessível
curl http://104.248.219.200:7687
```

### 2. Redis (Opcional)
```bash
# Se tiver Redis local:
redis-cli ping
# Deve retornar: PONG

# Se não tiver, o teste vai rodar em modo degradado (sem cache L2)
```

### 3. Dados de Teste no Neo4j
O script assume que você já rodou o `seed_neo4j.go` para popular os Eneatipos.

## Como Executar

### Opção 1: Executar Diretamente
```bash
cd d:\dev\EVA\EVA-Mind
go run cmd/test_fzpn.go
```

### Opção 2: Compilar e Executar
```bash
cd d:\dev\EVA\EVA-Mind
go build -o test_fzpn.exe cmd/test_fzpn.go
./test_fzpn.exe
```

## O Que Esperar

### Output Esperado
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🧪 FZPN VALIDATION TEST SUITE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 Running Test Suites...

🔬 TEST 1: FDPN Priming Latency
   Objetivo: Verificar se priming é < 10ms (com cache)

   ✅ PASS | FDPN Cold Query (Neo4j direto)
      └─ Latência: 45ms ✅ Excelente (< 100ms)

   ✅ PASS | FDPN Hot Query (Redis cache)
      └─ Latência: 3ms 🚀 PERFEITO (< 10ms)

   ✅ PASS | FDPN Parallel Priming (5 keywords)
      └─ Latência: 12ms 🚀 Goroutines funcionando!

🔬 TEST 2: Zeta Personality Routing
   Objetivo: Verificar mudanças de tipo por emoção

   ✅ PASS | Zeta Stress Path (9 → 6)
      └─ Base: 9, Emoção: anxiety → Tipo: 6, Modo: stress ✅ Correto!

   ✅ PASS | Zeta Growth Path (9 → 3)
      └─ Base: 9, Emoção: joy → Tipo: 3, Modo: growth ✅ Correto!

   ✅ PASS | Zeta Attention Weights (Tipo 6)
      └─ RISCO: 2.2, SEGURANÇA: 2.0, AMBIGUIDADE: 0.5 ✅ Zeros corretos!

🔬 TEST 3: Lacan Signifier Detection
   Objetivo: Rastrear significantes emocionais

   ✅ PASS | Lacan Track Signifier
      └─ 5 textos processados em 234ms

   ✅ PASS | Lacan Retrieve Signifiers
      └─ Encontrados 1 significantes. Top: 'solidão' (freq: 5)

🔬 TEST 4: Anti-Sycophancy (Co-Intelligence)
   Objetivo: Verificar se prompts bloqueiam concordância perigosa

   ✅ PASS | Anti-Sycophancy Prompt Check
      └─ Prompt contém 'DISCORDE IMEDIATAMENTE' ✅

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 FINAL REPORT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Total Tests: 9
✅ Passed: 9
❌ Failed: 0
Pass Rate: 100.0%

📈 Telemetry Snapshot:
   Enneatype: 0
   Priming Latency: 3ms
   Switches: 0

🎉 FZPN ARCHITECTURE VALIDATED!
   Sistema pronto para produção.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Interpretação dos Resultados

### ✅ PASS (100%)
Todos os componentes funcionando conforme especificação. Sistema pronto.

### ⚠️ PASS (80-99%)
Maioria dos testes passou. Revisar falhas específicas.

### ❌ FAIL (< 80%)
Componentes críticos com problemas. Não deploy em produção.

## Troubleshooting

### Erro: "Neo4j connection failed"
```bash
# Verificar se Neo4j está rodando
# Verificar credenciais em .env
NEO4J_URI=neo4j://104.248.219.200:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=Debian23
```

### Erro: "Redis not available"
```bash
# Não é crítico. Sistema roda sem Redis, mas mais lento.
# Para instalar Redis localmente:
# Windows: https://github.com/microsoftarchive/redis/releases
# Linux: sudo apt install redis-server
```

### Teste "FDPN Hot Query" falhou (> 10ms)
- Redis não está rodando ou está lento
- Verificar latência de rede para Redis
- Considerar Redis local em vez de remoto

### Teste "Zeta Stress Path" falhou
- Verificar se `personality_router.go` tem as rotas corretas
- Revisar mapeamento de emoções para tipos

## Próximos Passos Após Validação

1. **Se 100% Pass:** Deploy em staging
2. **Se 80-99% Pass:** Corrigir falhas específicas
3. **Se < 80% Pass:** Revisar arquitetura

## Métricas de Sucesso

| Métrica | Target | Crítico? |
|---------|--------|----------|
| FDPN Cold Query | < 100ms | Não |
| FDPN Hot Query | < 10ms | **Sim** |
| Zeta Routing | 100% correto | **Sim** |
| Lacan Tracking | Funcional | Sim |
| Anti-Sycophancy | Presente | **Sim** |

---

**Desenvolvido por:** EVA-Mind Team  
**Arquitetura:** FZPN (Fractal Zeta Priming Network)  
**Data:** 2026-01-16
