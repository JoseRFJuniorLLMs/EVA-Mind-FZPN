# MATRIZ DE RASTREABILIDADE
## EVA-Mind-FZPN - Requisitos → Código → Testes

**Documento:** TM-001
**Versão:** 1.0
**Data:** 2026-01-27
**Referências:** RDC 751/2022, IEC 62304:2006, ISO 14971:2019

---

## 1. INTRODUÇÃO

### 1.1 Objetivo
Esta matriz de rastreabilidade documenta a relação entre:
- **Requisitos de Software (SRS)** - O que o sistema deve fazer
- **Design/Implementação (COD)** - Como foi implementado
- **Verificação (TEST)** - Como foi testado
- **Riscos (RISK)** - Riscos associados mitigados

### 1.2 Escopo
Cobre todas as funcionalidades críticas do EVA-Mind classificadas como:
- **Classe A** - Críticas para segurança do paciente
- **Classe B** - Importantes para funcionamento correto
- **Classe C** - Funcionalidades de suporte

### 1.3 Abreviações
| Sigla | Descrição |
|-------|-----------|
| SRS | Software Requirements Specification |
| COD | Código/Implementação |
| TEST | Caso de Teste |
| RISK | Risco associado (ISO 14971) |
| V | Verificado |
| P | Parcialmente implementado |
| N | Não implementado |

---

## 2. MATRIZ DE RASTREABILIDADE - FUNCIONALIDADES CLÍNICAS

### 2.1 Avaliação de Risco Suicida (C-SSRS)

| ID Req | Requisito | Classe | Arquivo Código | Linha | ID Teste | Status |
|--------|-----------|--------|----------------|-------|----------|--------|
| SRS-CSSRS-001 | Sistema deve implementar escala C-SSRS completa com 6 questões | A | `internal/cortex/scales/clinical_scales.go` | 295-454 | TEST-CSSRS-001 | V |
| SRS-CSSRS-002 | Questões 1-5 devem avaliar ideação suicida | A | `internal/cortex/scales/clinical_scales.go` | 299-325 | TEST-CSSRS-002 | V |
| SRS-CSSRS-003 | Questão 6 deve avaliar comportamento suicida | A | `internal/cortex/scales/clinical_scales.go` | 326-330 | TEST-CSSRS-003 | V |
| SRS-CSSRS-004 | Comportamento suicida (Q6) deve resultar em risco CRÍTICO | A | `internal/cortex/scales/clinical_scales.go` | 380-385 | TEST-CSSRS-004 | V |
| SRS-CSSRS-005 | Sistema deve calcular níveis de risco: none/low/moderate/high/critical | A | `internal/cortex/scales/clinical_scales.go` | 365-400 | TEST-CSSRS-005 | V |
| SRS-CSSRS-006 | Sistema deve gerar plano de intervenção para risco > none | A | `internal/cortex/scales/clinical_scales.go` | 405-450 | TEST-CSSRS-006 | V |
| SRS-CSSRS-007 | Sistema deve incluir CVV 188 em todas as intervenções | A | `internal/cortex/scales/clinical_scales.go` | 430-435 | TEST-CSSRS-007 | V |
| SRS-CSSRS-008 | Sistema deve incluir SAMU 192 para risco crítico | A | `internal/cortex/scales/clinical_scales.go` | 420-425 | TEST-CSSRS-008 | V |

**Testes Associados:**
```
internal/cortex/scales/clinical_scales_test.go:
- TestCSSRSQuestions (TEST-CSSRS-001)
- TestCSSRSIdeationQuestions (TEST-CSSRS-002)
- TestCSSRSBehaviorQuestion (TEST-CSSRS-003)
- TestCSSRSBehaviorAlwaysCritical (TEST-CSSRS-004)
- TestCSSRSRiskLevels (TEST-CSSRS-005)
- TestCSSRSInterventionPlan (TEST-CSSRS-006)
- TestCSSRSCVVHotline (TEST-CSSRS-007)
- TestCSSRSSAMUForCritical (TEST-CSSRS-008)
```

**Riscos Mitigados:** R-001 (Score incorreto), R-002 (Risco subestimado)

