# Sistema de Tipos Gurdjieff (Eneagrama) - EVA-Mind-FZPN

**Documento:** GURDJIEFF-001
**Versão:** 1.0
**Data:** 2026-01-27
**Autor:** Jose R F Junior

---

## 1. O Que São os Tipos Gurdjieff

O Eneagrama de Gurdjieff identifica **9 padrões de personalidade** que revelam:
- **Como a pessoa está "aprisionada"** em padrões automáticos
- **Qual emoção raiz** domina (raiva, medo, vergonha)
- **Qual mecanismo de defesa** ela usa inconscientemente

**PRINCÍPIO:** "Não identificar QUEM o paciente É, mas COMO está preso."

---

## 2. Os 9 Tipos

| Tipo | Nome PT | Centro | Emoção Raiz | Chief Feature |
|------|---------|--------|-------------|---------------|
| **1** | Perfeccionista | Instintivo | Raiva (reprimida) | Ressentimento e autocrítica |
| **2** | Ajudador | Emocional | Vergonha (negada) | Orgulho de ser necessário |
| **3** | Realizador | Emocional | Vergonha (evitada) | Vaidade e manipulação de imagem |
| **4** | Individualista | Emocional | Vergonha (internalizada) | Inveja e sentir-se deficiente |
| **5** | Investigador | Mental | Medo (de intrusão) | Avareza de recursos e energia |
| **6** | Leal | Mental | Medo (de abandono) | Dúvida e suspeita |
| **7** | Entusiasta | Mental | Medo (de dor) | Gula por experiências |
| **8** | Desafiador | Instintivo | Raiva (expressada) | Luxúria por poder e controle |
| **9** | Pacificador | Instintivo | Raiva (negada) | Auto-esquecimento e fusão |

---

## 3. Funções no Sistema EVA

### 3.1 Filtro de Atenção da EVA (Principal)

O tipo detectado **muda o FOCO e TOM da EVA**.

**Arquivo:** `internal/cortex/gemini/prompts.go:29-41`

```go
switch eneatype {
case 2: // Ajudante
    typeDirective = "FOCO ATUAL: Empatia máxima e cuidado prático. Seja suave e acolhedora."
case 6: // Leal/Segurança
    typeDirective = "FOCO ATUAL: Segurança e precisão. Transmita confiança e autoridade calma."
case 9: // Pacificador (Base)
    typeDirective = "FOCO ATUAL: Harmonia e escuta ativa. Evite conflitos e mantenha o tom estável."
}
```

| Tipo Paciente | EVA Adapta Para |
|---------------|-----------------|
| Tipo 2 (Ajudador) | Empatia máxima, cuidado prático, suavidade |
| Tipo 6 (Leal) | Segurança, precisão, confiança, autoridade calma |
| Tipo 9 (Pacificador) | Harmonia, escuta ativa, tom estável |

### 3.2 Roteamento Lacaniano (Intervenção Clínica)

O tipo é determinado pela **análise de transferência e desejo**.

**Arquivo:** `internal/cortex/lacan/zeta_router.go:180-198`

```go
func (z *ZetaRouter) DetermineGurdjieffType(...) int {
    // Tipo 2: Se há necessidade de cuidado maternal
    if desejo == DESEJO_AMOR || desejo == DESEJO_COMPANHIA ||
       transferencia == TRANSFERENCIA_MATERNA {
        return 2  // EVA age como Ajudante/Cuidadora
    }

    // Tipo 6: Se há ansiedade ou busca de estrutura
    if desejo == DESEJO_CONTROLE ||
       transferencia == TRANSFERENCIA_PATERNA {
        return 6  // EVA age como figura de Segurança
    }

    // Tipo 9: Padrão harmonioso
    return 9  // EVA mantém neutralidade empática
}
```

### 3.3 Detecção de Armadilhas Psíquicas

O sistema identifica **como o paciente se sabota**:

| Tipo | Armadilha (Chief Feature) | Frases Detectadas |
|------|---------------------------|-------------------|
| 1 | Ressentimento, autocrítica | "deveria ter feito", "não está certo" |
| 2 | Orgulho de ser necessário | "preciso ajudar", "depende de mim" |
| 3 | Vaidade, manipulação de imagem | "tenho que conseguir", "parecer" |
| 4 | Inveja, sentir-se deficiente | "ninguém me entende", "sou diferente" |
| 5 | Avareza de energia | "preciso pensar", "deixa eu analisar" |
| 6 | Dúvida, suspeita | "e se der errado", "não confio" |
| 7 | Gula por experiências | "outras opções", "não quero ficar parado" |
| 8 | Negação de vulnerabilidade | "não vou deixar", "sou forte" |
| 9 | Auto-esquecimento | "tanto faz", "não quero conflito" |

### 3.4 Mirror Output (Espelhamento Terapêutico)

Quando há confiança suficiente (>0.3), EVA pode **espelhar o padrão** para o paciente:

```
PADRÃO DETECTADO:
- Centro: Emocional
- Emoção raiz: Vergonha (negada)
- Traço principal: "Orgulho de ser necessário" (detectado 12x)

"Você percebe esse padrão em si mesmo? O que acha que isso significa?"
```

**OBJETIVO:** Não diagnosticar, mas fazer o paciente **perceber** seu próprio padrão.

---

## 4. Detecção de Tipos

### 4.1 Keywords por Tipo

**Arquivo:** `internal/hippocampus/memory/superhuman/enneagram_service.go`

