# Infraestrutura e Operações
## EVA-Mind-FZPN - Companion IA para Idosos

**Documento:** INFRA-EVA-001
**Versão:** 1.0
**Data:** 2025-01-27

---

## 1. Hospedagem e Servidores

### 1.1 Provedor de Nuvem

| Item | Especificação |
|------|---------------|
| **Provedor** | AWS (Amazon Web Services) |
| **Região primária** | sa-east-1 (São Paulo) |
| **Região DR** | us-east-1 (N. Virginia) |
| **Conta** | Organization com múltiplas contas |
| **Suporte** | Business Support |

### 1.2 Certificações do Datacenter

| Certificação | Status |
|--------------|--------|
| ISO 27001 | ✅ Válido |
| SOC 1/2/3 | ✅ Válido |
| PCI DSS Level 1 | ✅ Válido |
| HIPAA Eligible | ✅ Válido |

### 1.3 Arquitetura de Alta Disponibilidade

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    ARQUITETURA DE INFRAESTRUTURA                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│                           ┌─────────────┐                               │
│                           │  CloudFront │                               │
│                           │    (CDN)    │                               │
│                           └──────┬──────┘                               │
│                                  │                                      │
│                           ┌──────▼──────┐                               │
│                           │     WAF     │                               │
│                           │   Shield    │                               │
│                           └──────┬──────┘                               │
│                                  │                                      │
│     ┌────────────────────────────┼────────────────────────────┐        │
│     │                   VPC (10.0.0.0/16)                     │        │
│     │                            │                            │        │
│     │  ┌─────────────────────────┼─────────────────────────┐ │        │
│     │  │              Public Subnets                       │ │        │
│     │  │     ┌───────────────────┴───────────────────┐     │ │        │
│     │  │     │                                       │     │ │        │
│     │  │  ┌──▼──┐                              ┌──▼──┐    │ │        │
│     │  │  │ ALB │                              │ ALB │    │ │        │
│     │  │  │ AZ-a│                              │ AZ-b│    │ │        │
│     │  │  └──┬──┘                              └──┬──┘    │ │        │
│     │  └─────┼────────────────────────────────────┼──────┘ │        │
│     │        │                                    │        │        │
│     │  ┌─────┼────────────────────────────────────┼──────┐ │        │
│     │  │     │       Private Subnets              │      │ │        │
│     │  │  ┌──▼────────┐              ┌────────────▼──┐   │ │        │
│     │  │  │   EKS     │              │     EKS       │   │ │        │
│     │  │  │ Node Pool │              │  Node Pool    │   │ │        │
│     │  │  │   AZ-a    │              │    AZ-b       │   │ │        │
│     │  │  └───────────┘              └───────────────┘   │ │        │
│     │  │        │                           │            │ │        │
│     │  │  ┌─────▼───────────────────────────▼─────┐      │ │        │
│     │  │  │              Services                 │      │ │        │
│     │  │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ │      │ │        │
│     │  │  │  │ Cortex  │ │Hippoc.  │ │  Motor  │ │      │ │        │
│     │  │  │  └─────────┘ └─────────┘ └─────────┘ │      │ │        │
│     │  │  └───────────────────────────────────────┘      │ │        │
│     │  └─────────────────────────────────────────────────┘ │        │
│     │                                                      │        │
│     │  ┌───────────────────────────────────────────────┐  │        │
│     │  │              Data Subnets                     │  │        │
│     │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐         │  │        │
│     │  │  │   RDS   │ │ Qdrant  │ │  Redis  │         │  │        │
│     │  │  │ Primary │ │ Cluster │ │ Cluster │         │  │        │
│     │  │  └────┬────┘ └─────────┘ └─────────┘         │  │        │
│     │  │       │                                       │  │        │
│     │  │  ┌────▼────┐                                  │  │        │
│     │  │  │   RDS   │                                  │  │        │
│     │  │  │ Standby │                                  │  │        │
│     │  │  └─────────┘                                  │  │        │
│     │  └───────────────────────────────────────────────┘  │        │
│     └──────────────────────────────────────────────────────┘        │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

