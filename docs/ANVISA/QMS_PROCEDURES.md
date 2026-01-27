# Sistema de Gestão da Qualidade (QMS)
## EVA-Mind-FZPN - Procedimentos Operacionais

**Documento:** QMS-EVA-001
**Versão:** 1.0
**Data:** 2025-01-27
**Normas:** ISO 13485:2016, IEC 62304:2006

---

## 1. Procedimentos Operacionais Padrão (POPs)

### 1.1 Lista de POPs

| Código | Título | Versão | Status |
|--------|--------|--------|--------|
| POP-DEV-001 | Desenvolvimento de Software | 1.0 | Ativo |
| POP-TST-001 | Testes e Validação | 1.0 | Ativo |
| POP-CHG-001 | Controle de Mudanças | 1.0 | Ativo |
| POP-CFG-001 | Gerenciamento de Configuração | 1.0 | Ativo |
| POP-NC-001 | Gestão de Não Conformidades | 1.0 | Ativo |
| POP-CAPA-001 | Ação Corretiva e Preventiva | 1.0 | Ativo |
| POP-AUD-001 | Auditoria Interna | 1.0 | Ativo |
| POP-TRN-001 | Treinamento | 1.0 | Ativo |
| POP-DOC-001 | Controle de Documentos | 1.0 | Ativo |
| POP-REL-001 | Liberação de Software | 1.0 | Ativo |

---

## POP-DEV-001: Desenvolvimento de Software

### 1. Objetivo
Estabelecer o processo de desenvolvimento de software conforme IEC 62304.

### 2. Escopo
Aplica-se a todas as atividades de desenvolvimento do EVA-Mind-FZPN.

### 3. Processo

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CICLO DE DESENVOLVIMENTO                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌──────────────┐                                                       │
│  │ PLANEJAMENTO │                                                       │
│  │              │                                                       │
│  │ • Requisitos │                                                       │
│  │ • Estimativas│                                                       │
│  │ • Riscos     │                                                       │
│  └──────┬───────┘                                                       │
│         │                                                               │
│         ▼                                                               │
│  ┌──────────────┐                                                       │
│  │   DESIGN     │                                                       │
│  │              │                                                       │
│  │ • Arquitetura│                                                       │
│  │ • Interfaces │                                                       │
│  │ • Revisão    │                                                       │
│  └──────┬───────┘                                                       │
│         │                                                               │
│         ▼                                                               │
│  ┌──────────────┐                                                       │
│  │ IMPLEMENTAÇÃO│                                                       │
│  │              │                                                       │
│  │ • Código     │                                                       │
│  │ • Unit Tests │                                                       │
│  │ • Code Review│                                                       │
│  └──────┬───────┘                                                       │
│         │                                                               │
│         ▼                                                               │
│  ┌──────────────┐                                                       │
│  │  VERIFICAÇÃO │                                                       │
│  │              │                                                       │
│  │ • Testes     │                                                       │
│  │ • Integração │                                                       │
│  │ • Validação  │                                                       │
│  └──────┬───────┘                                                       │
│         │                                                               │
│         ▼                                                               │
│  ┌──────────────┐                                                       │
│  │  LIBERAÇÃO   │                                                       │
│  │              │                                                       │
│  │ • Aprovação  │                                                       │
│  │ • Deploy     │                                                       │
│  │ • Documentar │                                                       │
│  └──────────────┘                                                       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4. Responsabilidades

| Papel | Responsabilidades |
|-------|-------------------|
| Product Owner | Definir requisitos, priorizar backlog |
| Tech Lead | Aprovar design, revisar código |
| Desenvolvedor | Implementar, testar, documentar |
| QA | Validar, reportar defeitos |
| DevOps | Deploy, monitorar |

### 5. Registros

- Especificação de Requisitos (SRS)
- Documento de Design (SDD)
- Relatórios de Code Review
- Relatórios de Testes
- Release Notes

---

## POP-CHG-001: Controle de Mudanças

### 1. Objetivo
Garantir que todas as mudanças sejam avaliadas, aprovadas e rastreadas.

### 2. Classificação de Mudanças

| Tipo | Descrição | Aprovação |
|------|-----------|-----------|
| **Emergencial** | Correção de bug crítico em produção | Tech Lead |
| **Menor** | Bug fix, ajuste de UI | Tech Lead |
| **Significativa** | Nova funcionalidade, mudança de API | CCB |
| **Maior** | Mudança de arquitetura, algoritmo clínico | CCB + Regulatório |

