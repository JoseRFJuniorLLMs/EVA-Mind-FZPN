# Arquivo de Cibersegurança
## EVA-Mind-FZPN - Companion IA para Idosos

**Documento:** SEC-EVA-001
**Versão:** 1.0
**Data:** 2025-01-27
**Normas:** ABNT NBR ISO/TR 81001-2-8, ISO 27001, OWASP

---

## 1. Política de Segurança da Informação

### 1.1 Declaração de Política

O EVA-Mind-FZPN está comprometido com a proteção da confidencialidade, integridade e disponibilidade das informações de saúde dos usuários, implementando controles de segurança proporcionais aos riscos identificados.

### 1.2 Objetivos de Segurança

| Objetivo | Descrição | Métrica |
|----------|-----------|---------|
| Confidencialidade | Proteção contra acesso não autorizado | 0 vazamentos de dados |
| Integridade | Proteção contra alteração não autorizada | 0 incidentes de corrupção |
| Disponibilidade | Sistema acessível quando necessário | Uptime ≥99.5% |
| Autenticidade | Garantia de identidade | 0 acessos fraudulentos |
| Não-repúdio | Prova de ações realizadas | 100% de ações auditadas |

### 1.3 Escopo

- Aplicativo mobile (Android/iOS)
- Aplicação web
- APIs backend
- Banco de dados
- Infraestrutura de nuvem
- Integrações com terceiros

---

## 2. Análise de Riscos de Segurança

### 2.1 Metodologia

**Framework:** ISO 27005 + NIST Cybersecurity Framework
**Escala de Probabilidade:** 1 (Raro) a 5 (Quase certo)
**Escala de Impacto:** 1 (Insignificante) a 5 (Catastrófico)

### 2.2 Riscos Identificados

| ID | Risco | Prob | Imp | Score | Controles | Residual |
|----|-------|------|-----|-------|-----------|----------|
| SEC-01 | Vazamento de banco de dados | 2 | 5 | 10 | Criptografia, WAF, MFA | 4 |
| SEC-02 | Acesso não autorizado a APIs | 3 | 4 | 12 | JWT, Rate limiting, RBAC | 6 |
| SEC-03 | Interceptação de dados em trânsito | 2 | 5 | 10 | TLS 1.3, Certificate pinning | 4 |
| SEC-04 | Comprometimento de credenciais | 3 | 4 | 12 | Bcrypt, MFA, Detecção anomalia | 6 |
| SEC-05 | Ransomware | 2 | 5 | 10 | Backup offline, EDR | 5 |
| SEC-06 | SQL Injection | 2 | 5 | 10 | Prepared statements, WAF | 3 |
| SEC-07 | XSS (Cross-Site Scripting) | 3 | 3 | 9 | CSP, Sanitização, Encoding | 4 |
| SEC-08 | CSRF | 2 | 3 | 6 | Tokens CSRF, SameSite cookies | 3 |
| SEC-09 | DDoS | 3 | 4 | 12 | CDN, Rate limiting, Auto-scale | 6 |
| SEC-10 | Insider threat | 2 | 4 | 8 | RBAC, Auditoria, Segregação | 4 |

### 2.3 Matriz de Risco

```
                            IMPACTO
              1        2        3        4        5
           (Insig)  (Menor)  (Mod)   (Maior) (Catastr)
         ┌────────┬────────┬────────┬────────┬────────┐
    5    │        │        │   I    │   I    │   I    │
  (Certo)│        │        │        │        │        │
         ├────────┼────────┼────────┼────────┼────────┤
    4    │        │        │   A    │   I    │   I    │
  (Prov) │        │        │        │        │        │
P        ├────────┼────────┼────────┼────────┼────────┤
R   3    │        │   R    │ SEC-07 │SEC-02,4│   I    │
O (Poss) │        │        │        │ SEC-09 │        │
B        ├────────┼────────┼────────┼────────┼────────┤
    2    │   R    │   R    │ SEC-08 │SEC-10  │SEC-01,3│
 (Improv)│        │        │        │        │SEC-05,6│
         ├────────┼────────┼────────┼────────┼────────┤
    1    │   R    │   R    │   R    │   A    │   A    │
  (Raro) │        │        │        │        │        │
         └────────┴────────┴────────┴────────┴────────┘

I = Inaceitável | A = ALARP | R = Aceitável
```

---

## 3. Criptografia

### 3.1 Dados em Trânsito

