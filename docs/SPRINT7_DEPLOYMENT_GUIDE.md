# 🚀 SPRINT 7 - Integration Layer Deployment Guide

## 📋 O Que Foi Criado

### 1. Database Schema (SQL)
**Arquivo**: `migrations/010_integration_layer.sql`
- 8 tabelas para gerenciamento de API
- 4 views para analytics
- 3 triggers automáticos
- Sistema completo de OAuth2, audit logs, webhooks, FHIR mappings e data exports

### 2. Go Integration Helpers
**Arquivos criados**:
```
internal/integration/
├── serializers.go      (~400 linhas) - 30+ DTOs para JSON
├── fhir_adapter.go     (~600 linhas) - FHIR R4 conversions
├── webhooks.go         (~400 linhas) - 15+ tipos de eventos
└── export.go           (~600 linhas) - LGPD/GDPR exports
```

### 3. Documentação Completa
- `SPRINT7_COMPLETED.md` - Guia técnico completo
- `SPRINT7_DEPLOYMENT_GUIDE.md` - Este guia de deployment

---

## 🔧 Pré-Requisitos de Deployment

### Verificar PostgreSQL

**Opção 1: PostgreSQL Local**
```bash
# Verificar se está instalado
psql --version

# Verificar se está rodando (Windows)
sc query postgresql-x64-16

# Se não estiver rodando, instalar:
# https://www.postgresql.org/download/windows/
```

**Opção 2: PostgreSQL Remoto**
```bash
# Seu .env atual aponta para: 127.0.0.1:5432
# Mas você tem conexões ativas com: 104.248.219.200:5432

# Se quiser usar o servidor remoto, atualize .env:
DATABASE_URL=postgres://postgres:Debian23@104.248.219.200:5432/eva-db?sslmode=disable
```

### Verificar Go Build Environment
```bash
cd D:\dev\EVA\EVA-Mind-FZPN
go version  # Deve ser 1.21+
go mod download
go build -o eva-mind-fzpn.exe
```

---

## 📦 Passos de Deployment

### Passo 1: Escolher Configuração de Database

**Se usar PostgreSQL LOCAL:**
```bash
# 1. Instalar PostgreSQL 16
# 2. Criar database
psql -U postgres -c "CREATE DATABASE \"eva-db\";"

# 3. Habilitar pgvector extension
psql -U postgres -d eva-db -c "CREATE EXTENSION IF NOT EXISTS vector;"

# 4. Executar migração SPRINT 7
psql -U postgres -d eva-db -f migrations/010_integration_layer.sql

# 5. Verificar tabelas criadas
psql -U postgres -d eva-db -c "\dt api_*"
```

**Se usar PostgreSQL REMOTO:**
```bash
# 1. Atualizar .env
DATABASE_URL=postgres://postgres:Debian23@104.248.219.200:5432/eva-db?sslmode=disable

# 2. Conectar ao remoto
set PGPASSWORD=Debian23
psql -U postgres -h 104.248.219.200 -d eva-db -f migrations/010_integration_layer.sql

# 3. Verificar tabelas criadas
psql -U postgres -h 104.248.219.200 -d eva-db -c "SELECT tablename FROM pg_tables WHERE tablename LIKE 'api_%';"
```

---

### Passo 2: Compilar Go Helpers

Você tem 4 opções de integração. Escolha a que melhor se adequa:

#### Opção A: HTTP Microservice (Recomendado)
```bash
cd D:\dev\EVA\EVA-Mind-FZPN

# Criar arquivo main para microservice
# (Ver exemplo em SPRINT7_COMPLETED.md - Seção "Integration Method 3")

go build -o eva_integration_service.exe cmd/integration_service/main.go

# Rodar serviço na porta 8081
./eva_integration_service.exe
```

**Python FastAPI chama assim:**
```python
import requests

# Serializar patient
response = requests.get("http://localhost:8081/serialize/patient/1")
patient_json = response.json()

# Converter para FHIR
response = requests.post("http://localhost:8081/fhir/patient", json={"patient_id": 1})
fhir_bundle = response.json()
```

#### Opção B: gRPC Service (Melhor Performance)
```bash
# 1. Instalar protobuf compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 2. Criar .proto file
# (Ver exemplo em SPRINT7_COMPLETED.md - Seção "Integration Method 2")

# 3. Compilar gRPC service
go build -o eva_grpc_service.exe cmd/grpc_service/main.go

# 4. Rodar serviço
./eva_grpc_service.exe
```

**Python chama via grpcio:**
```python
import grpc
import integration_pb2
import integration_pb2_grpc

channel = grpc.insecure_channel('localhost:50051')
stub = integration_pb2_grpc.IntegrationServiceStub(channel)

response = stub.SerializePatient(integration_pb2.PatientRequest(patient_id=1))
```

#### Opção C: Subprocess (Mais Simples)
```bash
# Criar CLI tool
go build -o eva_helpers.exe cmd/helpers/main.go
```