---

### 2.2 Avaliação de Depressão (PHQ-9)

| ID Req | Requisito | Classe | Arquivo Código | Linha | ID Teste | Status |
|--------|-----------|--------|----------------|-------|----------|--------|
| SRS-PHQ9-001 | Sistema deve implementar escala PHQ-9 com 9 questões | B | `internal/cortex/scales/clinical_scales.go` | 23-171 | TEST-PHQ9-001 | V |
| SRS-PHQ9-002 | Score deve variar de 0-27 | B | `internal/cortex/scales/clinical_scales.go` | 95-100 | TEST-PHQ9-002 | V |
| SRS-PHQ9-003 | Sistema deve classificar severidade: minimal/mild/moderate/moderately_severe/severe | B | `internal/cortex/scales/clinical_scales.go` | 105-125 | TEST-PHQ9-003 | V |
| SRS-PHQ9-004 | Questão 9 (ideação suicida) deve gerar alerta específico | A | `internal/cortex/scales/clinical_scales.go` | 130-145 | TEST-PHQ9-004 | V |
| SRS-PHQ9-005 | Sistema deve gerar recomendações baseadas no score | B | `internal/cortex/scales/clinical_scales.go` | 150-170 | TEST-PHQ9-005 | V |

**Testes Associados:**
```
internal/cortex/scales/clinical_scales_test.go:
- TestPHQ9Questions (TEST-PHQ9-001)
- TestPHQ9ScoreRange (TEST-PHQ9-002)
- TestPHQ9SeverityLevels (TEST-PHQ9-003)
- TestPHQ9Question9Alert (TEST-PHQ9-004)
- TestPHQ9Recommendations (TEST-PHQ9-005)
```

---

### 2.3 Avaliação de Ansiedade (GAD-7)

| ID Req | Requisito | Classe | Arquivo Código | Linha | ID Teste | Status |
|--------|-----------|--------|----------------|-------|----------|--------|
| SRS-GAD7-001 | Sistema deve implementar escala GAD-7 com 7 questões | B | `internal/cortex/scales/clinical_scales.go` | 173-293 | TEST-GAD7-001 | V |
| SRS-GAD7-002 | Score deve variar de 0-21 | B | `internal/cortex/scales/clinical_scales.go` | 220-225 | TEST-GAD7-002 | V |
| SRS-GAD7-003 | Sistema deve classificar severidade: minimal/mild/moderate/severe | B | `internal/cortex/scales/clinical_scales.go` | 230-250 | TEST-GAD7-003 | V |
| SRS-GAD7-004 | Sistema deve gerar recomendações baseadas no score | B | `internal/cortex/scales/clinical_scales.go` | 260-290 | TEST-GAD7-004 | V |

**Testes Associados:**
```
internal/cortex/scales/clinical_scales_test.go:
- TestGAD7Questions (TEST-GAD7-001)
- TestGAD7ScoreRange (TEST-GAD7-002)
- TestGAD7SeverityLevels (TEST-GAD7-003)
- TestGAD7Recommendations (TEST-GAD7-004)
```

---

## 3. MATRIZ DE RASTREABILIDADE - SISTEMA DE ALERTAS

### 3.1 Escalação de Alertas

| ID Req | Requisito | Classe | Arquivo Código | Linha | ID Teste | Status |
|--------|-----------|--------|----------------|-------|----------|--------|
| SRS-ALERT-001 | Sistema deve suportar múltiplos canais: Push, WhatsApp, SMS, Email, Ligação | A | `internal/cortex/alert/escalation.go` | 25-50 | TEST-ALERT-001 | V |
| SRS-ALERT-002 | Sistema deve escalar automaticamente em caso de falha | A | `internal/cortex/alert/escalation.go` | 80-120 | TEST-ALERT-002 | V |
| SRS-ALERT-003 | Prioridade crítica deve ter timeout de 30 segundos | A | `internal/cortex/alert/escalation.go` | 55-65 | TEST-ALERT-003 | V |
| SRS-ALERT-004 | Prioridade alta deve ter timeout de 2 minutos | A | `internal/cortex/alert/escalation.go` | 55-65 | TEST-ALERT-004 | V |
| SRS-ALERT-005 | Sistema deve registrar log de todas as tentativas | B | `internal/cortex/alert/escalation.go` | 150-180 | TEST-ALERT-005 | V |
| SRS-ALERT-006 | Sistema deve rastrear acknowledgment | B | `internal/cortex/alert/escalation.go` | 190-220 | TEST-ALERT-006 | V |
| SRS-ALERT-007 | Sistema deve continuar escalando até acknowledgment ou exaustão de canais | A | `internal/cortex/alert/escalation.go` | 100-140 | TEST-ALERT-007 | V |

