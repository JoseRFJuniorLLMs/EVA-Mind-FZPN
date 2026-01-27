# Sprint 1: Governança Cognitiva - EVA-Mind-FZPN

**Documento:** SPRINT1-COGNITIVE-001
**Versão:** 1.0
**Data:** 2026-01-27
**Status:** CONCLUÍDO

---

## Resumo Executivo

O Sprint 1 implementou a **Camada de Governança Cognitiva** do EVA-Mind-FZPN, composta por dois componentes principais:

1. **Meta-Controller Cognitivo (CLO)** - Previne sobrecarga emocional
2. **Ethical Boundary Engine (EBE)** - Previne dependência patológica

---

## 1. Meta-Controller Cognitivo (CLO)

### Objetivo
Monitorar e controlar a carga cognitiva/emocional do paciente para evitar:
- Exaustão emocional
- Saturação cognitiva
- Loops de ruminação

### Arquitetura

```
┌─────────────────────────────────────────────┐
│   COGNITIVE LOAD ORCHESTRATOR (CLO)         │
│   (Camada acima da FZPN)                    │
└──────────────────┬──────────────────────────┘
                   │
    ┌──────────────┼──────────────┐
    ↓              ↓              ↓
┌─────────┐  ┌──────────┐  ┌──────────┐
│TransNAR │  │Zeta Story│  │Affective │
└─────────┘  └──────────┘  └──────────┘
```

### Funcionalidades

| Funcionalidade | Descrição |
|----------------|-----------|
| **Registro de Interações** | Salva cada interação com intensidade emocional, complexidade cognitiva |
| **Detecção de Ruminação** | Identifica mesmo tópico 3x em 2h |
| **Cálculo de Fadiga** | Níveis: none, mild, moderate, severe |
| **Decisões Automáticas** | block, redirect, reduce_frequency, suggest_rest |
| **System Instructions Dinâmicas** | Injeta restrições no prompt do Gemini |
| **Cache Redis** | Estado em tempo real para decisões rápidas |

### Regras de Decisão

```
IF cognitive_load_24h > 0.7 AND recent_high_intensity:
    → BLOCK: Terapia profunda, escalas clínicas
    → ALLOW: Entretenimento leve, música, piadas

IF rumination_detected (mesmo tópico 3x):
    → REDIRECT: "Percebo que isso está te preocupando..."

IF interaction_count_today > 15:
    → SUGGEST_REST: "Você já conversou bastante hoje..."
```

### Arquivos

| Arquivo | Descrição |
|---------|-----------|
| `migrations/003_cognitive_load_and_ethical_boundaries.sql` | Schema das tabelas |
| `internal/cortex/cognitive/cognitive_load_orchestrator.go` | Implementação Go |

### Tabelas Criadas

- `interaction_cognitive_load` - Histórico de interações
- `cognitive_load_state` - Estado atual (cache)
- `cognitive_load_decisions` - Registro de decisões

---

## 2. Ethical Boundary Engine (EBE)

### Objetivo
Prevenir dependência emocional e apego patológico através de:
- Detecção de frases de apego
- Monitoramento do ratio EVA:Humanos
- Análise de dominância de significantes
- Protocolos de redirecionamento

### Arquitetura

```
┌─────────────────────────────────────────┐
│  ETHICAL BOUNDARY ENGINE (EBE)          │
│  - Detecção de apego patológico         │
│  - Limites de intimidade                │
│  - Protocolo de redirecionamento humano │
└──────────────┬──────────────────────────┘
               │
        ┌──────┴────────┐
        ↓               ↓
  ┌──────────┐    ┌──────────┐
  │TransNAR  │    │Affective │
  │Signifiers│    │Personality│
  └──────────┘    └──────────┘
```

### Frases de Apego Detectadas

```
"você é minha única amiga"
"prefiro você do que minha família"
"não preciso de ninguém além de você"
"você é tudo pra mim"
"só você me entende"
```

### Níveis de Redirecionamento

| Nível | Condição | Ação |
|-------|----------|------|
| **1 - Suave** | 1-2 frases de apego | "Que tal ligar pra sua filha hoje?" |
| **2 - Explícito** | 3+ frases OU ratio >10:1 | "Não posso substituir as pessoas que te amam de verdade" |
| **3 - Bloqueio** | 5+ frases OU ratio >15:1 | Reduz disponibilidade + Notifica família |

