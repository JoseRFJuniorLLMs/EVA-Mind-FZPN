# Sistema de Armazenamento de Memórias - EVA-Mind-FZPN

**Documento:** MEMORY-STORAGE-001
**Versão:** 1.0
**Data:** 2026-01-27
**Autor:** Jose R F Junior

---

## 1. Visão Geral

O EVA-Mind-FZPN implementa um sistema de memória multi-camadas para armazenar conversas dos idosos com a IA EVA.

### Bancos de Dados Utilizados

| Banco | Função | Dados |
|-------|--------|-------|
| **PostgreSQL** | Principal | Mensagens, embeddings, metadados |
| **Qdrant** | Busca semântica | Vetores 768-dim |
| **Neo4j** | Grafo causal | Relações paciente→tópicos→emoções |
| **Redis** | Cache | Estado em tempo real |

---

## 2. O Que É Armazenado

### 2.1 Tabela Principal: `episodic_memories`

```sql
CREATE TABLE episodic_memories (
    id INTEGER PRIMARY KEY,
    idoso_id INTEGER NOT NULL,          -- ID do paciente
    timestamp TIMESTAMP WITH TIME ZONE, -- Quando falou
    speaker VARCHAR(20),                -- 'user' ou 'assistant'
    content TEXT NOT NULL,              -- Texto da mensagem
    emotion VARCHAR(50),                -- feliz, triste, ansioso, confuso, neutro
    importance DOUBLE PRECISION,        -- 0.0 a 1.0
    topics TEXT[],                      -- [saúde, família, medicamento, lazer]
    session_id VARCHAR(100),            -- ID da sessão
    embedding VECTOR(768)               -- Vetor Gemini text-embedding-004
);
```

### 2.2 Cálculo de Importance

| Conteúdo da Mensagem | Importance |
|---------------------|------------|
| "emergência", "socorro", "caí" | **1.0** |
| "dor", "médico", "remédio" | **0.8** |
| Qualquer outro texto | **0.5** |
| "tempo", "hora" | **0.3** |

**Arquivo:** `internal/hippocampus/memory/analyzer.go`

### 2.3 Sistemas de Memória Avançada

O Sprint 12 implementou 12 sistemas de memória do paciente:

1. **Eneagrama Gurdjieff** - Tipo de personalidade
2. **Self-Core** - Auto-descrições do paciente
3. **Padrões Comportamentais** - Gatilhos e respostas
4. **Intenções vs Realizações** - Promessas feitas
5. **Contrafactuais** - Ruminações "e se..."
6. **Metáforas Pessoais** - "peso no peito", "num buraco"
7. **Padrões Transgeracionais** - Traumas familiares
8. **Correlações Somáticas** - Corpo × emoção
9. **Contexto Histórico-Cultural** - Eventos vividos
10. **Aprendizado Terapêutico** - O que funciona
11. **Preditores de Crise** - Marcadores de alerta
12. **Mapa do Mundo** - Pessoas, lugares, objetos

---

## 3. Quando É Armazenado

```
T=0ms      User fala
T=10ms     INSERT PostgreSQL (síncrono)
T=100ms    Gera embedding (Gemini API)
T=1100ms   Upsert Qdrant (async)
T=2000ms+  Store Neo4j (async)
```

**Código:** `internal/cortex/brain/memory.go`

```go
func (s *Service) ProcessUserSpeech(ctx context.Context, idosoID int64, text string) {
    go s.SaveEpisodicMemory(idosoID, "user", text)  // Fire and forget
}
```

---

## 4. Política de Retenção

### 4.1 Regra Principal

**MEMÓRIAS NUNCA SÃO DELETADAS AUTOMATICAMENTE.**

A EVA precisa das memórias para:
- Manter contexto das conversas
- Personalizar interações
- Detectar padrões de comportamento
- Prever crises

### 4.2 LGPD Compliance

A LGPD **NÃO obriga** deleção automática. Exige apenas:

| Direito | Implementação |
|---------|---------------|
| Acesso | `ExportPersonalData()` |
| Correção | `RectifyPersonalData()` |
| Eliminação | `DeletePersonalData()` - **sob demanda** |
| Portabilidade | `ExportToJSON()` |

**Arquivo:** `internal/audit/data_rights.go`

### 4.3 Dados Clínicos

Dados clínicos (PHQ-9, GAD-7, etc.) devem ser mantidos por **mínimo 20 anos** (Resolução CFM 1.821/2007).

---

## 5. Funções Administrativas (Restritas)

### 5.1 Acesso Restrito

As funções de deleção são **RESTRITAS** ao criador da EVA:

```go
// CPF do Criador - Jose R F Junior
const CREATOR_CPF = "64525430249"
```

**Arquivo:** `internal/hippocampus/memory/storage.go`

### 5.2 Funções Disponíveis

#### `DeleteOld()` - Deletar memórias antigas

```go
func (m *MemoryStore) DeleteOld(
    ctx context.Context,
    requesterCPF string,    // Deve ser 64525430249
    idosoID int64,          // ID do paciente (0 = todos)
    olderThanDays int,      // Dias de idade
    minImportance float64,  // Deletar apenas < este valor
) (int64, error)
```

**Exemplo:**
```go
// Deletar memórias > 90 dias com importance < 0.5 do paciente 123
deleted, err := memoryStore.DeleteOld(ctx, "64525430249", 123, 90, 0.5)
```

