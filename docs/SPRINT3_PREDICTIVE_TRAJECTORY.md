# Sprint 3: Predictive Life Trajectory Engine - EVA-Mind-FZPN

**Documento:** SPRINT3-TRAJECTORY-001
**Versão:** 1.0
**Data:** 2026-01-27
**Status:** CONCLUÍDO

---

## Resumo Executivo

O Sprint 3 implementou o **Motor de Trajetória Preditiva**, que utiliza simulações Monte Carlo para prever a evolução do estado de saúde mental dos pacientes e recomendar intervenções preventivas.

---

## 1. Predictive Life Trajectory Engine

### Objetivo

Simular o futuro do paciente para:
- Prever probabilidade de crise em 7 e 30 dias
- Identificar fatores críticos de risco
- Recomendar intervenções preventivas
- Comparar cenários "what-if"

### Arquitetura

```
┌────────────────────────────────────────────────────────────────┐
│              PREDICTIVE LIFE TRAJECTORY ENGINE                 │
│                                                                │
│  ┌─────────────┐     ┌──────────────────┐     ┌─────────────┐ │
│  │ Patient     │────▶│ Monte Carlo      │────▶│ Aggregated  │ │
│  │ State       │     │ Simulations      │     │ Results     │ │
│  └─────────────┘     │ (N=1000)         │     └─────────────┘ │
│                      └──────────────────┘            │        │
│                             │                        ▼        │
│                             ▼                 ┌─────────────┐ │
│                      ┌──────────────────┐     │ Risk        │ │
│                      │ Bayesian Network │     │ Assessment  │ │
│                      │ Parameters       │     └─────────────┘ │
│                      └──────────────────┘            │        │
│                                                      ▼        │
│  ┌─────────────┐     ┌──────────────────┐     ┌─────────────┐ │
│  │ What-If     │◀────│ Intervention     │◀────│ Recommen-   │ │
│  │ Scenarios   │     │ Simulation       │     │ dations     │ │
│  └─────────────┘     └──────────────────┘     └─────────────┘ │
└────────────────────────────────────────────────────────────────┘
```

### Funcionalidades

| Funcionalidade | Descrição |
|----------------|-----------|
| **Simulação Monte Carlo** | 1000 trajetórias possíveis por paciente |
| **Predição de Crise** | Probabilidade em 7 e 30 dias |
| **Fatores Críticos** | Identifica principais drivers de risco |
| **Cenários What-If** | Simula impacto de intervenções |
| **Recomendações** | Gera intervenções priorizadas |

---

## 2. Modelo de Predição

### Estado do Paciente

| Variável | Descrição | Range |
|----------|-----------|-------|
| `PHQ9Score` | Score de depressão | 0-27 |
| `GAD7Score` | Score de ansiedade | 0-21 |
| `MedicationAdherence` | Adesão medicamentosa | 0-1 |
| `SleepHours` | Horas de sono/noite | 2-10 |
| `SocialIsolationDays` | Dias sem contato humano | 0+ |
| `VoiceEnergyScore` | Energia vocal | 0-1 |

### Regras de Transição

```go
// PHQ-9 aumenta se:
// - Adesão medicamentosa < 50%
// - Sono < 5 horas/noite
if currentAdherence < 0.5 {
    phq9Delta += 0.2
}
if currentSleep < 5 {
    phq9Delta += 0.15
}

// Probabilidade de crise (modelo logístico):
logit = -5.0
logit += (phq9 - 10) * 0.15      // PHQ-9 alto = mais risco
logit += (0.7 - adherence) * 2.0 // Baixa adesão = mais risco
logit += (6 - sleep) * 0.3       // Pouco sono = mais risco
prob = 1 / (1 + exp(-logit))
```

---

## 3. Intervenções Recomendadas

| Tipo | Gatilho | Impacto Esperado |
|------|---------|------------------|
| `psychiatric_consultation` | Risco >60% | -25% risco |
| `medication_adherence_boost` | Adesão <60% | -15% risco, -3 PHQ-9 |
| `sleep_hygiene_protocol` | Sono <5h | -10% risco |
| `family_engagement` | Isolamento >5 dias | -12% risco |
| `therapy_intensification` | PHQ-9 >15 | -18% risco, -4.5 PHQ-9 |

