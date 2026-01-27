# EVA-Mind-FZPN - Guia de Testes

**Documento:** TEST-GUIDE-001
**Versão:** 1.0
**Data:** 2026-01-27
**Status:** ATIVO

---

## Resumo

Este documento descreve a estrutura de testes do EVA-Mind-FZPN, incluindo testes unitários, testes de integração e métricas de cobertura.

---

## 1. Estrutura de Testes

```
EVA-Mind-FZPN/
├── internal/
│   ├── cortex/
│   │   ├── scales/clinical_scales_test.go        # Escalas clínicas (C-SSRS, PHQ-9, GAD-7)
│   │   ├── cognitive/cognitive_load_orchestrator_test.go  # Carga cognitiva
│   │   ├── ethics/ethical_boundary_engine_test.go         # Fronteiras éticas
│   │   ├── alert/escalation_test.go              # Escalação de alertas
│   │   ├── medgemma/service_test.go              # Análise médica
│   │   ├── learning/continuous_learning_test.go  # Aprendizado contínuo
│   │   └── llm/thinking/detector_test.go         # Detector de pensamento
│   ├── persona/persona_manager_test.go           # Sistema de personas
│   ├── hippocampus/memory/
│   │   ├── pattern_miner_test.go                 # Mineração de padrões
│   │   └── superhuman/consciousness_service_test.go  # Consciência
│   ├── audit/
│   │   ├── lgpd_audit_test.go                    # Auditoria LGPD
│   │   └── data_rights_test.go                   # Direitos de dados
│   ├── metrics/metrics_test.go                   # Métricas Prometheus
│   ├── mocks/mocks_test.go                       # Validação de mocks
│   └── motor/workers/pattern_worker_test.go      # Workers de padrões
│
├── test/
│   └── integration/
│       ├── suite_test.go                         # Setup de integração
│       ├── alert_integration_test.go             # Testes de alertas
│       ├── cognitive_integration_test.go         # Testes cognitivos
│       ├── trajectory_integration_test.go        # Testes de trajetória
│       ├── research_integration_test.go          # Testes de pesquisa
│       └── e2e_integration_test.go               # Testes end-to-end
│
└── Makefile                                      # Comandos de teste
```

---

## 2. Executando Testes

### 2.1 Testes Unitários (Rápidos)

```bash
# Todos os testes unitários
make test-unit

# Testes críticos apenas
make test-critical

# Testes específicos
go test ./internal/cortex/scales/... -v
```

### 2.2 Testes de Integração (Requer DB)

```bash
# Configurar variável de ambiente
export DATABASE_URL="postgres://user:pass@host:5432/eva_test"

# Executar testes de integração
make test-integration

# Ou diretamente
go test ./test/integration/... -v -timeout 5m
```

### 2.3 Cobertura de Código

```bash
# Gerar relatório de cobertura
make test-coverage

# Relatório detalhado por pacote
make test-coverage-detailed

# Abrir relatório HTML
open coverage/coverage.html
```

### 2.4 Benchmarks

```bash
make test-bench
```

---

## 3. Testes Críticos (Prioridade Máxima)

### 3.1 C-SSRS (Columbia Suicide Severity Rating Scale)

**Arquivo:** `internal/cortex/scales/clinical_scales_test.go`

| Teste | Descrição | Criticidade |
|-------|-----------|-------------|
| `TestGetCSSRSQuestions` | Verifica 6 perguntas | 🔴 CRÍTICA |
| `TestCSSRSCalculation_NoRisk` | Todas negativas = sem risco | 🔴 CRÍTICA |
| `TestCSSRSCalculation_LowRisk` | Q1 positiva = risco baixo | 🔴 CRÍTICA |
| `TestCSSRSCalculation_ModerateRisk` | Q1+Q2 = risco moderado | 🔴 CRÍTICA |
| `TestCSSRSCalculation_HighRisk_WithPlan` | Q1-Q4 = risco alto | 🔴 CRÍTICA |
| `TestCSSRSCalculation_CriticalRisk_Behavior` | Q6 = CRÍTICO | 🔴 CRÍTICA |
| `TestCSSRSInterventions_ContainEmergencyInfo` | SAMU 192, CVV 188 | 🔴 CRÍTICA |

### 3.2 Sistema de Alertas

**Arquivo:** `internal/cortex/alert/escalation_test.go`

| Teste | Descrição | Criticidade |
|-------|-----------|-------------|
| `TestEscalationChain` | Push → SMS → Email → Call | 🔴 CRÍTICA |
| `TestAlertTimeout` | Escalonamento após timeout | 🔴 CRÍTICA |
| `TestCriticalAlertCreation` | Criação de alerta crítico | 🔴 CRÍTICA |

### 3.3 Escalas PHQ-9 e GAD-7

| Teste | Descrição | Criticidade |
|-------|-----------|-------------|
| `TestPHQ9Calculation_*` | Scores 0-27 | 🟠 ALTA |
| `TestPHQ9Calculation_SuicideRisk` | Q9 positiva | 🔴 CRÍTICA |
| `TestGAD7Calculation_*` | Scores 0-21 | 🟠 ALTA |