### 1.4 Especificações de Recursos

| Serviço | Tipo | Especificação | Quantidade |
|---------|------|---------------|------------|
| EKS Nodes | c6i.xlarge | 4 vCPU, 8 GB RAM | 4-8 (auto-scale) |
| RDS PostgreSQL | db.r6g.large | 2 vCPU, 16 GB RAM | 2 (primary + standby) |
| ElastiCache Redis | cache.r6g.large | 2 vCPU, 13 GB RAM | 3 (cluster) |
| Qdrant | r6i.xlarge | 4 vCPU, 32 GB RAM | 3 (cluster) |

### 1.5 SLA e Uptime

| Serviço | SLA AWS | SLA EVA | Medição |
|---------|---------|---------|---------|
| Aplicação | - | 99.5% | Mensal |
| API | - | 99.9% | Mensal |
| Banco de dados | 99.95% | 99.9% | Mensal |
| CDN | 99.9% | 99.9% | Mensal |

---

## 2. Escalabilidade

### 2.1 Capacidade

| Métrica | Atual | Máximo (auto-scale) |
|---------|-------|---------------------|
| Usuários simultâneos | 500 | 5.000 |
| Requests/segundo | 200 | 2.000 |
| Mensagens/dia | 50.000 | 500.000 |
| Storage (DB) | 100 GB | 1 TB |
| Storage (Vectors) | 50 GB | 500 GB |

### 2.2 Auto-Scaling

```yaml
# Kubernetes HPA
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: eva-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: eva-api
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 60
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
```

---

## 3. Backup e Recuperação de Desastres

### 3.1 Política de Backup

| Dados | Frequência | Retenção | Localização |
|-------|------------|----------|-------------|
| PostgreSQL (full) | Diário 02:00 | 30 dias | S3 sa-east-1 |
| PostgreSQL (incremental) | 6 horas | 7 dias | S3 sa-east-1 |
| PostgreSQL (WAL) | Contínuo | 7 dias | S3 sa-east-1 |
| Redis (snapshot) | Diário | 7 dias | S3 sa-east-1 |
| Qdrant (snapshot) | Diário | 14 dias | S3 sa-east-1 |
| Configurações | A cada mudança | 90 dias | S3 + Git |

### 3.2 Recuperação de Desastres (DR)

| Métrica | Objetivo | Atual |
|---------|----------|-------|
| **RTO** (Recovery Time Objective) | 4 horas | 2 horas |
| **RPO** (Recovery Point Objective) | 1 hora | 15 minutos |

### 3.3 Procedimento de DR

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    PROCEDIMENTO DE DR                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  CENÁRIO: Falha total da região sa-east-1                              │
│                                                                         │
│  1. DETECÇÃO (0-15 min)                                                │
│     • Alertas automáticos de indisponibilidade                        │
│     • Verificação por equipe de plantão                                │
│     • Decisão de ativar DR                                             │
│                                                                         │
│  2. ATIVAÇÃO DR (15-60 min)                                            │
│     • Promover RDS standby em us-east-1                                │
│     • Restaurar último snapshot de Qdrant                              │
│     • Redirecionar DNS para us-east-1                                  │
│     • Escalar EKS nodes em us-east-1                                   │
│                                                                         │
│  3. VALIDAÇÃO (60-120 min)                                             │
│     • Testes de smoke em ambiente DR                                   │
│     • Verificar integridade de dados                                   │
│     • Monitorar métricas de saúde                                      │
│                                                                         │
│  4. COMUNICAÇÃO                                                         │
│     • Notificar stakeholders internos                                  │
│     • Atualizar status page                                            │
│     • Comunicar usuários se necessário                                 │
│                                                                         │
│  5. RETORNO À NORMALIDADE (quando região primária disponível)          │
│     • Sincronizar dados de volta para sa-east-1                        │
│     • Testar ambiente primário                                         │
│     • Failback gradual                                                 │
│     • Atualizar DNS                                                    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.4 Testes de Backup