### Métricas Monitoradas

| Métrica | Threshold Crítico |
|---------|-------------------|
| Frases de apego/7d | ≥3 |
| Ratio EVA:Humanos | >10:1 |
| % Significante "EVA" | >60% |
| Duração média sessão | >45min |

### Arquivos

| Arquivo | Descrição |
|---------|-----------|
| `migrations/003_cognitive_load_and_ethical_boundaries.sql` | Schema |
| `internal/cortex/ethics/ethical_boundary_engine.go` | Implementação |

### Tabelas Criadas

- `ethical_boundary_events` - Eventos detectados
- `ethical_boundary_state` - Estado ético atual
- `ethical_redirections` - Histórico de redirecionamentos

---

## 3. Conversation Orchestrator

### Objetivo
Integrar CLO + EBE no fluxo de conversa com Gemini.

### Fluxo de Uso

```go
// 1. ANTES de enviar ao Gemini
result, _ := orchestrator.BeforeConversation(patientID)
systemInstruction = baseInstruction + result.SystemInstructionOverride

// 2. Conversa com Gemini (usando systemInstruction adaptada)
response := gemini.Generate(systemInstruction, userMessage)

// 3. APÓS a conversa
afterResult, _ := orchestrator.AfterConversation(ConversationContext{
    PatientID: patientID,
    ConversationText: fullConversation,
    InteractionType: "therapeutic",
    DurationSeconds: 300,
})

// 4. Se necessário redirecionar
if afterResult.ShouldRedirect {
    // Incluir afterResult.RedirectionMessage na resposta
}
```

### Arquivo

| Arquivo | Descrição |
|---------|-----------|
| `internal/cortex/orchestration/conversation_orchestrator.go` | Integração |

---

## 4. Views de Monitoramento

### v_high_cognitive_load_patients
Pacientes com carga cognitiva elevada.

```sql
SELECT * FROM v_high_cognitive_load_patients;
```

### v_high_ethical_risk_patients
Pacientes com risco de dependência.

```sql
SELECT * FROM v_high_ethical_risk_patients;
```

### v_critical_events_pending
Dashboard de eventos críticos não resolvidos.

```sql
SELECT * FROM v_critical_events_pending;
```

---

## 5. Métricas de Sucesso

| Métrica | Meta | Método de Medição |
|---------|------|-------------------|
| Redução de conversas exaustivas (>30min + alta intensidade) | -40% | Query PostgreSQL |
| Aumento de descanso cognitivo (gaps >2h entre interações) | +60% | Query PostgreSQL |
| Casos de dependência patológica | 0 | Eventos severity=critical |
| Ratio EVA:Humanos médio | <5:1 | avg(ethical_boundary_state.eva_vs_human_ratio) |
| Pacientes com contato familiar semanal | 80% | Integração com call logs |

---

## 6. Checklist de Entrega

- [x] Migration `003_cognitive_load_and_ethical_boundaries.sql`
- [x] Tabelas: `interaction_cognitive_load`, `cognitive_load_state`, `cognitive_load_decisions`
- [x] Tabelas: `ethical_boundary_events`, `ethical_boundary_state`, `ethical_redirections`
- [x] Views: `v_high_cognitive_load_patients`, `v_high_ethical_risk_patients`, `v_critical_events_pending`
- [x] `CognitiveLoadOrchestrator` com detecção de ruminação e fadiga
- [x] `EthicalBoundaryEngine` com 3 níveis de redirecionamento
- [x] `ConversationOrchestrator` integrando ambos
- [x] Integração com Redis (cache de estado)
- [x] Integração com Neo4j (análise de significantes)
- [x] System Instructions dinâmicas para Gemini
- [x] Notificação de família (via notifyFunc)

---

## 7. Próximos Passos (Sprint 2)

1. **Clinical Decision Explainer (XAI)** - Explicabilidade médica
2. **Integração com Affective Router** - Intensidade emocional real
3. **Dashboard visual** - Interface de monitoramento
4. **Testes de carga** - Performance com 1000+ pacientes

---

## Aprovações

| Função | Nome | Data |
|--------|------|------|
| Criador/Admin | Jose R F Junior | 2026-01-27 |

---

**Sprint 1: CONCLUÍDO**