#### `DeleteAllMemories()` - Deletar TODAS memórias

```go
func (m *MemoryStore) DeleteAllMemories(
    ctx context.Context,
    requesterCPF string,    // Deve ser 64525430249
    idosoID int64,          // ID do paciente
) (int64, error)
```

**Exemplo:**
```go
// Deletar TODAS memórias do paciente 123 (CUIDADO!)
deleted, err := memoryStore.DeleteAllMemories(ctx, "64525430249", 123)
```

#### `GetMemoryStats()` - Estatísticas

```go
func (m *MemoryStore) GetMemoryStats(
    ctx context.Context,
    requesterCPF string,    // Deve ser 64525430249
) (map[string]interface{}, error)
```

**Retorna:**
```json
{
  "total_memories": 15420,
  "total_patients_with_memories": 87,
  "avg_memories_per_patient": 177.24,
  "by_importance": {
    "critical (>=0.9)": 234,
    "important (0.7-0.9)": 1567,
    "normal (0.5-0.7)": 8945,
    "low (<0.5)": 4674
  },
  "oldest_memory": "2025-06-15T10:30:00Z",
  "newest_memory": "2026-01-27T14:25:00Z"
}
```

### 5.3 Segurança

Tentativas não autorizadas são **logadas**:

```
🚫 [SECURITY] Tentativa não autorizada de DeleteOld por CPF: 111.111.111-11
```

Uso autorizado:
```
🔧 [ADMIN] DeleteOld autorizado para criador Jose R F Junior
🔧 [ADMIN] Parâmetros: idosoID=123, olderThanDays=90, minImportance=0.50
✅ [ADMIN] DeleteOld concluído: 150 memórias removidas
```

### 5.4 Erro de Autorização

```go
var ErrUnauthorized = errors.New("acesso negado: apenas o criador pode executar esta função")
```

---

## 6. Fluxo de Armazenamento

```
┌─────────────────────────────────────────────────────────────────┐
│                    IDOSO FALA COM EVA                           │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────────┐
│  ProcessUserSpeech() → go SaveEpisodicMemory()                  │
└──────────────────────────┬───────────────────────────────────────┘
                           │
           ┌───────────────┼───────────────┬───────────────┐
           │               │               │               │
           ▼               ▼               ▼               ▼
    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
    │PostgreSQL│    │  Gemini  │    │  Qdrant  │    │  Neo4j   │
    │          │    │Embedding │    │ (async)  │    │ (async)  │
    │ INSERT   │◄───│  768-dim │───►│  Upsert  │    │ MERGE    │
    └──────────┘    └──────────┘    └──────────┘    └──────────┘
         │                                               │
         ▼                                               ▼
    ┌──────────────────────────┐         ┌─────────────────────────┐
    │ episodic_memories        │         │ (Person)-[:EXPERIENCED] │
    │ patient_behavioral_*     │         │ (Event)-[:RELATED_TO]   │
    │ patient_self_core        │         │ (Topic)                 │
    └──────────────────────────┘         └─────────────────────────┘
```

---

## 7. Busca Semântica (Hybrid Search)

```go
// Query: "Como ele se sente sobre medicamentos?"

func (r *RetrievalService) Retrieve(ctx context.Context, idosoID int64, query string, k int) {
    // 1. Gerar embedding da query
    queryEmbedding, _ := r.embedder.GenerateEmbedding(ctx, query)

    // 2. Busca no PostgreSQL (pgvector)
    results := search_similar_memories(idosoID, queryEmbedding, k, 0.5)

    // 3. Busca no Qdrant (se disponível)
    qdrantResults := r.qdrant.Search(ctx, "memories", queryEmbedding, k)

    // 4. Merge e deduplica
    return mergeResults(results, qdrantResults)
}
```

---

## 8. Índices para Performance

```sql
-- Busca por paciente
CREATE INDEX idx_episodic_memories_idoso_id ON episodic_memories(idoso_id);

-- Busca temporal
CREATE INDEX idx_episodic_memories_timestamp ON episodic_memories(timestamp DESC);

-- Busca por speaker
CREATE INDEX idx_episodic_memories_speaker ON episodic_memories(speaker);

-- Busca por importância
CREATE INDEX idx_episodic_memories_importance ON episodic_memories(importance DESC);

-- Busca semântica (pgvector)
CREATE INDEX idx_episodic_memories_embedding
  ON episodic_memories USING ivfflat(embedding vector_cosine_ops)
  WITH (lists = 100);

-- Full-text search
CREATE INDEX idx_episodic_memories_content_gin
  ON episodic_memories USING GIN(to_tsvector('portuguese', content));
```

---

## 9. Resumo

| Aspecto | Comportamento |
|---------|---------------|
| Armazenamento | Automático, toda conversa |
| Retenção | Indefinida (nunca deleta automaticamente) |
| Deleção manual | Apenas criador (CPF 64525430249) |
| LGPD | Export/Delete sob demanda do usuário |
| Dados clínicos | Mínimo 20 anos |

---

## Aprovações

| Função | Nome | Data |
|--------|------|------|
| Criador/Admin | Jose R F Junior | 2026-01-27 |

---

**Documento controlado - Versão 1.0**