| Teste | Frequência | Último Teste | Resultado |
|-------|------------|--------------|-----------|
| Restauração de DB | Trimestral | 2025-01-15 | ✅ Sucesso (45 min) |
| Failover de RDS | Trimestral | 2025-01-15 | ✅ Sucesso (2 min) |
| DR completo | Anual | 2024-11-20 | ✅ Sucesso (1h 45min) |
| Recuperação de arquivo | Mensal | 2025-01-20 | ✅ Sucesso |

---

## 4. Monitoramento e Alertas

### 4.1 Stack de Monitoramento

| Ferramenta | Uso |
|------------|-----|
| **Prometheus** | Coleta de métricas |
| **Grafana** | Visualização e dashboards |
| **CloudWatch** | Logs e métricas AWS |
| **PagerDuty** | Alertas e on-call |
| **Datadog APM** | Tracing distribuído |

### 4.2 Métricas Monitoradas

| Categoria | Métrica | Threshold Warning | Threshold Critical |
|-----------|---------|-------------------|-------------------|
| **Disponibilidade** | Uptime | <99.9% | <99.5% |
| **Latência** | P50 | >200ms | >500ms |
| **Latência** | P99 | >1s | >2s |
| **Erros** | Taxa de erro 5xx | >0.5% | >2% |
| **CPU** | Utilização | >70% | >90% |
| **Memória** | Utilização | >75% | >90% |
| **Disco** | Utilização | >70% | >85% |
| **DB** | Conexões | >80% | >95% |
| **DB** | Replication lag | >10s | >60s |

### 4.3 Dashboard Principal

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    EVA-Mind Operations Dashboard                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  STATUS GERAL: 🟢 Operacional                                          │
│                                                                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │
│  │   Uptime    │  │   RPS       │  │  Latência   │  │   Erros     │   │
│  │   99.98%    │  │    156      │  │   P99: 420ms│  │    0.02%    │   │
│  │   🟢        │  │   🟢        │  │   🟢        │  │   🟢        │   │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │
│                                                                         │
│  RECURSOS:                                                              │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │ CPU    ████████████░░░░░░░░  45%                                │  │
│  │ Memory ██████████████░░░░░░  62%                                │  │
│  │ Disk   ████████░░░░░░░░░░░░  38%                                │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  PODS (EKS):                                                            │
│  ├── eva-api:        6/6 Running                                       │
│  ├── eva-cortex:     4/4 Running                                       │
│  ├── eva-hippocampus:3/3 Running                                       │
│  └── eva-motor:      2/2 Running                                       │
│                                                                         │
│  ALERTAS ATIVOS: 0                                                      │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.4 Configuração de Alertas

| Alerta | Condição | Severidade | Canal |
|--------|----------|------------|-------|
| HighLatency | P99 > 2s por 5min | Critical | PagerDuty |
| HighErrorRate | 5xx > 2% por 5min | Critical | PagerDuty |
| PodCrashLoop | Restart > 5 em 10min | Critical | PagerDuty |
| DBConnectionHigh | Conexões > 90% | Warning | Slack |
| DiskSpaceLow | Uso > 85% | Warning | Slack |
| CertificateExpiring | Expira em < 30 dias | Warning | E-mail |

---

## 5. CI/CD

### 5.1 Pipeline

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         PIPELINE CI/CD                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐             │
│  │  COMMIT │───▶│  BUILD  │───▶│  TEST   │───▶│  SCAN   │             │
│  │         │    │         │    │         │    │         │             │
│  │ • Push  │    │ • Go    │    │ • Unit  │    │ • SAST  │             │
│  │ • PR    │    │   build │    │ • Integ │    │ • SCA   │             │
│  │         │    │ • Docker│    │ • Lint  │    │ • DAST  │             │
│  └─────────┘    └─────────┘    └─────────┘    └────┬────┘             │
│                                                     │                   │
│                                     ┌───────────────┴───────────────┐   │
│                                     │                               │   │
│                                     ▼                               ▼   │
│                              ┌─────────────┐                ┌───────────┐
│                              │   STAGING   │                │   FAIL    │
│                              │             │                │           │
│                              │ • Deploy    │                │ • Alert   │
│                              │ • Smoke     │                │ • Block   │
│                              │ • E2E       │                └───────────┘
│                              └──────┬──────┘                            │
│                                     │                                   │
│                                     ▼                                   │
│                              ┌─────────────┐                            │
│                              │  APPROVAL   │                            │
│                              │  (Manual)   │                            │
│                              └──────┬──────┘                            │
│                                     │                                   │
│                                     ▼                                   │
│                              ┌─────────────┐                            │
│                              │ PRODUCTION  │                            │
│                              │             │                            │
│                              │ • Canary    │                            │
│                              │ • Monitor   │                            │
│                              │ • Rollback  │                            │
│                              └─────────────┘                            │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Ferramentas

