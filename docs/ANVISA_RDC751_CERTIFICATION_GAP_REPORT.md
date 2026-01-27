# RELATÓRIO DE ANÁLISE DE GAPS - CERTIFICAÇÃO ANVISA RDC 751/2022
## EVA-Mind-FZPN - Software como Dispositivo Médico (SaMD) Classe II

**Data:** 2026-01-27
**Versão:** 1.0
**Autor:** Análise automatizada + documentação técnica
**Referências:** EVA-Mind_Cap3_Analise_Tecnica.docx, RDC 751/2022 (Anexo II, Cap. 3)

---

## 1. RESUMO EXECUTIVO

### 1.1 Classificação do Produto
| Aspecto | Valor |
|---------|-------|
| **Tipo** | Software como Dispositivo Médico (SaMD) |
| **Classe de Risco** | II (Risco Médio) |
| **Uso Pretendido** | Acompanhamento e suporte emocional para idosos |
| **Funcionalidades Críticas** | Avaliação de risco suicida (C-SSRS), depressão (PHQ-9), ansiedade (GAD-7) |
| **Regulamentação** | RDC 751/2022 (Anexo II, Capítulo 3) |

### 1.2 Status Geral de Conformidade

| Área | Status | Progresso |
|------|--------|-----------|
| Funcionalidades Clínicas | ✅ Implementado | 95% |
| Sistema de Alertas | ✅ Implementado | 90% |
| Trilha de Auditoria LGPD | ✅ Implementado | 95% |
| Direitos do Titular (Art. 18) | ✅ Implementado | 90% |
| Métricas e Monitoramento | ✅ Implementado | 85% |
| Interoperabilidade FHIR | ✅ Implementado | 80% |
| **Gestão de Riscos (ISO 14971)** | ⚠️ Parcial | 40% |
| **Usabilidade (IEC 62366-1)** | ⚠️ Parcial | 30% |
| **Documentação Técnica ANVISA** | ❌ Pendente | 20% |
| **Validação Clínica** | ❌ Pendente | 15% |
| **Testes de Software (IEC 62304)** | ⚠️ Parcial | 50% |

---

## 2. ANÁLISE DETALHADA POR REQUISITO

### 2.1 RDC 751/2022 - Anexo II, Capítulo 3

#### Art. 47 - Validação de Escalas Clínicas

| Requisito | Status | Evidência |
|-----------|--------|-----------|
| PHQ-9 validado cientificamente | ✅ | Escala internacional validada |
| GAD-7 validado cientificamente | ✅ | Escala internacional validada |
| C-SSRS validado cientificamente | ✅ | Escala Columbia - padrão ouro |
| Algoritmos de scoring documentados | ✅ | `clinical_scales.go:23-454` |
| Níveis de risco estratificados | ✅ | none→low→moderate→high→critical |
| Intervenções por nível de risco | ✅ | CVV 188, SAMU 192, hospitalização |

**Implementação atual:**
```go
// internal/cortex/scales/clinical_scales.go
type CSSRSResult struct {
    Score             int
    RiskLevel         string  // none, low, moderate, high, critical
    HasSuicidalIdeation bool
    HasSuicidalBehavior bool  // Comportamento = CRÍTICO automático
    RequiresIntervention bool
    Recommendations   []string
    InterventionPlan  *InterventionPlan
}
```

**Gap identificado:** ⚠️ Falta documentação formal da validação clínica das escalas no contexto brasileiro (estudos de validação semântica/cultural).

---

#### Art. 54 - Dados Clínicos Estruturados

| Requisito | Status | Evidência |
|-----------|--------|-----------|
| Armazenamento estruturado | ✅ | `clinical_assessments` table |
| Histórico de avaliações | ✅ | Respostas individuais persistidas |
| Codificação padronizada | ⚠️ | FHIR parcial, falta ICD-10/CID-10 |
| Exportação de dados | ✅ | JSON, FHIR R4 Bundle |

**Schema implementado:**
```sql
-- migrations/002_clinical_and_vision_features.sql
CREATE TABLE clinical_assessments (
    id SERIAL PRIMARY KEY,
    patient_id BIGINT NOT NULL,
    assessment_type VARCHAR(20), -- 'PHQ-9', 'GAD-7', 'C-SSRS', 'MMSE', 'MoCA'
    total_score INTEGER,
    risk_level VARCHAR(20),
    completed_at TIMESTAMP,
    metadata JSONB
);
```

