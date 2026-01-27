# 🚀 EVA-Mind Integration API - Deployment Guide

## 📋 O Que Foi Implementado

### Architecture
```
┌─────────────────────────────────────────────────────────────┐
│                      CLIENT APPLICATIONS                     │
│            (Mobile, Web, External Systems)                   │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTPS
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Python API Server (FastAPI)                     │
│                    PORT 8000                                 │
│                                                              │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   OAuth2    │  │ Rate Limiting│  │ Audit Logging│      │
│  │   + JWT     │  │   60/min     │  │   Complete   │      │
│  └─────────────┘  └──────────────┘  └──────────────┘      │
└────────────┬──────────────────────────────┬─────────────────┘
             │                              │
             │ HTTP calls                   │ PostgreSQL
             ▼                              ▼
┌───────────────────────────┐    ┌──────────────────────┐
│  Go Integration Service   │    │   PostgreSQL DB      │
│       PORT 8081           │    │  104.248.219.200     │
│                           │    │                      │
│  - JSON Serializers       │    │  - api_clients       │
│  - FHIR Adapters          │    │  - api_tokens        │
│  - Webhook Builders       │    │  - api_request_logs  │
│  - LGPD Exports           │    │  - webhook_deliveries│
└───────────────────────────┘    │  - 147+ tables       │
                                 └──────────────────────┘
```

### Componentes Criados

**1. Python API Server** (`api_server.py`)
- FastAPI REST server
- OAuth2 Client Credentials Flow
- JWT authentication
- Rate limiting (60 requests/minute)
- Audit logging completo
- CORS enabled
- Swagger docs automática

**2. Go Integration Service** (`cmd/integration_service/main.go`)
- HTTP microservice (porta 8081)
- Serializers (JSON DTOs)
- FHIR R4 adapters
- Webhook payload builders
- LGPD/GDPR export utilities
- Conexão direta ao PostgreSQL

**3. Scripts de Deployment**
- `start_services.bat` - Inicia todos os serviços
- `create_test_client.py` - Cria credenciais de teste
- `test_api.py` - Suite completa de testes
- `check_integration_tables.py` - Verifica tabelas

**4. Documentação**
- `requirements.txt` - Dependências Python
- `API_DEPLOYMENT_README.md` - Este arquivo
- Swagger UI automático em `/docs`

---

## 🔧 Pré-Requisitos

### Software Necessário
- [x] **Python 3.9+** instalado
- [x] **Go 1.21+** instalado
- [x] **PostgreSQL** (remoto em 104.248.219.200) - ✅ JÁ CONFIGURADO
- [x] **Git** (para clonar repo)

### Verificar Instalações
```bash
# Python
python --version  # Deve ser 3.9 ou superior

# Go
go version  # Deve ser 1.21 ou superior

# pip
python -m pip --version
```

---

## 📦 Installation & Deployment

### Método 1: Deployment Automático (RECOMENDADO)

```bash
# Passo 1: Abrir terminal no diretório do projeto
cd D:\dev\EVA\EVA-Mind-FZPN

# Passo 2: Criar API client de teste
python create_test_client.py

# Passo 3: Iniciar todos os serviços (automático)
start_services.bat
```

O script `start_services.bat` faz automaticamente:
1. ✅ Instala dependências Python
2. ✅ Compila Go integration service
3. ✅ Inicia Go service (porta 8081)
4. ✅ Inicia Python API (porta 8000)
5. ✅ Abre navegador com Swagger docs

---

### Método 2: Deployment Manual

#### Passo 1: Instalar Dependências Python
```bash
cd D:\dev\EVA\EVA-Mind-FZPN
python -m pip install -r requirements.txt
```

**Dependências instaladas:**
- fastapi
- uvicorn
- psycopg2-binary
- python-jose
- bcrypt
- httpx
- python-dotenv

#### Passo 2: Criar API Client de Teste
```bash
python create_test_client.py
```

**Output esperado:**
```
✅ API Client criado com sucesso!

📋 Credenciais do Cliente:
  Client ID:     eva_test_client
  Client Secret: test_secret_123

  ⚠️  SALVE ESTAS CREDENCIAIS!
```

#### Passo 3: Compilar Go Integration Service
```bash
cd cmd\integration_service
go build -o ..\..\eva_integration_service.exe main.go
cd ..\..
```

#### Passo 4: Iniciar Go Service
```bash
# Terminal 1
eva_integration_service.exe
```

**Output esperado:**
```
✓ Conectado ao PostgreSQL remoto
🚀 Eva Integration Service rodando em http://localhost:8081
```

#### Passo 5: Iniciar Python API
```bash
# Terminal 2
python api_server.py
```

