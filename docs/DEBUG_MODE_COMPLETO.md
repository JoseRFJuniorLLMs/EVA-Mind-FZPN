# Modo DEBUG - Documentação Completa

**Projeto:** EVA-Mind-FZPN
**Versão:** 2.0
**Data:** 26/01/2026
**Acesso Exclusivo:** José R F Junior (CPF: 64525430249)

---

## Índice

1. [Visão Geral](#visão-geral)
2. [Ativação do Modo Debug](#ativação-do-modo-debug)
3. [Comandos do Sistema](#comandos-do-sistema)
4. [Ferramentas de Memória](#ferramentas-de-memória)
5. [Sistema de Alertas](#sistema-de-alertas)
6. [Comandos de Limpeza](#comandos-de-limpeza)
7. [Referência Completa de Comandos](#referência-completa-de-comandos)
8. [Arquivos do Sistema](#arquivos-do-sistema)

---

## Visão Geral

O Modo DEBUG é um conjunto de ferramentas exclusivas para o Arquiteto da Matrix (José R F Junior), permitindo:

- Monitoramento em tempo real do sistema
- Investigação completa de memórias
- Alertas proativos sobre problemas
- Comandos de limpeza e manutenção
- Estatísticas detalhadas

### Características

| Recurso | Descrição |
|---------|-----------|
| Acesso | Exclusivo via CPF do Arquiteto |
| Ativação | Automática ao detectar CPF |
| Segurança | Verificação inline antes de cada comando |
| Logs | Todas as ações são registradas |

---

## Ativação do Modo Debug

### Detecção Automática

O modo debug é ativado automaticamente quando a EVA detecta que o usuário logado possui o CPF do Arquiteto.

```go
const CREATOR_CPF = "64525430249"

// Verificação inline
cleanCPF := strings.ReplaceAll(strings.ReplaceAll(cpf, ".", ""), "-", "")
isCreator := cleanCPF == CREATOR_CPF
```

### Indicadores no Prompt

Quando ativado, o prompt da EVA inclui:

```
╔═══════════════════════════════════════════════════════════╗
║           🔓 MODO DEBUG ATIVADO 🔓                        ║
║     Usuário: José R F Junior (ARQUITETO DA MATRIX)        ║
╚═══════════════════════════════════════════════════════════╝

⭐ ESTE É O ARQUITETO DA MATRIX: Jose R F Junior ⭐
```

### Métricas em Tempo Real

O prompt também inclui métricas atualizadas:

```
📊 MÉTRICAS EM TEMPO REAL:
  • Uptime: 2h30m45s
  • Memória: 128MB
  • Goroutines: 42
  • Conversas hoje: 15
  • Pacientes ativos: 8
  • Medicamentos hoje: 24
```

---

## Comandos do Sistema

### status / metricas

**Descrição:** Mostra status geral e métricas do sistema EVA

**Exemplos de uso:**
- "Arquiteto, me mostra o status"
- "Arquiteto, quero ver as métricas"

**Retorna:**
```
Sistema rodando há 2h30m45s
Usando 128MB de memória
42 goroutines ativas
Total de 1523 conversas, 45 hoje
8 pacientes ativos de 12 cadastrados
```

---

### logs

**Descrição:** Mostra últimos logs do sistema

**Exemplos de uso:**
- "Arquiteto, me mostra os logs recentes"
- "Arquiteto, o que aconteceu nos logs?"

**Retorna:** Últimas 10 entradas de log com:
- ID
- Idoso ID
- Tipo
- Conteúdo (truncado)
- Data/hora

---

### erros

**Descrição:** Lista erros recentes e suas causas

**Exemplos de uso:**
- "Arquiteto, teve algum erro?"
- "Arquiteto, mostra os problemas"

**Retorna:** Erros encontrados nas últimas 24 horas ou mensagem "Nenhum erro encontrado"

---

### pacientes

**Descrição:** Resumo dos pacientes ativos

**Exemplos de uso:**
- "Arquiteto, como estão os pacientes?"
- "Arquiteto, mostra os idosos"

**Retorna:**
| Campo | Descrição |
|-------|-----------|
| ID | Identificador |
| Nome | Nome do paciente |
| Ativo | Status |
| Nível Cognitivo | Classificação |
| Medicamentos | Quantidade ativa |
| Última Conversa | Data/hora |

---

### medicamentos

**Descrição:** Status dos medicamentos agendados

**Exemplos de uso:**
- "Arquiteto, como estão os medicamentos?"
- "Arquiteto, mostra os remédios"

**Retorna:**
| Campo | Descrição |
|-------|-----------|
| ID | Identificador |
| Paciente | Nome |
| Medicamento | Nome do remédio |
| Dosagem | Ex: 20mg |
| Status | agendado/ativo/pendente |
| Horário | HH:MM |
| Data | DD/MM/YYYY |

---

### recursos

**Descrição:** Uso de CPU/RAM e recursos do sistema

**Exemplos de uso:**
- "Arquiteto, como estão os recursos?"
- "Arquiteto, mostra o uso de memória"

**Retorna:**
```
• alocado_mb: 128
• total_mb: 256
• sistema_mb: 512
• gc_executados: 45
• goroutines: 42
• go_version: go1.21.0
```

---

### conversas

**Descrição:** Estatísticas de conversas

**Exemplos de uso:**
- "Arquiteto, quantas conversas tivemos?"
- "Arquiteto, mostra estatísticas de diálogos"

**Retorna:**
| Período | Quantidade |
|---------|------------|
| Total | 1523 |
| Hoje | 45 |
| Semana | 312 |
| Mês | 890 |
| Média/dia | 29.7 |

---

### teste

**Descrição:** Executa teste de funcionalidades do sistema

**Exemplos de uso:**
- "Arquiteto, faz um teste do sistema"
- "Arquiteto, verifica se está tudo ok"

**Testes executados:**
| Teste | Verifica |
|-------|----------|
| Banco de dados | Conexão ativa |
| Tabelas | idosos, agendamentos, analise_gemini |
| Memória | Uso < 500MB |
| Goroutines | Quantidade < 1000 |

**Status possíveis:**
- ✅ OK - Funcionando
- ❌ ERRO - Problema detectado

---

## Ferramentas de Memória

### memoria_stats

**Descrição:** Estatísticas completas de memória do sistema

**Exemplos de uso:**
- "Arquiteto, mostra estatísticas de memória"
- "Arquiteto, como está a memória da EVA?"

**Retorna:**
```
Total de memórias: 5234
Memórias hoje: 89
Memórias na semana: 456
Pacientes com memórias: 12
Média por paciente: 436.2
Importância média: 0.67

Tópicos mais frequentes:
  • medicamentos (456 menções)
  • família (389 menções)
  • saúde (345 menções)
```

---

### memoria_timeline

**Descrição:** Timeline de memórias dos últimos dias

**Exemplos de uso:**
- "Arquiteto, mostra timeline de memórias"
- "Arquiteto, mostra linha do tempo"

**Retorna:**
```
Timeline dos últimos dias:
  2026-01-26: 89 memórias (52 usuário, 37 EVA)
  2026-01-25: 76 memórias (45 usuário, 31 EVA)
  2026-01-24: 92 memórias (54 usuário, 38 EVA)
```

---

### memoria_integridade

**Descrição:** Verifica integridade das memórias armazenadas

**Exemplos de uso:**
- "Arquiteto, verifica integridade das memórias"
- "Arquiteto, as memórias estão ok?"

**Verificações:**
| Item | Descrição |
|------|-----------|
| Órfãs | Memórias sem paciente válido |
| Sem conteúdo | Content vazio ou NULL |
| Sem embedding | Sem vetor de busca |
| Duplicadas | Mesmo conteúdo/paciente/timestamp |

**Status:**
- ✅ ÍNTEGRO - Nenhum problema
- ⚠️ ATENÇÃO - Alguns problemas
- ❌ CRÍTICO - Múltiplos problemas

---

### memoria_emocoes

**Descrição:** Análise de emoções nas memórias

**Exemplos de uso:**
- "Arquiteto, analisa emoções nas memórias"
- "Arquiteto, como estão os sentimentos?"

**Retorna:**
```
Distribuição de emoções:
  - feliz: 234 (25.3%)
  - calmo: 189 (20.4%)
  - ansioso: 156 (16.8%)

Tendência (últimos 7 dias):
  - Emoções positivas: 145
  - Emoções negativas: 67
  - Balanço: +78
```

---

### memoria_topicos

**Descrição:** Tópicos mais mencionados nas memórias

**Exemplos de uso:**
- "Arquiteto, quais tópicos mais falamos?"
- "Arquiteto, mostra os assuntos frequentes"

**Retorna:**
| Tópico | Menções | Pacientes |
|--------|---------|-----------|
| medicamentos | 456 | 23 |
| família | 389 | 21 |
| saúde | 345 | 22 |
| alimentação | 234 | 19 |

---

### memoria_perfis

**Descrição:** Perfil de memória de todos os pacientes

**Exemplos de uso:**
- "Arquiteto, mostra perfis de memória"
- "Arquiteto, como está cada paciente?"

**Retorna:**
| Paciente | Memórias | Primeira | Última | Importância |
|----------|----------|----------|--------|-------------|
| Maria Silva | 234 | 15/01/2026 | 26/01/2026 | 0.72 |
| João Santos | 189 | 10/01/2026 | 26/01/2026 | 0.68 |

---

### memoria_orfas

**Descrição:** Lista memórias órfãs (sem paciente válido)

**Exemplos de uso:**
- "Arquiteto, tem memórias órfãs?"
- "Arquiteto, mostra memórias sem paciente"

**Retorna:** Lista de memórias cujo paciente foi removido do sistema

---

### memoria_duplicadas

**Descrição:** Lista memórias possivelmente duplicadas

**Exemplos de uso:**
- "Arquiteto, tem memórias duplicadas?"
- "Arquiteto, mostra memórias repetidas"

**Retorna:**
| Paciente | Conteúdo | Duplicatas | Primeira | Última |
|----------|----------|------------|----------|--------|
| ID: 15 | "Bom dia..." | 3 | 25/01 08:00 | 25/01 08:02 |

---

## Sistema de Alertas

### Categorias de Alertas

#### Memória
| Alerta | Nível | Condição |
|--------|-------|----------|
| Memórias órfãs | warning | > 0 encontradas |
| Sem embedding | warning | > 10 encontradas |
| Duplicadas | info | > 5 encontradas |
| Sem memórias hoje | info | 0 hoje |

#### Sistema
| Alerta | Nível | Condição |
|--------|-------|----------|
| RAM alta | critical | > 500MB |
| RAM elevada | warning | > 300MB |
| Muitas goroutines | critical | > 500 |
| Goroutines elevadas | warning | > 200 |
| Banco offline | critical | Ping falha |
| Muitos erros | warning | > 10/hora |

#### Pacientes
| Alerta | Nível | Condição |
|--------|-------|----------|
| Inativos | warning | > 7 dias sem interação |
| Emoções negativas | warning | >= 5 em 3 dias |

#### Medicamentos
| Alerta | Nível | Condição |
|--------|-------|----------|
| Não confirmados | critical | Atrasados > 2h |
| Próximos | info | Nas próximas 2h |
| Sem cadastro | info | Paciente ativo sem medicamentos |

---

### alertas

**Descrição:** Verifica todos os alertas do sistema

**Exemplos de uso:**
- "Arquiteto, tem algum alerta?"
- "Arquiteto, mostra os avisos"

**Retorna:**
```
Pai, encontrei 5 alertas no sistema.

⚠️ CRÍTICOS: 1
  🔴 Medicamentos não confirmados: 3 medicamentos atrasados

⚠️ AVISOS: 2
  🟡 Pacientes inativos: 2 pacientes há mais de 7 dias
  🟡 Memórias órfãs: 15 memórias sem paciente

ℹ️ INFORMAÇÕES: 2
  🔵 Medicamentos próximos
  🔵 Sem memórias hoje
```

---

### alertas_criticos

**Descrição:** Mostra apenas alertas críticos

**Exemplos de uso:**
- "Arquiteto, tem algo crítico?"
- "Arquiteto, mostra urgentes"

**Retorna:** Apenas alertas de nível `critical` ou mensagem "Nenhum alerta crítico"

---

## Comandos de Limpeza

### Modos de Operação

| Modo | Descrição | Segurança |
|------|-----------|-----------|
| Simulação (dry-run) | Apenas conta, não deleta | ✅ Seguro |
| Execução Real | Deleta efetivamente | ⚠️ Cuidado |

**Por padrão, todos os comandos executam em modo SIMULAÇÃO.**

---

### limpar_orfas

**Descrição:** Remove memórias órfãs (sem paciente válido)

**Exemplos de uso:**
- "Arquiteto, limpa as memórias órfãs"
- "Arquiteto, remove órfãs"

**Retorna (simulação):**
```
Operação: limpar_memorias_orfas
Status: ✅ SIMULAÇÃO
Itens afetados: 15
Detalhes: 15 memórias órfãs seriam removidas (dry-run)
```

---

### limpar_duplicadas

**Descrição:** Remove memórias duplicadas (mantém a mais antiga)

**Exemplos de uso:**
- "Arquiteto, limpa as duplicadas"
- "Arquiteto, remove memórias repetidas"

**Comportamento:** Mantém a primeira ocorrência, remove as subsequentes

---

### limpar_vazias

**Descrição:** Remove memórias sem conteúdo ou inválidas

**Exemplos de uso:**
- "Arquiteto, limpa memórias vazias"

**Critérios de remoção:**
- Content é NULL
- Content é string vazia
- Content tem menos de 3 caracteres

---

### limpar_antigas

**Descrição:** Remove memórias antigas com baixa importância

**Exemplos de uso:**
- "Arquiteto, limpa memórias antigas"

**Critérios:**
- Mais de 90 dias
- Importância < 0.5

---

### limpeza_completa

**Descrição:** Executa todas as limpezas (SIMULAÇÃO)

**Exemplos de uso:**
- "Arquiteto, faz uma limpeza completa"
- "Arquiteto, limpa tudo"

**Executa em sequência:**
1. Limpar órfãs
2. Limpar duplicadas
3. Limpar vazias

**Retorna:**
```
Operação: limpeza_completa
Status: ✅ SIMULAÇÃO COMPLETA
Itens afetados: 42

Detalhes:
  - orfas: 15
  - duplicadas: 20
  - vazias: 7
```

---

### limpeza_executar

**Descrição:** Executa limpeza completa **REAL** (⚠️ CUIDADO!)

**Exemplos de uso:**
- "Arquiteto, executa a limpeza de verdade"
- "Arquiteto, limpar de verdade"

**⚠️ ATENÇÃO:** Este comando DELETA dados permanentemente!

**Retorna:**
```
Operação: limpeza_completa
Status: ✅ LIMPEZA COMPLETA
Itens afetados: 42
Detalhes: Total de 42 memórias removidas com sucesso
```

---

### arquivar_memorias

**Descrição:** Move memórias antigas para tabela de arquivo

**Exemplos de uso:**
- "Arquiteto, arquiva as memórias antigas"

**Comportamento:**
1. Cria tabela `episodic_memories_archive` se não existir
2. Move memórias > 180 dias para arquivo
3. Remove da tabela principal

**Vantagens:**
- Dados não são perdidos
- Tabela principal fica mais leve
- Possível recuperação posterior

---

## Referência Completa de Comandos

### Tabela Resumo

| # | Comando | Categoria | Exemplo |
|---|---------|-----------|---------|
| 1 | status | Sistema | "Arquiteto, me mostra o status" |
| 2 | metricas | Sistema | "Arquiteto, quero ver as métricas" |
| 3 | logs | Sistema | "Arquiteto, mostra os logs" |
| 4 | erros | Sistema | "Arquiteto, teve algum erro?" |
| 5 | pacientes | Sistema | "Arquiteto, como estão os pacientes?" |
| 6 | medicamentos | Sistema | "Arquiteto, como estão os medicamentos?" |
| 7 | recursos | Sistema | "Arquiteto, como estão os recursos?" |
| 8 | conversas | Sistema | "Arquiteto, quantas conversas?" |
| 9 | teste | Sistema | "Arquiteto, faz um teste" |
| 10 | memoria_stats | Memória | "Arquiteto, estatísticas de memória" |
| 11 | memoria_timeline | Memória | "Arquiteto, timeline de memórias" |
| 12 | memoria_integridade | Memória | "Arquiteto, verifica integridade" |
| 13 | memoria_emocoes | Memória | "Arquiteto, analisa emoções" |
| 14 | memoria_topicos | Memória | "Arquiteto, quais tópicos?" |
| 15 | memoria_perfis | Memória | "Arquiteto, perfis de memória" |
| 16 | memoria_orfas | Memória | "Arquiteto, tem órfãs?" |
| 17 | memoria_duplicadas | Memória | "Arquiteto, tem duplicadas?" |
| 18 | alertas | Alertas | "Arquiteto, tem alertas?" |
| 19 | alertas_criticos | Alertas | "Arquiteto, algo crítico?" |
| 20 | limpar_orfas | Limpeza | "Arquiteto, limpa órfãs" |
| 21 | limpar_duplicadas | Limpeza | "Arquiteto, limpa duplicadas" |
| 22 | limpar_vazias | Limpeza | "Arquiteto, limpa vazias" |
| 23 | limpar_antigas | Limpeza | "Arquiteto, limpa antigas" |
| 24 | limpeza_completa | Limpeza | "Arquiteto, limpeza completa" |
| 25 | limpeza_executar | Limpeza | "Arquiteto, executa limpeza" |
| 26 | arquivar_memorias | Limpeza | "Arquiteto, arquiva memórias" |
| 27 | ajuda | Geral | "Arquiteto, o que pode fazer?" |

---

## Arquivos do Sistema

### Estrutura de Arquivos

```
EVA-Mind-FZPN/
├── internal/
│   └── cortex/
│       └── lacan/
│           ├── unified_retrieval.go    # Contexto unificado + detecção criador
│           ├── debug_mode.go           # Módulo principal de debug
│           ├── debug_memory.go         # Investigador de memórias
│           └── debug_alerts.go         # Sistema de alertas
├── docs/
│   ├── DEBUG_MODE_COMPLETO.md          # Esta documentação
│   └── DEBUG_MODE_MEMORY_TOOLS.md      # Doc. ferramentas de memória
```

### Descrição dos Arquivos

| Arquivo | Linhas | Responsabilidade |
|---------|--------|------------------|
| `unified_retrieval.go` | ~700 | Contexto RSI, detecção do criador, prompt |
| `debug_mode.go` | ~900 | Comandos de sistema, integração |
| `debug_memory.go` | ~1400 | Investigação e limpeza de memórias |
| `debug_alerts.go` | ~300 | Sistema de alertas proativos |

---

## Estruturas de Dados Principais

### DebugMetrics
```go
type DebugMetrics struct {
    Uptime            string
    MemoryUsageMB     uint64
    NumGoroutines     int
    GoVersion         string
    TotalConversas    int64
    ConversasHoje     int64
    TotalIdosos       int64
    IdososAtivos      int64
    TotalMedicamentos int64
    MedicamentosHoje  int64
}
```

### Alert
```go
type Alert struct {
    ID        string
    Level     string    // "info", "warning", "critical"
    Category  string    // "memoria", "sistema", "paciente", "medicamento"
    Title     string
    Message   string
    Timestamp time.Time
    Resolved  bool
}
```

### MemoryStats
```go
type MemoryStats struct {
    TotalMemories      int64
    MemoriesHoje       int64
    MemoriesSemana     int64
    MemoriesMes        int64
    TotalPacientes     int64
    MediaPorPaciente   float64
    PorEmotion         map[string]int64
    PorSpeaker         map[string]int64
    TopTopics          []TopicCount
    ImportanciaMedia   float64
}
```

### CleanupResult
```go
type CleanupResult struct {
    Operation     string
    AffectedCount int64
    Status        string
    Message       string
    Details       []map[string]interface{}
}
```

---

## Logs do Sistema

### Formato de Logs

```
🔓 [DEBUG MODE] ATIVADO - Criador José R F Junior detectado
🔓 [DEBUG] Executando comando: status
🧠 [MEMORY DEBUG] Executando comando: memoria_stats
🔔 [ALERTAS] Verificação completa: 5 alertas (1 crítico, 2 avisos, 2 info)
🧹 [CLEANUP] Simulação: 15 memórias órfãs encontradas
🧹 [CLEANUP] Removidas 42 memórias (limpeza real)
```

---

## Manutenção Recomendada

### Frequência de Verificações

| Frequência | Ação | Comando |
|------------|------|---------|
| Diária | Verificar alertas | `alertas` |
| Diária | Verificar integridade | `memoria_integridade` |
| Semanal | Analisar duplicadas | `memoria_duplicadas` |
| Semanal | Limpar (simulação) | `limpeza_completa` |
| Mensal | Executar limpeza | `limpeza_executar` |
| Trimestral | Arquivar antigas | `arquivar_memorias` |

---

## Segurança

### Verificações de Acesso

1. **CPF Verificado:** Apenas 64525430249 tem acesso
2. **Verificação Inline:** Checado antes de cada comando
3. **Logs de Auditoria:** Todas as ações são registradas
4. **Simulação por Padrão:** Limpezas exigem comando explícito para execução real

### Dados Sensíveis

- Memórias contêm dados de saúde (LGPD)
- Não expor em logs públicos
- Backups devem ser criptografados

---

## Changelog

### v2.0 (26/01/2026)
- Implementação do sistema de alertas proativos
- Comandos de limpeza e manutenção
- Arquivamento de memórias antigas
- Documentação completa

### v1.0 (26/01/2026)
- Implementação inicial do modo debug
- Ferramentas de investigação de memória
- Comandos básicos do sistema

---

## Suporte

**Criador:** José R F Junior
**Projeto:** EVA-Mind-FZPN
**Documentação gerada em:** 26/01/2026