**Gap identificado:** ⚠️ Integração com terminologias CID-10/CID-11 para classificação de condições identificadas.

---

#### Art. 58 - Gestão de Riscos

| Requisito ISO 14971:2019 | Status | Evidência |
|--------------------------|--------|-----------|
| Análise de risco documentada | ❌ | Não existe documento formal |
| Matriz de risco (probabilidade x severidade) | ❌ | Não existe |
| FMEA (Failure Mode Effects Analysis) | ❌ | Não existe |
| Identificação de perigos | ⚠️ | Implícito no código, não documentado |
| Controles de risco implementados | ✅ | Alertas, escalação, intervenções |
| Risco residual aceitável | ❌ | Não documentado |
| Rastreabilidade req→risco | ❌ | Não existe |

**Controles implementados (não documentados formalmente):**
- Detecção automática de risco suicida (C-SSRS Q6 = CRÍTICO)
- Escalação multi-canal (Push→WhatsApp→SMS→Email→Ligação)
- Timeouts por prioridade (Crítica=30s, Alta=2min)
- Fornecimento de recursos de emergência (CVV 188, SAMU 192)

**Gap crítico:** ❌ **Documento de Gestão de Riscos ISO 14971 não existe.** Prioridade máxima.

---

#### Art. 59 - Segurança da Informação

| Requisito | Status | Evidência |
|-----------|--------|-----------|
| Criptografia de senhas | ✅ | bcrypt cost 14 |
| Tokens de acesso seguros | ✅ | JWT HS256, 15min expiry |
| Controle de acesso | ✅ | Middleware JWT, roles |
| Trilha de auditoria | ✅ | `lgpd_audit_log` completo |
| Criptografia em trânsito | ✅ | HTTPS (configurável) |
| Criptografia em repouso | ⚠️ | Depende do PostgreSQL |
| Política de retenção | ✅ | Auto-expiração por categoria |
| Backup e recuperação | ❌ | Não documentado |

**Implementação atual:**
```go
// internal/brainstem/auth/service.go
func HashPassword(password string) (string, error) {
    return bcrypt.GenerateFromPassword([]byte(password), 14)
}

// internal/audit/lgpd_audit.go
type AuditEvent struct {
    EventType     AuditEventType  // DATA_ACCESS, DATA_CREATE, etc.
    DataCategory  DataCategory    // PERSONAL, SENSITIVE, CLINICAL
    LegalBasis    LegalBasis      // CONSENT, HEALTH_PROTECTION, etc.
    RetentionDays int             // Auto-expiration
}
```

**Gap identificado:** ⚠️ Falta documentação de política de backup/recuperação de dados clínicos.

---

#### Art. 60 - Proteção de Dados Pessoais (LGPD)

| Requisito LGPD | Status | Evidência |
|----------------|--------|-----------|
| Art. 7 - Base legal documentada | ✅ | `legal_basis` em cada evento |
| Art. 8 - Gestão de consentimento | ✅ | `lgpd_consents` table |
| Art. 18, I - Acesso aos dados | ✅ | `GetDataAccessReport()` |
| Art. 18, II - Correção | ✅ | `RectifyPersonalData()` |
| Art. 18, III - Anonimização | ✅ | Implementado na deleção |
| Art. 18, V - Portabilidade | ✅ | `ExportPersonalData()` JSON |
| Art. 18, VI - Eliminação | ✅ | `DeletePersonalData()` |
| Art. 37 - Registro de operações | ✅ | Audit trail completo |
| DPO designado | ❌ | Não documentado |
| RIPD (Relatório de Impacto) | ❌ | Não existe |

**Implementação atual:**
```go
// internal/audit/data_rights.go
func (s *DataRightsService) ExportPersonalData(ctx context.Context, subjectID int64, format string) (*DataExportResult, error)
func (s *DataRightsService) DeletePersonalData(ctx context.Context, subjectID int64, retainAuditLog bool) (*DeletionResult, error)
func (s *DataRightsService) RectifyPersonalData(ctx context.Context, subjectID int64, field, oldValue, newValue string) error
```

