# ✅ SPRINT 7: Integration Layer (Go Helpers) - COMPLETO

**Status:** ✅ IMPLEMENTADO
**Data:** 2026-01-24
**Complexidade:** 🟡 ALTA
**Impacto:** 🔥🔥🔥🔥 MUITO ALTO
**Abordagem:** Go Helpers → Python API (hybrid approach)

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Arquitetura Híbrida](#arquitetura-híbrida)
3. [O Que Foi Criado](#o-que-foi-criado)
4. [Schema SQL](#schema-sql)
5. [Go Helpers](#go-helpers)
6. [Como Integrar com Python](#como-integrar-com-python)
7. [Exemplos de Uso](#exemplos-de-uso)
8. [FHIR R4 Mapping](#fhir-r4-mapping)
9. [Webhooks](#webhooks)
10. [Data Export (LGPD/GDPR)](#data-export-lgpdgdpr)
11. [Próximos Passos](#próximos-passos)

---

## 🎯 Visão Geral

O **Integration Layer** fornece helpers Go para facilitar a integração com APIs REST externas (especialmente sua API Python).

### O Que NÃO Foi Criado
❌ Servidor HTTP em Go (você já tem Python para isso)
❌ Rotas/endpoints Go
❌ Controllers/Views Go

### O Que FOI Criado
✅ **JSON Serializers** - DTOs otimizados para API
✅ **FHIR R4 Adapters** - Converter EVA → FHIR padrão HL7
✅ **Webhook Payload Builders** - Criar eventos estruturados
✅ **Data Export Utilities** - LGPD/GDPR compliance, research exports
✅ **Schema SQL** - Tracking de API calls, OAuth2 tokens, webhooks

---

## 🏗️ Arquitetura Híbrida

```
┌──────────────────────────────────────────────────────────────────┐
│                         FRONTEND                                  │
│              (Web App, Mobile App, Hospital System)               │
└────────────────────────┬─────────────────────────────────────────┘
                         │ HTTP/REST
                         ▼
┌──────────────────────────────────────────────────────────────────┐
│                    PYTHON API (FastAPI/Flask)                     │
│  - Endpoints REST                                                 │
│  - OAuth2 authentication                                          │
│  - Rate limiting                                                  │
│  - Request/response handling                                      │
└────────────────────────┬─────────────────────────────────────────┘
                         │ Calls Go functions via subprocess/gRPC/HTTP
                         ▼
┌──────────────────────────────────────────────────────────────────┐
│                    GO BUSINESS LOGIC                              │
│  - internal/persona/persona_manager.go                            │
│  - internal/exit/exit_protocol_manager.go                         │
│  - internal/research/research_engine.go                           │
│  - internal/trajectory/trajectory_engine.go                       │
└────────────────────────┬─────────────────────────────────────────┘
                         │ Uses helpers
                         ▼
┌──────────────────────────────────────────────────────────────────┐
│                GO INTEGRATION HELPERS (NOVO!)                     │
│  - internal/integration/serializers.go                            │
│  - internal/integration/fhir_adapter.go                           │
│  - internal/integration/webhooks.go                               │
│  - internal/integration/export.go                                 │
└────────────────────────┬─────────────────────────────────────────┘
                         │
                         ▼
                   ┌─────────────┐
                   │ PostgreSQL  │
                   └─────────────┘
```

---

## 📦 O Que Foi Criado

### 1. **migrations/010_integration_layer.sql** (~800 linhas)

**8 Tabelas:**
- `api_clients` - Aplicações autorizadas (OAuth2 clients)
- `api_tokens` - Access tokens (Bearer tokens)
- `api_request_logs` - Auditoria de todas as chamadas
- `webhook_deliveries` - Fila de webhooks assíncronos
- `rate_limit_tracking` - Proteção contra abuso
- `fhir_resource_mappings` - Sync com sistemas FHIR
- `external_system_credentials` - Credenciais de EHRs/labs
- `data_export_jobs` - Jobs de exportação (LGPD, research)

**4 Views:**
- `v_api_usage_stats` - Estatísticas por cliente
- `v_pending_webhooks` - Webhooks prontos para retry
- `v_fhir_sync_issues` - Recursos com problemas de sync
- `v_top_api_endpoints` - Endpoints mais usados

**3 Triggers:**
- Auto-update `updated_at`
- Increment token usage counter
- Cleanup de logs antigos (90 dias)

---

### 2. **internal/integration/serializers.go** (~400 linhas)

**30+ DTOs (Data Transfer Objects):**

```go
// Pacientes
type PatientDTO struct { ... }
type PatientListDTO struct { ... }

// Assessments
type AssessmentDTO struct { ... }
type AssessmentSummaryDTO struct { ... }

// Voice Analysis
type VoiceAnalysisDTO struct { ... }

// Trajectory
type TrajectoryDTO struct { ... }
type PredictionPointDTO struct { ... }

// Research
type ResearchStudyDTO struct { ... }
type FindingDTO struct { ... }

// Personas
type PersonaDTO struct { ... }
type PersonaSessionDTO struct { ... }

// Exit Protocol
type LastWishesDTO struct { ... }
type QualityOfLifeDTO struct { ... }
type PainLogDTO struct { ... }
type LegacyMessageDTO struct { ... }

// Responses padrão
type PaginatedResponse struct { ... }
type ErrorResponse struct { ... }
type SuccessResponse struct { ... }
```

**Utility Functions:**
```go
ToJSON(v interface{}) (string, error)
ToJSONCompact(v interface{}) (string, error)
FromJSON(jsonStr string, v interface{}) error
NewPaginatedResponse(...) *PaginatedResponse
NewErrorResponse(...) *ErrorResponse
NewSuccessResponse(...) *SuccessResponse
```

---

### 3. **internal/integration/fhir_adapter.go** (~600 linhas)

**FHIR R4 Resources Implementados:**

```go
// Pacientes
type FHIRPatient struct { ... }
PatientDTOToFHIR(*PatientDTO) *FHIRPatient

// Observações (Assessments, sinais vitais)
type FHIRObservation struct { ... }
PHQ9ToFHIR(*AssessmentDTO) *FHIRObservation

// Questionários
type FHIRQuestionnaireResponse struct { ... }

// Condições/Diagnósticos
type FHIRCondition struct { ... }

// Medicações
type FHIRMedicationRequest struct { ... }

// Bundles (coleções)
type FHIRBundle struct { ... }
CreatePatientBundle(*FHIRPatient, []*FHIRObservation) *FHIRBundle
```

**FHIR Mapping Completo:**

| EVA Resource | FHIR Resource | LOINC/SNOMED Codes |
|--------------|---------------|---------------------|
| Patient      | Patient       | N/A |
| PHQ-9        | Observation   | LOINC 44249-1 |
| GAD-7        | Observation   | LOINC 69737-5 |
| C-SSRS       | Observation   | LOINC 73831-0 |
| Medication   | MedicationRequest | RxNorm codes |
| Voice Analysis | Observation | Custom extension |
| Pain Log     | Observation   | LOINC 38208-5 |

---

### 4. **internal/integration/webhooks.go** (~400 linhas)

**15+ Event Types:**

```go
// Pacientes
PatientCreatedEvent(*PatientDTO) *WebhookEvent
PatientUpdatedEvent(patientID, changes) *WebhookEvent

// Assessments
AssessmentCompletedEvent(*AssessmentDTO) *WebhookEvent
SuicideRiskDetectedEvent(patientID, assessmentID, score) *WebhookEvent

// Crises
CrisisDetectedEvent(patientID, crisisType, severity, details) *WebhookEvent

// Personas
PersonaTransitionEvent(patientID, from, to, reason) *WebhookEvent

// Exit Protocol
PainAlertEvent(patientID, *PainLogDTO) *WebhookEvent
QualityOfLifeChangedEvent(patientID, oldScore, newScore, trend) *WebhookEvent

// Research
ResearchFindingEvent(studyID, *FindingDTO) *WebhookEvent

// Medicações
MedicationAdherenceAlertEvent(patientID, medication, missedDoses) *WebhookEvent

// Trajectory
TrajectoryRiskIncreasedEvent(patientID, riskType, oldRisk, newRisk) *WebhookEvent

// Custom
CustomEvent(eventType string, data map[string]interface{}) *WebhookEvent
```

**Webhook Security (HMAC-SHA256):**
```go
SignWebhookPayload(payload, secret) string
VerifyWebhookSignature(payload, signature, secret) bool
event.AddSignature(secret) error
```

---

### 5. **internal/integration/export.go** (~600 linhas)

**Export Types:**

```go
// LGPD/GDPR Portabilidade
type LGPDPortabilityExport struct { ... }
NewLGPDPortabilityExport(patientID) *LGPDPortabilityExport

// Research Dataset
type ResearchDatasetExport struct { ... }

// Clinical Summary
type ClinicalSummaryExport struct { ... }

// Export Job (tracking)
type ExportJob struct { ... }
```

**Anonymization Utilities:**
```go
AnonymizePatientID(patientID int64) string // SHA-256
AnonymizeName(name string) string         // "João" → "Jo*****"
AnonymizeEmail(email string) string       // "user@domain.com" → "us*****@domain.com"
RemoveSensitiveFields(data, sensitiveFields) map[string]interface{}
```

**CSV Export:**
```go
type CSVExportConfig struct { ... }
GenerateCSVHeader(columns, delimiter) string
RowToCSV(row, columns, delimiter) string
```

**FHIR Bundle Export:**
```go
ExportPatientAsFHIRBundle(patient, assessments) (*FHIRBundle, error)
```

**Compliance Checks:**
```go
RunLGPDComplianceChecks(*LGPDPortabilityExport) []ComplianceCheck
```

**Templates:**
```go
DefaultLGPDExportConfig
DefaultResearchExportConfig
DefaultClinicalSummaryConfig
```

---

## 🗄️ Schema SQL

### Tabela: `api_clients`

Registra aplicações autorizadas a usar a API.

```sql
CREATE TABLE api_clients (
    id UUID PRIMARY KEY,
    client_name VARCHAR(200),
    client_type VARCHAR(50), -- 'web_app', 'mobile_app', 'hospital_system', 'ehr_system'

    -- OAuth2
    client_id VARCHAR(100) UNIQUE,
    client_secret_hash VARCHAR(256), -- bcrypt

    -- Permissões
    scopes TEXT[], -- ['read:patients', 'write:assessments']
    allowed_endpoints TEXT[],

    -- Rate limiting
    rate_limit_per_minute INTEGER DEFAULT 60,
    rate_limit_per_hour INTEGER DEFAULT 1000,

    -- Webhook callback
    webhook_url VARCHAR(500),
    webhook_secret VARCHAR(256),
    webhook_events TEXT[],

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    is_approved BOOLEAN DEFAULT FALSE, -- Requer aprovação manual

    -- Metadados
    organization VARCHAR(200),
    contact_email VARCHAR(200),

    created_at TIMESTAMP DEFAULT NOW()
);
```

**Exemplo:**
```sql
INSERT INTO api_clients (
    client_name, client_type, client_id, client_secret_hash,
    scopes, rate_limit_per_minute, webhook_url, webhook_events
) VALUES (
    'Hospital ABC EHR',
    'hospital_system',
    'hosp_abc_12345',
    '$2a$10$...',
    ARRAY['read:patients', 'read:assessments', 'write:observations'],
    120,
    'https://hospital-abc.com/webhooks/eva-mind',
    ARRAY['assessment.completed', 'crisis.detected']
);
```

---

### Tabela: `api_tokens`

Access tokens OAuth2.

```sql
CREATE TABLE api_tokens (
    id UUID PRIMARY KEY,
    client_id UUID REFERENCES api_clients(id),

    -- Token
    access_token VARCHAR(256) UNIQUE, -- JWT
    refresh_token VARCHAR(256),

    -- Escopo
    scopes TEXT[],

    -- Expiração
    expires_at TIMESTAMP,

    -- Revogação
    is_revoked BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT NOW()
);
```

---

### Tabela: `api_request_logs`

Auditoria completa de chamadas.

```sql
CREATE TABLE api_request_logs (
    id BIGSERIAL PRIMARY KEY,
    request_id UUID UNIQUE,

    client_id UUID REFERENCES api_clients(id),
    token_id UUID REFERENCES api_tokens(id),

    -- Request
    http_method VARCHAR(10), -- GET, POST, PUT, DELETE
    endpoint VARCHAR(500),
    query_params JSONB,

    -- Headers
    user_agent TEXT,
    ip_address INET,

    -- Response
    http_status_code INTEGER,
    response_time_ms INTEGER,
    error_message TEXT,

    -- Rate limiting
    rate_limit_hit BOOLEAN DEFAULT FALSE,

    timestamp TIMESTAMP DEFAULT NOW()
);
```

**Partitioning por mês recomendado** para performance em grandes volumes.

---

### Tabela: `webhook_deliveries`

Fila de webhooks assíncronos.

```sql
CREATE TABLE webhook_deliveries (
    id UUID PRIMARY KEY,
    client_id UUID REFERENCES api_clients(id),

    -- Evento
    event_type VARCHAR(100), -- 'patient.created', 'crisis.detected'
    event_data JSONB,

    -- Delivery
    webhook_url VARCHAR(500),

    -- Status
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'sent', 'failed'
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,

    -- Resultado
    last_attempt_at TIMESTAMP,
    last_http_status INTEGER,
    last_error_message TEXT,

    -- Retry
    next_retry_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT NOW()
);
```

**Retry Strategy:** Exponential backoff (1min, 5min, 15min).

---

### View: `v_api_usage_stats`

Estatísticas de uso por cliente.

```sql
SELECT
    client_name,
    total_requests,
    successful_requests,
    failed_requests,
    avg_response_time_ms,
    p95_response_time_ms,
    rate_limit_hits,
    last_request_at
FROM v_api_usage_stats
WHERE client_id = 'hosp_abc_12345';
```

---

## 🔧 Go Helpers

### Serializers

**Uso típico:**

```go
package main

import "eva-mind/internal/integration"

func main() {
    // Criar DTO
    patient := &integration.PatientDTO{
        ID:          1,
        Name:        "João Silva",
        Age:         72,
        Gender:      "M",
        DateOfBirth: "1952-03-15",
    }

    // Converter para JSON
    json, _ := integration.ToJSON(patient)
    fmt.Println(json)

    // Criar resposta paginada
    patients := []integration.PatientListDTO{ /* ... */ }
    response := integration.NewPaginatedResponse(patients, 1, 20, 150)

    // Criar resposta de erro
    errorResp := integration.NewErrorResponse("validation_error", "Email inválido")
}
```

**Output JSON:**
```json
{
  "id": 1,
  "name": "João Silva",
  "age": 72,
  "gender": "M",
  "date_of_birth": "1952-03-15",
  "created_at": "2026-01-24T10:30:00Z",
  "updated_at": "2026-01-24T10:30:00Z"
}
```

---

### FHIR Adapter

**Uso típico:**

```go
package main

import "eva-mind/internal/integration"

func main() {
    // Paciente EVA
    patient := &integration.PatientDTO{
        ID:          1,
        Name:        "João Silva",
        Gender:      "M",
        DateOfBirth: "1952-03-15",
    }

    // Converter para FHIR R4
    fhirPatient := integration.PatientDTOToFHIR(patient)

    // Assessment PHQ-9
    assessment := &integration.AssessmentDTO{
        ID:             "a123",
        PatientID:      1,
        AssessmentType: "PHQ-9",
        TotalScore:     &score, // score = 15
        Status:         "completed",
    }

    // Converter para FHIR Observation
    fhirObs := integration.PHQ9ToFHIR(assessment)

    // Criar bundle
    bundle := integration.CreatePatientBundle(fhirPatient, []*integration.FHIRObservation{fhirObs})

    // Serializar para JSON FHIR
    fhirJSON, _ := integration.ToFHIRJSON(bundle)
    fmt.Println(fhirJSON)
}
```

**Output FHIR JSON:**
```json
{
  "resourceType": "Bundle",
  "type": "collection",
  "total": 2,
  "entry": [
    {
      "fullUrl": "Patient/1",
      "resource": {
        "resourceType": "Patient",
        "id": "1",
        "identifier": [{
          "system": "https://eva-mind.com/patient-id",
          "value": "1"
        }],
        "name": [{
          "use": "official",
          "text": "João Silva"
        }],
        "gender": "male",
        "birthDate": "1952-03-15"
      }
    },
    {
      "fullUrl": "Observation/a123",
      "resource": {
        "resourceType": "Observation",
        "id": "a123",
        "status": "final",
        "code": {
          "coding": [{
            "system": "http://loinc.org",
            "code": "44249-1",
            "display": "PHQ-9 quick depression assessment panel"
          }]
        },
        "subject": {
          "reference": "Patient/1"
        },
        "valueInteger": 15,
        "interpretation": [{
          "text": "Moderately severe depression"
        }]
      }
    }
  ]
}
```

Este JSON FHIR pode ser enviado para qualquer sistema hospitalar compatível com HL7 FHIR R4!

---

### Webhooks

**Criar evento:**

```go
package main

import "eva-mind/internal/integration"

func main() {
    // Assessment completado
    assessment := &integration.AssessmentDTO{ /* ... */ }
    event := integration.AssessmentCompletedEvent(assessment)

    // Crise detectada
    event2 := integration.SuicideRiskDetectedEvent(1, "a123", 4)

    // Dor severa
    painLog := &integration.PainLogDTO{
        PatientID:     1,
        PainIntensity: 8,
        PainLocation:  []string{"abdomen"},
    }
    event3 := integration.PainAlertEvent(1, painLog)

    // Adicionar assinatura HMAC
    secret := "webhook_secret_key_123"
    event.AddSignature(secret)

    // Serializar para enviar via HTTP POST
    payload, _ := event.ToJSON()

    // No receptor (Python API do hospital):
    // 1. Verificar assinatura
    // 2. Processar evento
}
```

**Payload JSON:**
```json
{
  "id": "20260124103045",
  "type": "assessment.completed",
  "timestamp": "2026-01-24T10:30:45Z",
  "source": "EVA-Mind",
  "data": {
    "assessment_id": "a123",
    "patient_id": 1,
    "assessment_type": "PHQ-9",
    "total_score": 15,
    "severity": "moderate",
    "flags": ["sleep_disturbance", "fatigue"]
  },
  "signature": "abc123def456..." // HMAC-SHA256
}
```

---

### Data Export

**LGPD Portability Export:**

```go
package main

import "eva-mind/internal/integration"

func main() {
    // Criar export LGPD
    export := integration.NewLGPDPortabilityExport(1)

    // Adicionar dados
    export.Patient = &integration.PatientDTO{ /* ... */ }
    export.Assessments = []integration.AssessmentDTO{ /* ... */ }
    export.LastWishes = &integration.LastWishesDTO{ /* ... */ }

    // Adicionar consentimento
    export.ConsentHistory = []integration.ConsentRecord{
        {
            Purpose:   "treatment",
            Granted:   true,
            GrantedAt: time.Now(),
        },
    }

    // Adicionar logs de processamento
    export.DataProcessing = []integration.DataProcessingLog{
        {
            Activity:    "phq9_administration",
            ProcessedAt: time.Now(),
            ProcessedBy: "EVA-Clinical",
            LegalBasis:  "consent",
        },
    }

    // Verificar compliance LGPD
    checks := integration.RunLGPDComplianceChecks(export)
    for _, check := range checks {
        fmt.Printf("%s: %v\n", check.Rule, check.Compliant)
    }

    // Serializar para JSON
    json, _ := integration.ToJSON(export)

    // Salvar em arquivo
    // ioutil.WriteFile("patient_1_lgpd_export.json", []byte(json), 0644)
}
```

**Research Export Anonimizado:**

```go
// Criar datapoints anonimizados
datapoint := integration.AnonymizedDatapoint{
    AnonymousID: integration.AnonymizePatientID(1), // SHA-256
    DayOffset:   7, // 7 dias após baseline
    Variables: map[string]interface{}{
        "phq9":             12,
        "voice_pitch_mean": 142.5,
        "medication_adherence": 0.85,
    },
}

// Adicionar a dataset
dataset := &integration.ResearchDatasetExport{
    StudyID:    "study123",
    StudyCode:  "EVA-VOICE-PHQ9-001",
    Datapoints: []integration.AnonymizedDatapoint{datapoint},
    Metadata: integration.ResearchMetadata{
        AnonymizationMethod: "SHA-256",
        KAnonymity:          5,
        Variables:           []string{"phq9", "voice_pitch_mean", "medication_adherence"},
    },
}

// Export como CSV
csv := ""
csv += integration.GenerateCSVHeader(columns, ",")
for _, dp := range dataset.Datapoints {
    csv += integration.RowToCSV(dp.Variables, columns, ",")
}
```

---

## 🔗 Como Integrar com Python

### Opção 1: Subprocess (Mais Simples)

**Go:** Compile helpers em executável
```bash
go build -o eva_helpers cmd/api_helpers/main.go
```

**Python:** Chame via subprocess
```python
import subprocess
import json

# Converter PatientDTO para JSON
patient_json = subprocess.run(
    ["./eva_helpers", "serialize", "patient", "1"],
    capture_output=True,
    text=True
).stdout

patient = json.loads(patient_json)

# Converter para FHIR
fhir_json = subprocess.run(
    ["./eva_helpers", "fhir", "patient", "1"],
    capture_output=True,
    text=True
).stdout
```

---

### Opção 2: gRPC (Mais Performático)

**Go:** Criar servidor gRPC
```go
// cmd/grpc_server/main.go
service IntegrationService {
    rpc SerializePatient(PatientRequest) returns (PatientResponse);
    rpc ConvertToFHIR(FHIRRequest) returns (FHIRResponse);
    rpc CreateWebhookEvent(EventRequest) returns (EventResponse);
}
```

**Python:** Cliente gRPC
```python
import grpc
import integration_pb2
import integration_pb2_grpc

channel = grpc.insecure_channel('localhost:50051')
stub = integration_pb2_grpc.IntegrationServiceStub(channel)

response = stub.SerializePatient(integration_pb2.PatientRequest(id=1))
patient_json = response.json_data
```

---

### Opção 3: HTTP Microservice (Mais Flexível)

**Go:** HTTP server simples
```go
// cmd/integration_service/main.go
http.HandleFunc("/serialize/patient", serializePatientHandler)
http.HandleFunc("/fhir/patient", fhirPatientHandler)
http.HandleFunc("/webhook/create", webhookHandler)
http.ListenAndServe(":8081", nil)
```

**Python:** HTTP client
```python
import requests

response = requests.get("http://localhost:8081/serialize/patient/1")
patient_json = response.json()

fhir_response = requests.get("http://localhost:8081/fhir/patient/1")
fhir_bundle = fhir_response.json()
```

---

### Opção 4: Shared Library (C-compatible)

**Go:** Compile como .so/.dll
```bash
go build -buildmode=c-shared -o libeva_integration.so
```

**Python:** Usar ctypes
```python
import ctypes
import json

lib = ctypes.CDLL('./libeva_integration.so')

# Chamar função Go
result = lib.SerializePatient(1)
patient_json = ctypes.string_at(result).decode('utf-8')
```

---

## 📚 Exemplos Completos de Uso

### Exemplo 1: API Python (FastAPI) + Go Helpers via HTTP

**Go Service (port 8081):**
```go
// cmd/integration_service/main.go
package main

import (
    "encoding/json"
    "eva-mind/internal/integration"
    "net/http"
    "strconv"
)

func serializePatientHandler(w http.ResponseWriter, r *http.Request) {
    patientIDStr := r.URL.Query().Get("id")
    patientID, _ := strconv.ParseInt(patientIDStr, 10, 64)

    // Buscar paciente do DB
    patient := getPatientFromDB(patientID) // Sua função

    // Converter para DTO
    dto := &integration.PatientDTO{
        ID:   patient.ID,
        Name: patient.Nome,
        // ... mapping
    }

    // Serializar
    jsonData, _ := integration.ToJSON(dto)

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(jsonData))
}

func main() {
    http.HandleFunc("/serialize/patient", serializePatientHandler)
    http.HandleFunc("/fhir/patient", fhirPatientHandler)
    http.ListenAndServe(":8081", nil)
}
```

**Python API (FastAPI, port 8000):**
```python
# main.py
from fastapi import FastAPI
import requests

app = FastAPI()
GO_SERVICE = "http://localhost:8081"

@app.get("/api/v1/patients/{patient_id}")
async def get_patient(patient_id: int):
    # Buscar do Go service (que busca do DB e serializa)
    response = requests.get(f"{GO_SERVICE}/serialize/patient?id={patient_id}")
    patient = response.json()
    return patient

@app.get("/api/v1/fhir/patients/{patient_id}")
async def get_patient_fhir(patient_id: int):
    # Buscar versão FHIR
    response = requests.get(f"{GO_SERVICE}/fhir/patient?id={patient_id}")
    fhir_bundle = response.json()
    return fhir_bundle
```

---

### Exemplo 2: Webhook Sender (Go) → Webhook Receiver (Python)

**Go (envia webhooks):**
```go
package main

import (
    "bytes"
    "eva-mind/internal/integration"
    "net/http"
)

func sendWebhook(event *integration.WebhookEvent, url string, secret string) error {
    // Adicionar assinatura
    event.AddSignature(secret)

    // Serializar
    payload, _ := event.ToJSON()

    // Enviar HTTP POST
    resp, err := http.Post(url, "application/json", bytes.NewBufferString(payload))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return fmt.Errorf("webhook failed: %d", resp.StatusCode)
    }

    return nil
}

// Quando assessment completa
func onAssessmentCompleted(assessment *Assessment) {
    dto := toAssessmentDTO(assessment)
    event := integration.AssessmentCompletedEvent(dto)

    // Buscar clientes que se inscreveram neste evento
    clients := getClientsSubscribedTo("assessment.completed")

    for _, client := range clients {
        sendWebhook(event, client.WebhookURL, client.WebhookSecret)
    }
}
```

**Python (recebe webhooks):**
```python
# webhook_receiver.py
from fastapi import FastAPI, Request, HTTPException
import hmac
import hashlib

app = FastAPI()

WEBHOOK_SECRET = "webhook_secret_key_123"

def verify_signature(payload: bytes, signature: str) -> bool:
    expected = hmac.new(
        WEBHOOK_SECRET.encode(),
        payload,
        hashlib.sha256
    ).hexdigest()
    return hmac.compare_digest(signature, expected)

@app.post("/webhooks/eva-mind")
async def receive_webhook(request: Request):
    # Ler payload
    body = await request.body()
    event = await request.json()

    # Verificar assinatura
    signature = event.get("signature", "")
    if not verify_signature(body, signature):
        raise HTTPException(status_code=401, detail="Invalid signature")

    # Processar evento
    event_type = event["type"]
    data = event["data"]

    if event_type == "assessment.completed":
        patient_id = data["patient_id"]
        total_score = data["total_score"]

        # Atualizar sistema do hospital
        update_hospital_system(patient_id, total_score)

    elif event_type == "crisis.suicide_risk_detected":
        patient_id = data["patient_id"]

        # Alerta urgente
        send_urgent_alert(patient_id)

    return {"status": "received"}
```

---

## 🔐 OAuth2 Flow (Conceitual)

Já que você vai implementar em Python, aqui está o fluxo:

### 1. Client Credentials Flow (Server-to-Server)

```
┌─────────────┐                               ┌─────────────┐
│  Hospital   │                               │ EVA-Mind    │
│  System     │                               │ API (Python)│
└──────┬──────┘                               └──────┬──────┘
       │                                             │
       │ POST /oauth/token                           │
       │ client_id=hosp_abc_12345                    │
       │ client_secret=secret                        │
       │ grant_type=client_credentials               │
       ├────────────────────────────────────────────>│
       │                                             │
       │                                             │ Verify credentials
       │                                             │ (check api_clients table)
       │                                             │
       │ 200 OK                                      │
       │ {                                           │
       │   "access_token": "eyJhbGc...",            │
       │   "token_type": "Bearer",                   │
       │   "expires_in": 3600                        │
       │ }                                           │
       │<────────────────────────────────────────────┤
       │                                             │
       │ GET /api/v1/patients/1                      │
       │ Authorization: Bearer eyJhbGc...            │
       ├────────────────────────────────────────────>│
       │                                             │
       │                                             │ Verify token
       │                                             │ Check scopes
       │                                             │ Check rate limit
       │                                             │ Log request
       │                                             │
       │ 200 OK                                      │
       │ { patient data }                            │
       │<────────────────────────────────────────────┤
```

**Python (FastAPI) Implementation:**
```python
from fastapi import FastAPI, Depends, HTTPException
from fastapi.security import OAuth2PasswordBearer
import jwt

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="oauth/token")

@app.post("/oauth/token")
async def login(client_id: str, client_secret: str):
    # Verificar credenciais no DB (api_clients table)
    client = db.query("SELECT * FROM api_clients WHERE client_id = %s", [client_id])

    if not client or not bcrypt.verify(client_secret, client['client_secret_hash']):
        raise HTTPException(status_code=401, detail="Invalid credentials")

    if not client['is_active'] or not client['is_approved']:
        raise HTTPException(status_code=403, detail="Client not approved")

    # Gerar JWT
    token = jwt.encode({
        "client_id": client['id'],
        "scopes": client['scopes'],
        "exp": datetime.utcnow() + timedelta(hours=1)
    }, SECRET_KEY, algorithm="HS256")

    # Salvar token no DB (api_tokens table)
    db.insert("api_tokens", {
        "client_id": client['id'],
        "access_token": token,
        "scopes": client['scopes'],
        "expires_at": datetime.utcnow() + timedelta(hours=1)
    })

    return {
        "access_token": token,
        "token_type": "Bearer",
        "expires_in": 3600
    }

@app.get("/api/v1/patients/{patient_id}")
async def get_patient(patient_id: int, token: str = Depends(oauth2_scheme)):
    # Verificar token
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
    except jwt.ExpiredSignatureError:
        raise HTTPException(status_code=401, detail="Token expired")

    # Verificar scopes
    if "read:patients" not in payload["scopes"]:
        raise HTTPException(status_code=403, detail="Insufficient permissions")

    # Rate limiting (check rate_limit_tracking table)
    check_rate_limit(payload["client_id"])

    # Log request (api_request_logs table)
    log_api_request(payload["client_id"], "GET", f"/api/v1/patients/{patient_id}", 200)

    # Buscar paciente (via Go service)
    patient = requests.get(f"http://localhost:8081/serialize/patient?id={patient_id}").json()

    return patient
```

---

## 📊 Métricas e Monitoramento

### Queries Úteis

**Top endpoints:**
```sql
SELECT * FROM v_top_api_endpoints LIMIT 10;
```

**Clientes com mais erros:**
```sql
SELECT
    client_name,
    COUNT(*) FILTER (WHERE http_status_code >= 400) AS error_count
FROM api_request_logs arl
JOIN api_clients ac ON ac.id = arl.client_id
WHERE timestamp > NOW() - INTERVAL '24 hours'
GROUP BY client_name
ORDER BY error_count DESC;
```

**Webhooks falhando:**
```sql
SELECT * FROM v_pending_webhooks
WHERE attempts >= 2;
```

**Rate limit violations:**
```sql
SELECT
    client_name,
    COUNT(*) AS violations
FROM api_request_logs arl
JOIN api_clients ac ON ac.id = arl.client_id
WHERE rate_limit_hit = TRUE
  AND timestamp > NOW() - INTERVAL '1 hour'
GROUP BY client_name;
```

---

## 🚀 Próximos Passos

### Curto Prazo (1 semana)
1. **Executar migration SQL**
2. **Compilar Go helpers** (HTTP service ou gRPC)
3. **Implementar Python API** (FastAPI recomendado)
4. **Testar OAuth2 flow**

### Médio Prazo (1 mês)
5. **Implementar webhook sender** (Go worker que processa `webhook_deliveries`)
6. **Conectar com hospital test** (FHIR bundle export)
7. **LGPD export completo** (interface para pacientes solicitarem)

### Longo Prazo (3 meses)
8. **HL7 FHIR sync em tempo real** (bidirectional)
9. **API Gateway** (Kong, Tyk) para produção
10. **Monitoring** (Prometheus, Grafana)

---

## 📝 Conclusão

O **Integration Layer** fornece tudo que você precisa para conectar EVA-Mind com o mundo externo:

✅ **Schema SQL** completo para API management
✅ **DTOs otimizados** para JSON
✅ **FHIR R4 adapters** para interoperabilidade hospitalar
✅ **Webhook system** para notificações assíncronas
✅ **Data export utilities** para LGPD/GDPR compliance

**Agora você pode:**
- Implementar API REST em Python (FastAPI/Flask)
- Chamar helpers Go via HTTP/gRPC/subprocess
- Integrar com sistemas hospitalares via FHIR
- Enviar webhooks para clientes externos
- Exportar dados com compliance LGPD

---

**Arquivo:** `SPRINT7_COMPLETED.md`
**Última Atualização:** 2026-01-24
**Versão:** 1.0
**Status:** ✅ COMPLETO