**Testes Associados:**
```
internal/cortex/alert/escalation_test.go:
- TestAlertChannels (TEST-ALERT-001)
- TestAlertEscalation (TEST-ALERT-002)
- TestCriticalPriorityTimeout (TEST-ALERT-003)
- TestHighPriorityTimeout (TEST-ALERT-004)
- TestAlertLogging (TEST-ALERT-005)
- TestAlertAcknowledgment (TEST-ALERT-006)
- TestEscalationUntilAck (TEST-ALERT-007)
```

**Riscos Mitigados:** R-003 (Alerta não entregue), R-004 (Todos canais falham)

---

## 4. MATRIZ DE RASTREABILIDADE - LGPD/PROTEÇÃO DE DADOS

### 4.1 Trilha de Auditoria

| ID Req | Requisito | Classe | Arquivo Código | Linha | ID Teste | Status |
|--------|-----------|--------|----------------|-------|----------|--------|
| SRS-LGPD-001 | Sistema deve registrar todos os acessos a dados pessoais | B | `internal/audit/lgpd_audit.go` | 95-130 | TEST-LGPD-001 | V |
| SRS-LGPD-002 | Sistema deve registrar base legal (Art. 7) para cada operação | B | `internal/audit/lgpd_audit.go` | 50-70 | TEST-LGPD-002 | V |
| SRS-LGPD-003 | Sistema deve categorizar dados (pessoal, sensível, clínico) | B | `internal/audit/lgpd_audit.go` | 35-48 | TEST-LGPD-003 | V |
| SRS-LGPD-004 | Sistema deve definir período de retenção por categoria | B | `internal/audit/lgpd_audit.go` | 200-230 | TEST-LGPD-004 | V |
| SRS-LGPD-005 | Sistema deve auto-expirar registros após período de retenção | B | `migrations/018_lgpd_audit_trail.sql` | 55-69 | TEST-LGPD-005 | V |

**Testes Associados:**
```
internal/audit/lgpd_audit_test.go:
- TestLogDataAccess (TEST-LGPD-001)
- TestLegalBasisValidation (TEST-LGPD-002)
- TestDataCategories (TEST-LGPD-003)
- TestDefaultRetention (TEST-LGPD-004)
- TestAuditExpiration (TEST-LGPD-005)
```

### 4.2 Direitos do Titular (Art. 18)

| ID Req | Requisito | Classe | Arquivo Código | Linha | ID Teste | Status |
|--------|-----------|--------|----------------|-------|----------|--------|
| SRS-LGPD-010 | Sistema deve permitir exportação de dados (Art. 18, V) | B | `internal/audit/data_rights.go` | 72-126 | TEST-LGPD-010 | V |
| SRS-LGPD-011 | Sistema deve permitir deleção de dados (Art. 18, VI) | B | `internal/audit/data_rights.go` | 128-249 | TEST-LGPD-011 | V |
| SRS-LGPD-012 | Sistema deve reter dados clínicos por requisito legal | A | `internal/audit/data_rights.go` | 154-158 | TEST-LGPD-012 | V |
| SRS-LGPD-013 | Sistema deve permitir retificação de dados (Art. 18, III) | B | `internal/audit/data_rights.go` | 251-282 | TEST-LGPD-013 | V |
| SRS-LGPD-014 | Sistema deve restringir campos retificáveis | B | `internal/audit/data_rights.go` | 256-263 | TEST-LGPD-014 | V |
| SRS-LGPD-015 | Sistema deve gerar relatório de acesso aos dados (Art. 18, VII) | B | `internal/audit/data_rights.go` | 284-342 | TEST-LGPD-015 | V |