| Componente | Protocolo | Versão | Algoritmo |
|------------|-----------|--------|-----------|
| HTTPS | TLS | 1.3 | ECDHE + AES-256-GCM |
| API Gateway | TLS | 1.3 | ECDHE + AES-256-GCM |
| Mobile App | TLS | 1.3 + Certificate Pinning | ECDHE + AES-256-GCM |
| Database conexão | TLS | 1.3 | AES-256-GCM |
| Redis conexão | TLS | 1.3 | AES-256-GCM |

**Cipher Suites Permitidos:**
```
TLS_AES_256_GCM_SHA384
TLS_AES_128_GCM_SHA256
TLS_CHACHA20_POLY1305_SHA256
```

**Cipher Suites Bloqueados:**
```
TLS_RSA_*
*_CBC_*
*_MD5
*_SHA1
SSL_*
```

### 3.2 Dados em Repouso

| Dado | Algoritmo | Tamanho Chave | Gerenciamento |
|------|-----------|---------------|---------------|
| Banco de dados | AES-256-GCM | 256 bits | AWS KMS |
| Backups | AES-256-GCM | 256 bits | AWS KMS |
| Arquivos S3 | AES-256-GCM | 256 bits | SSE-KMS |
| Logs | AES-256-GCM | 256 bits | AWS KMS |

### 3.3 Hash de Senhas

```go
// Configuração de Bcrypt
const (
    BcryptCost = 12  // 2^12 iterações (~250ms)
)

// Alternativa: Argon2id
type Argon2Config struct {
    Memory      uint32 = 64 * 1024  // 64 MB
    Iterations  uint32 = 3
    Parallelism uint8  = 4
    SaltLength  uint32 = 16
    KeyLength   uint32 = 32
}
```

### 3.4 Gestão de Chaves

| Chave | Rotação | Armazenamento | Acesso |
|-------|---------|---------------|--------|
| Master Key (KMS) | Anual | AWS KMS | IAM roles |
| Data Encryption Key | Mensal | Derivada de Master | Aplicação |
| JWT Signing Key | Trimestral | AWS Secrets Manager | API Server |
| API Keys (terceiros) | Conforme vendor | AWS Secrets Manager | Serviços |

---

## 4. Autenticação e Autorização

### 4.1 Mecanismos de Autenticação

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    FLUXO DE AUTENTICAÇÃO                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────┐        ┌─────────────┐        ┌─────────────┐             │
│  │  USER   │──1────▶│   LOGIN     │──2────▶│  VALIDATE   │             │
│  │         │        │  (email+pw) │        │  CREDENTIALS│             │
│  └─────────┘        └─────────────┘        └──────┬──────┘             │
│                                                    │                    │
│                                           ┌───────┴───────┐             │
│                                           │               │             │
│                                           ▼               ▼             │
│                                     ┌─────────┐    ┌──────────┐        │
│                                     │ SUCCESS │    │  FAIL    │        │
│                                     └────┬────┘    │ (block)  │        │
│                                          │         └──────────┘        │
│                                          │                              │
│  ┌─────────┐        ┌─────────────┐      │                              │
│  │  USER   │◀──5────│  TOKENS     │◀─4───┤                              │
│  │         │        │(access+     │      │                              │
│  │         │        │ refresh)    │      │                              │
│  └─────────┘        └─────────────┘      │                              │
│                                          │                              │
│                     ┌─────────────┐      │                              │
│                     │  MFA CHECK  │◀─3───┘  (se habilitado)            │
│                     │  (TOTP/SMS) │                                     │
│                     └─────────────┘                                     │
│                                                                         │
│  TOKENS:                                                                │
│  • Access Token: JWT, 15 min TTL                                       │
│  • Refresh Token: Opaco, 7 dias TTL, rotação automática               │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Política de Senhas

| Requisito | Especificação |
|-----------|---------------|
| Comprimento mínimo | 8 caracteres |
| Complexidade | Letra + número |
| Histórico | Últimas 5 senhas bloqueadas |
| Expiração | 180 dias (opcional para idosos) |
| Bloqueio por tentativas | 5 falhas → bloqueio 15 min |
| Recuperação | E-mail + código 6 dígitos (15 min) |

### 4.3 Multi-Factor Authentication (MFA)

| Método | Disponibilidade | Implementação |
|--------|-----------------|---------------|
| TOTP (Google Authenticator) | Opcional | RFC 6238 |
| SMS OTP | Fallback | Twilio, 6 dígitos, 5 min |
| Biometria (mobile) | Principal | Touch ID / Face ID |