### Priorização

| Prioridade | Critério | Prazo |
|------------|----------|-------|
| **Crítica** | Risco >60% ou PHQ-9 ≥20 | 24-48h |
| **Alta** | Risco 40-60% ou PHQ-9 15-19 | 3-5 dias |
| **Média** | Risco 20-40% | 1 semana |
| **Baixa** | Risco <20% | 2 semanas |

---

## 4. Uso do Sistema

### Executar Simulação

```go
engine := NewTrajectoryEngine(db)

// Simular 30 dias com 1000 trajetórias
simulation, err := engine.SimulateTrajectory(patientID, 30, 1000)

fmt.Printf("Risco 7d: %.1f%%\n", simulation.CrisisProbability7d*100)
fmt.Printf("Risco 30d: %.1f%%\n", simulation.CrisisProbability30d*100)
fmt.Printf("Fatores críticos: %v\n", simulation.CriticalFactors)
```

### Simular Intervenção

```go
// O que acontece se aumentarmos a adesão em 15%?
interventions := []Intervention{
    {Type: "medication_adherence_boost", ImpactAdherence: 0.15},
}

scenario, err := engine.SimulateWithIntervention(simulation.ID, interventions)

fmt.Printf("Redução de risco: %.1f%%\n", scenario.RiskReduction30d*100)
```

### Gerar Recomendações

```go
recommendations, err := engine.GenerateRecommendations(simulation.ID)

for _, rec := range recommendations {
    fmt.Printf("[%s] %s\n", rec.Priority, rec.Title)
    fmt.Printf("  Impacto: -%.0f%% risco\n", rec.ExpectedRiskReduction*100)
    fmt.Printf("  Ações: %v\n", rec.ActionSteps)
}
```

---

## 5. Banco de Dados

### Tabelas

| Tabela | Descrição |
|--------|-----------|
| `trajectory_simulations` | Resultados de simulações |
| `intervention_scenarios` | Cenários what-if |
| `recommended_interventions` | Intervenções recomendadas |
| `trajectory_prediction_accuracy` | Acurácia do modelo |
| `bayesian_network_parameters` | Parâmetros da rede Bayesiana |

### Views

```sql
-- Última simulação por paciente
SELECT * FROM v_latest_trajectory_simulations;

-- Pacientes de alto risco com intervenções pendentes
SELECT * FROM v_high_risk_patients_pending_interventions;

-- Acurácia do modelo
SELECT * FROM v_model_accuracy_by_version;

-- Relatório completo de um paciente
SELECT get_patient_trajectory_report(123);
```

---

## 6. Arquivos Implementados

| Arquivo | Descrição |
|---------|-----------|
| `migrations/005_predictive_trajectory.sql` | Schema completo |
| `internal/cortex/predictive/trajectory_engine.go` | Engine principal |

---

## 7. Checklist de Entrega

- [x] Migration `005_predictive_trajectory.sql`
- [x] Tabelas: `trajectory_simulations`, `intervention_scenarios`, `recommended_interventions`
- [x] Views: `v_latest_trajectory_simulations`, `v_high_risk_patients_pending_interventions`
- [x] `TrajectoryEngine` com Monte Carlo
- [x] Simulação de intervenções (what-if)
- [x] Geração de recomendações priorizadas
- [x] Função SQL `get_patient_trajectory_report()`
- [x] Documentação completa

---

## 8. Próximos Passos

1. **Treinar Rede Bayesiana** com dados reais
2. **Validar modelo** com outcomes históricos
3. **Integrar com dashboard** para visualização
4. **Adicionar mais variáveis** (voz, atividade física, clima)

---

## Aprovações

| Função | Nome | Data |
|--------|------|------|
| Criador/Admin | Jose R F Junior | 2026-01-27 |

---

**Sprint 3: CONCLUÍDO**