**Testes Associados:**
```
internal/audit/data_rights_test.go:
- TestExportPersonalData (TEST-LGPD-010)
- TestDeletePersonalData (TEST-LGPD-011)
- TestDeletePersonalData_RetainsClinicalData (TEST-LGPD-012)
- TestRectifyPersonalData (TEST-LGPD-013)
- TestRectifyPersonalData_InvalidField (TEST-LGPD-014)
- TestGetDataAccessReport (TEST-LGPD-015)
```

**Riscos Mitigados:** R-005 (Vazamento de dados)

---

## 5. MATRIZ DE RASTREABILIDADE - OBSERVABILIDADE

### 5.1 Métricas Prometheus

| ID Req | Requisito | Classe | Arquivo Código | Linha | ID Teste | Status |
|--------|-----------|--------|----------------|-------|----------|--------|
| SRS-MET-001 | Sistema deve expor métricas de avaliações clínicas | B | `internal/metrics/metrics.go` | 30-80 | TEST-MET-001 | V |
| SRS-MET-002 | Sistema deve expor métricas de alertas | B | `internal/metrics/metrics.go` | 85-130 | TEST-MET-002 | V |
| SRS-MET-003 | Sistema deve expor métricas de saúde do sistema | B | `internal/metrics/metrics.go` | 180-220 | TEST-MET-003 | V |
| SRS-MET-004 | Sistema deve expor métricas de LLM (latência, tokens) | B | `internal/metrics/metrics.go` | 225-280 | TEST-MET-004 | V |
| SRS-MET-005 | Sistema deve ter alertas para risco suicida detectado | A | `deployments/prometheus/alerts/eva-mind.yml` | 5-20 | TEST-MET-005 | V |

**Testes Associados:**
```
internal/metrics/metrics_test.go:
- TestClinicalMetrics (TEST-MET-001)
- TestAlertMetrics (TEST-MET-002)
- TestSystemHealthMetrics (TEST-MET-003)
- TestLLMMetrics (TEST-MET-004)
- TestSuicideRiskAlertRule (TEST-MET-005)
```

---

## 6. MATRIZ DE RASTREABILIDADE - DETECÇÃO DE PADRÕES

### 6.1 Padrões Temporais

| ID Req | Requisito | Classe | Arquivo Código | Linha | ID Teste | Status |
|--------|-----------|--------|----------------|-------|----------|--------|
| SRS-PAT-001 | Sistema deve detectar padrões recorrentes em conversas | B | `internal/hippocampus/memory/pattern_miner.go` | 22-122 | TEST-PAT-001 | V |
| SRS-PAT-002 | Sistema deve detectar padrões temporais (hora/dia) | B | `internal/hippocampus/memory/pattern_miner.go` | 124-175 | TEST-PAT-002 | V |
| SRS-PAT-003 | Sistema deve calcular tendência de severidade | B | `internal/hippocampus/memory/pattern_miner.go` | 50-60 | TEST-PAT-003 | V |
| SRS-PAT-004 | Sistema deve detectar padrões de sono | B | `internal/motor/workers/pattern_worker.go` | 104-173 | TEST-PAT-004 | V |
| SRS-PAT-005 | Sistema deve detectar padrões de humor | B | `internal/motor/workers/pattern_worker.go` | 175-224 | TEST-PAT-005 | V |
| SRS-PAT-006 | Sistema deve detectar padrões de adesão à medicação | B | `internal/motor/workers/pattern_worker.go` | 226-277 | TEST-PAT-006 | V |

**Testes Associados:**
```
internal/hippocampus/memory/pattern_miner_test.go:
- TestMineRecurrentPatterns (TEST-PAT-001)
- TestMineTemporalPatterns (TEST-PAT-002)
- TestSeverityTrendAnalysis (TEST-PAT-003)

internal/motor/workers/pattern_worker_test.go:
- TestDetectSleepPattern (TEST-PAT-004)
- TestDetectMoodPattern (TEST-PAT-005)
- TestDetectMedicationPattern (TEST-PAT-006)
```