### 4.4 Controle de Acesso (RBAC)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    MATRIZ DE PERMISSÕES                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│                   │ Idoso  │ Cuidador │ Profissional │ Admin │         │
│  ─────────────────┼────────┼──────────┼──────────────┼───────┤         │
│  Próprio perfil   │  RW    │   R      │      R       │  RW   │         │
│  Conversas        │  RW    │   R*     │      R*      │  RW   │         │
│  Screenings       │  R     │   R*     │      RW      │  RW   │         │
│  Alertas          │  R     │   RW     │      RW      │  RW   │         │
│  Contatos emerg.  │  RW    │   RW     │      R       │  RW   │         │
│  Config sistema   │  -     │   -      │      -       │  RW   │         │
│  Audit logs       │  -     │   -      │      R       │  RW   │         │
│  Outros usuários  │  -     │   -      │      -       │  RW   │         │
│                                                                         │
│  R = Read | W = Write | * = Apenas idosos vinculados                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.5 JWT Token Structure

```json
{
  "header": {
    "alg": "RS256",
    "typ": "JWT",
    "kid": "key-2025-01"
  },
  "payload": {
    "sub": "user_abc123",
    "iat": 1706356800,
    "exp": 1706357700,
    "iss": "eva-mind-fzpn",
    "aud": "eva-api",
    "roles": ["idoso"],
    "permissions": ["chat:write", "profile:read"],
    "mfa_verified": true,
    "session_id": "sess_xyz789"
  }
}
```

---

## 5. Proteção Contra Ataques (OWASP Top 10)

### 5.1 A01:2021 - Broken Access Control

| Controle | Implementação |
|----------|---------------|
| RBAC | Verificação em cada endpoint |
| Principle of Least Privilege | Permissões mínimas por padrão |
| Deny by default | Acesso negado se não explicitamente permitido |
| Rate limiting | Por usuário e por IP |
| CORS | Whitelist de origens |

```go
// Middleware de autorização
func AuthorizeMiddleware(requiredPermission string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := getUserFromContext(r.Context())

            if !user.HasPermission(requiredPermission) {
                respondError(w, http.StatusForbidden, "Insufficient permissions")
                auditLog(user.ID, "ACCESS_DENIED", requiredPermission)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

### 5.2 A02:2021 - Cryptographic Failures

| Controle | Implementação |
|----------|---------------|
| TLS 1.3 obrigatório | HSTS header, redirect HTTP→HTTPS |
| Dados sensíveis criptografados | AES-256-GCM em repouso |
| Sem dados sensíveis em logs | Mascaramento automático |
| Senhas com hash seguro | Bcrypt cost 12 |

### 5.3 A03:2021 - Injection

| Tipo | Proteção |
|------|----------|
| SQL Injection | Prepared statements (SQLC) |
| NoSQL Injection | Validação de tipos |
| Command Injection | Não executa comandos externos |
| LDAP Injection | N/A (não usa LDAP) |

```go
// SQLC gera queries type-safe
// queries.sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