| Etapa | Ferramenta |
|-------|------------|
| Source Control | GitHub |
| CI/CD | GitHub Actions |
| Container Registry | Amazon ECR |
| Infrastructure as Code | Terraform |
| Kubernetes Deploy | ArgoCD |
| Secrets | AWS Secrets Manager |

### 5.3 Workflow de Deploy

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build and push Docker image
        run: |
          docker build -t eva-api:${{ github.sha }} .
          docker push $ECR_REGISTRY/eva-api:${{ github.sha }}

  test:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Check coverage
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
          if [ ${coverage%\%} -lt 80 ]; then exit 1; fi

  security-scan:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Run Snyk
        uses: snyk/actions/golang@master

      - name: Run Trivy
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: ${{ env.ECR_REGISTRY }}/eva-api:${{ github.sha }}

  deploy-staging:
    needs: [test, security-scan]
    runs-on: ubuntu-latest
    environment: staging
    steps:
      - name: Deploy to staging
        run: kubectl apply -f k8s/staging/

      - name: Run smoke tests
        run: ./scripts/smoke-test.sh staging

  deploy-production:
    needs: deploy-staging
    runs-on: ubuntu-latest
    environment: production
    steps:
      - name: Deploy canary (10%)
        run: kubectl apply -f k8s/production/canary.yaml

      - name: Monitor canary
        run: ./scripts/monitor-canary.sh --duration 15m

      - name: Promote to 100%
        run: kubectl apply -f k8s/production/full.yaml
```

### 5.4 Rollback

| Método | Tempo | Uso |
|--------|-------|-----|
| Kubernetes rollback | <1 min | Problemas de aplicação |
| ArgoCD sync anterior | <2 min | Configuração errada |
| Database restore | 15-60 min | Dados corrompidos |

---

## 6. Ambientes

### 6.1 Matriz de Ambientes

| Ambiente | Propósito | Dados | Acesso |
|----------|-----------|-------|--------|
| **Development** | Desenvolvimento local | Sintéticos | Devs |
| **Staging** | Testes de integração | Anonimizados | QA + Devs |
| **Production** | Produção real | Reais | Operações |
| **DR** | Disaster Recovery | Réplica | Emergência |

### 6.2 Paridade de Ambientes

| Componente | Development | Staging | Production |
|------------|-------------|---------|------------|
| Kubernetes | Minikube | EKS (menor) | EKS (full) |
| PostgreSQL | Docker | RDS (menor) | RDS Multi-AZ |
| Redis | Docker | ElastiCache | ElastiCache Cluster |
| Qdrant | Docker | EC2 single | EC2 Cluster |

---

## 7. Conclusão

A infraestrutura do EVA-Mind-FZPN foi projetada para:

- **Alta disponibilidade**: Multi-AZ, auto-scaling, failover automático
- **Segurança**: VPC isolada, WAF, criptografia em trânsito e repouso
- **Escalabilidade**: Suporta 10x da carga atual
- **Recuperação**: RTO de 4h, RPO de 1h
- **Observabilidade**: Monitoramento completo, alertas proativos

---

## Aprovações

| Função | Nome | Assinatura | Data |
|--------|------|------------|------|
| DevOps Lead | | | |
| SRE | | | |
| CTO | José R F Junior | | 2025-01-27 |

---

**Documento controlado - Versão 1.0**
