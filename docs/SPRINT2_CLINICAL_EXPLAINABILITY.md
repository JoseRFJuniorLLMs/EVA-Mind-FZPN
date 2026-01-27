# Sprint 2: Explicabilidade Clínica (XAI) - EVA-Mind-FZPN

**Documento:** SPRINT2-XAI-001
**Versão:** 1.0
**Data:** 2026-01-27
**Status:** CONCLUÍDO

---

## Resumo Executivo

O Sprint 2 implementou a **Camada de Explicabilidade Clínica** do EVA-Mind-FZPN, permitindo que decisões da IA sejam explicadas para profissionais de saúde de forma transparente e auditável.

---

## 1. Clinical Decision Explainer

### Objetivo

Gerar explicações detalhadas e compreensíveis para:
- Predições de risco (crise, depressão, suicídio)
- Alertas clínicos
- Recomendações de intervenção

### Arquitetura

```
┌─────────────────────────────────────────────────────────┐
│              CLINICAL DECISION EXPLAINER                │
│                                                         │
│  ┌───────────────┐  ┌───────────────┐  ┌─────────────┐ │
│  │ Feature       │  │ SHAP-like     │  │ Natural     │ │
│  │ Analysis      │→ │ Contributions │→ │ Language    │ │
│  └───────────────┘  └───────────────┘  └─────────────┘ │
│           │                                     │       │
│           ▼                                     ▼       │
│  ┌───────────────┐                     ┌─────────────┐ │
│  │ Recommendations│                    │ PDF Report  │ │
│  └───────────────┘                     └─────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### Tipos de Decisão Explicados

| Tipo | Descrição |
|------|-----------|
| `crisis_prediction` | Risco de crise mental nas próximas 24-48h |
| `depression_alert` | Deterioração de sintomas depressivos |
| `anxiety_alert` | Aumento de sintomas ansiosos |
| `medication_alert` | Problema de adesão medicamentosa |
| `suicide_risk` | Ideação ou risco suicida |
| `hospitalization_risk` | Risco de internação |
| `fall_risk` | Risco de queda |

### Features Analisadas

| Feature | Peso | Descrição |
|---------|------|-----------|
| `medication_adherence` | 1.5x | % de doses tomadas |
| `voice_biomarkers` | 1.3x | Pitch, energia, variabilidade |
| `phq9_score` | 1.2x | Escala de depressão |
| `gad7_score` | 1.2x | Escala de ansiedade |
| `sleep_quality` | 1.1x | Horas e qualidade de sono |
| `activity_level` | 1.0x | Nível de atividade física |

---

## 2. Geração de Explicações

### Fluxo de Geração

```go
// 1. Criar predição com features
prediction := ClinicalPrediction{
    PatientID:           123,
    DecisionType:        "depression_alert",
    PredictionScore:     0.75,
    PredictionTimeframe: "7-14 dias",
    Severity:            "high",
    Features: map[string]Feature{
        "medication_adherence": {CurrentValue: 0.42, BaselineValue: 0.85, Status: "critical"},
        "phq9_score":           {CurrentValue: 18, BaselineValue: 8, Status: "concerning"},
        "sleep_hours":          {CurrentValue: 4.2, BaselineValue: 7, Status: "warning"},
    },
}

// 2. Gerar explicação
explainer := NewClinicalDecisionExplainer(db)
explanation, err := explainer.ExplainDecision(prediction)

// 3. Gerar PDF (opcional)
pdfGen := NewPDFGenerator(db, explainer)
report, err := pdfGen.GenerateExplanationPDF(explanation)
```

### Saída da Explicação

```
🚨 ALERTA: Alerta de Depressão

Probabilidade: 75% (alto)
Janela temporal: 7-14 dias

📊 FATORES PRINCIPAIS (por ordem de importância):

1. Medication Adherence (contribuição: 42%)
   Status: 🔴 Crítico
   Adesão medicamentosa crítica: apenas 42% das doses tomadas
   Comparação: ↓ 50.6% abaixo da baseline

