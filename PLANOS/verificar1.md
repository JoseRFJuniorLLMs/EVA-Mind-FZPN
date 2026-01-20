# 🚀 EVA-Mind: Plano de Ação Imediato
## *Da Documentação à Execução - Próximas 2 Horas*

---

## ✅ CHECKLIST PRÉ-EXECUÇÃO

Antes de começar, verifique:

```bash
# 1. Conexões dos bancos
psql -h 34.175.224.36 -U postgres -d eva_db -c "SELECT 1"
# ✅ Deve retornar 1

# 2. Neo4j rodando
curl http://localhost:7474
# ✅ Deve retornar HTML do Neo4j Browser

# 3. Redis ativo
redis-cli ping
# ✅ Deve retornar PONG

# 4. Qdrant online
curl http://localhost:6333/collections
# ✅ Deve retornar JSON
```

---

## 🔴 FIX #1: Embedding Dimension (15 minutos)

### Passo 1: Verificar Problema Atual

```bash
cd /root/eva-mind-fzpn

# Ver schema atual do Postgres
psql -h 34.175.224.36 -U postgres -d eva_db -c "\d episodic_memories"
```

**Output esperado (ERRADO):**
```
embedding | vector(1536) |   <-- PROBLEMA: Deveria ser 768
```

### Passo 2: Criar Arquivo de Migração

```bash
nano migrations/004_fix_embedding_dimension.sql
```

**Cole este conteúdo:**

```sql
-- migrations/004_fix_embedding_dimension.sql
-- Correção: Gemini text-embedding-004 retorna 768 dimensões, não 1536

BEGIN;

-- 1. Criar nova coluna com dimensão correta
ALTER TABLE episodic_memories 
ADD COLUMN embedding_new vector(768);

-- 2. Marcar dados antigos como inválidos (serão re-embeddados)
UPDATE episodic_memories 
SET embedding_new = NULL;

-- 3. Remover coluna antiga
ALTER TABLE episodic_memories 
DROP COLUMN embedding;

-- 4. Renomear
ALTER TABLE episodic_memories 
RENAME COLUMN embedding_new TO embedding;

-- 5. Atualizar função de busca (com aspas para palavras reservadas)
DROP FUNCTION IF EXISTS search_similar_memories(BIGINT, vector, INT, FLOAT);

CREATE OR REPLACE FUNCTION search_similar_memories(
    p_idoso_id BIGINT,
    p_query_embedding vector(768),  -- ✅ 768 dimensões
    p_limit INT DEFAULT 10,
    p_min_similarity FLOAT DEFAULT 0.5
)
RETURNS TABLE (
    id BIGINT,
    content TEXT,
    speaker TEXT,
    "timestamp" TIMESTAMPTZ,  -- ✅ Aspas para palavra reservada
    emotion TEXT,
    importance FLOAT,
    topics TEXT,
    similarity FLOAT
) AS $
BEGIN
    RETURN QUERY
    SELECT 
        em.id,
        em.content,
        em.speaker,
        em."timestamp",  -- ✅ Aspas aqui também
        em.emotion,
        em.importance,
        em.topics::TEXT,
        1 - (em.embedding <=> p_query_embedding) AS similarity
    FROM episodic_memories em
    WHERE em.idoso_id = p_idoso_id
      AND em.embedding IS NOT NULL
      AND (1 - (em.embedding <=> p_query_embedding)) >= p_min_similarity
    ORDER BY em.embedding <=> p_query_embedding
    LIMIT p_limit;
END;
$ LANGUAGE plpgsql;

COMMIT;

-- Verificação
SELECT 
    COUNT(*) as total_memories,
    COUNT(embedding) as with_embedding,
    COUNT(*) - COUNT(embedding) as need_reembedding
FROM episodic_memories;
```

### Passo 3: Executar Migração

```bash
psql -h 34.175.224.36 -U postgres -d eva_db -f migrations/004_fix_embedding_dimension.sql
```

**Output esperado:**
```
BEGIN
ALTER TABLE
UPDATE 150  (ou quantas memórias você tem)
ALTER TABLE
ALTER TABLE
DROP FUNCTION
CREATE FUNCTION
COMMIT
 total_memories | with_embedding | need_reembedding 
----------------+----------------+------------------
            150 |              0 |              150
```

**✅ MIGRAÇÃO CONCLUÍDA COM SUCESSO!**

### Passo 4: Validar no Código Go

Abra `internal/memory/embeddings.go` e adicione validação:

```bash
nano internal/memory/embeddings.go
```

**Adicione após o método GenerateEmbedding (linha ~50):**