**Python chama assim:**
```python
import subprocess
import json

result = subprocess.run(
    ["./eva_helpers.exe", "serialize", "patient", "1"],
    capture_output=True,
    text=True
)
patient_json = json.loads(result.stdout)
```

#### Opção D: Shared Library (C-compatible)
```bash
# Compilar como shared library
go build -buildmode=c-shared -o libeva_integration.so

# Python chama via ctypes
```

---

### Passo 3: Implementar Python API (FastAPI)

Criar `api_server.py`:

```python
from fastapi import FastAPI, Depends, HTTPException
from fastapi.security import OAuth2PasswordBearer
import psycopg2
import jwt
import bcrypt
import requests  # Para chamar Go microservice

app = FastAPI(title="EVA-Mind Integration API")

# Configurar conexão PostgreSQL
DB_CONFIG = {
    "host": "127.0.0.1",  # ou 104.248.219.200
    "port": 5432,
    "database": "eva-db",
    "user": "postgres",
    "password": "Debian23@"
}

# URL do Go microservice
GO_SERVICE_URL = "http://localhost:8081"

# OAuth2 Token Endpoint
@app.post("/oauth/token")
async def login(client_id: str, client_secret: str):
    conn = psycopg2.connect(**DB_CONFIG)
    cur = conn.cursor()

    cur.execute("""
        SELECT id, client_secret_hash, scopes, is_active, is_approved
        FROM api_clients
        WHERE client_id = %s
    """, (client_id,))

    client = cur.fetchone()
    if not client or not bcrypt.checkpw(client_secret.encode(), client[1].encode()):
        raise HTTPException(status_code=401, detail="Invalid credentials")

    if not client[3] or not client[4]:  # is_active, is_approved
        raise HTTPException(status_code=403, detail="Client not active or not approved")

    # Gerar JWT token
    token = jwt.encode({
        "client_id": str(client[0]),
        "scopes": client[2],
        "exp": datetime.utcnow() + timedelta(hours=1)
    }, SECRET_KEY)

    # Salvar token no database
    cur.execute("""
        INSERT INTO api_tokens (id, client_id, access_token, scopes, expires_at)
        VALUES (gen_random_uuid(), %s, %s, %s, %s)
    """, (client[0], token, client[2], datetime.utcnow() + timedelta(hours=1)))

    conn.commit()
    cur.close()
    conn.close()

    return {"access_token": token, "token_type": "Bearer"}

# Protected Endpoint: Get Patient
@app.get("/api/v1/patients/{patient_id}")
async def get_patient(patient_id: int, token: str = Depends(oauth2_scheme)):
    # Validar token
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
        client_id = payload["client_id"]
    except:
        raise HTTPException(status_code=401)

    # Rate limiting check
    conn = psycopg2.connect(**DB_CONFIG)
    cur = conn.cursor()

    cur.execute("""
        SELECT COUNT(*) FROM api_request_logs
        WHERE client_id = %s
        AND timestamp > NOW() - INTERVAL '1 minute'
    """, (client_id,))

    request_count = cur.fetchone()[0]

    cur.execute("SELECT rate_limit_per_minute FROM api_clients WHERE id = %s", (client_id,))
    rate_limit = cur.fetchone()[0]

    if request_count >= rate_limit:
        raise HTTPException(status_code=429, detail="Rate limit exceeded")

    # Chamar Go microservice para serializar patient
    response = requests.get(f"{GO_SERVICE_URL}/serialize/patient/{patient_id}")

    # Log request
    cur.execute("""
        INSERT INTO api_request_logs (id, client_id, http_method, endpoint, http_status_code, response_time_ms)
        VALUES (DEFAULT, %s, 'GET', %s, %s, %s)
    """, (client_id, f"/api/v1/patients/{patient_id}", response.status_code, response.elapsed.total_seconds() * 1000))

    conn.commit()
    cur.close()
    conn.close()

    return response.json()

# FHIR Export Endpoint
@app.get("/api/v1/fhir/patients/{patient_id}")
async def get_patient_fhir(patient_id: int, token: str = Depends(oauth2_scheme)):
    # Chamar Go microservice
    response = requests.post(f"{GO_SERVICE_URL}/fhir/patient", json={"patient_id": patient_id})
    return response.json()

# Webhook Sender (Exemplo)
@app.post("/internal/send_webhook")
async def send_webhook(event_type: str, data: dict):
    # Buscar clientes com webhooks configurados
    conn = psycopg2.connect(**DB_CONFIG)
    cur = conn.cursor()

    cur.execute("SELECT id, webhook_url, webhook_secret FROM api_clients WHERE webhook_url IS NOT NULL AND is_active = TRUE")
    clients = cur.fetchall()

    for client in clients:
        client_id, webhook_url, webhook_secret = client

        # Chamar Go helper para criar payload assinado
        response = requests.post(f"{GO_SERVICE_URL}/webhook/build", json={
            "event_type": event_type,
            "data": data,
            "secret": webhook_secret
        })

        webhook_payload = response.json()

        # Enviar webhook
        try:
            webhook_response = requests.post(webhook_url, json=webhook_payload, timeout=5)

            # Registrar delivery
            cur.execute("""
                INSERT INTO webhook_deliveries (id, client_id, event_type, event_data, status, http_status_code, response_time_ms)
                VALUES (gen_random_uuid(), %s, %s, %s, 'delivered', %s, %s)
            """, (client_id, event_type, data, webhook_response.status_code, webhook_response.elapsed.total_seconds() * 1000))
        except Exception as e:
            # Falha no webhook
            cur.execute("""
                INSERT INTO webhook_deliveries (id, client_id, event_type, event_data, status, attempts, error_message)
                VALUES (gen_random_uuid(), %s, %s, %s, 'failed', 1, %s)
            """, (client_id, event_type, data, str(e)))

    conn.commit()
    cur.close()
    conn.close()

    return {"sent_to": len(clients)}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
```

