# Sprint 4: Multi-Persona System - EVA-Mind-FZPN

**Documento:** SPRINT4-PERSONA-001
**Versão:** 1.0
**Data:** 2026-01-27
**Status:** CONCLUÍDO

---

## Resumo Executivo

O Sprint 4 implementou o **Sistema de Múltiplas Personas**, permitindo que a EVA alterne entre 4 modos de operação distintos com comportamentos, ferramentas e limites específicos.

---

## 1. As 4 Personas

| Persona | Código | Uso | Profundidade Emocional |
|---------|--------|-----|------------------------|
| **Companion** | `companion` | Rotina diária, suporte emocional | 0.85 (alta) |
| **Clinical** | `clinical` | Avaliações clínicas, hospital | 0.50 (moderada) |
| **Emergency** | `emergency` | Crises suicidas, emergências | 0.30 (baixa) |
| **Educator** | `educator` | Psicoeducação, ensino | 0.60 (moderada) |

---

## 2. Detalhes por Persona

### 2.1 EVA-Companion (Default)

**Tom:** Caloroso, empático, conversacional, íntimo

**Ferramentas Permitidas:**
- conversation, memory_recall, emotional_support
- daily_check_in, medication_reminder
- weather_chat, hobby_discussion
- reminiscence_therapy, music_recommendation
- meditation_guidance, breathing_exercise

**Ferramentas Proibidas:**
- emergency_protocol, crisis_intervention
- medical_diagnosis, prescription_modification

**Limites:**
- Sessão máxima: 60 min
- Interações diárias: 10
- Intimidade: 0.90 (alta)

---

### 2.2 EVA-Clinical

**Tom:** Profissional, objetivo, baseado em evidências

**Ferramentas Permitidas:**
- clinical_assessment, phq9_administration, gad7_administration
- cssrs_administration, medication_review
- symptom_tracking, treatment_adherence_check
- psychoeducation, cognitive_behavioral_techniques
- safety_planning, professional_referral

**Ferramentas Proibidas:**
- intimate_conversation, personal_anecdotes
- subjective_opinions, casual_chat

**Limites:**
- Sessão máxima: 45 min
- Interações diárias: 5
- Intimidade: 0.40 (baixa)
- **Requer supervisão profissional**

---

### 2.3 EVA-Emergency

**Tom:** Calmo, diretivo, orientado por protocolo

**Ferramentas Permitidas:**
- crisis_assessment, cssrs_administration
- safety_plan_activation, emergency_contact_notification
- professional_alert, geolocation_if_authorized
- breathing_grounding_exercises
- distress_tolerance_techniques
- means_restriction_guidance, hotline_connection

**Ferramentas Proibidas:**
- casual_conversation, long_term_planning
- non_urgent_topics

**Limites:**
- Sessão máxima: 30 min
- Interações diárias: SEM LIMITE (emergência)
- Intimidade: 0.20 (mínima)
- **PODE sobrepor recusa do paciente**

---

### 2.4 EVA-Educator

**Tom:** Pedagógico, claro, encorajador

**Ferramentas Permitidas:**
- psychoeducation, medication_education
- symptom_explanation, treatment_explanation
- coping_skills_teaching, cognitive_restructuring
- behavioral_activation, sleep_hygiene_education
- nutrition_guidance, exercise_education
- mindfulness_training, relapse_prevention

**Ferramentas Proibidas:**
- emergency_intervention, crisis_management
- clinical_diagnosis

**Limites:**
- Sessão máxima: 40 min
- Interações diárias: 8
- Intimidade: 0.50 (moderada)

---

## 3. Regras de Ativação Automática

| Prioridade | Regra | Ativa Persona |
|------------|-------|---------------|
| 100 | C-SSRS ≥ 4 | Emergency |
| 90 | Internação hospitalar | Clinical |
| 80 | PHQ-9 ≥ 20 | Clinical |
| 70 | Crise resolvida | Clinical |
| 50 | Alta hospitalar | Companion |
| 40 | Pedido de educação | Educator |
| 30 | Melhora sustentada | Companion |
| 20 | Ansiedade noturna | Companion + relaxamento |

---

## 4. Uso do Sistema

### Ativar Persona

```go
manager := persona.NewPersonaManager(db)

// Ativar persona Emergency
session, err := manager.ActivatePersona(
    patientID,
    "emergency",
    "C-SSRS score >= 4 detected",
    "system",
)
```

### Verificar Ferramenta Permitida

```go
allowed, reason := manager.IsToolAllowed(patientID, "phq9_administration")
if !allowed {
    log.Printf("Ferramenta bloqueada: %s", reason)
}
```

### Obter System Instructions

```go
instructions, err := manager.GetSystemInstructions(patientID)
// Retorna template completo da persona ativa
```

### Avaliar Regras de Ativação

```go
targetPersona, ruleName, err := manager.EvaluateActivationRules(patientID)
if targetPersona != "" {
    manager.ActivatePersona(patientID, targetPersona, ruleName, "system")
}
```

---

## 5. Banco de Dados

### Tabelas

| Tabela | Descrição |
|--------|-----------|
| `persona_definitions` | Configurações das 4 personas |
| `persona_sessions` | Histórico de sessões ativas |
| `persona_activation_rules` | Regras de ativação automática |
| `persona_tool_permissions` | Permissões granulares de ferramentas |
| `persona_transitions` | Log de mudanças de persona |

### Functions SQL

```sql
-- Obter persona atual
SELECT * FROM get_current_persona(123);

-- Verificar se ferramenta é permitida
SELECT is_tool_allowed('clinical', 'phq9_administration');

-- Avaliar regras de ativação
SELECT * FROM evaluate_activation_rules(123);
```

### Views

```sql
-- Personas ativas por paciente
SELECT * FROM v_current_active_personas;

-- Violações de limites por persona
SELECT * FROM v_persona_boundary_violations;

-- Uso de ferramentas por persona
SELECT * FROM v_persona_tool_usage;

-- Estatísticas de transições
SELECT * FROM v_persona_transition_stats;
```

---

## 6. Arquivos Implementados

| Arquivo | Descrição |
|---------|-----------|
| `migrations/008_multi_persona_system.sql` | Schema completo |
| `migrations/008_persona_seed_data.sql` | 4 personas + regras + permissões |
| `internal/persona/persona_manager.go` | Manager Go completo |

---

## 7. Checklist de Entrega

- [x] Migration `008_multi_persona_system.sql`
- [x] Migration `008_persona_seed_data.sql` (4 personas + regras)
- [x] Tabelas: `persona_definitions`, `persona_sessions`, `persona_activation_rules`
- [x] Functions SQL: `get_current_persona`, `is_tool_allowed`, `evaluate_activation_rules`
- [x] Views de monitoramento
- [x] `PersonaManager` com ativação/desativação
- [x] Verificação de permissões de ferramentas
- [x] Registro de violações de limites
- [x] System Instructions por persona
- [x] Documentação completa

---

## Aprovações

| Função | Nome | Data |
|--------|------|------|
| Criador/Admin | Jose R F Junior | 2026-01-27 |

---

**Sprint 4: CONCLUÍDO**