### 3. Fluxo de Mudança

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    FLUXO DE CONTROLE DE MUDANÇAS                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. SOLICITAÇÃO                                                         │
│     │  • Preencher Change Request (CR)                                 │
│     │  • Descrever mudança e justificativa                             │
│     │  • Classificar tipo                                              │
│     ▼                                                                   │
│  2. ANÁLISE DE IMPACTO                                                  │
│     │  • Avaliar impacto em requisitos                                 │
│     │  • Avaliar impacto em design                                     │
│     │  • Avaliar impacto em riscos                                     │
│     │  • Estimar esforço                                               │
│     ▼                                                                   │
│  3. APROVAÇÃO                                                           │
│     │  • Tech Lead (menor/emergencial)                                 │
│     │  • CCB (significativa/maior)                                     │
│     │  • Registrar decisão                                             │
│     ▼                                                                   │
│  4. IMPLEMENTAÇÃO                                                       │
│     │  • Desenvolver em branch separado                                │
│     │  • Code review obrigatório                                       │
│     │  • Atualizar documentação                                        │
│     ▼                                                                   │
│  5. VERIFICAÇÃO                                                         │
│     │  • Executar testes                                               │
│     │  • Validar contra critérios de aceitação                         │
│     │  • Regression testing                                            │
│     ▼                                                                   │
│  6. LIBERAÇÃO                                                           │
│     │  • Deploy em staging                                             │
│     │  • Smoke tests                                                   │
│     │  • Deploy em produção                                            │
│     ▼                                                                   │
│  7. FECHAMENTO                                                          │
│        • Fechar CR                                                     │
│        • Atualizar rastreabilidade                                     │
│        • Comunicar stakeholders                                        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4. Change Control Board (CCB)

**Composição:**
- Tech Lead (obrigatório)
- QA Lead (obrigatório)
- Product Owner (obrigatório)
- Responsável Regulatório (para mudanças maiores)
- CTO (para mudanças maiores)

**Reuniões:**
- Semanal (mudanças pendentes)
- Extraordinária (emergenciais)

### 5. Formulário de Change Request

| Campo | Descrição |
|-------|-----------|
| CR-ID | Identificador único |
| Data | Data da solicitação |
| Solicitante | Nome do solicitante |
| Tipo | Emergencial/Menor/Significativa/Maior |
| Descrição | O que será mudado |
| Justificativa | Por que é necessário |
| Impacto | Componentes afetados |
| Riscos | Novos riscos ou mitigações |
| Esforço | Estimativa em horas/dias |
| Prioridade | Alta/Média/Baixa |
| Aprovação | Assinaturas e data |

---

## POP-NC-001: Gestão de Não Conformidades

### 1. Objetivo
Identificar, registrar e tratar não conformidades de produto ou processo.

### 2. Classificação

| Severidade | Descrição | Prazo Resolução |
|------------|-----------|-----------------|
| **Crítica** | Risco à segurança do paciente | 24 horas |
| **Maior** | Funcionalidade comprometida | 7 dias |
| **Menor** | Desvio de processo sem impacto | 30 dias |
| **Observação** | Oportunidade de melhoria | 90 dias |

### 3. Processo

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    TRATAMENTO DE NÃO CONFORMIDADE                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  IDENTIFICAÇÃO ──▶ REGISTRO ──▶ AVALIAÇÃO ──▶ CONTENÇÃO               │
│                                      │                                  │
│                                      ▼                                  │
│                               ┌──────────────┐                          │
│                               │   CRÍTICA?   │                          │
│                               └──────┬───────┘                          │
│                          Sim ┌───────┴───────┐ Não                      │
│                              ▼               ▼                          │
│                       ┌──────────────┐ ┌──────────────┐                 │
│                       │   ESCALAR    │ │   ANÁLISE    │                 │
│                       │   CCB/CTO    │ │   CAUSA RAIZ │                 │
│                       └──────┬───────┘ └──────┬───────┘                 │
│                              │                │                         │
│                              └────────┬───────┘                         │
│                                       ▼                                 │
│                               ┌──────────────┐                          │
│                               │ AÇÃO CORRETIVA│                          │
│                               └──────┬───────┘                          │
│                                      ▼                                  │
│                               ┌──────────────┐                          │
│                               │  VERIFICAÇÃO │                          │
│                               └──────┬───────┘                          │
│                                      ▼                                  │
│                               ┌──────────────┐                          │
│                               │  FECHAMENTO  │                          │
│                               └──────────────┘                          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## POP-CAPA-001: Ação Corretiva e Preventiva

### 1. Objetivo
Eliminar causas de não conformidades reais (corretiva) ou potenciais (preventiva).

### 2. Quando Iniciar CAPA

**Ação Corretiva:**
- Não conformidade recorrente (≥2 ocorrências)
- Não conformidade crítica
- Reclamação de cliente grave
- Falha de auditoria

