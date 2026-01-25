# 🔌 Integration Quick Reference - Go ↔️ Python

## 📊 Comparação dos Métodos

| Método | Performance | Complexidade | Uso Recomendado |
|--------|-------------|--------------|-----------------|
| **HTTP Microservice** | ⭐⭐⭐⭐ | ⭐⭐ | Produção, múltiplos clientes |
| **gRPC** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | Alta performance, streaming |
| **Subprocess** | ⭐⭐ | ⭐ | Desenvolvimento, scripts simples |
| **Shared Library** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | Integração C/Python nativa |

---

## 🚀 Método 1: HTTP Microservice (RECOMENDADO)

### Go Service (`cmd/integration_service/main.go`)

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "eva-mind-fzpn/internal/integration"
)

func main() {
    r := mux.NewRouter()

    // Serialize endpoints
    r.HandleFunc("/serialize/patient/{id}", serializePatientHandler).Methods("GET")
    r.HandleFunc("/serialize/assessment/{id}", serializeAssessmentHandler).Methods("GET")

    // FHIR endpoints
    r.HandleFunc("/fhir/patient", convertPatientToFHIRHandler).Methods("POST")
    r.HandleFunc("/fhir/bundle", createFHIRBundleHandler).Methods("POST")

    // Webhook endpoints
    r.HandleFunc("/webhook/build", buildWebhookHandler).Methods("POST")
    r.HandleFunc("/webhook/sign", signWebhookHandler).Methods("POST")

    // Export endpoints
    r.HandleFunc("/export/lgpd/{patient_id}", exportLGPDHandler).Methods("GET")
    r.HandleFunc("/export/research", exportResearchHandler).Methods("POST")

    log.Println("Integration service running on :8081")
    log.Fatal(http.ListenAndServe(":8081", r))
}

func serializePatientHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    patientID, _ := strconv.ParseInt(vars["id"], 10, 64)

    // TODO: Buscar patient do database
    patient := &integration.PatientDTO{
        ID:          patientID,
        Name:        "João Silva",
        DateOfBirth: "1980-05-15",
        Age:         45,
        Gender:      "M",
    }

    json, err := integration.ToJSON(patient)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(json))
}

func convertPatientToFHIRHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        PatientID int64 `json:"patient_id"`
    }

    json.NewDecoder(r.Body).Decode(&req)

    // TODO: Buscar patient do database
    patient := &integration.PatientDTO{
        ID:          req.PatientID,
        Name:        "João Silva",
        DateOfBirth: "1980-05-15",
        Age:         45,
        Gender:      "M",
    }

    fhirPatient := integration.PatientDTOToFHIR(patient)

    jsonStr, err := integration.ToFHIRJSON(fhirPatient)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/fhir+json")
    w.Write([]byte(jsonStr))
}

func buildWebhookHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        EventType string                 `json:"event_type"`
        Data      map[string]interface{} `json:"data"`
        Secret    string                 `json:"secret"`
    }

    json.NewDecoder(r.Body).Decode(&req)

    event := integration.CustomEvent(req.EventType, req.Data)
    event.AddSignature(req.Secret)

    jsonStr, _ := event.ToJSON()

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(jsonStr))
}

func exportLGPDHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    patientID, _ := strconv.ParseInt(vars["patient_id"], 10, 64)

    export := integration.NewLGPDPortabilityExport(patientID)

    // TODO: Popular com dados do database

    jsonStr, _ := integration.ToJSON(export)

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=lgpd_export_%d.json", patientID))
    w.Write([]byte(jsonStr))
}
```

### Compilar e Rodar
```bash
cd D:\dev\EVA\EVA-Mind-FZPN
go get github.com/gorilla/mux
go build -o eva_integration_service.exe cmd/integration_service/main.go
./eva_integration_service.exe
```

### Python Client
```python
import requests

BASE_URL = "http://localhost:8081"

# 1. Serializar Patient
response = requests.get(f"{BASE_URL}/serialize/patient/1")
patient = response.json()
print(patient)

# 2. Converter para FHIR
response = requests.post(f"{BASE_URL}/fhir/patient", json={"patient_id": 1})
fhir_patient = response.json()
print(fhir_patient)

