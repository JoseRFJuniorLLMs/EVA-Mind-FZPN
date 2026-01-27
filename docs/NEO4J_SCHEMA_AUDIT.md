# Auditoria do Schema Neo4j - EVA Mind

**Data**: 2026-01-27
**Versão**: 1.0

---

## 1. Nodes (Labels)

| Label | Descrição | Propriedades | Arquivo de Origem |
|-------|-----------|--------------|-------------------|
| `Person` | Idoso/Paciente | `id` (int64) | graph_store.go, significante.go, fdpn_engine.go |
| `Patient` | Alias de Person (usado em ethics) | `id` | ethical_boundary_engine.go |
| `Paciente` | Alias de Person (usado em knowledge) | `id` | graph_reasoning.go |
| `Event` | Evento/Conversa | `id`, `content`, `timestamp`, `speaker`, `emotion`, `importance`, `sessionId`, `type` | graph_store.go, significante.go |
| `Topic` | Tópico mencionado | `name`, `created` | graph_store.go, pattern_miner.go |
| `Emotion` | Estado emocional | `name` | graph_store.go |
| `Significante` | Palavra significativa (Lacan) | `word`, `idoso_id`, `emotional_valence`, `frequency`, `created_at`, `context`, `last_interpellation` | significante.go |
| `Demand` | Demanda/Pedido (FDPN) | `type`, `text`, `timestamp`, `urgency` | fdpn_engine.go |
| `Addressee` | Destinatário de demanda | `type` | fdpn_engine.go |
| `Pattern` | Padrão comportamental | `idoso_id`, `name`, `type`, `first_seen`, `last_seen`, `occurrences`, `avg_hour`, `confidence` | pattern_miner.go |
| `Condition` | Condição médica | - | unified_retrieval.go |
| `Medication` | Medicamento | - | unified_retrieval.go |
| `Symptom` | Sintoma | - | unified_retrieval.go |
| `Phrase` | Frase dita | - | ethical_boundary_engine.go |
| `Eneatipo` | Tipo de personalidade | - | fdpn_engine.go |

---

## 2. Relationships

| Tipo | De | Para | Propriedades | Arquivo |
|------|-----|------|--------------|---------|
| `EXPERIENCED` | Person | Event | - | graph_store.go, significante.go |
| `RELATED_TO` | Event | Topic | - | graph_store.go |
| `MENTIONED` | Person | Topic | `count`, `first_mention`, `last_mention` | graph_store.go |
| `FEELS` | Person | Emotion | `count`, `first_felt`, `last_felt` | graph_store.go |
| `EVOCA` | Event | Significante | - | significante.go |
| `DEMANDS` | Person | Demand | - | fdpn_engine.go |
| `ADDRESSED_TO` | Demand | Addressee | - | fdpn_engine.go |
| `FREQUENTLY_ADDRESSES` | Person | Addressee | `count`, `last_time` | fdpn_engine.go |
| `HAS_PATTERN` | Person | Pattern | - | pattern_miner.go |
| `REPRESENTS` | Pattern | Topic | - | pattern_miner.go |
| `HAS_CONDITION` | Person | Condition | - | unified_retrieval.go |
| `TAKES_MEDICATION` | Person | Medication | - | unified_retrieval.go |
| `SAID` | Patient | Phrase | - | ethical_boundary_engine.go |

---

## 3. Queries por Funcionalidade

### 3.1 Armazenamento de Memórias (graph_store.go)

```cypher
-- Criar evento de conversa
MERGE (p:Person {id: $idosoId})
CREATE (e:Event {
    id: $id,
    content: $content,
    timestamp: datetime($timestamp),
    speaker: $speaker,
    emotion: $emotion,
    importance: $importance,
    sessionId: $sessionId
})
CREATE (p)-[:EXPERIENCED]->(e)
```

