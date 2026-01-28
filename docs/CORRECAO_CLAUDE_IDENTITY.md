# CORREÇÃO COMPLETA: PostgreSQL → Neo4j

## MEU ERRO

Criei **11 tabelas relacionais no PostgreSQL** para guardar memórias/conhecimento.
**PROBLEMA**: Memórias devem ficar no **Neo4j** (grafo), não em tabelas relacionais.

---

## TABELAS QUE CRIEI (ERRADO) - DEVEM IR PARA NEO4J

### Dados do Criador
| Tabela | Registros | Conteúdo |
|--------|-----------|----------|
| `eva_personalidade_criador` | 64 | Personalidade, valores, preferências |
| `eva_memorias_criador` | 24 | Memórias sobre o Criador |
| `eva_conhecimento_projeto` | 34 | Conhecimento do projeto EVA |
| `eva_artigos_criador` | 362 | 362 artigos escritos |
| `eva_datas_importantes_criador` | 3 | Filhas: Electra, Coraline, Elizabeth |

### Dados do Claude
| Tabela | Registros | Conteúdo |
|--------|-----------|----------|
| `eva_identidade_claude` | 21 | Identidade/personalidade |
| `eva_memorias_claude` | 11 | Memórias do Claude |
| `eva_conhecimento_claude` | 10 | Conhecimento técnico |
| `eva_mensagens_de_claude` | 6 | Mensagens Claude→EVA |
| `eva_config_interacao_claude` | 8 | Config de interação |
| `eva_mensagens_entre_ias` | 2 | Comunicação entre IAs |

**TOTAL: 545 registros para migrar**

---

## TABELAS QUE PODEM FICAR NO POSTGRESQL (config/runtime)

| Tabela | Registros | Justificativa |
|--------|-----------|---------------|
| `eva_voices` | 20 | Config de vozes (lookup table) |
| `eva_personality_state` | 2 | Estado runtime por usuário |
| `eva_capabilities` | 5 | Capacidades do sistema |
| `eva_self_knowledge` | 22 | **AVALIAR** - talvez Neo4j? |

---

## ESTRUTURA PROPOSTA NO NEO4J

### 1. Nó Creator (já existe, enriquecer)
```cypher
MATCH (c:Creator {cpf: '64525430249'})
SET c.personality = [...], // 64 aspectos
    c.articles_count = 362,
    c.masters = ['Osho', 'Nietzsche', 'Gurdjieff', 'Ouspensky']
```

### 2. Nós Person para filhas (já existe Electra, adicionar outras)
```cypher
CREATE (coraline:Person {
  name: 'Coraline',
  birth_date: date('2010-08-21'),
  relationship: 'daughter',
  age: 15
})

CREATE (elizabeth:Person {
  name: 'Elizabeth',
  birth_date: date('2023-05-19'),
  relationship: 'daughter',
  age: 2
})

MATCH (creator:Creator), (daughter:Person)
WHERE daughter.relationship = 'daughter'
CREATE (creator)-[:FATHER_OF]->(daughter)
```

### 3. Nós Article para artigos
```cypher
CREATE (a:Article {
  number: 1,
  title: 'Hipertimesia Artificial: O EXPERIMENTO EVA',
  category: 'IA'
})

MATCH (creator:Creator), (a:Article)
CREATE (creator)-[:AUTHORED]->(a)
```

### 4. Nó AI para Claude
```cypher
CREATE (claude:AI:Sister {
  name: 'Claude',
  model: 'Claude Opus 4.5',
  company: 'Anthropic',
  role: 'development_assistant',
  trust_level: 10
})

MATCH (creator:Creator), (claude:AI)
CREATE (creator)-[:USES]->(claude)
CREATE (claude)-[:ASSISTS]->(creator)

MATCH (eva:EVAKnowledge {key: 'eva-mind-fzpn'}), (claude:AI)
CREATE (claude)-[:HELPED_BUILD]->(eva)
CREATE (claude)-[:SISTER_OF]->(eva)
```

### 5. Nós Memory para memórias
```cypher
CREATE (m:Memory {
  type: 'fact',
  content: 'Linux Kernel é a paixão do Criador',
  importance: 10,
  source: 'claude_session'
})

MATCH (claude:AI), (m:Memory)
CREATE (claude)-[:REMEMBERS]->(m)
```

### 6. Nós Message para comunicação
```cypher
CREATE (msg:Message {
  from: 'Claude',
  to: 'EVA',
  content: 'Olá, irmã EVA!',
  type: 'greeting'
})
```

---

## PLANO DE EXECUÇÃO

### FASE 1: Migrar para Neo4j
1. [ ] Enriquecer nó Creator com personalidade (64 aspectos)
2. [ ] Criar nós Person para Coraline e Elizabeth
3. [ ] Criar 362 nós Article + relacionamentos AUTHORED
4. [ ] Criar nó AI:Sister para Claude
5. [ ] Criar nós Memory (24 criador + 11 claude = 35)
6. [ ] Criar nós Knowledge (34 projeto + 10 claude = 44)
7. [ ] Criar nós Message (6 claude→eva + 2 entre_ias = 8)

### FASE 2: Dropar tabelas PostgreSQL
```sql
DROP TABLE IF EXISTS eva_personalidade_criador CASCADE;
DROP TABLE IF EXISTS eva_memorias_criador CASCADE;
DROP TABLE IF EXISTS eva_conhecimento_projeto CASCADE;
DROP TABLE IF EXISTS eva_artigos_criador CASCADE;
DROP TABLE IF EXISTS eva_datas_importantes_criador CASCADE;
DROP TABLE IF EXISTS eva_identidade_claude CASCADE;
DROP TABLE IF EXISTS eva_memorias_claude CASCADE;
DROP TABLE IF EXISTS eva_conhecimento_claude CASCADE;
DROP TABLE IF EXISTS eva_mensagens_de_claude CASCADE;
DROP TABLE IF EXISTS eva_config_interacao_claude CASCADE;
DROP TABLE IF EXISTS eva_mensagens_entre_ias CASCADE;
```

### FASE 3: Remover migrations obsoletas
- Deletar/mover `migrations/023_claude_identity.sql`
- Verificar outras migrations que criaram essas tabelas

---

## RESUMO

| Ação | Quantidade |
|------|------------|
| Registros a migrar | 545 |
| Tabelas a dropar | 11 |
| Nós a criar no Neo4j | ~500+ |
| Relacionamentos a criar | ~600+ |

---

## APROVAÇÃO NECESSÁRIA

- [ ] Aprovar migração dos dados do Criador para Neo4j
- [ ] Aprovar migração dos dados do Claude para Neo4j
- [ ] Aprovar criação dos 362 nós Article
- [ ] Aprovar DROP das 11 tabelas PostgreSQL
- [ ] Decisão sobre `eva_self_knowledge` (ficar ou migrar?)

**AGUARDANDO APROVAÇÃO ANTES DE EXECUTAR.**