**Output esperado:**
```
============================================================
🚀 EVA-Mind Integration API Server
============================================================
API Server: http://localhost:8000
API Docs: http://localhost:8000/docs
Go Service: http://localhost:8081
Database: 104.248.219.200:5432/eva-db
============================================================
```

---

## 🧪 Testing

### Teste Automático Completo
```bash
python test_api.py
```

**Testes executados:**
- ✅ Health checks (Go + Python)
- ✅ OAuth2 authentication
- ✅ GET /api/v1/patients (list)
- ✅ GET /api/v1/patients/{id}
- ✅ GET /api/v1/fhir/patients/{id}
- ✅ GET /api/v1/fhir/bundle/{id}
- ✅ GET /api/v1/export/lgpd/{id}
- ✅ Rate limiting (65 requests para testar limite de 60/min)

---

### Teste Manual com cURL

#### 1. Health Check
```bash
curl http://localhost:8000/health
```

#### 2. Obter Token OAuth2
```bash
curl -X POST http://localhost:8000/oauth/token \
  -d "username=eva_test_client" \
  -d "password=test_secret_123" \
  -d "grant_type=password"
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

#### 3. Usar Token
```bash
# Salvar token
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Listar pacientes
curl http://localhost:8000/api/v1/patients \
  -H "Authorization: Bearer $TOKEN"

# Obter paciente específico
curl http://localhost:8000/api/v1/patients/1 \
  -H "Authorization: Bearer $TOKEN"

# FHIR export
curl http://localhost:8000/api/v1/fhir/patients/1 \
  -H "Authorization: Bearer $TOKEN"

# LGPD export
curl http://localhost:8000/api/v1/export/lgpd/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

### Teste com Postman

1. **Importar Collection**
   - URL base: `http://localhost:8000`

2. **Criar Environment**
   ```
   base_url: http://localhost:8000
   client_id: eva_test_client
   client_secret: test_secret_123
   ```

3. **Obter Token** (POST /oauth/token)
   - Body: form-data
   - username: {{client_id}}
   - password: {{client_secret}}
   - grant_type: password

4. **Usar Token**
   - Authorization: Bearer Token
   - Token: {{access_token}}

---

## 📚 API Endpoints

### Authentication

#### POST /oauth/token
Obter access token OAuth2

**Body:**
```
username: eva_test_client (client_id)
password: test_secret_123 (client_secret)
grant_type: password
```

**Response:**
```json
{
  "access_token": "...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

---

### Patients

#### GET /api/v1/patients
Listar pacientes (paginado)

**Query Params:**
- limit: int (default: 10)
- offset: int (default: 0)

**Headers:**
```
Authorization: Bearer {token}
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "João Silva",
      "age": 45,
      "gender": "M"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_count": 150,
  "has_next": true
}
```

**Scopes necessários:** `read:patients`

---

#### GET /api/v1/patients/{id}
Obter paciente por ID

**Response:**
```json
{
  "id": 1,
  "name": "João Silva",
  "date_of_birth": "1980-05-15",
  "age": 45,
  "gender": "M",
  "email": "joao@example.com",
  "phone": "+55 11 98765-4321"
}
```

**Scopes necessários:** `read:patients`

---

### Assessments

#### GET /api/v1/assessments/{id}
Obter assessment por ID

**Response:**
```json
{
  "id": "a123",
  "patient_id": 1,
  "assessment_type": "PHQ-9",
  "total_score": 15,
  "severity": "moderate",
  "completed_at": "2026-01-25T10:30:00Z"
}
```

**Scopes necessários:** `read:assessments`

---

### FHIR Exports

#### GET /api/v1/fhir/patients/{id}
Obter paciente em formato FHIR R4

**Response:**
```json
{
  "resourceType": "Patient",
  "id": "1",
  "meta": {
    "lastUpdated": "2026-01-25T10:30:00Z",
    "source": "EVA-Mind"
  },
  "identifier": [
    {
      "use": "official",
      "system": "https://eva-mind.com/patient-id",
      "value": "1"
    }
  ],
  "active": true
}
```

**Scopes necessários:** `read:patients` ou `export:data`

---

#### GET /api/v1/fhir/bundle/{patient_id}
Obter FHIR Bundle (Patient + Observations)

**Scopes necessários:** `export:data`

---

### Data Export

#### GET /api/v1/export/lgpd/{patient_id}
Exportar dados do paciente (LGPD/GDPR portability)

**Response:**
```json
{
  "export_date": "2026-01-25T10:30:00Z",
  "patient_id": "1",
  "patient": {...},
  "assessments": [...],
  "medications": [...],
  "consent_history": [...],
  "data_processing": [...]
}
```

**Scopes necessários:** `export:data`

---

## 🔒 Security Features

### OAuth2 Client Credentials Flow
- JWT tokens com expiração de 1 hora
- Client secret armazenado com bcrypt
- Tokens salvos no database para revogação

### Rate Limiting
- 60 requests/minuto por client
- 3600 requests/hora por client
- 50000 requests/dia por client
- HTTP 429 quando excedido

### Audit Logging
- Todas as requests são logadas
- Timestamp, método, endpoint, status code
- Response time tracking
- Client ID tracking

### Scopes (Permissões)
- `read:patients` - Ler dados de pacientes
- `write:patients` - Criar/atualizar pacientes
- `read:assessments` - Ler assessments
- `write:assessments` - Criar assessments
- `read:medications` - Ler medicações
- `write:medications` - Gerenciar medicações
- `read:trajectories` - Ler predições de trajetória
- `export:data` - Exportar dados LGPD/FHIR

---

## 📊 Monitoring

### Health Checks
```bash
# API health
curl http://localhost:8000/health