# 3. Criar webhook assinado
response = requests.post(f"{BASE_URL}/webhook/build", json={
    "event_type": "patient.created",
    "data": {"patient_id": 1, "name": "João Silva"},
    "secret": "my_webhook_secret"
})
webhook_payload = response.json()
print(webhook_payload['signature'])  # HMAC-SHA256

# 4. Exportar dados LGPD
response = requests.get(f"{BASE_URL}/export/lgpd/1")
lgpd_export = response.json()
print(lgpd_export)
```

---

## ⚡ Método 2: gRPC (MELHOR PERFORMANCE)

### Proto Definition (`proto/integration.proto`)
```protobuf
syntax = "proto3";

package integration;

option go_package = "eva-mind-fzpn/proto";

service IntegrationService {
    rpc SerializePatient(PatientRequest) returns (PatientResponse);
    rpc ConvertToFHIR(FHIRRequest) returns (FHIRResponse);
    rpc BuildWebhook(WebhookRequest) returns (WebhookResponse);
    rpc ExportLGPD(LGPDRequest) returns (LGPDResponse);
}

message PatientRequest {
    int64 patient_id = 1;
}

message PatientResponse {
    string json_data = 1;
}

message FHIRRequest {
    int64 patient_id = 1;
}

message FHIRResponse {
    string fhir_json = 1;
}

message WebhookRequest {
    string event_type = 1;
    string data_json = 2;
    string secret = 3;
}

message WebhookResponse {
    string webhook_json = 1;
}

message LGPDRequest {
    int64 patient_id = 1;
}

message LGPDResponse {
    string export_json = 1;
}
```

### Compilar Proto
```bash
# Instalar protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Compilar
protoc --go_out=. --go-grpc_out=. proto/integration.proto
```

### Go gRPC Server
```go
package main

import (
    "context"
    "log"
    "net"
    "google.golang.org/grpc"
    pb "eva-mind-fzpn/proto"
    "eva-mind-fzpn/internal/integration"
)

type server struct {
    pb.UnimplementedIntegrationServiceServer
}

func (s *server) SerializePatient(ctx context.Context, req *pb.PatientRequest) (*pb.PatientResponse, error) {
    patient := &integration.PatientDTO{
        ID:   req.PatientId,
        Name: "João Silva",
        Age:  45,
    }

    json, _ := integration.ToJSON(patient)
    return &pb.PatientResponse{JsonData: json}, nil
}

func main() {
    lis, _ := net.Listen("tcp", ":50051")
    s := grpc.NewServer()
    pb.RegisterIntegrationServiceServer(s, &server{})

    log.Println("gRPC server running on :50051")
    s.Serve(lis)
}
```

### Python gRPC Client
```python
import grpc
import integration_pb2
import integration_pb2_grpc
import json

# Conectar
channel = grpc.insecure_channel('localhost:50051')
stub = integration_pb2_grpc.IntegrationServiceStub(channel)

# Chamar RPC
response = stub.SerializePatient(integration_pb2.PatientRequest(patient_id=1))
patient = json.loads(response.json_data)
print(patient)
```

---

## 🛠️ Método 3: Subprocess (MAIS SIMPLES)

### Go CLI (`cmd/helpers/main.go`)
```go
package main

import (
    "fmt"
    "os"
    "strconv"
    "eva-mind-fzpn/internal/integration"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: eva_helpers <command> <args...>")
        os.Exit(1)
    }

    command := os.Args[1]

    switch command {
    case "serialize":
        resource := os.Args[2]
        id, _ := strconv.ParseInt(os.Args[3], 10, 64)

        if resource == "patient" {
            patient := &integration.PatientDTO{ID: id, Name: "João Silva"}
            json, _ := integration.ToJSON(patient)
            fmt.Println(json)
        }

    case "fhir":
        id, _ := strconv.ParseInt(os.Args[2], 10, 64)
        patient := &integration.PatientDTO{ID: id}
        fhir := integration.PatientDTOToFHIR(patient)
        json, _ := integration.ToFHIRJSON(fhir)
        fmt.Println(json)

    case "webhook":
        eventType := os.Args[2]
        event := integration.CustomEvent(eventType, map[string]interface{}{})
        json, _ := event.ToJSON()
        fmt.Println(json)

    default:
        fmt.Println("Unknown command:", command)
    }
}
```

### Compilar
```bash
go build -o eva_helpers.exe cmd/helpers/main.go
```

### Python Subprocess
```python
import subprocess
import json