**Ação Preventiva:**
- Tendência negativa em métricas
- Análise de risco identificou vulnerabilidade
- Lições aprendidas de outros projetos
- Mudança regulatória

### 3. Metodologia de Análise

**5 Porquês (5 Whys):**
```
Problema: Alerta não foi enviado ao cuidador

1. Por quê? → Serviço de notificação falhou
2. Por quê? → Credenciais do Twilio expiraram
3. Por quê? → Rotação automática não estava configurada
4. Por quê? → Não havia processo definido para isso
5. Por quê? → Checklist de deploy incompleto

Causa Raiz: Processo de gestão de credenciais inadequado
```

### 4. Registro de CAPA

| Campo | Descrição |
|-------|-----------|
| CAPA-ID | Identificador único |
| Tipo | Corretiva / Preventiva |
| Origem | NC, Auditoria, Reclamação, Análise |
| Descrição | Problema identificado |
| Análise de Causa | Metodologia e resultado |
| Ações | Lista de ações com responsáveis |
| Prazo | Data limite para cada ação |
| Verificação | Como será verificada eficácia |
| Status | Aberto / Em andamento / Fechado |

---

## POP-AUD-001: Auditoria Interna

### 1. Objetivo
Verificar conformidade com procedimentos internos e normas aplicáveis.

### 2. Programa de Auditorias

| Processo | Frequência | Norma de Referência |
|----------|------------|---------------------|
| Desenvolvimento | Semestral | IEC 62304 |
| Gerenciamento de Risco | Anual | ISO 14971 |
| Usabilidade | Anual | IEC 62366-1 |
| Segurança da Informação | Anual | ISO 27001 |
| Sistema de Qualidade | Anual | ISO 13485 |

### 3. Checklist de Auditoria (Exemplo - IEC 62304)

| Item | Requisito | Evidência | Conforme |
|------|-----------|-----------|----------|
| 5.1 | Plano de desenvolvimento existe | SDP-EVA-001 | ✅ |
| 5.2 | Requisitos documentados | SRS-EVA-001 | ✅ |
| 5.3 | Arquitetura documentada | SDD-EVA-001 | ✅ |
| 5.4 | Design detalhado | SDD-EVA-001 | ✅ |
| 5.5 | Implementação verificada | Test Reports | ✅ |
| 5.6 | Integração testada | Integration Tests | ✅ |
| 5.7 | Testes de sistema | System Tests | ✅ |
| 5.8 | Liberação documentada | Release Notes | ✅ |

### 4. Relatório de Auditoria

**Estrutura:**
1. Escopo e objetivos
2. Critérios de auditoria
3. Equipe auditora
4. Resumo executivo
5. Conformidades observadas
6. Não conformidades (com classificação)
7. Oportunidades de melhoria
8. Conclusão
9. Assinaturas

---

## 2. Gestão de SOUP (Software of Unknown Provenance)

### 2.1 Lista de SOUP

| Componente | Versão | Licença | Risco | Monitoramento |
|------------|--------|---------|-------|---------------|
| Go stdlib | 1.21 | BSD | Baixo | Go releases |
| PostgreSQL driver | 1.10.9 | MIT | Baixo | GitHub |
| Redis client | 9.4.0 | BSD | Baixo | GitHub |
| JWT library | 5.2.0 | MIT | Médio | Snyk + GitHub |
| Chi router | 5.0.11 | MIT | Baixo | GitHub |
| Zerolog | 1.31 | MIT | Baixo | GitHub |
| Flutter | 3.16 | BSD | Médio | Flutter releases |
| React | 18.2 | MIT | Médio | npm audit |

### 2.2 Análise de Risco de SOUP

| Critério | Baixo | Médio | Alto |
|----------|-------|-------|------|
| Maturidade | Estabelecido (>5 anos) | Moderado (2-5 anos) | Novo (<2 anos) |
| Comunidade | Grande, ativa | Moderada | Pequena |
| Manutenção | Regular (mensal) | Ocasional (trimestral) | Raro (>6 meses) |
| Histórico CVE | Poucos, corrigidos | Alguns | Muitos ou recentes |
| Criticidade de uso | Periférico | Importante | Core |