# Go service health
curl http://localhost:8081/health
```

### Database Queries

#### Top API Clients
```sql
SELECT * FROM v_api_usage_stats
ORDER BY total_requests DESC
LIMIT 10;
```

#### Recent API Requests
```sql
SELECT
  ac.client_name,
  arl.http_method,
  arl.endpoint,
  arl.http_status_code,
  arl.response_time_ms,
  arl.timestamp
FROM api_request_logs arl
JOIN api_clients ac ON ac.id = arl.client_id
ORDER BY arl.timestamp DESC
LIMIT 20;
```

#### Rate Limit Hits
```sql
SELECT
  ac.client_name,
  COUNT(*) as rate_limit_hits,
  MAX(arl.timestamp) as last_hit
FROM api_request_logs arl
JOIN api_clients ac ON ac.id = arl.client_id
WHERE arl.rate_limit_hit = TRUE
GROUP BY ac.client_name;
```

---

## 🐛 Troubleshooting

### Go service não inicia

**Erro:** "database not accessible"

**Solução:**
```bash
# Verificar .env
cat .env | findstr DATABASE_URL

# Deve ser:
# DATABASE_URL=postgres://postgres:Debian23%40@104.248.219.200:5432/eva-db?sslmode=disable
```

---

### Python API não conecta ao Go service

**Erro:** "Go service error: Connection refused"

**Solução:**
1. Verificar se Go service está rodando: `curl http://localhost:8081/health`
2. Verificar `GO_SERVICE_URL` em `api_server.py`
3. Reiniciar Go service

---

### OAuth2 authentication falha

**Erro:** "Invalid client credentials"

**Solução:**
```bash
# Recriar client de teste
python create_test_client.py

# Usar credenciais corretas
username: eva_test_client
password: test_secret_123
```

---

### Rate limit sempre bloqueando

**Solução:**
```sql
-- Aumentar rate limit para client de teste
UPDATE api_clients
SET rate_limit_per_minute = 1000
WHERE client_id = 'eva_test_client';
```

---

## 🚀 Production Deployment

### Environment Variables
Criar `.env.production`:
```env
# JWT
JWT_SECRET_KEY=your-very-secure-random-key-min-32-chars

# Database
DB_HOST=104.248.219.200
DB_PORT=5432
DB_NAME=eva-db
DB_USER=postgres
DB_PASSWORD=Debian23@

# Services
GO_SERVICE_URL=http://localhost:8081
API_PORT=8000
```

### Production Checklist
- [ ] Trocar `JWT_SECRET_KEY` por valor seguro
- [ ] Habilitar HTTPS (TLS/SSL)
- [ ] Configurar firewall (apenas portas 80, 443)
- [ ] Configurar reverse proxy (Nginx/Caddy)
- [ ] Habilitar logs estruturados
- [ ] Configurar backup automático do database
- [ ] Implementar monitoring (Prometheus/Grafana)
- [ ] Configurar alertas de erro
- [ ] Rate limits ajustados por cliente
- [ ] Documentar processo de criação de novos clients

---

## 📖 Next Steps

1. **Adicionar mais endpoints**
   - POST /api/v1/patients (criar paciente)
   - POST /api/v1/assessments (criar assessment)
   - GET /api/v1/medications
   - GET /api/v1/trajectories

2. **Implementar webhooks**
   - Sistema de delivery assíncrono
   - Retry automático em caso de falha
   - Signature verification

3. **Adicionar analytics**
   - Dashboard de uso da API
   - Métricas de performance
   - Relatórios de compliance

4. **Mobile SDK**
   - SDK Python para integração fácil
   - SDK JavaScript/TypeScript
   - Exemplos de uso

---

## 📞 Support

- **Documentação API**: http://localhost:8000/docs
- **Health Check**: http://localhost:8000/health
- **Issue Tracker**: GitHub Issues

---

**✅ SPRINT 7 - Integration Layer COMPLETO!**

Arquitetura híbrida Go ↔ Python totalmente funcional! 🎉