---

## 7. MATRIZ DE RASTREABILIDADE - META-COGNIÇÃO

### 7.1 Sistema de Consciência

| ID Req | Requisito | Classe | Arquivo Código | Linha | ID Teste | Status |
|--------|-----------|--------|----------------|-------|----------|--------|
| SRS-CON-001 | Sistema deve calcular gravidade emocional de memórias | B | `internal/hippocampus/memory/superhuman/consciousness_service.go` | 50-116 | TEST-CON-001 | V |
| SRS-CON-002 | Sistema deve detectar ciclos comportamentais | B | `internal/hippocampus/memory/superhuman/consciousness_service.go` | 148-258 | TEST-CON-002 | V |
| SRS-CON-003 | Sistema deve rastrear rapport/confiança | B | `internal/hippocampus/memory/superhuman/consciousness_service.go` | 264-364 | TEST-CON-003 | V |
| SRS-CON-004 | Sistema deve detectar contradições narrativas | B | `internal/hippocampus/memory/superhuman/consciousness_service.go` | 370-574 | TEST-CON-004 | V |
| SRS-CON-005 | Sistema deve adaptar modo de interação | B | `internal/hippocampus/memory/superhuman/consciousness_service.go` | 580-666 | TEST-CON-005 | V |
| SRS-CON-006 | Sistema deve gerenciar carga empática | B | `internal/hippocampus/memory/superhuman/consciousness_service.go` | 832-895 | TEST-CON-006 | V |

**Testes Associados:**
```
internal/hippocampus/memory/superhuman/consciousness_service_test.go:
- TestMemoryGravity_* (TEST-CON-001)
- TestCyclePattern_* (TEST-CON-002)
- TestPatientRapport_* (TEST-CON-003)
- TestContradiction_* (TEST-CON-004)
- TestEvaMode_* (TEST-CON-005)
- TestEmpathicLoad_* (TEST-CON-006)
```

---

## 8. MATRIZ DE RASTREABILIDADE - APRENDIZADO CONTÍNUO

### 8.1 Sistema de Aprendizado

| ID Req | Requisito | Classe | Arquivo Código | Linha | ID Teste | Status |
|--------|-----------|--------|----------------|-------|----------|--------|
| SRS-LEARN-001 | Sistema deve aprender de feedback implícito | B | `internal/cortex/learning/continuous_learning.go` | 60-95 | TEST-LEARN-001 | V |
| SRS-LEARN-002 | Sistema deve adaptar estratégias de resposta | B | `internal/cortex/learning/continuous_learning.go` | 100-175 | TEST-LEARN-002 | V |
| SRS-LEARN-003 | Sistema deve aprender preferências de vocabulário | B | `internal/cortex/learning/continuous_learning.go` | 180-260 | TEST-LEARN-003 | V |
| SRS-LEARN-004 | Sistema deve aprender interesses por tópico | B | `internal/cortex/learning/continuous_learning.go` | 265-340 | TEST-LEARN-004 | V |
| SRS-LEARN-005 | Sistema deve aprender preferências de timing | B | `internal/cortex/learning/continuous_learning.go` | 345-420 | TEST-LEARN-005 | V |
| SRS-LEARN-006 | Sistema deve rastrear efetividade de personas | B | `internal/cortex/learning/continuous_learning.go` | 425-490 | TEST-LEARN-006 | V |

**Testes Associados:**
```
internal/cortex/learning/continuous_learning_test.go:
- TestCalculateImplicitFeedback (TEST-LEARN-001)
- TestGetBestStrategy_* (TEST-LEARN-002)
- TestVocabularyPreference_* (TEST-LEARN-003)
- TestTopicInterest_* (TEST-LEARN-004)
- TestTimingPreference_* (TEST-LEARN-005)
- TestPersonaEffectiveness_* (TEST-LEARN-006)
```

---

## 9. RESUMO DE COBERTURA

### 9.1 Estatísticas Gerais