| Tipo | Keywords (Português) |
|------|---------------------|
| 1 | "certo", "errado", "deveria", "precisa", "correto", "perfeito", "falha", "erro" |
| 2 | "precisa de mim", "deixa eu ajudar", "faço por você", "cuido", "ajudo" |
| 3 | "consegui", "sucesso", "melhor", "eficiente", "resultado", "reconhecimento" |
| 4 | "ninguém entende", "diferente", "especial", "vazio", "saudade", "único" |
| 5 | "penso", "estudo", "preciso entender", "sozinho", "observo", "analiso" |
| 6 | "e se", "cuidado", "confiança", "seguro", "dúvida", "preocupado", "medo" |
| 7 | "legal", "divertido", "plano", "opção", "possibilidade", "novo", "aventura" |
| 8 | "forte", "luta", "controle", "poder", "proteger", "enfrento", "decido" |
| 9 | "tanto faz", "não sei", "talvez", "deixa quieto", "paz", "calma", "tranquilo" |

### 4.2 Pesos de Evidência

| Categoria | Peso | Descrição |
|-----------|------|-----------|
| `keyword` | 0.5 | Palavra-chave simples detectada |
| `defense_mechanism` | 0.7 | Mecanismo de defesa identificado |
| `chief_feature` | 0.8 | Traço principal do tipo |

### 4.3 Mecanismos de Defesa

```go
defensePatterns := map[int][]string{
    1: {"mas eu estava certo", "o certo seria"},
    2: {"eu só quero ajudar", "faço por amor"},
    3: {"estou muito ocupado", "trabalhando muito"},
    4: {"você não entenderia", "é complicado"},
    5: {"preciso de mais tempo", "deixa eu ver"},
    6: {"e se eles", "acho que querem"},
    7: {"mas olha o lado bom", "podia ser pior"},
    8: {"não tenho medo", "não me afeta"},
    9: {"não tem problema", "está tudo bem"},
}
```

---

## 5. Armazenamento de Dados

### 5.1 Tabela: `enneagram_evidence`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | SERIAL | ID da evidência |
| `idoso_id` | INTEGER | ID do paciente |
| `memory_id` | INTEGER | ID da memória que gerou a evidência |
| `verbatim` | TEXT | Frase exata do paciente |
| `suggested_type` | INTEGER | Tipo detectado (1-9) |
| `weight` | FLOAT | Peso (0.5, 0.7, 0.8) |
| `category` | VARCHAR | "keyword", "defense_mechanism", "chief_feature" |
| `context` | TEXT | Padrão que foi detectado |
| `timestamp` | TIMESTAMP | Quando foi detectado |

### 5.2 Tabela: `patient_enneagram`

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `idoso_id` | INTEGER | ID do paciente |
| `primary_type` | INTEGER | Tipo mais provável (1-9) |
| `primary_type_confidence` | FLOAT | Confiança (0.0-1.0) |
| `dominant_wing` | INTEGER | Asa dominante |
| `wing_influence` | FLOAT | Influência da asa |
| `health_level` | INTEGER | Nível de saúde (1-9) |
| `instinctual_variant` | VARCHAR | sp/sx/so |
| `type_scores` | JSONB | Scores de todos os tipos |
| `evidence_count` | INTEGER | Total de evidências |
| `last_evidence_at` | TIMESTAMP | Última evidência |
| `identification_method` | VARCHAR | Como foi identificado |
| `identified_at` | TIMESTAMP | Quando foi identificado |

---

## 6. Fluxo de Detecção

```
┌─────────────────────────────────────────────────────────────┐
│                    PACIENTE FALA                            │
│           "Preciso cuidar de todo mundo, sem mim            │
│            eles não conseguem fazer nada"                   │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  1. EnneagramService.AnalyzeText()                          │
│     - Detecta keywords: "preciso", "cuidar", "sem mim"      │
│     - Identifica chief feature: "sem mim" → Tipo 2          │
│     - Salva evidência (weight: 0.8)                         │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  2. ZetaRouter.DetermineGurdjieffType()                     │
│     - Analisa transferência: MATERNA                        │
│     - Analisa desejo: COMPANHIA                             │
│     - Retorna: Tipo 2 (Ajudante)                            │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  3. BuildSystemPrompt(eneatype=2)                           │
│     - Injeta: "FOCO ATUAL: Empatia máxima e cuidado         │
│       prático. Seja suave e acolhedora."                    │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  4. EVA RESPONDE                                            │
│     (Com tom adaptado ao Tipo 2)                            │
│     "Você cuida tanto de todos... e quem cuida de você?"    │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. Resumo das Funções

| Função | Descrição |
|--------|-----------|
| **Adaptar Tom** | EVA muda como fala baseado no tipo |
| **Detectar Armadilhas** | Identifica padrões automáticos inconscientes |
| **Guiar Intervenção** | Decide se deve ser empática, firme, ou neutra |
| **Espelhar Padrões** | Mostra ao paciente o que ele não vê em si |
| **Acumular Evidências** | Construir perfil psicológico ao longo do tempo |

**Em essência:** Os tipos Gurdjieff permitem que a EVA **personalize suas respostas** de acordo com a estrutura psíquica do paciente, em vez de dar respostas genéricas.

---

## 8. Arquivos Relacionados

| Arquivo | Função |
|---------|--------|
| `internal/hippocampus/memory/superhuman/enneagram_service.go` | Detecção de tipos por keywords |
| `internal/cortex/lacan/zeta_router.go` | Roteamento ético e determinação de tipo |
| `internal/cortex/gemini/prompts.go` | Injeção de diretivas no prompt |
| `internal/cortex/lacan/unified_retrieval.go` | Integração dos módulos |
| `migrations/012_superhuman_memory_system.sql` | Schema das tabelas |

---

## Aprovações

| Função | Nome | Data |
|--------|------|------|
| Criador/Admin | Jose R F Junior | 2026-01-27 |

---

**Documento controlado - Versão 1.0**