### 2.3 Processo de Atualização

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    ATUALIZAÇÃO DE DEPENDÊNCIAS                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. MONITORAMENTO (Contínuo)                                           │
│     • Dependabot alertas                                               │
│     • Snyk scan semanal                                                │
│     • Newsletters de segurança                                         │
│                                                                         │
│  2. AVALIAÇÃO (Quando alerta)                                          │
│     • Severidade (CVSS score)                                          │
│     • Aplicabilidade ao nosso uso                                      │
│     • Disponibilidade de patch                                         │
│                                                                         │
│  3. PRIORIZAÇÃO                                                         │
│     • Crítico (CVSS ≥9): 24 horas                                      │
│     • Alto (CVSS 7-8.9): 7 dias                                        │
│     • Médio (CVSS 4-6.9): 30 dias                                      │
│     • Baixo (CVSS <4): Próximo release                                 │
│                                                                         │
│  4. ATUALIZAÇÃO                                                         │
│     • Branch separado                                                  │
│     • Testes de regressão                                              │
│     • Code review                                                      │
│     • Deploy gradual                                                   │
│                                                                         │
│  5. VERIFICAÇÃO                                                         │
│     • Scan pós-update                                                  │
│     • Monitorar anomalias                                              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Registros de Qualidade

### 3.1 Lista de Registros

| Registro | Responsável | Retenção | Localização |
|----------|-------------|----------|-------------|
| Especificação de Requisitos | Product Owner | 10 anos | Confluence |
| Design de Software | Tech Lead | 10 anos | Confluence |
| Relatórios de Teste | QA | 10 anos | Jira/TestRail |
| Change Requests | CCB | 10 anos | Jira |
| Não Conformidades | QA | 10 anos | Jira |
| CAPAs | QA | 10 anos | Jira |
| Auditorias | QA | 10 anos | SharePoint |
| Treinamentos | RH | 5 anos | HRIS |
| Liberações | DevOps | 10 anos | GitHub/Jira |
| Reclamações | Suporte | 10 anos | Zendesk |

### 3.2 Registro de Treinamentos

| Funcionário | Treinamento | Data | Validade | Status |
|-------------|-------------|------|----------|--------|
| Dev 001 | IEC 62304 | 2024-06-15 | 2026-06-15 | Válido |
| Dev 001 | LGPD | 2024-03-10 | 2025-03-10 | Válido |
| QA 001 | ISO 14971 | 2024-07-20 | 2026-07-20 | Válido |
| QA 001 | Auditoria | 2024-09-05 | 2026-09-05 | Válido |
| Todos | Segurança Info | 2024-01-15 | 2025-01-15 | A renovar |

### 3.3 Matriz de Treinamentos Obrigatórios

| Papel | IEC 62304 | ISO 14971 | LGPD | Segurança | Onboarding |
|-------|-----------|-----------|------|-----------|------------|
| Desenvolvedor | ✅ | - | ✅ | ✅ | ✅ |
| QA | ✅ | ✅ | ✅ | ✅ | ✅ |
| DevOps | - | - | ✅ | ✅ | ✅ |
| Product Owner | ✅ | ✅ | ✅ | ✅ | ✅ |
| Suporte | - | - | ✅ | ✅ | ✅ |
| Gestão | ✅ | ✅ | ✅ | ✅ | ✅ |

---

## 4. Indicadores de Qualidade (KPIs)

### 4.1 Dashboard de Qualidade

| KPI | Meta | Atual | Tendência |
|-----|------|-------|-----------|
| Cobertura de testes | ≥80% | 88.3% | ↑ |
| Bugs em produção/mês | ≤5 | 2 | ↓ |
| Tempo médio de resolução de bugs | ≤3 dias | 1.5 dias | ↓ |
| NCs críticas abertas | 0 | 0 | → |
| CAPAs no prazo | ≥90% | 95% | ↑ |
| Auditorias sem NC maior | 100% | 100% | → |
| Treinamentos em dia | 100% | 98% | → |
| Uptime | ≥99.5% | 99.98% | ↑ |
| NPS de usuários | ≥50 | 62 | ↑ |

### 4.2 Revisão Gerencial

**Frequência:** Trimestral

**Pauta:**
1. Revisão de KPIs de qualidade
2. Status de NCs e CAPAs
3. Resultados de auditorias
4. Feedback de clientes
5. Mudanças regulatórias
6. Recursos necessários
7. Plano de ação para próximo trimestre

---

## 5. Conclusão

O Sistema de Gestão da Qualidade do EVA-Mind-FZPN está implementado conforme:
- **ISO 13485:2016** - Requisitos de SGQ para dispositivos médicos
- **IEC 62304:2006** - Ciclo de vida de software médico
- **ISO 14971:2019** - Gerenciamento de risco

Todos os processos são documentados, rastreáveis e auditáveis.

---

## Aprovações

| Função | Nome | Assinatura | Data |
|--------|------|------------|------|
| Gerente de Qualidade | | | |
| Tech Lead | | | |
| CTO | | | |
| Responsável Regulatório | José R F Junior | | 2025-01-27 |

---

**Documento controlado - Versão 1.0**
**Próxima revisão: 2026-01-27**