---

## 4. Testes de Integração

### 4.1 Suite de Testes

**Arquivo:** `test/integration/suite_test.go`

- Conexão automática com banco de dados
- Criação de paciente de teste
- Limpeza automática após testes

### 4.2 Fluxo Completo de Crise

**Arquivo:** `test/integration/e2e_integration_test.go`

```go
TestE2E_CrisisDetectionToIntervention
├── Step1_InitialAssessment      // PHQ-9 moderado
├── Step2_CognitiveLoadIncrease  // Carga aumentando
├── Step3_VoiceProsodyChange     // Voz alterada
├── Step4_CSSRSTriggered         // C-SSRS ativado
├── Step5_TrajectoryPrediction   // Simulação Monte Carlo
├── Step6_InterventionGeneration // Recomendações geradas
├── Step7_AlertEscalation        // Alerta escalado
└── Step8_VerifyWorkflow         // Verificação completa
```

### 4.3 Troca de Personas

```go
TestE2E_PersonaSwitchingOnContext
├── StartWithCompanion   // Início padrão
├── SwitchToEducator     // Pedido de educação
├── SwitchToEmergency    // Crise detectada
├── SwitchToClinical     // Pós-crise
└── VerifyTransitions    // 3 transições corretas
```

---

## 5. Mocks

### 5.1 Interfaces Mockadas

**Arquivo:** `internal/mocks/interfaces.go`

| Interface | Descrição |
|-----------|-----------|
| `PushService` | Firebase push notifications |
| `SMSService` | Twilio SMS |
| `VoiceService` | Twilio chamadas |
| `EmailService` | SMTP email |
| `AlertService` | Serviço de alertas |
| `CSSRSService` | Escala C-SSRS |
| `PHQ9Service` | Escala PHQ-9 |
| `GAD7Service` | Escala GAD-7 |

### 5.2 Mocks Disponíveis

**Arquivos em `internal/mocks/`:**

- `firebase_mock.go` - Mock Firebase
- `twilio_mock.go` - Mock Twilio (SMS + Voice)
- `email_mock.go` - Mock SMTP
- `alert_mock.go` - Mock de alertas

---

## 6. Métricas de Cobertura

### 6.1 Targets de Cobertura

| Pacote | Target | Prioridade |
|--------|--------|------------|
| `cortex/scales` | ≥90% | 🔴 CRÍTICA |
| `cortex/alert` | ≥90% | 🔴 CRÍTICA |
| `cortex/cognitive` | ≥80% | 🟠 ALTA |
| `cortex/ethics` | ≥80% | 🟠 ALTA |
| `persona` | ≥80% | 🟠 ALTA |
| `research` | ≥70% | 🟡 MÉDIA |
| `audit` | ≥80% | 🟠 ALTA |

### 6.2 Verificando Cobertura

```bash
# Cobertura total
go tool cover -func=coverage/coverage.out | tail -1

# Cobertura por arquivo
go tool cover -func=coverage/coverage.out | grep -E "scales|alert"

# Relatório HTML
go tool cover -html=coverage/coverage.out -o coverage/coverage.html
```

---

## 7. CI/CD

### 7.1 GitHub Actions (Sugerido)

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_DB: eva_test
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Run unit tests
        run: make test-unit

      - name: Run integration tests
        run: make test-integration
        env:
          TEST_DATABASE_URL: postgres://postgres:postgres@localhost:5432/eva_test

      - name: Generate coverage
        run: make test-coverage

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage/coverage.out
```

---

## 8. Boas Práticas

### 8.1 Escrevendo Testes

```go
func TestFeature_Scenario(t *testing.T) {
    // Arrange
    input := createTestInput()

    // Act
    result := functionUnderTest(input)

    // Assert
    assert.Equal(t, expected, result)
}
```

### 8.2 Table-Driven Tests

```go
func TestFeature_AllScenarios(t *testing.T) {
    testCases := []struct {
        name     string
        input    int
        expected string
    }{
        {"scenario1", 10, "result1"},
        {"scenario2", 20, "result2"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := functionUnderTest(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

### 8.3 Cleanup

```go
func TestWithCleanup(t *testing.T) {
    // Setup
    resource := createResource()

    // Cleanup garantido
    t.Cleanup(func() {
        resource.Close()
    })

    // Test
    // ...
}
```

---

## 9. Troubleshooting

### 9.1 Database Connection Failed

```bash
# Verificar conectividade
psql "$DATABASE_URL" -c "SELECT 1"

# Usar banco local
export TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/eva_test"
```

### 9.2 Tests Timeout

```bash
# Aumentar timeout
go test ./... -timeout 10m
```

### 9.3 Race Conditions

```bash
# Detectar race conditions
go test ./... -race
```

---

## Aprovações

| Função | Nome | Data |
|--------|------|------|
| Criador/Admin | Jose R F Junior | 2026-01-27 |

---

**Testes são fundamentais para a segurança dos pacientes. Execute-os antes de cada deploy.**