```cypher
-- Conectar tópicos
MATCH (e:Event {id: $eventId})
MATCH (p:Person {id: $idosoId})
MERGE (t:Topic {name: $topic})
ON CREATE SET t.created = datetime()
MERGE (e)-[:RELATED_TO]->(t)
MERGE (p)-[r:MENTIONED]->(t)
ON CREATE SET r.count = 1, r.first_mention = datetime()
ON MATCH SET r.count = r.count + 1, r.last_mention = datetime()
```

```cypher
-- Conectar emoções
MATCH (p:Person {id: $idosoId})
MERGE (em:Emotion {name: $emotion})
MERGE (p)-[r:FEELS]->(em)
ON CREATE SET r.count = 1, r.first_felt = datetime()
ON MATCH SET r.count = r.count + 1, r.last_felt = datetime()
```

### 3.2 Significantes Lacanianos (significante.go)

```cypher
-- Registrar significante
MERGE (p:Person {id: $idosoId})
ON CREATE SET p.created = datetime()

MERGE (s:Significante {word: $word, idoso_id: $idosoId})
ON CREATE SET
    s.emotional_valence = $valence,
    s.frequency = 1,
    s.created_at = datetime(),
    s.context = $context
ON MATCH SET
    s.emotional_valence = (s.emotional_valence + $valence) / 2,
    s.frequency = s.frequency + 1

CREATE (e:Event {type: 'utterance', content: $context, timestamp: datetime()})
MERGE (e)-[:EVOCA]->(s)
MERGE (p)-[:EXPERIENCED]->(e)
```

```cypher
-- Buscar significantes
MATCH (s:Significante {idoso_id: $idosoId})
RETURN s.word AS word,
       s.emotional_valence AS valence,
       s.frequency AS frequency,
       s.context AS lastContext
ORDER BY s.frequency DESC
LIMIT 50
```

### 3.3 Motor FDPN (fdpn_engine.go)

```cypher
-- Registrar demanda
MERGE (p:Person {id: $idosoId})
ON CREATE SET p.created = datetime()

MERGE (a:Addressee {type: $addressee})
ON CREATE SET a.created = datetime()

CREATE (d:Demand {
    type: $demandType,
    text: $text,
    timestamp: datetime(),
    urgency: $urgency
})
MERGE (p)-[:DEMANDS]->(d)
MERGE (d)-[:ADDRESSED_TO]->(a)

MERGE (p)-[r:FREQUENTLY_ADDRESSES]->(a)
ON CREATE SET r.count = 1, r.last_time = datetime()
ON MATCH SET r.count = r.count + 1, r.last_time = datetime()
```

```cypher
-- Buscar destinatários frequentes
MATCH (p:Person {id: $idosoId})-[r:FREQUENTLY_ADDRESSES]->(a:Addressee)
RETURN a.type AS addressee, r.count AS count
ORDER BY r.count DESC
LIMIT 5
```

### 3.4 Priming Engine (priming_engine.go)

```cypher
-- Buscar eventos recentes para priming
MATCH (p:Person {id: $idosoId})-[:EXPERIENCED]->(e:Event)
WHERE e.timestamp > datetime() - duration({days: 7})
WITH e ORDER BY e.timestamp DESC LIMIT 10
OPTIONAL MATCH (e)-[:RELATED_TO|EVOCA]->(related)
RETURN e.content AS content,
       e.emotion AS emotion,
       collect(related) AS associations
```

### 3.5 Pattern Miner (pattern_miner.go)

```cypher
-- Descobrir padrões
MATCH (p:Person {id: $idosoId})-[:EXPERIENCED]->(e:Event)-[:RELATED_TO]->(t:Topic)
WHERE e.timestamp > datetime() - duration({days: 30})
WITH t.name AS topic,
     e.timestamp.hour AS hour,
     count(*) AS mentions
WHERE mentions >= 3
RETURN topic, hour, mentions
ORDER BY mentions DESC
```