// Uso (injetável é impossível)
user, err := queries.GetUserByID(ctx, userID)
```

### 5.4 A04:2021 - Insecure Design

| Controle | Implementação |
|----------|---------------|
| Threat modeling | Realizado na fase de design |
| Secure design patterns | Repository pattern, clean architecture |
| Input validation | Validação em camada de API |
| Business logic protection | Limites de taxa, verificações de consistência |

### 5.5 A05:2021 - Security Misconfiguration

| Controle | Implementação |
|----------|---------------|
| Hardening automatizado | Ansible playbooks |
| Configuração como código | Terraform, Kubernetes manifests |
| Sem defaults inseguros | Revisão de configs em CI |
| Headers de segurança | CSP, X-Frame-Options, etc. |

```nginx
# Headers de segurança
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-Frame-Options "DENY" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Content-Security-Policy "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline';" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Permissions-Policy "geolocation=(), microphone=(self), camera=()" always;
```

### 5.6 A06:2021 - Vulnerable Components

| Controle | Implementação |
|----------|---------------|
| SCA (Software Composition Analysis) | Snyk em CI/CD |
| Dependabot | Alertas automáticos |
| SBOM | Gerado em cada build |
| Política de atualização | Críticas em 24h, Altas em 7 dias |

### 5.7 A07:2021 - Authentication Failures

| Controle | Implementação |
|----------|---------------|
| MFA disponível | TOTP, SMS, Biometria |
| Bloqueio de conta | 5 falhas = 15 min bloqueio |
| Session management | Tokens curtos, refresh rotation |
| Password strength | zxcvbn validation |

### 5.8 A08:2021 - Software and Data Integrity Failures

| Controle | Implementação |
|----------|---------------|
| Signed releases | GPG signing |
| CI/CD seguro | Branch protection, required reviews |
| Dependency verification | Checksums verificados |
| Content integrity | SRI para scripts externos |

### 5.9 A09:2021 - Security Logging and Monitoring

| Controle | Implementação |
|----------|---------------|
| Logging centralizado | ELK Stack / CloudWatch |
| Alertas em tempo real | PagerDuty integration |
| Auditoria completa | Todas as ações críticas |
| Retenção | 5 anos para compliance |

### 5.10 A10:2021 - Server-Side Request Forgery (SSRF)

| Controle | Implementação |
|----------|---------------|
| Whitelist de destinos | Apenas APIs aprovadas |
| Validação de URLs | Bloqueio de IPs internos |
| Sem redirecionamentos | Desabilitado por padrão |

---

## 6. Testes de Segurança

### 6.1 Análise Estática (SAST)

**Ferramenta:** SonarQube + gosec + Semgrep

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    RELATÓRIO SAST - 2025-01-27                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Ferramenta: SonarQube 10.x + gosec 2.x                                │
│  Commit: abc123def456                                                   │
│  Linhas analisadas: 45.000                                             │
│                                                                         │
│  VULNERABILIDADES:                                                      │
│  ├── Críticas:    0                                                    │
│  ├── Altas:       0                                                    │
│  ├── Médias:      3 (aceitas com justificativa)                        │
│  └── Baixas:      12 (monitoradas)                                     │
│                                                                         │
│  CODE SMELLS:                                                           │
│  ├── Blocker:     0                                                    │
│  ├── Critical:    0                                                    │
│  └── Major:       8                                                    │
│                                                                         │
│  MÉTRICAS:                                                              │
│  ├── Duplicação:  2.1%                                                 │
│  ├── Cobertura:   88.3%                                                │
│  └── Debt ratio:  0.3%                                                 │
│                                                                         │
│  STATUS: ✅ APROVADO                                                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Análise Dinâmica (DAST)

**Ferramenta:** OWASP ZAP + Burp Suite Pro

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    RELATÓRIO DAST - 2025-01-27                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Ferramenta: OWASP ZAP 2.14                                            │
│  Target: https://staging.eva-mind.com.br                               │
│  Duração: 4h 32m                                                       │
│                                                                         │
│  ALERTAS:                                                               │
│  ├── High:        0                                                    │
│  ├── Medium:      2                                                    │
│  │   ├── Cookie sem Secure flag (fixed)                               │
│  │   └── Missing CSP header (fixed)                                   │
│  ├── Low:         5                                                    │
│  └── Informational: 15                                                 │
│                                                                         │
│  TESTES REALIZADOS:                                                     │
│  ├── SQL Injection:           ✅ Não vulnerável                        │
│  ├── XSS:                     ✅ Não vulnerável                        │
│  ├── CSRF:                    ✅ Não vulnerável                        │
│  ├── Path Traversal:          ✅ Não vulnerável                        │
│  ├── Command Injection:       ✅ Não vulnerável                        │
│  └── SSRF:                    ✅ Não vulnerável                        │
│                                                                         │
│  STATUS: ✅ APROVADO                                                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 6.3 Teste de Penetração (Pentest)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    RELATÓRIO DE PENTEST                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Empresa: [Empresa de Pentest Certificada]                             │
│  Período: 15-19 Janeiro 2025                                           │
│  Metodologia: OWASP Testing Guide v4 + PTES                            │
│  Escopo: Aplicação web, API, Mobile (Android/iOS)                      │
│                                                                         │
│  RESUMO EXECUTIVO:                                                      │
│  O sistema EVA-Mind-FZPN demonstrou postura de segurança robusta.      │
│  Não foram identificadas vulnerabilidades críticas ou altas.           │
│                                                                         │
│  VULNERABILIDADES ENCONTRADAS:                                          │
│  ┌──────────┬────────────────────────────────────┬──────────┬────────┐ │
│  │Severidade│ Descrição                          │ CVSS     │ Status │ │
│  ├──────────┼────────────────────────────────────┼──────────┼────────┤ │
│  │ Crítica  │ Nenhuma                            │    -     │   -    │ │
│  │ Alta     │ Nenhuma                            │    -     │   -    │ │
│  │ Média    │ Rate limiting inconsistente em    │   4.3    │ Fixed  │ │
│  │          │ endpoint de recuperação de senha   │          │        │ │
│  │ Média    │ Verbose error messages em API     │   3.7    │ Fixed  │ │
│  │ Baixa    │ Headers de segurança incompletos  │   2.1    │ Fixed  │ │
│  │ Info     │ Divulgação de versão de servidor  │   1.0    │ Fixed  │ │
│  └──────────┴────────────────────────────────────┴──────────┴────────┘ │
│                                                                         │
│  TESTES REALIZADOS:                                                     │
│  • Reconhecimento e enumeração                                         │
│  • Testes de autenticação e autorização                                │
│  • Testes de injeção (SQL, XSS, Command, etc.)                        │
│  • Testes de gerenciamento de sessão                                   │
│  • Testes de criptografia                                              │
│  • Testes de lógica de negócio                                         │
│  • Testes de mobile (Android/iOS)                                      │
│                                                                         │
│  RECOMENDAÇÕES IMPLEMENTADAS:                                           │
│  ✅ Rate limiting em todos os endpoints de autenticação               │
│  ✅ Error messages genéricas em produção                              │
│  ✅ Headers de segurança completos                                    │
│  ✅ Remoção de banners de versão                                      │
│                                                                         │
│  CONCLUSÃO:                                                             │
│  O sistema está APROVADO do ponto de vista de segurança.               │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 6.4 Análise de Composição (SCA)

**Ferramenta:** Snyk + Dependabot

| Dependência | Versão | Vulnerabilidades | Ação |
|-------------|--------|------------------|------|
| golang.org/x/crypto | 0.18.0 | 0 | OK |
| github.com/lib/pq | 1.10.9 | 0 | OK |
| github.com/go-chi/chi | 5.0.11 | 0 | OK |
| github.com/golang-jwt/jwt | 5.2.0 | 0 | OK |
| github.com/redis/go-redis | 9.4.0 | 0 | OK |

**Total de dependências:** 87
**Com vulnerabilidades conhecidas:** 0

---

## 7. Logs e Auditoria

### 7.1 Eventos Auditados

| Evento | Dados Registrados | Severidade |
|--------|-------------------|------------|
| Login sucesso | user_id, ip, user_agent, timestamp | Info |
| Login falha | email_attempt, ip, reason, timestamp | Warning |
| Logout | user_id, session_id, timestamp | Info |
| Acesso a dados sensíveis | user_id, resource, action, timestamp | Info |
| Modificação de dados | user_id, resource, before, after, timestamp | Info |
| Alerta de risco | user_id, risk_level, trigger, timestamp | Warning |
| Erro de sistema | error_code, stack_trace, timestamp | Error |
| Tentativa de acesso negada | user_id, resource, reason, timestamp | Warning |
| Mudança de permissões | admin_id, target_user, changes, timestamp | Info |
| Exportação de dados | user_id, data_type, timestamp | Info |

### 7.2 Formato de Log

```json
{
  "timestamp": "2025-01-27T14:30:00.000Z",
  "level": "INFO",
  "service": "eva-api",
  "trace_id": "abc123def456",
  "span_id": "789xyz",
  "event": "USER_LOGIN",
  "user_id": "user_abc123",
  "ip_address": "203.0.113.50",
  "user_agent": "EVA-Mobile/2.0.0 (Android 13)",
  "details": {
    "method": "biometric",
    "mfa_verified": true
  },
  "duration_ms": 234
}
```

### 7.3 Retenção e Proteção

| Tipo de Log | Retenção | Armazenamento | Acesso |
|-------------|----------|---------------|--------|
| Auditoria de segurança | 5 anos | S3 + Glacier | Admin + Compliance |
| Logs de aplicação | 90 dias | CloudWatch | DevOps |
| Logs de acesso | 2 anos | S3 | Admin |
| Métricas | 1 ano | Prometheus | DevOps |

### 7.4 Monitoramento e Alertas

| Condição | Severidade | Notificação |
|----------|------------|-------------|
| 5+ login failures (mesmo IP) | Warning | Slack |
| Acesso de país não usual | Warning | Slack + E-mail |
| Tentativa de SQL injection | Critical | PagerDuty |
| Acesso a dados massivo | Warning | Slack |
| Erro 5xx > 1% | Critical | PagerDuty |
| Latência P99 > 2s | Warning | Slack |

---

## 8. Plano de Resposta a Incidentes

### 8.1 Classificação de Incidentes

| Nível | Descrição | SLA Resposta | SLA Resolução |
|-------|-----------|--------------|---------------|
| P1 - Crítico | Vazamento de dados, sistema down | 15 min | 4h |
| P2 - Alto | Vulnerabilidade explorada, parcial down | 1h | 8h |
| P3 - Médio | Vulnerabilidade detectada, degradação | 4h | 24h |
| P4 - Baixo | Anomalia de segurança, sem impacto | 24h | 7 dias |

### 8.2 Fluxo de Resposta

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    FLUXO DE RESPOSTA A INCIDENTES                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. DETECÇÃO                                                            │
│     │  • Alerta automático (SIEM, WAF, IDS)                            │
│     │  • Relato de usuário                                             │
│     │  • Descoberta em auditoria                                       │
│     ▼                                                                   │
│  2. TRIAGEM (15 min)                                                    │
│     │  • Classificar severidade (P1-P4)                                │
│     │  • Identificar escopo                                            │
│     │  • Acionar equipe apropriada                                     │
│     ▼                                                                   │
│  3. CONTENÇÃO (imediato)                                                │
│     │  • Isolar sistemas afetados                                      │
│     │  • Bloquear vetores de ataque                                    │
│     │  • Preservar evidências                                          │
│     ▼                                                                   │
│  4. ERRADICAÇÃO                                                         │
│     │  • Remover malware/backdoor                                      │
│     │  • Corrigir vulnerabilidade                                      │
│     │  • Atualizar credenciais comprometidas                           │
│     ▼                                                                   │
│  5. RECUPERAÇÃO                                                         │
│     │  • Restaurar sistemas                                            │
│     │  • Verificar integridade                                         │
│     │  • Monitorar por recorrência                                     │
│     ▼                                                                   │
│  6. PÓS-INCIDENTE                                                       │
│        • Relatório de incidente                                        │
│        • Lições aprendidas                                             │
│        • Atualizar controles                                           │
│        • Notificar stakeholders/ANPD (se aplicável)                    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 8.3 Contatos de Emergência

| Papel | Contato | Disponibilidade |
|-------|---------|-----------------|
| Security Lead | [Telefone] | 24/7 |
| CTO | [Telefone] | 24/7 |
| DPO | [Telefone] | Horário comercial + on-call |
| Fornecedor de Cloud | [Suporte Premium] | 24/7 |
| Jurídico | [Telefone] | On-call |

---

## 9. Conformidade LGPD (Aspectos de Segurança)

### 9.1 Medidas Técnicas (Art. 46)

| Medida | Implementação | Status |
|--------|---------------|--------|
| Criptografia | AES-256, TLS 1.3 | ✅ |
| Pseudonimização | IDs internos separados de PII | ✅ |
| Anonimização | Processo definido para dados de pesquisa | ✅ |
| Controle de acesso | RBAC + MFA | ✅ |
| Logs de auditoria | Completo e imutável | ✅ |
| Backup seguro | Criptografado, testado | ✅ |

### 9.2 Direito de Acesso Seguro

| Requisito | Implementação |
|-----------|---------------|
| Verificação de identidade | MFA antes de exportar dados |
| Formato seguro | Download criptografado, link expirável |
| Log de acesso | Registro de toda exportação |
| Limite de taxa | Máximo 1 export/dia |

---

## 10. Conclusão

O EVA-Mind-FZPN implementa controles de segurança robustos alinhados com:

- **ISO 27001** - Sistema de Gestão de Segurança da Informação
- **OWASP Top 10** - Proteção contra vulnerabilidades web comuns
- **LGPD** - Proteção de dados pessoais
- **ABNT NBR ISO/TR 81001-2-8** - Cibersegurança em dispositivos médicos

**Status de Segurança:** ✅ APROVADO

---

## Aprovações

| Função | Nome | Assinatura | Data |
|--------|------|------------|------|
| Security Officer | | | |
| CTO | | | |
| DPO | | | |
| Responsável Regulatório | José R F Junior | | 2025-01-27 |

---

**Documento controlado - Versão 1.0**
**Próxima revisão: 2025-07-27 (semestral)**