**Gap identificado:** ❌ RIPD (Relatório de Impacto à Proteção de Dados) não elaborado.

---

#### Art. 73 - Usabilidade (IEC 62366-1:2015)

| Requisito | Status | Evidência |
|-----------|--------|-----------|
| Análise de uso pretendido | ⚠️ | Documentado parcialmente |
| Perfil de usuários | ⚠️ | Implícito (idosos, cuidadores) |
| Tarefas críticas identificadas | ❌ | Não documentado |
| Estudo de usabilidade formativo | ❌ | Não realizado |
| Estudo de usabilidade somativo | ❌ | Não realizado |
| Erros de uso identificados | ❌ | Não documentado |
| Interface adaptada para idosos | ✅ | Interface por voz, personalizada |
| Relatório de engenharia de usabilidade | ❌ | Não existe |

**Gap crítico:** ❌ **Arquivo de Engenharia de Usabilidade IEC 62366-1 não existe.** Prioridade alta.

---

### 2.2 IEC 62304 - Ciclo de Vida de Software para Dispositivos Médicos

| Requisito | Status | Evidência |
|-----------|--------|-----------|
| Classificação de segurança do SW | ⚠️ | Classe B implícita |
| Plano de desenvolvimento | ❌ | Não formalizado |
| Requisitos de software | ⚠️ | Parcial em docs/ |
| Arquitetura documentada | ✅ | Estrutura cortex/brainstem/hippocampus |
| Design detalhado | ⚠️ | Código auto-documentado |
| Testes unitários | ✅ | 109 testes passando |
| Testes de integração | ⚠️ | Parcial |
| Testes de sistema | ❌ | Não formalizados |
| Validação de software | ❌ | Não realizada |
| Gestão de configuração | ✅ | Git |
| Gestão de problemas | ⚠️ | GitHub Issues (informal) |
| Rastreabilidade | ❌ | Não implementada |

**Testes implementados:**
```
internal/mocks/              - 12 testes (mocks Firebase, Twilio, Email)
internal/cortex/scales/      - 25 testes (PHQ-9, GAD-7, C-SSRS)
internal/cortex/alert/       - 17 testes (escalação, prioridades)
internal/metrics/            - 18 testes (Prometheus)
internal/audit/              - 37 testes (LGPD audit, data rights)
TOTAL: 109 testes passando
```

**Gap identificado:** ❌ Matriz de rastreabilidade requisitos→código→testes não existe.

---

## 3. MAPA DE GAPS vs. IMPLEMENTAÇÃO

### 3.1 O que ESTÁ implementado no código:

| Funcionalidade | Arquivo | Status |
|----------------|---------|--------|
| Escalas clínicas (C-SSRS, PHQ-9, GAD-7) | `internal/cortex/scales/clinical_scales.go` | ✅ Completo |
| Estratificação de risco | `internal/cortex/scales/clinical_scales.go` | ✅ Completo |
| Alertas de emergência | `internal/cortex/alert/escalation.go` | ✅ Completo |
| Escalação multi-canal | `internal/cortex/alert/escalation.go` | ✅ Completo |
| Predição de crise | `internal/cortex/prediction/crisis_predictor.go` | ✅ Completo |
| Explicabilidade clínica | `internal/cortex/explainability/clinical_decision_explainer.go` | ✅ Completo |
| Audit trail LGPD | `internal/audit/lgpd_audit.go` | ✅ Completo |
| Direitos do titular | `internal/audit/data_rights.go` | ✅ Completo |
| Métricas Prometheus | `internal/metrics/metrics.go` | ✅ Completo |
| Interoperabilidade FHIR | `internal/integration/fhir_adapter.go` | ✅ Parcial |
| Autenticação segura | `internal/brainstem/auth/` | ✅ Completo |
| Consentimento | `migrations/018_lgpd_audit_trail.sql` | ✅ Completo |

### 3.2 O que FALTA para certificação:

| Documento/Artefato | Prioridade | Esforço Estimado |
|--------------------|------------|------------------|
| **Arquivo de Gestão de Riscos ISO 14971** | 🔴 Crítica | 40-60h |
| **FMEA (Failure Mode Effects Analysis)** | 🔴 Crítica | 20-30h |
| **Arquivo de Usabilidade IEC 62366-1** | 🔴 Crítica | 60-80h |
| **Plano de Validação de Software** | 🔴 Crítica | 20-30h |
| **Dossiê Técnico ANVISA** | 🔴 Crítica | 80-100h |
| RIPD (LGPD) | 🟡 Alta | 16-24h |
| Matriz de rastreabilidade | 🟡 Alta | 16-24h |
| Protocolo de validação clínica | 🟡 Alta | 40-60h |
| Plano de gestão de configuração | 🟠 Média | 8-16h |
| Política de backup/recuperação | 🟠 Média | 8-16h |
| Manual do usuário | 🟠 Média | 24-40h |
| Instruções de uso (IFU) | 🟠 Média | 16-24h |

---

## 4. PLANO DE AÇÃO PRIORITIZADO

### Fase 1: Documentação Regulatória Crítica (4-6 semanas)

#### 4.1.1 Gestão de Riscos ISO 14971:2019

**Entregáveis:**
1. [ ] **Arquivo de Gestão de Riscos** contendo:
   - Escopo e contexto de uso
   - Identificação de perigos (hazards)
   - Situações perigosas
   - Estimativa de risco (probabilidade × severidade)
   - Avaliação de risco (aceitabilidade)
   - Controles de risco implementados
   - Risco residual

2. [ ] **Matriz de Risco** com categorias:
   - Probabilidade: Muito improvável → Frequente
   - Severidade: Insignificante → Catastrófica
   - Aceitabilidade: Aceitável / ALARP / Inaceitável

3. [ ] **FMEA (Análise de Modo e Efeito de Falha)**:

   | Componente | Modo de Falha | Efeito | Severidade | Causa | Ocorrência | Controle | Detecção | RPN |
   |------------|---------------|--------|------------|-------|------------|----------|----------|-----|
   | C-SSRS | Score incorreto | Risco subestimado | 5 | Bug no algoritmo | 1 | Testes unitários | 2 | 10 |
   | Alertas | Não entregue | Atraso na resposta | 5 | Falha de rede | 2 | Multi-canal | 2 | 20 |
   | Auth | Token vazado | Acesso não autorizado | 4 | Vulnerabilidade | 1 | JWT+HTTPS | 2 | 8 |

#### 4.1.2 Usabilidade IEC 62366-1:2015

**Entregáveis:**
1. [ ] **Especificação de Uso**:
   - Uso pretendido: Acompanhamento emocional de idosos
   - Perfil de usuários: Idosos (65+), cuidadores, profissionais de saúde
   - Ambiente de uso: Domiciliar, via smartphone/tablet

2. [ ] **Análise de Tarefas**:
   - Tarefas críticas para segurança
   - Erros de uso potenciais
   - Cenários de uso relacionados a risco

3. [ ] **Plano de Validação de Usabilidade**:
   - Estudo formativo (N≥5 por perfil)
   - Estudo somativo (N≥15 por perfil)
   - Métricas: sucesso da tarefa, erros, tempo, satisfação

4. [ ] **Relatório de Engenharia de Usabilidade**

#### 4.1.3 Dossiê Técnico ANVISA

**Estrutura requerida:**
```
1. Informações Gerais do Produto
   1.1 Nome comercial e técnico
   1.2 Modelo e versão
   1.3 Classificação de risco (Classe II)
   1.4 Regra de classificação aplicável
   1.5 Uso pretendido

2. Descrição do Produto
   2.1 Princípios de funcionamento
   2.2 Algoritmos utilizados (PHQ-9, GAD-7, C-SSRS)
   2.3 Arquitetura de software
   2.4 Integrações e interoperabilidade

3. Gestão de Riscos
   3.1 Arquivo de gestão de riscos ISO 14971
   3.2 FMEA
   3.3 Controles implementados

4. Verificação e Validação
   4.1 Plano de verificação
   4.2 Resultados de testes (unitários, integração, sistema)
   4.3 Validação clínica (se aplicável)

5. Usabilidade
   5.1 Arquivo de usabilidade IEC 62366-1
   5.2 Resultados de estudos de usabilidade

6. Segurança da Informação
   6.1 Controles de segurança
   6.2 Proteção de dados (LGPD)
   6.3 Trilha de auditoria

7. Rotulagem
   7.1 Instruções de uso
   7.2 Manual do usuário
   7.3 Informações de segurança
```

