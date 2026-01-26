# Ferramentas de Investigação de Memória - Modo DEBUG

**Versão:** 1.0
**Data:** 26/01/2026
**Acesso:** Exclusivo para o Arquiteto da Matrix (José R F Junior - CPF: 64525430249)

---

## Visão Geral

O Modo DEBUG da EVA inclui ferramentas avançadas de investigação de memória que permitem ao Arquiteto analisar, verificar e manter a integridade das memórias armazenadas no sistema.

### Arquivos Relacionados

| Arquivo | Descrição |
|---------|-----------|
| `internal/cortex/lacan/debug_mode.go` | Módulo principal do modo debug |
| `internal/cortex/lacan/debug_memory.go` | Investigador de memórias |
| `internal/cortex/lacan/unified_retrieval.go` | Integração com contexto unificado |

---

## Comandos de Memória Disponíveis

### Comandos Básicos

| Comando | Descrição | Exemplo de Uso |
|---------|-----------|----------------|
| `memoria_stats` | Estatísticas completas de memória | "Arquiteto, mostra estatísticas de memória" |
| `memoria_timeline` | Timeline de memórias (últimos 14 dias) | "Arquiteto, mostra timeline de memórias" |
| `memoria_integridade` | Verifica integridade das memórias | "Arquiteto, verifica integridade das memórias" |
| `memoria_emocoes` | Análise de emoções nas memórias | "Arquiteto, analisa emoções nas memórias" |
| `memoria_topicos` | Tópicos mais mencionados | "Arquiteto, quais tópicos mais falamos?" |
| `memoria_perfis` | Perfil de memória de todos pacientes | "Arquiteto, mostra perfis de memória" |
| `memoria_orfas` | Lista memórias órfãs (sem paciente) | "Arquiteto, tem memórias órfãs?" |
| `memoria_duplicadas` | Lista memórias duplicadas | "Arquiteto, tem memórias duplicadas?" |

---

## Funcionalidades Detalhadas

### 1. Estatísticas de Memória (`memoria_stats`)

Retorna estatísticas completas do sistema de memórias:

```
- Total de memórias armazenadas
- Memórias criadas hoje/semana/mês
- Total de pacientes com memórias
- Média de memórias por paciente
- Memória mais antiga e mais recente
- Distribuição por emoção
- Distribuição por speaker (user/assistant)
- Top 10 tópicos mais frequentes
- Importância média das memórias
- Tamanho médio em bytes
```

**Estrutura de Dados:**
```go
type MemoryStats struct {
    TotalMemories      int64
    MemoriesHoje       int64
    MemoriesSemana     int64
    MemoriesMes        int64
    TotalPacientes     int64
    MediaPorPaciente   float64
    MemoriasMaisAntiga time.Time
    MemoriaMaisRecente time.Time
    PorEmotion         map[string]int64
    PorSpeaker         map[string]int64
    TopTopics          []TopicCount
    ImportanciaMedia   float64
    TamanhoMedioBytes  int64
}
```

---

### 2. Timeline de Memórias (`memoria_timeline`)

Mostra a linha do tempo de memórias dos últimos dias:

```
- Data
- Total de memórias no dia
- Mensagens do usuário
- Mensagens da EVA
- Emoções detectadas no dia
```

**Exemplo de Saída:**
```
Timeline dos últimos dias:
  2026-01-26: 45 memórias (28 usuário, 17 EVA)
  2026-01-25: 38 memórias (22 usuário, 16 EVA)
  2026-01-24: 52 memórias (30 usuário, 22 EVA)
```

---

### 3. Verificação de Integridade (`memoria_integridade`)

Verifica a saúde das memórias armazenadas:

| Verificação | Descrição |
|-------------|-----------|
| Memórias órfãs | Memórias sem paciente válido associado |
| Sem conteúdo | Memórias com content vazio ou NULL |
| Sem embedding | Memórias sem vetor de embedding |
| Duplicadas | Mesmo conteúdo, paciente e timestamp |

**Status Possíveis:**
- ✅ ÍNTEGRO - Nenhum problema encontrado
- ⚠️ ATENÇÃO - Alguns problemas detectados
- ❌ CRÍTICO - Múltiplos problemas detectados

---

### 4. Análise de Emoções (`memoria_emocoes`)

Analisa as emoções presentes nas memórias:

```
Distribuição de emoções:
  - feliz: 234 (25.3%)
  - calmo: 189 (20.4%)
  - ansioso: 156 (16.8%)
  - triste: 98 (10.6%)
  ...

Tendência (últimos 7 dias):
  - Emoções positivas: 145
  - Emoções negativas: 67
  - Balanço: +78
```

**Emoções Positivas Rastreadas:**
- feliz, alegre, satisfeito, calmo, esperançoso

**Emoções Negativas Rastreadas:**
- triste, ansioso, irritado, preocupado, frustrado

---

### 5. Análise de Tópicos (`memoria_topicos`)

Lista os tópicos mais mencionados nas memórias:

```
Tópicos mais frequentes:
  1. medicamentos (456 menções, 23 pacientes)
  2. família (389 menções, 21 pacientes)
  3. saúde (345 menções, 22 pacientes)
  4. alimentação (234 menções, 19 pacientes)
  5. sono (198 menções, 18 pacientes)
```

---

### 6. Perfis de Memória (`memoria_perfis`)

Mostra resumo de memórias por paciente:

| Paciente | Memórias | Primeira | Última | Importância |
|----------|----------|----------|--------|-------------|
| Maria Silva | 234 | 15/01/2026 | 26/01/2026 | 0.72 |
| João Santos | 189 | 10/01/2026 | 26/01/2026 | 0.68 |
| Ana Costa | 156 | 20/01/2026 | 25/01/2026 | 0.75 |

---

### 7. Memórias Órfãs (`memoria_orfas`)

Lista memórias que não têm paciente válido associado (paciente foi removido):

```
Memórias órfãs encontradas: 12

  [ID: 4523] PACIENTE REMOVIDO - "Tomei o remédio às 8h..."
  [ID: 4524] PACIENTE REMOVIDO - "Estou me sentindo bem..."
```

**Ação Recomendada:** Revisar e possivelmente excluir essas memórias.

---

### 8. Memórias Duplicadas (`memoria_duplicadas`)

Lista memórias possivelmente duplicadas:

```
Possíveis duplicatas encontradas: 8

  Paciente ID: 15
  Conteúdo: "Bom dia, como você está?"
  Duplicatas: 3
  Primeira: 25/01/2026 08:00
  Última: 25/01/2026 08:02
```

---

## Estruturas de Dados Principais

### MemoryDetail
```go
type MemoryDetail struct {
    ID            int64
    IdosoID       int64
    IdosoNome     string
    Timestamp     time.Time
    Speaker       string    // "user" ou "assistant"
    Content       string
    ContentLength int
    Emotion       string
    Importance    float64
    Topics        []string
    SessionID     string
    HasEmbedding  bool
}
```

### PatientMemoryProfile
```go
type PatientMemoryProfile struct {
    IdosoID           int64
    Nome              string
    TotalMemories     int64
    PrimeiraMemoria   time.Time
    UltimaMemoria     time.Time
    EmocoesMaisComuns []string
    TopicosFrequentes []string
    ImportanciaMedia  float64
    SessoesUnicas     int64
    MemoriasPorMes    map[string]int64
}
```

---

## Como Usar

### Ativação Automática

O modo debug é ativado automaticamente quando a EVA detecta que o usuário logado é o Arquiteto (CPF: 64525430249).

### Comandos por Voz

Basta falar naturalmente com a EVA:

```
"Arquiteto, me mostra as estatísticas de memória"
"Arquiteto, verifica se tem memórias com problemas"
"Arquiteto, quais são os tópicos mais falados?"
"Arquiteto, tem memórias duplicadas no sistema?"
```

### Detecção de Comandos

O sistema detecta automaticamente palavras-chave:

| Palavras-Chave | Comando Ativado |
|----------------|-----------------|
| "estatísticas de memória", "stats de memória" | `memoria_stats` |
| "timeline", "linha do tempo" | `memoria_timeline` |
| "integridade", "verificar memórias" | `memoria_integridade` |
| "emoções", "sentimentos" | `memoria_emocoes` |
| "tópicos", "assuntos" | `memoria_topicos` |
| "perfis de memória" | `memoria_perfis` |
| "órfãs", "sem paciente" | `memoria_orfas` |
| "duplicadas", "repetidas" | `memoria_duplicadas` |

---

## Exportação de Dados

### Exportar Memórias de um Paciente

```go
json, err := memoryInvestigator.ExportPatientMemories(ctx, idosoID)
```

Retorna JSON completo com:
- Data de exportação
- Perfil do paciente
- Todas as memórias
- Total de memórias

---

## Tabelas do Banco de Dados

### episodic_memories
```sql
CREATE TABLE episodic_memories (
    id SERIAL PRIMARY KEY,
    idoso_id BIGINT REFERENCES idosos(id),
    timestamp TIMESTAMP DEFAULT NOW(),
    speaker VARCHAR(20),        -- 'user' ou 'assistant'
    content TEXT,
    embedding VECTOR(768),      -- pgvector
    emotion VARCHAR(50),
    importance FLOAT,
    topics TEXT[],
    session_id VARCHAR(100),
    call_history_id BIGINT
);
```

---

## Logs do Sistema

Quando comandos de memória são executados, logs são gerados:

```
🧠 [MEMORY DEBUG] Executando comando: memoria_stats
✅ [MemoryInvestigator] Estatísticas coletadas com sucesso
🔓 [DEBUG MODE] Resposta formatada para o Arquiteto
```

---

## Considerações de Segurança

1. **Acesso Restrito:** Apenas o CPF do Arquiteto (64525430249) pode acessar essas funcionalidades
2. **Verificação Inline:** CPF é verificado antes de cada comando
3. **Logs de Auditoria:** Todos os comandos de debug são logados
4. **Dados Sensíveis:** Memórias contêm dados de saúde - manter confidencialidade

---

## Manutenção Recomendada

| Frequência | Ação |
|------------|------|
| Diária | Verificar `memoria_integridade` |
| Semanal | Analisar `memoria_duplicadas` e limpar |
| Mensal | Exportar backup de memórias importantes |
| Trimestral | Revisar e arquivar memórias antigas |

---

## Changelog

### v1.0 (26/01/2026)
- Implementação inicial do MemoryInvestigator
- 8 comandos de memória disponíveis
- Integração com modo debug
- Documentação completa

---

**Desenvolvido para:** EVA-Mind-FZPN
**Arquiteto:** José R F Junior