def call_go_helper(command, *args):
    result = subprocess.run(
        ["./eva_helpers.exe", command] + list(args),
        capture_output=True,
        text=True
    )
    return json.loads(result.stdout)

# Usar
patient = call_go_helper("serialize", "patient", "1")
print(patient)

fhir = call_go_helper("fhir", "1")
print(fhir)
```

---

## 📦 Método 4: Shared Library (C-COMPATIBLE)

### Go Shared Library
```go
package main

import "C"
import (
    "eva-mind-fzpn/internal/integration"
)

//export SerializePatient
func SerializePatient(patientID int64) *C.char {
    patient := &integration.PatientDTO{ID: patientID, Name: "João Silva"}
    json, _ := integration.ToJSON(patient)
    return C.CString(json)
}

//export ConvertToFHIR
func ConvertToFHIR(patientID int64) *C.char {
    patient := &integration.PatientDTO{ID: patientID}
    fhir := integration.PatientDTOToFHIR(patient)
    json, _ := integration.ToFHIRJSON(fhir)
    return C.CString(json)
}

func main() {}
```

### Compilar
```bash
go build -buildmode=c-shared -o libeva_integration.so
```

### Python ctypes
```python
import ctypes
import json

# Carregar library
lib = ctypes.CDLL('./libeva_integration.so')

# Configurar tipos
lib.SerializePatient.argtypes = [ctypes.c_int64]
lib.SerializePatient.restype = ctypes.c_char_p

# Chamar função
json_str = lib.SerializePatient(1).decode('utf-8')
patient = json.loads(json_str)
print(patient)
```

---

## 🎯 Escolher Método

**Use HTTP Microservice se:**
- ✅ Precisa de produção estável
- ✅ Múltiplos clientes (mobile, web, outros serviços)
- ✅ RESTful é mais familiar
- ✅ Debugging fácil (curl, Postman)

**Use gRPC se:**
- ✅ Performance crítica
- ✅ Streaming bidirecional necessário
- ✅ Tipagem forte entre serviços
- ✅ Equipe experiente com gRPC

**Use Subprocess se:**
- ✅ Desenvolvimento/testes rápidos
- ✅ Scripts simples
- ✅ Não precisa de alta performance
- ✅ Menos complexidade

**Use Shared Library se:**
- ✅ Performance máxima
- ✅ Integração C/Python nativa
- ✅ Sem overhead de rede
- ⚠️ Equipe avançada

---

## 📚 Recursos Disponíveis

### DTOs (JSON Serialization)
```go
integration.PatientDTO
integration.AssessmentDTO
integration.VoiceAnalysisDTO
integration.TrajectoryDTO
integration.PersonaDTO
integration.LastWishesDTO
integration.QualityOfLifeDTO
integration.PainLogDTO
integration.LegacyMessageDTO
integration.ResearchStudyDTO
```

### FHIR Adapters
```go
integration.PatientDTOToFHIR(patient)
integration.PHQ9ToFHIR(assessment)
integration.CreatePatientBundle(patient, observations)
integration.ToFHIRJSON(resource)
```

### Webhook Builders
```go
integration.PatientCreatedEvent(patient)
integration.AssessmentCompletedEvent(assessment)
integration.CrisisDetectedEvent(patientID, type, severity, details)
integration.SuicideRiskDetectedEvent(patientID, assessmentID, score)
integration.PersonaTransitionEvent(patientID, from, to, reason)
integration.SignWebhookPayload(payload, secret)
```

### Export Utilities
```go
integration.NewLGPDPortabilityExport(patientID)
integration.AnonymizePatientID(patientID)
integration.AnonymizeName(name)
integration.AnonymizeEmail(email)
integration.GenerateCSVHeader(columns, delimiter)
integration.ExportPatientAsFHIRBundle(patient, assessments)
```

---

## 🔧 Troubleshooting

### Port já em uso
```bash
# Windows
netstat -ano | findstr "8081"
taskkill /PID <PID> /F
```

### Go module errors
```bash
go mod tidy
go mod download
```

### CORS errors (HTTP)
```go
// Adicionar middleware CORS
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
```

---

**Integration Quick Reference - Pronto para Uso! 🔌**