```go
const (
    expectedDimension = 768 // ✅ CRÍTICO
)

// Dentro de GenerateEmbedding, ANTES de return:
actualDim := len(result.Embedding.Values)
if actualDim != expectedDimension {
    return nil, fmt.Errorf(
        "❌ DIMENSION MISMATCH!\n"+
        "   Expected: %d (Postgres schema)\n"+
        "   Got: %d (Gemini API)\n",
        expectedDimension,
        actualDim,
    )
}

log.Printf("✅ [EMBEDDING] Generated %d dimensions (validated)", actualDim)
```

### Passo 5: Re-embedar Memórias Existentes

```bash
# Criar script de re-embedding
nano cmd/reembed/main.go
```

**Cole este código completo:**

```go
package main

import (
    "context"
    "database/sql"
    "eva-mind/internal/config"
    "eva-mind/internal/memory"
    "fmt"
    "log"
    "time"

    _ "github.com/lib/pq"
)

func main() {
    log.Println("🔄 Re-embedding Script Started")

    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("❌ Config error: %v", err)
    }

    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("❌ DB connection error: %v", err)
    }
    defer db.Close()

    embedder := memory.NewEmbeddingService(cfg.GoogleAPIKey)

    query := `
        SELECT id, content 
        FROM episodic_memories 
        WHERE embedding IS NULL 
          AND content IS NOT NULL
          AND LENGTH(content) > 10
        ORDER BY id
    `

    rows, err := db.Query(query)
    if err != nil {
        log.Fatalf("❌ Query error: %v", err)
    }
    defer rows.Close()

    ctx := context.Background()
    success := 0
    failed := 0

    for rows.Next() {
        var id int64
        var content string

        if err := rows.Scan(&id, &content); err != nil {
            continue
        }

        embedding, err := embedder.GenerateEmbedding(ctx, content)
        if err != nil {
            log.Printf("❌ ID=%d failed: %v", id, err)
            failed++
            continue
        }

        // Converter para formato Postgres
        vectorStr := "["
        for i, v := range embedding {
            if i > 0 {
                vectorStr += ","
            }
            vectorStr += fmt.Sprintf("%f", v)
        }
        vectorStr += "]"

        _, err = db.Exec("UPDATE episodic_memories SET embedding = $1 WHERE id = $2", vectorStr, id)
        if err != nil {
            log.Printf("❌ ID=%d update failed: %v", id, err)
            failed++
            continue
        }

        success++
        if success%10 == 0 {
            log.Printf("✅ Progress: %d embeddings created", success)
        }

        time.Sleep(100 * time.Millisecond) // Rate limit
    }

    log.Printf("📊 Complete: %d success, %d failed", success, failed)
}
```

**Execute:**

```bash
go run cmd/reembed/main.go
```

---

## 🔴 FIX #2: Qdrant Upsert (5 minutos)

### Problema
Memórias salvam no Postgres mas nunca chegam no Qdrant.

### Solução: Adicionar Upsert no save_memory_helper

Abra `main.go` e localize a função `saveAsMemory` (linha ~650):

```bash
nano main.go
```

**Adicione APÓS salvar no Postgres (linha ~720):**

```go
// ✅ NOVO: UPSERT NO QDRANT (Assíncrono)
if s.qdrantClient != nil {
    go func() {
        qctx, qcancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer qcancel()

        point := &qdrant.PointStruct{
            Id: &qdrant.PointId{
                PointIdOptions: &qdrant.PointId_Num{Num: uint64(mem.ID)},
            },
            Vectors: &qdrant.Vectors{
                VectorsOptions: &qdrant.Vectors_Vector{
                    Vector: &qdrant.Vector{Data: mem.Embedding},
                },
            },
            Payload: map[string]*qdrant.Value{
                "content": {Kind: &qdrant.Value_StringValue{StringValue: mem.Content}},
                "speaker": {Kind: &qdrant.Value_StringValue{StringValue: mem.Speaker}},
                "idoso_id": {Kind: &qdrant.Value_IntegerValue{IntegerValue: mem.IdosoID}},
            },
        }

        for attempt := 1; attempt <= 3; attempt++ {
            err := s.qdrantClient.Upsert(qctx, "memories", []*qdrant.PointStruct{point})
            if err == nil {
                log.Printf("✅ [QDRANT] Memory %d indexed", mem.ID)
                break
            }
            if attempt < 3 {
                time.Sleep(time.Second * time.Duration(attempt))
            } else {
                log.Printf("❌ [QDRANT] Failed after 3 attempts: %v", err)
            }
        }
    }()
}
```

---

## 🔴 FIX #3: FDPN Hook (2 minutos)

### Problema
`StreamingPrime` existe mas nunca é chamado.

### Solução: Adicionar Hook na Transcrição

No `client.go` do Gemini, a transcrição já chama um callback. Precisamos garantir que o FDPN seja ativado lá.

Abra `main.go` e localize `setupGeminiSession` (linha ~380):

```bash
nano main.go
```

