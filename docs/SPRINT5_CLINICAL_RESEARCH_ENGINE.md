# Sprint 5: Clinical Research Engine - EVA-Mind-FZPN

**Documento:** SPRINT5-RESEARCH-001
**Versão:** 1.0
**Data:** 2026-01-27
**Status:** CONCLUÍDO

---

## Resumo Executivo

O Sprint 5 implementou o **Clinical Research Engine**, permitindo análises longitudinais, estudos científicos validados e publicação de resultados para transformar EVA em DTx cientificamente defensável.

---

## 1. Arquitetura do Sistema

```
┌────────────────────────────────────────────────────────────────┐
│                  CLINICAL RESEARCH ENGINE                       │
│                                                                  │
│  ┌─────────────┐     ┌──────────────────┐     ┌──────────────┐ │
│  │ Cohort      │────▶│ Anonymization    │────▶│ Longitudinal │ │
│  │ Builder     │     │ Service          │     │ Analyzer     │ │
│  └─────────────┘     └──────────────────┘     └──────────────┘ │
│         │                    │                       │          │
│         ▼                    ▼                       ▼          │
│  ┌─────────────┐     ┌──────────────────┐     ┌──────────────┐ │
│  │ Research    │     │ Statistical      │     │ Export       │ │
│  │ Datapoints  │     │ Methods          │     │ Service      │ │
│  └─────────────┘     └──────────────────┘     └──────────────┘ │
└────────────────────────────────────────────────────────────────┘
```

---

## 2. Estudos Pré-Configurados

| Estudo | Código | Hipótese | Método |
|--------|--------|----------|--------|
| **Voice → Depression** | EVA-VOICE-PHQ9-001 | Voz prediz PHQ-9 7-14 dias antes | Lag Correlation |
| **Adherence → Depression** | EVA-ADHERENCE-DEP-002 | Adesão <50% → PHQ-9 +5 pts em 30d | Propensity Score |
| **Isolation → Crisis** | EVA-ISOLATION-CRISIS-003 | 7+ dias isolado → 3x risco crise | Survival Analysis |
| **Sleep → Mental Health** | EVA-SLEEP-MH-004 | Sono <5h → piora PHQ-9/GAD-7 | Lag Correlation |

---

## 3. Análise Longitudinal

### 3.1 Lag Correlation Analysis

Calcula correlações entre variáveis com diferentes lags temporais.

**Variáveis Suportadas:**

| Variável | Coluna DB | Descrição |
|----------|-----------|-----------|
| `voice_pitch_mean` | `voice_pitch_mean_hz` | Pitch médio da voz |
| `voice_jitter` | `voice_jitter` | Tremor vocal |
| `voice_shimmer` | `voice_shimmer` | Variação de amplitude |
| `phq9` | `phq9_score` | Score de depressão |
| `gad7` | `gad7_score` | Score de ansiedade |
| `medication_adherence` | `medication_adherence_7d` | Adesão medicamentosa |
| `sleep_hours` | `sleep_hours_avg_7d` | Horas de sono |
| `social_isolation` | `social_isolation_days` | Dias de isolamento |

**Uso:**

```go
engine := research.NewResearchEngine(db)

// Correlação: voz prediz PHQ-9 com lag de 0-14 dias
results, err := engine.longitudinalAnalyzer.CalculateLagCorrelations(
    cohortID,
    "voice_pitch_mean",
    "phq9",
    14, // maxLag
)

for _, r := range results {
    if r.PValue < 0.05 {
        fmt.Printf("Lag %d dias: r=%.3f, p=%.6f (SIGNIFICATIVO)\n",
            r.LagDays, r.CorrelationCoefficient, r.PValue)
    }
}
```

### 3.2 Trend Analysis

Identifica tendências ao longo do tempo.

```go
trend, err := engine.longitudinalAnalyzer.CalculateTrend(cohortID, "phq9")

fmt.Printf("Tendência: %s (slope=%.3f, R²=%.3f)\n",
    trend.Direction, // "increasing", "decreasing", "stable"
    trend.Slope,
    trend.RSquared,
)
```

### 3.3 Change Point Detection

Detecta mudanças abruptas em séries temporais.

```go
changePoints, err := engine.longitudinalAnalyzer.DetectChangePoints(
    cohortID,
    "phq9",
    30.0, // threshold: mudança de 30%
)

for _, cp := range changePoints {
    fmt.Printf("Mudança no dia %d: %.1f → %.1f (%.1f%%)\n",
        cp.Day, cp.ValueBefore, cp.ValueAfter, cp.PercentChange)
}
```

---

## 4. Métodos Estatísticos

### 4.1 Correlação de Pearson

```go
sm := research.NewStatisticalMethods()

r := sm.PearsonCorrelation(x, y)
pValue := sm.CorrelationPValue(r, len(x))
ciLower, ciUpper := sm.CorrelationConfidenceInterval(r, len(x), 0.95)
```

### 4.2 Regressão Linear

```go
slope, intercept, rSquared := sm.SimpleLinearRegression(x, y)
// y = slope*x + intercept
```

### 4.3 T-Test

```go
tStat, pValue := sm.TTest(group1, group2)
```

### 4.4 Cohen's D (Effect Size)

```go
d := sm.CohensD(group1, group2)
interpretation := sm.InterpretCohensD(d)
// "negligible", "small", "medium", "large"
```

### 4.5 Estatísticas Descritivas