---

### Fase 2: Validação e Testes (4-6 semanas)

#### 4.2.1 Validação de Software IEC 62304

**Entregáveis:**
1. [ ] **Plano de Verificação e Validação**:
   - Estratégia de testes
   - Critérios de aceitação
   - Ambiente de teste

2. [ ] **Matriz de Rastreabilidade**:
   ```
   REQ-001 → SRS-001 → COD-001 → TEST-001
   REQ-002 → SRS-002 → COD-002 → TEST-002
   ...
   ```

3. [ ] **Testes de Sistema** para cenários críticos:
   - Paciente com risco suicida crítico → alerta entregue em <30s
   - Falha de canal primário → escalação funciona
   - Dados exportados são completos e corretos
   - Deleção de dados é efetiva

4. [ ] **Relatório de Testes** com:
   - 109 testes unitários existentes
   - Cobertura de código (meta: >80%)
   - Testes de integração
   - Testes de regressão

#### 4.2.2 Validação Clínica

**Protocolo sugerido:**
1. [ ] Estudo piloto (N=30-50 pacientes)
2. [ ] Comparação com avaliação profissional
3. [ ] Métricas: sensibilidade, especificidade, VPP, VPN
4. [ ] Aprovação por comitê de ética (CEP)

---

### Fase 3: Documentação de Suporte (2-3 semanas)

#### 4.3.1 Documentação LGPD

1. [ ] **RIPD (Relatório de Impacto à Proteção de Dados)**:
   - Descrição do tratamento
   - Necessidade e proporcionalidade
   - Riscos aos titulares
   - Medidas de mitigação

2. [ ] **Política de Privacidade** atualizada

3. [ ] **Termo de Consentimento Livre e Esclarecido** (TCLE)

#### 4.3.2 Manuais e Instruções

1. [ ] **Instruções de Uso (IFU)**:
   - Indicações de uso
   - Contraindicações
   - Advertências e precauções
   - Instruções de operação

2. [ ] **Manual do Usuário**:
   - Instalação
   - Operação
   - Resolução de problemas
   - Suporte técnico

---

## 5. CHECKLIST DE CONFORMIDADE

### 5.1 RDC 751/2022 - Anexo II, Capítulo 3

- [x] 3.1 - Identificação do produto
- [x] 3.2 - Descrição do produto
- [ ] 3.3 - Referência a normas aplicadas
- [x] 3.4 - Análise de risco (parcial)
- [ ] 3.5 - Verificação e validação
- [x] 3.6 - Biocompatibilidade (N/A - software)
- [x] 3.7 - Desempenho do produto
- [ ] 3.8 - Segurança elétrica (N/A)
- [ ] 3.9 - Proteção contra radiação (N/A)
- [x] 3.10 - Desempenho clínico (parcial)
- [x] 3.11 - Software (IEC 62304 parcial)
- [ ] 3.12 - Usabilidade (IEC 62366-1)
- [x] 3.13 - Rotulagem (parcial)
- [ ] 3.14 - Relatório de avaliação clínica

### 5.2 ISO 14971:2019

- [ ] 4.1 - Processo de gestão de risco
- [ ] 4.2 - Responsabilidades da alta direção
- [ ] 4.3 - Competência do pessoal
- [ ] 5.1 - Análise de risco
- [ ] 5.2 - Identificação de perigos
- [ ] 5.3 - Estimativa de risco
- [ ] 6 - Avaliação de risco
- [x] 7 - Controle de risco (implementado)
- [ ] 8 - Avaliação de risco residual geral
- [ ] 9 - Revisão da gestão de risco
- [ ] 10 - Atividades de produção e pós-produção

### 5.3 IEC 62366-1:2015

- [x] 5.1 - Preparar especificação de uso
- [ ] 5.2 - Identificar características relacionadas à segurança
- [ ] 5.3 - Identificar perigos e situações perigosas
- [ ] 5.4 - Selecionar tarefas para avaliação
- [ ] 5.5 - Elaborar especificação de interface
- [ ] 5.6 - Estabelecer plano de avaliação
- [ ] 5.7 - Realizar avaliação formativa
- [ ] 5.8 - Realizar avaliação somativa
- [ ] 5.9 - Documentar arquivo de engenharia de usabilidade