```cypher
-- Salvar padrão descoberto
MATCH (p:Person {id: $idosoId})
MERGE (pat:Pattern {
    idoso_id: $idosoId,
    name: $patternName,
    type: $patternType
})
ON CREATE SET
    pat.first_seen = datetime(),
    pat.last_seen = datetime(),
    pat.occurrences = 1,
    pat.avg_hour = $avgHour,
    pat.confidence = $confidence
ON MATCH SET
    pat.last_seen = datetime(),
    pat.occurrences = pat.occurrences + 1,
    pat.confidence = $confidence
MERGE (p)-[:HAS_PATTERN]->(pat)
```

### 3.6 Unified Retrieval (unified_retrieval.go)

```cypher
-- Buscar contexto médico
MATCH (p:Person {id: $idosoId})
OPTIONAL MATCH (p)-[:HAS_CONDITION]->(c:Condition)
OPTIONAL MATCH (p)-[:TAKES_MEDICATION]->(m:Medication)
OPTIONAL MATCH (p)-[:EXPERIENCED]->(s:Symptom)
RETURN p, collect(DISTINCT c) AS conditions,
       collect(DISTINCT m) AS medications,
       collect(DISTINCT s) AS symptoms
```

---

## 4. Índices Recomendados

```cypher
-- Criar índices para performance
CREATE INDEX person_id IF NOT EXISTS FOR (p:Person) ON (p.id);
CREATE INDEX event_timestamp IF NOT EXISTS FOR (e:Event) ON (e.timestamp);
CREATE INDEX event_id IF NOT EXISTS FOR (e:Event) ON (e.id);
CREATE INDEX significante_idoso IF NOT EXISTS FOR (s:Significante) ON (s.idoso_id);
CREATE INDEX topic_name IF NOT EXISTS FOR (t:Topic) ON (t.name);
CREATE INDEX pattern_idoso IF NOT EXISTS FOR (pat:Pattern) ON (pat.idoso_id);
```

---

## 5. Queries de Diagnóstico

### Ver todos os dados de todos os idosos:
```cypher
MATCH (n)
RETURN labels(n) AS tipo, count(n) AS quantidade
ORDER BY quantidade DESC
```

### Ver todas as conversas:
```cypher
MATCH (p:Person)-[:EXPERIENCED]->(e:Event)
RETURN p.id AS idoso_id,
       e.content AS mensagem,
       e.speaker AS quem_falou,
       e.timestamp AS quando,
       e.emotion AS emocao
ORDER BY e.timestamp DESC
LIMIT 100
```

### Ver grafo completo de um idoso:
```cypher
MATCH (p:Person {id: 1})-[r*1..3]-(n)
RETURN p, r, n
```

### Verificar se banco está vazio:
```cypher
MATCH (n)
RETURN count(n) AS total_nodes
```

---

## 6. Problemas Identificados

### 6.1 Inconsistência de Labels
- `Person` vs `Patient` vs `Paciente` - Deveriam ser unificados
- Recomendação: Usar apenas `Person` com propriedade `type`

### 6.2 Dados Não Sendo Salvos
Se o banco está vazio, verificar:
1. Se `GraphStore.StoreCausalMemory()` está sendo chamado
2. Se a conexão Neo4j está funcionando
3. Se há erros nos logs

### 6.3 Queries Sem Dados
As queries que retornaram vazio podem indicar:
- Nenhuma conversa foi processada ainda
- O fluxo de salvamento não está ativo

---

## 7. Como Testar

```bash
# Conectar ao Neo4j
cypher-shell -u neo4j -p Debian23 -a bolt://104.248.219.200:7687

# Verificar nodes existentes
MATCH (n) RETURN labels(n), count(n);

# Criar dados de teste
CREATE (p:Person {id: 1})
CREATE (e:Event {id: 'test-1', content: 'Olá, como vai?', speaker: 'user', timestamp: datetime()})
CREATE (p)-[:EXPERIENCED]->(e);

# Verificar se foi criado
MATCH (p:Person)-[:EXPERIENCED]->(e:Event) RETURN p, e;
```