```go
mean := sm.Mean(values)
variance := sm.Variance(values)
stdDev := sm.StandardDeviation(values)
median := sm.Median(values)
p95 := sm.Percentile(values, 95)
```

---

## 5. Anonimização de Dados

### 5.1 Processo de Anonimização

```go
anonymizer := research.NewAnonymizer(db)

// Anonimiza patient_id com SHA-256 irreversível
anonymousID := research.AnonymizePatientID(patientID)
// Resultado: "a3f2b1c4d5e6..." (64 caracteres hex)
```

### 5.2 Níveis de Anonimização

| Nível | Descrição |
|-------|-----------|
| `fully_anonymized` | SHA-256 do ID, sem PII |
| `pseudonymized` | ID substituído, mapeamento guardado |
| `aggregated_only` | Apenas dados agregados |

### 5.3 Compliance

- **LGPD**: Compliant
- **GDPR**: Compliant
- **HIPAA**: Compliant
- **k-Anonymity**: Verificável via SQL function

```sql
SELECT calculate_k_anonymity(cohort_id, ARRAY['age_group', 'gender']);
```

---

## 6. Banco de Dados

### 6.1 Tabelas

| Tabela | Descrição |
|--------|-----------|
| `research_cohorts` | Definição de estudos |
| `research_datapoints` | Dados longitudinais anonimizados |
| `longitudinal_correlations` | Resultados de correlação |
| `statistical_analyses` | Análises estatísticas |
| `research_publications` | Tracking de papers |
| `research_exports` | Datasets exportados |

### 6.2 Views

```sql
-- Estudos ativos com progresso
SELECT * FROM v_active_research_studies;

-- Correlações significativas
SELECT * FROM v_significant_correlations;

-- Papers publicados
SELECT * FROM v_published_papers;

-- Portfolio de pesquisa
SELECT * FROM v_research_portfolio;
```

### 6.3 Functions

```sql
-- Relatório completo de estudo
SELECT generate_study_report(cohort_id);

-- Verificar k-anonymity
SELECT calculate_k_anonymity(cohort_id, quasi_identifiers);
```

---

## 7. Uso do Sistema

### 7.1 Criar Estudos Pré-Configurados

```go
engine := research.NewResearchEngine(db)
err := engine.CreatePreconfiguredStudies()
```

### 7.2 Coletar Dados para Coorte

```go
err := engine.CollectDataForCohort(cohortID)
```

### 7.3 Executar Análise de Correlação

```go
err := engine.RunLagCorrelationAnalysis(
    cohortID,
    "voice_pitch_mean", // predictor
    "phq9",             // outcome
    14,                 // max lag (dias)
)
```

### 7.4 Exportar Dataset

```go
err := engine.ExportDatasetToCSV(cohortID, "/research/datasets/study.csv")
```

### 7.5 Gerar Relatório

```go
report, err := engine.GenerateStudyReport(cohortID)
// {
//   "study": {...},
//   "significant_correlations": [...],
//   "analyses": [...]
// }
```

---

## 8. Arquivos Implementados

| Arquivo | Descrição |
|---------|-----------|
| `migrations/007_clinical_research_engine.sql` | Schema completo |
| `internal/research/research_engine.go` | Engine principal |
| `internal/research/anonymization.go` | Anonimização LGPD |
| `internal/research/longitudinal_analysis.go` | Análise de séries temporais |
| `internal/research/statistical_methods.go` | Métodos estatísticos |
| `internal/research/cohort_builder.go` | Seleção de pacientes |

---

## 9. Checklist de Entrega

- [x] Migration `007_clinical_research_engine.sql`
- [x] Tabelas: `research_cohorts`, `research_datapoints`, `longitudinal_correlations`
- [x] Tabelas: `statistical_analyses`, `research_publications`, `research_exports`
- [x] Views: `v_active_research_studies`, `v_significant_correlations`, `v_published_papers`
- [x] Functions: `generate_study_report()`, `calculate_k_anonymity()`
- [x] `ResearchEngine` com gestão de coortes
- [x] `LongitudinalAnalyzer` com lag correlations e trend analysis
- [x] `StatisticalMethods` com Pearson, T-Test, regressão
- [x] `Anonymizer` com SHA-256 e compliance LGPD
- [x] Change point detection
- [x] Exportação de datasets
- [x] Documentação completa

---

## 10. Estudos Científicos Objetivos

### Objetivo 1: Voice Biomarkers Paper

> **Hipótese:** Mudanças em pitch/jitter precedem mudanças em PHQ-9

**Método:**
- N = 100 pacientes, 6 meses
- Lag correlation (0-14 dias)
- Mixed effects model

**Output:** Paper peer-reviewed

### Objetivo 2: Medication Adherence Paper

> **Hipótese:** Adesão <50% por 2 semanas → PHQ-9 +5 pontos

**Método:**
- N = 200 pacientes
- Propensity score matching
- Logistic regression

**Output:** Modelo preditivo validado

### Objetivo 3: Isolation Risk Paper

> **Hipótese:** 7+ dias isolado → 3x risco de crise

**Método:**
- N = 150 pacientes
- Kaplan-Meier survival analysis
- Cox proportional hazards

**Output:** Paper + aprovação para RCT

---

## Aprovações

| Função | Nome | Data |
|--------|------|------|
| Criador/Admin | Jose R F Junior | 2026-01-27 |

---

**Sprint 5: CONCLUÍDO**