---

## 6. CRONOGRAMA SUGERIDO

```
Semana 1-2:   Arquivo de Gestão de Riscos ISO 14971
Semana 3-4:   FMEA + Matriz de Rastreabilidade
Semana 5-8:   Arquivo de Usabilidade IEC 62366-1
Semana 9-10:  Validação de Software / Testes de Sistema
Semana 11-12: RIPD + Documentação LGPD
Semana 13-14: Dossiê Técnico ANVISA (compilação)
Semana 15-16: Revisão final + Submissão
```

**Tempo total estimado:** 4-5 meses
**Esforço total estimado:** 400-600 horas

---

## 7. RECURSOS NECESSÁRIOS

### 7.1 Expertise Requerida

| Área | Perfil | Dedicação |
|------|--------|-----------|
| Assuntos Regulatórios | Especialista ANVISA/FDA | 40-60h |
| Gestão de Riscos | Engenheiro de qualidade | 60-80h |
| Usabilidade | Especialista UX/Human Factors | 80-100h |
| Validação Clínica | Médico/Psicólogo | 40-60h |
| Segurança da Informação | Especialista LGPD/ISO 27001 | 20-40h |

### 7.2 Estudos Clínicos

- Comitê de Ética (CEP) - aprovação necessária
- Amostra de pacientes (N≥30 para validação)
- Profissionais de saúde para comparação

---

## 8. CONCLUSÃO

### 8.1 Pontos Fortes do EVA-Mind-FZPN

1. **Escalas clínicas robustas** - PHQ-9, GAD-7, C-SSRS implementados corretamente com estratificação de risco
2. **Sistema de alertas maduro** - Multi-canal com escalação e timeout por prioridade
3. **Conformidade LGPD avançada** - Audit trail, direitos do titular, consentimento
4. **Monitoramento completo** - Prometheus + Grafana com métricas clínicas
5. **Arquitetura bem documentada** - Estrutura cognitiva clara (cortex/brainstem/hippocampus)
6. **Testes unitários** - 109 testes cobrindo funcionalidades críticas

### 8.2 Gaps Críticos para Certificação

1. ❌ **Gestão de Riscos ISO 14971** - Documento formal inexistente
2. ❌ **Usabilidade IEC 62366-1** - Arquivo de engenharia inexistente
3. ❌ **Validação Clínica** - Sem estudo formal
4. ❌ **Dossiê Técnico** - Não compilado
5. ⚠️ **Rastreabilidade** - Req→Código→Teste não documentada

### 8.3 Recomendação

O EVA-Mind-FZPN possui **fundamentos técnicos sólidos** para certificação ANVISA Classe II. O código implementa corretamente as funcionalidades clínicas e de segurança requeridas. No entanto, a **documentação regulatória formal** necessária para submissão ainda não existe.

**Próximo passo recomendado:** Iniciar pela **Gestão de Riscos ISO 14971**, pois ela é pré-requisito para todas as outras atividades de documentação.

---

## ANEXOS

### A. Referências Normativas
- RDC 751/2022 - ANVISA
- NBR ISO 14971:2019 - Gestão de riscos
- IEC 62366-1:2015 - Engenharia de usabilidade
- IEC 62304:2006/Amd1:2015 - Ciclo de vida de software
- Lei 13.709/2018 - LGPD
- RDC 185/2001 - Registro de produtos (histórico)

### B. Arquivos do Projeto Relevantes
- `internal/cortex/scales/clinical_scales.go` - Escalas clínicas
- `internal/cortex/alert/escalation.go` - Sistema de alertas
- `internal/audit/lgpd_audit.go` - Trilha de auditoria
- `internal/audit/data_rights.go` - Direitos do titular
- `internal/metrics/metrics.go` - Métricas Prometheus
- `migrations/018_lgpd_audit_trail.sql` - Schema LGPD

### C. Contatos de Referência
- ANVISA: https://www.gov.br/anvisa
- INMETRO: https://www.gov.br/inmetro
- ABIMED: https://abimed.org.br