**Localize onde configura callbacks (linha ~430) e VERIFIQUE se tem isto:**

```go
// 🔍 3. Callback de Transcrição (Dual-Model + AUTO-SAVE)
func(role, text string) {
    if role == "user" {
        // ✅ DEVE TER ISTO - FDPN Priming
        if s.fdpnEngine != nil {
            go func() {
                err := s.fdpnEngine.StreamingPrime(
                    client.ctx, 
                    strconv.FormatInt(client.IdosoID, 10), 
                    text,
                )
                if err != nil {
                    log.Printf("⚠️ FDPN Error: %v", err)
                }
            }()
        }
        
        // Outros processamentos...
    }

    // AUTO-SAVE (ambos roles)
    go s.saveAsMemory(client.IdosoID, role, text)
},
```

**Se NÃO tiver, adicione.**

---

## ✅ TESTE COMPLETO (5 minutos)

### 1. Reiniciar o Sistema

```bash
# Compilar
go build -o eva-mind cmd/main.go

# Rodar
./eva-mind
```

### 2. Fazer Conversa de Teste

Via mobile ou WebSocket, diga algo como:
> "Estou com dor de cabeça e me sinto sozinho"

### 3. Verificar Logs

Você DEVE ver estas linhas:

```
✅ [EMBEDDING] Generated 768 dimensions (validated)
✅ [POSTGRES] Memory saved: ID=123, Speaker=user
✅ [QDRANT] Memory 123 indexed
✅ [FDPN] Primed in 45ms (user=1)
🔍 [FDPN] Keywords: [dor, cabeça, sinto, sozinho]
```

### 4. Validar Bancos

```bash
# Postgres - deve ter embedding
psql -h 34.175.224.36 -U postgres -d eva_db -c \
  "SELECT id, LEFT(content, 50), embedding IS NOT NULL FROM episodic_memories ORDER BY id DESC LIMIT 5"

# Qdrant - deve ter pontos
curl http://localhost:6333/collections/memories | jq '.result.points_count'

# Neo4j - deve ter eventos
# Abrir http://localhost:7474
# Executar: MATCH (e:Event) RETURN e ORDER BY e.timestamp DESC LIMIT 10
```

---

## 📊 CHECKLIST FINAL

```
[ ] Fix #1: Embedding = 768 dimensões (verificado no schema)
[ ] Fix #1: GenerateEmbedding retorna 768 (sem erros)
[ ] Fix #1: Re-embedding executado (X memórias processadas)
[ ] Fix #2: Qdrant recebe pontos (points_count > 0)
[ ] Fix #3: Logs FDPN aparecem ("✅ [FDPN] Primed")
[ ] Teste: Conversa completa funciona
[ ] Teste: Busca semântica retorna resultados
```

---

## 🚨 TROUBLESHOOTING

### Erro: "DIMENSION MISMATCH"
✅ Fix #1 não foi aplicado corretamente. Re-execute migração SQL.

### Erro: "Qdrant connection refused"
```bash
docker ps | grep qdrant
# Se não estiver rodando:
docker start qdrant
```

### Erro: "Neo4j unavailable"
```bash
sudo systemctl status neo4j
sudo systemctl start neo4j
```

### FDPN não ativa
- Verifique se `s.fdpnEngine != nil` no código
- Confirme que Neo4j está rodando
- Veja logs: `grep FDPN /var/log/eva-mind.log`

---

## 🎯 PRÓXIMOS PASSOS (Após 2h)

Após validar os 3 fixes:

**Semana 1:**
- [ ] Implementar `UnifiedRetrieval` (orquestra 4 DBs)
- [ ] Implementar `ReliableSaver` (retry + dead letter queue)
- [ ] Adicionar Neo4j Schema (constraints + índices)

**Semana 2:**
- [ ] Pattern Mining (detectar recorrências)
- [ ] Temporal Decay (spreading activation inteligente)
- [ ] Metadata Analyzer LLM (substituir regex)

**Semana 3:**
- [ ] Story Repository (Qdrant com histórias)
- [ ] Zeta Router (seleção por Eneagrama)
- [ ] Prompt System atualizado

**Semana 4:**
- [ ] Health Checks
- [ ] Observability
- [ ] Dashboard

---

## 💡 DICA FINAL

**Faça commit após cada fix:**

```bash
git add migrations/004_fix_embedding_dimension.sql
git commit -m "fix: embedding dimension 1536→768 (Gemini API)"

git add main.go
git commit -m "fix: adicionar Qdrant upsert em saveAsMemory"

git add main.go
git commit -m "fix: ativar FDPN hook na transcrição"

git push
```

---

**Tempo total estimado:** 22 minutos  
**Impacto:** Sistema 30% → 70% funcional  
**Próximo objetivo:** Sistema 70% → 90% (Semana 1)

Boa sorte, José! 🚀