2. Phq9 Score (contribuição: 35%)
   Status: ⚠️ Preocupante
   Depressão moderadamente severa (PHQ-9: 18)
   Comparação: ↑ 125.0% acima da baseline

3. Sleep Hours (contribuição: 23%)
   Status: ⚠️ Atenção
   Qualidade de sono ruim: 4.2 horas/noite
   Comparação: ↓ 40.0% abaixo da baseline
```

---

## 3. Relatórios PDF

### Tipos de Relatório

| Tipo | Descrição | Frequência |
|------|-----------|------------|
| Explicação Clínica | Detalhes de um alerta específico | Por evento |
| Resumo Semanal | Métricas e alertas da semana | Semanal |
| Relatório de Crise | Urgente para risco alto/crítico | Por evento |

### Estrutura do PDF

1. **Header**: Logo EVA, data, ID do relatório
2. **Info Paciente**: Nome, idade, médico responsável
3. **Alerta Principal**: Tipo, probabilidade, severidade
4. **Fatores Principais**: Top 3 contribuintes
5. **Fatores Secundários**: Demais fatores
6. **Recomendações**: Ações com urgência e prazo
7. **Footer**: Disclaimer LGPD, versão do modelo

---

## 4. Banco de Dados

### Tabelas Criadas

| Tabela | Descrição |
|--------|-----------|
| `clinical_decision_explanations` | Explicações completas |
| `decision_factors` | Fatores individuais (SHAP values) |
| `prediction_accuracy_log` | Histórico de acurácia |

### Views de Monitoramento

```sql
-- Predições de alto risco não revisadas
SELECT * FROM v_high_risk_predictions;

-- Acurácia do modelo por tipo
SELECT * FROM v_model_accuracy_by_type;

-- Alertas pendentes de revisão médica
SELECT * FROM v_pending_doctor_review;
```

---

## 5. Arquivos Implementados

| Arquivo | Descrição |
|---------|-----------|
| `migrations/004_clinical_decision_explainer.sql` | Schema das tabelas |
| `internal/cortex/explainability/clinical_decision_explainer.go` | Engine de explicação |
| `internal/cortex/explainability/pdf_generator.go` | Gerador de PDFs |

---

## 6. Integração com Sistema

### Uso no Fluxo Principal

```go
// Quando uma predição de risco é feita
if riskPrediction.Score > 0.6 {
    // 1. Gerar explicação
    explanation, _ := explainer.ExplainDecision(riskPrediction)

    // 2. Se severidade alta, gerar PDF
    if explanation.Severity == "high" || explanation.Severity == "critical" {
        report, _ := pdfGen.GenerateExplanationPDF(explanation)

        // 3. Notificar médico
        notifyDoctor(patientID, explanation, report.S3URL)
    }
}
```

---

## 7. Checklist de Entrega

- [x] Migration `004_clinical_decision_explainer.sql`
- [x] Tabelas: `clinical_decision_explanations`, `decision_factors`, `prediction_accuracy_log`
- [x] Views: `v_high_risk_predictions`, `v_model_accuracy_by_type`, `v_pending_doctor_review`
- [x] `ClinicalDecisionExplainer` com cálculo SHAP-like
- [x] Classificação de fatores (primary/secondary)
- [x] Geração de recomendações clínicas
- [x] Explicações em linguagem natural (português)
- [x] `PDFGenerator` para relatórios
- [x] Templates HTML para conversão PDF
- [x] Documentação completa

---

## 8. Próximos Passos (Sprint 3)

1. **Predictive Life Trajectory** - Trajetória de vida preditiva
2. **Integração com biblioteca PDF real** (wkhtmltopdf ou chromedp)
3. **API REST para consulta de explicações**
4. **Dashboard visual de monitoramento**

---

## Aprovações

| Função | Nome | Data |
|--------|------|------|
| Criador/Admin | Jose R F Junior | 2026-01-27 |

---

**Sprint 2: CONCLUÍDO**