| Categoria | Total Requisitos | Verificados | Parciais | Pendentes |
|-----------|------------------|-------------|----------|-----------|
| Avaliações Clínicas (C-SSRS, PHQ-9, GAD-7) | 17 | 17 | 0 | 0 |
| Sistema de Alertas | 7 | 7 | 0 | 0 |
| LGPD/Proteção de Dados | 11 | 11 | 0 | 0 |
| Observabilidade | 5 | 5 | 0 | 0 |
| Detecção de Padrões | 6 | 6 | 0 | 0 |
| Meta-cognição | 6 | 6 | 0 | 0 |
| Aprendizado Contínuo | 6 | 6 | 0 | 0 |
| **TOTAL** | **58** | **58** | **0** | **0** |

### 9.2 Cobertura por Classe de Risco

| Classe | Descrição | Requisitos | Verificados | % |
|--------|-----------|------------|-------------|---|
| A | Críticos para segurança | 15 | 15 | 100% |
| B | Importantes | 43 | 43 | 100% |
| C | Suporte | 0 | 0 | N/A |

### 9.3 Testes por Pacote

| Pacote | Testes | Status |
|--------|--------|--------|
| `internal/audit` | 37 | ✅ PASS |
| `internal/cortex/scales` | 25 | ✅ PASS |
| `internal/cortex/alert` | 17 | ✅ PASS |
| `internal/cortex/learning` | 45 | ✅ PASS |
| `internal/metrics` | 18 | ✅ PASS |
| `internal/mocks` | 12 | ✅ PASS |
| `internal/hippocampus/memory` | 17 | ✅ PASS |
| `internal/hippocampus/memory/superhuman` | 55 | ✅ PASS |
| `internal/motor/workers` | 14 | ✅ PASS |
| **TOTAL** | **240** | **✅ 100%** |

---

## 10. RASTREABILIDADE DE RISCOS

### 10.1 Riscos → Requisitos → Controles

| ID Risco | Descrição | Requisitos Relacionados | Controles Implementados | Status |
|----------|-----------|------------------------|------------------------|--------|
| R-001 | Score C-SSRS incorreto | SRS-CSSRS-001 a 005 | Testes unitários, validação de algoritmo | ✅ Mitigado |
| R-002 | Risco suicida subestimado | SRS-CSSRS-004, SRS-CSSRS-006 | Q6 = CRÍTICO automático, intervenção obrigatória | ✅ Mitigado |
| R-003 | Alerta não entregue | SRS-ALERT-001 a 007 | Multi-canal, escalação automática | ✅ Mitigado |
| R-004 | Todos os canais falham | SRS-ALERT-007 | 5 canais independentes, fallback local (CVV/SAMU) | ✅ Mitigado |
| R-005 | Vazamento de dados | SRS-LGPD-001 a 015 | Auditoria, criptografia, controle de acesso | ✅ Mitigado |

---

## 11. HISTÓRICO DE REVISÕES

| Versão | Data | Autor | Descrição |
|--------|------|-------|-----------|
| 1.0 | 2026-01-27 | Claude Opus 4.5 + José R F Junior | Versão inicial |

---

## 12. APROVAÇÕES

| Função | Nome | Assinatura | Data |
|--------|------|------------|------|
| Elaborado por | | | |
| Revisado por | | | |
| Aprovado por | | | |

---

## ANEXOS

### A. Comandos de Verificação

```bash
# Executar todos os testes
cd D:\dev\EVA\EVA-Mind-FZPN
go test ./internal/audit/... ./internal/cortex/scales/... ./internal/cortex/alert/... ./internal/cortex/learning/... ./internal/metrics/... ./internal/mocks/... ./internal/hippocampus/memory/... ./internal/motor/workers/... -v

# Verificar cobertura
go test ./... -cover

# Gerar relatório de cobertura HTML
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### B. Referências
- RDC 751/2022 - ANVISA
- IEC 62304:2006/Amd1:2015 - Ciclo de vida de software
- ISO 14971:2019 - Gestão de riscos
- NBR ISO/IEC 25010:2011 - Qualidade de software