**Rodar API:**
```bash
pip install fastapi uvicorn psycopg2-binary pyjwt bcrypt requests
python api_server.py
```

---

### Passo 4: Testar Integração End-to-End

```bash
# Terminal 1: Rodar Go microservice
cd D:\dev\EVA\EVA-Mind-FZPN
./eva_integration_service.exe

# Terminal 2: Rodar Python API
python api_server.py

# Terminal 3: Testar endpoints
# 1. Criar API client
psql -U postgres -d eva-db -c "
INSERT INTO api_clients (id, client_name, client_id, client_secret_hash, scopes, rate_limit_per_minute)
VALUES (
    gen_random_uuid(),
    'Test Client',
    'test_client',
    '\$2b\$12\$abcd...',  -- bcrypt hash de 'test_secret'
    ARRAY['read:patients', 'write:assessments'],
    60
);
"

# 2. Obter token
curl -X POST http://localhost:8000/oauth/token \
  -d "client_id=test_client" \
  -d "client_secret=test_secret"

# 3. Chamar endpoint protegido
curl http://localhost:8000/api/v1/patients/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# 4. Exportar FHIR
curl http://localhost:8000/api/v1/fhir/patients/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 📊 Verificar Deployment

### Verificar Tabelas Criadas
```sql
-- Conectar ao database
psql -U postgres -d eva-db

-- Listar tabelas de integração
\dt api_*
\dt webhook*
\dt fhir*
\dt data_export*

-- Verificar views
\dv v_api*
```

### Verificar Go Services Rodando
```bash
# Windows
netstat -an | findstr "8081"  # Go microservice
netstat -an | findstr "8000"  # Python API
```

### Verificar Logs
```sql
-- Top API clients por requests
SELECT * FROM v_api_usage_stats ORDER BY total_requests DESC LIMIT 10;

-- Webhooks pendentes
SELECT * FROM v_webhook_delivery_stats WHERE failed_deliveries > 0;

-- Rate limits atingidos
SELECT * FROM api_request_logs WHERE rate_limit_hit = TRUE ORDER BY timestamp DESC LIMIT 10;
```

---

## 🔒 Security Checklist

- [ ] PostgreSQL rodando com senha forte
- [ ] JWT secret key configurado (não usar padrão)
- [ ] Rate limits configurados para cada client
- [ ] Webhooks usando HMAC-SHA256 signature
- [ ] HTTPS habilitado em produção (não HTTP)
- [ ] Logs de audit habilitados
- [ ] Backup automático do database
- [ ] Firewall configurado (só portas necessárias abertas)

---

## 🎯 Próximos Passos

1. **Decidir configuração do database** (local ou remoto)
2. **Executar migração SQL** (criar tabelas)
3. **Escolher método de integração** (HTTP microservice recomendado)
4. **Compilar Go helpers**
5. **Implementar Python API** (FastAPI)
6. **Testar end-to-end**
7. **Deploy em produção**

---

## 📚 Referências

- **SPRINT7_COMPLETED.md** - Documentação técnica completa
- **migrations/010_integration_layer.sql** - Schema SQL
- **internal/integration/** - Código Go dos helpers
- **HL7 FHIR R4**: https://www.hl7.org/fhir/R4/
- **LOINC Codes**: https://loinc.org/
- **LGPD**: Lei Geral de Proteção de Dados

---

## ❓ Troubleshooting

### Erro: "database does not exist"
```bash
psql -U postgres -c "CREATE DATABASE \"eva-db\";"
```

### Erro: "extension vector does not exist"
```bash
psql -U postgres -d eva-db -c "CREATE EXTENSION vector;"
```

### Erro: "relation api_clients already exists"
```bash
# Tabelas já foram criadas antes. Para recriar:
psql -U postgres -d eva-db -c "DROP TABLE IF EXISTS api_clients CASCADE;"
# Depois executar migration novamente
```

### Go service não responde
```bash
# Verificar se porta está em uso
netstat -an | findstr "8081"

# Verificar logs
./eva_integration_service.exe > service.log 2>&1
```

---

**SPRINT 7 Integration Layer - Pronto para Deploy! 🚀**
