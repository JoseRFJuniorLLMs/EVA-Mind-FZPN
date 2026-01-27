# 🚀 EVA-Mind-FZPN - SPRINT PLAN 2025

## Visão Geral do Projeto

**Objetivo:** Elevar a EVA de "funcional" para "produção enterprise-ready"

**Duração Total:** 6 Sprints (6 semanas)

**Participantes:**
- Jose R F Junior (Arquiteto/Dev)
- EVA (Claude) - Pair Programming

---

## 📋 ÉPICOS

| # | Épico | Sprints | Prioridade |
|---|-------|---------|------------|
| E1 | Testes para Ferramentas Críticas | Sprint 1 | 🔴 CRÍTICA |
| E2 | Métricas e Observabilidade (Prometheus) | Sprint 2 | 🟠 ALTA |
| E3 | Auditoria LGPD | Sprint 3 | 🟠 ALTA |
| E4 | Memória de Longo Prazo Avançada | Sprint 4 | 🟡 MÉDIA |
| E5 | Autoconsciência e Meta-cognição | Sprint 5 | 🟡 MÉDIA |
| E6 | Sistema de Aprendizado Contínuo | Sprint 6 | 🟢 DESEJÁVEL |

---

# 🔴 SPRINT 1: TESTES CRÍTICOS
**Duração:** 1 semana
**Épico:** E1 - Testes para Ferramentas Críticas

## Objetivo
Garantir que ferramentas de vida-ou-morte funcionem 100% do tempo.

## User Stories

### US-1.1: Testes do C-SSRS (Risco Suicida)
**Como** sistema de saúde mental
**Quero** garantir que a escala C-SSRS nunca falhe silenciosamente
**Para** proteger vidas

**Critérios de Aceite:**
- [ ] Teste unitário para `apply_cssrs` handler
- [ ] Teste de integração: fluxo completo de aplicação
- [ ] Teste de edge cases: respostas inválidas, timeout, erro de DB
- [ ] Teste de alerta: verifica se notificação é enviada em risco positivo
- [ ] Cobertura mínima: 90%

**Arquivos a criar:**
```
internal/tools/handlers_test.go
internal/tools/cssrs_test.go
internal/cortex/scales/cssrs_test.go
```

### US-1.2: Testes do Sistema de Alertas
**Como** cuidador
**Quero** garantir que alertas sempre cheguem
**Para** responder emergências a tempo

**Critérios de Aceite:**
- [ ] Teste unitário para cada canal (Push, SMS, Email)
- [ ] Teste de fallback chain: Push falha → SMS → Email
- [ ] Teste de retry logic
- [ ] Teste de escalação por severidade
- [ ] Mock de serviços externos (Firebase, Twilio)

**Arquivos a criar:**
```
internal/brainstem/push/firebase_test.go
internal/motor/sms/twilio_test.go
internal/motor/email/smtp_test.go
internal/cortex/alert/escalation_test.go
```

### US-1.3: Testes do PHQ-9 e GAD-7
**Como** profissional de saúde
**Quero** garantir que escalas psicológicas calculem scores corretamente
**Para** não dar diagnósticos errados

**Critérios de Aceite:**
- [ ] Teste de cálculo de score (0-27 para PHQ-9, 0-21 para GAD-7)
- [ ] Teste de categorização (mínimo, leve, moderado, grave)
- [ ] Teste de persistência dos resultados
- [ ] Teste de fluxo conversacional

**Arquivos a criar:**
```
internal/cortex/scales/phq9_test.go
internal/cortex/scales/gad7_test.go
internal/tools/assessment_test.go
```

### US-1.4: Testes de Medicação Visual
**Como** idoso
**Quero** que a identificação de medicamentos seja precisa
**Para** não tomar remédio errado

**Critérios de Aceite:**
- [ ] Teste de parsing de resposta do Gemini Vision
- [ ] Teste de matching com banco de medicamentos
- [ ] Teste de falha graceful quando câmera não disponível

**Arquivos a criar:**
```
internal/motor/vision/medication_identifier_test.go
```

## Tarefas Técnicas

| ID | Tarefa | Estimativa | Responsável |
|----|--------|------------|-------------|
| T1.1 | Setup de framework de testes (testify) | 2h | Dev |
| T1.2 | Criar mocks para Firebase/Twilio | 4h | Dev |
| T1.3 | Implementar testes C-SSRS | 8h | Dev |
| T1.4 | Implementar testes de alertas | 8h | Dev |
| T1.5 | Implementar testes PHQ-9/GAD-7 | 4h | Dev |
| T1.6 | Implementar testes de medicação | 4h | Dev |
| T1.7 | Configurar CI (GitHub Actions) | 4h | Dev |
| T1.8 | Documentar cobertura de testes | 2h | Dev |

**Total Sprint 1:** ~36h (1 semana)

---

# 🟠 SPRINT 2: MÉTRICAS E OBSERVABILIDADE
**Duração:** 1 semana
**Épico:** E2 - Prometheus + Grafana

## Objetivo
Saber exatamente o que está acontecendo em produção, em tempo real.

## User Stories

### US-2.1: Métricas de Sistema
**Como** operador
**Quero** ver métricas de saúde do sistema
**Para** detectar problemas antes dos usuários

**Métricas a implementar:**
- [ ] `eva_requests_total` - Total de requests por endpoint
- [ ] `eva_request_duration_seconds` - Latência por endpoint
- [ ] `eva_active_sessions` - Sessões WebSocket ativas
- [ ] `eva_errors_total` - Erros por tipo e severidade
- [ ] `eva_db_connections` - Pool de conexões DB

### US-2.2: Métricas de Negócio
**Como** gestor de produto
**Quero** ver métricas de uso
**Para** entender como a EVA é utilizada

**Métricas a implementar:**
- [ ] `eva_conversations_total` - Conversas por dia/usuário
- [ ] `eva_tool_invocations_total` - Uso de cada ferramenta
- [ ] `eva_alerts_sent_total` - Alertas por severidade
- [ ] `eva_memory_operations_total` - Operações de memória
- [ ] `eva_llm_tokens_total` - Tokens consumidos do Gemini

### US-2.3: Métricas de Saúde Mental
**Como** profissional clínico
**Quero** ver tendências de saúde mental
**Para** identificar pacientes em risco

**Métricas a implementar:**
- [ ] `eva_phq9_scores` - Histogram de scores PHQ-9
- [ ] `eva_gad7_scores` - Histogram de scores GAD-7
- [ ] `eva_cssrs_triggers` - Gatilhos de risco suicida
- [ ] `eva_emotion_distribution` - Distribuição de emoções detectadas

### US-2.4: Dashboard Grafana
**Como** operador
**Quero** um dashboard visual
**Para** monitorar tudo em um lugar

**Painéis a criar:**
- [ ] Visão Geral (requests, erros, latência)
- [ ] Sessões Ativas (mapa, contagem)
- [ ] Ferramentas (uso por tipo)
- [ ] Alertas (enviados, falhas)
- [ ] Saúde Mental (tendências)

## Tarefas Técnicas

| ID | Tarefa | Estimativa | Responsável |
|----|--------|------------|-------------|
| T2.1 | Adicionar prometheus client ao Go | 2h | Dev |
| T2.2 | Criar pacote internal/metrics | 4h | Dev |
| T2.3 | Instrumentar endpoints HTTP | 4h | Dev |
| T2.4 | Instrumentar WebSocket | 4h | Dev |
| T2.5 | Instrumentar ferramentas | 4h | Dev |
| T2.6 | Instrumentar escalas psicológicas | 2h | Dev |
| T2.7 | Configurar Prometheus server | 2h | DevOps |
| T2.8 | Criar dashboards Grafana | 8h | Dev |
| T2.9 | Configurar alertas (AlertManager) | 4h | DevOps |
| T2.10 | Documentar métricas | 2h | Dev |

**Total Sprint 2:** ~36h (1 semana)

---

# 🟠 SPRINT 3: AUDITORIA LGPD
**Duração:** 1 semana
**Épico:** E3 - Compliance e Trilha de Auditoria

## Objetivo
Estar 100% em conformidade com LGPD/GDPR para dados de saúde.

## User Stories

### US-3.1: Trilha de Auditoria Universal
**Como** DPO (Data Protection Officer)
**Quero** log de todas as operações com dados pessoais
**Para** responder a auditorias e incidentes

**Critérios de Aceite:**
- [ ] Toda leitura de dados pessoais é logada
- [ ] Toda escrita/modificação é logada
- [ ] Toda exclusão é logada
- [ ] Logs incluem: quem, quando, o quê, de onde (IP)
- [ ] Logs são imutáveis (append-only)

**Tabela a criar:**
```sql
CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    actor_type VARCHAR(20) NOT NULL, -- user, system, admin
    actor_id VARCHAR(100),
    action VARCHAR(50) NOT NULL, -- read, write, delete, export
    resource_type VARCHAR(50) NOT NULL, -- memory, vital, alert
    resource_id VARCHAR(100),
    idoso_id INTEGER,
    ip_address INET,
    user_agent TEXT,
    details JSONB,
    -- Imutabilidade
    hash_chain VARCHAR(64) -- SHA256 do registro anterior
);
```

### US-3.2: Direito ao Esquecimento
**Como** titular dos dados
**Quero** poder solicitar exclusão dos meus dados
**Para** exercer meu direito LGPD

**Critérios de Aceite:**
- [ ] Endpoint para solicitar exclusão
- [ ] Processo de verificação de identidade
- [ ] Exclusão em cascata (todas as tabelas)
- [ ] Anonimização de dados para pesquisa (opcional)
- [ ] Confirmação por email
- [ ] Prazo máximo: 15 dias

**Arquivos a criar:**
```
internal/security/gdpr/deletion_service.go
internal/security/gdpr/anonymization.go
```

### US-3.3: Exportação de Dados (Portabilidade)
**Como** titular dos dados
**Quero** exportar todos os meus dados
**Para** exercer direito de portabilidade

**Critérios de Aceite:**
- [ ] Endpoint para solicitar exportação
- [ ] Formato: JSON estruturado + PDF legível
- [ ] Inclui: memórias, vitais, alertas, histórico
- [ ] Link de download temporário (24h)
- [ ] Notificação por email

**Arquivos a criar:**
```
internal/security/gdpr/export_service.go
internal/security/gdpr/pdf_generator.go
```

### US-3.4: Consentimento Granular
**Como** titular dos dados
**Quero** controlar quais dados compartilho
**Para** ter autonomia sobre minha privacidade

**Critérios de Aceite:**
- [ ] Tela de consentimentos no app
- [ ] Categorias: memória, vitais, localização, voz
- [ ] Histórico de alterações de consentimento
- [ ] Revogação a qualquer momento

**Tabela a criar:**
```sql
CREATE TABLE consent_records (
    id SERIAL PRIMARY KEY,
    idoso_id INTEGER NOT NULL,
    consent_type VARCHAR(50) NOT NULL,
    granted BOOLEAN NOT NULL,
    granted_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    ip_address INET,
    version INTEGER DEFAULT 1
);
```

### US-3.5: Retenção e Expiração de Dados
**Como** sistema
**Quero** deletar dados antigos automaticamente
**Para** minimizar riscos de vazamento

**Critérios de Aceite:**
- [ ] Política de retenção configurável por tipo de dado
- [ ] Worker de limpeza automática
- [ ] Dados de saúde: 5 anos (regulatório)
- [ ] Logs de auditoria: 10 anos
- [ ] Sessões: 30 dias
- [ ] Notificação antes de exclusão

## Tarefas Técnicas

| ID | Tarefa | Estimativa | Responsável |
|----|--------|------------|-------------|
| T3.1 | Criar tabela audit_log | 2h | Dev |
| T3.2 | Implementar middleware de auditoria | 4h | Dev |
| T3.3 | Instrumentar todas as operações CRUD | 8h | Dev |
| T3.4 | Implementar hash chain para imutabilidade | 2h | Dev |
| T3.5 | Implementar deletion_service | 4h | Dev |
| T3.6 | Implementar export_service | 4h | Dev |
| T3.7 | Criar tabela consent_records | 2h | Dev |
| T3.8 | Implementar API de consentimento | 4h | Dev |
| T3.9 | Criar worker de retenção | 4h | Dev |
| T3.10 | Documentar políticas LGPD | 2h | Dev |

**Total Sprint 3:** ~36h (1 semana)

---

# 🟡 SPRINT 4: MEMÓRIA DE LONGO PRAZO AVANÇADA
**Duração:** 1 semana
**Épico:** E4 - Padrões Temporais e Insights

## Objetivo
EVA entende padrões ao longo de meses/anos, não só conversas recentes.

## User Stories

### US-4.1: Detecção de Padrões Temporais
**Como** EVA
**Quero** detectar padrões que se repetem no tempo
**Para** antecipar necessidades do paciente

**Padrões a detectar:**
- [ ] Sazonais (Natal = saudade, Inverno = tristeza)
- [ ] Semanais (Domingo = solidão)
- [ ] Diários (Noite = ansiedade)
- [ ] Climáticos (Chuva = melancolia)
- [ ] Datas significativas (aniversário de morte)

**Arquivos a criar:**
```
internal/hippocampus/memory/temporal_patterns.go
internal/hippocampus/memory/pattern_detector.go
```

### US-4.2: Correlações Causais
**Como** EVA
**Quero** entender o que causa o quê
**Para** fazer intervenções preventivas

**Exemplos:**
- "Quando não toma medicação → fica irritado no dia seguinte"
- "Quando fala com a filha → fica feliz por 2 dias"
- "Quando não dorme bem → mais propenso a quedas"

**Arquivos a criar:**
```
internal/hippocampus/memory/causal_inference.go
```

### US-4.3: Linha do Tempo de Vida
**Como** profissional
**Quero** ver a história de vida do paciente
**Para** entender contexto das questões atuais

**Critérios de Aceite:**
- [ ] Timeline visual de eventos importantes
- [ ] Extração automática de marcos (nascimento, casamento, perdas)
- [ ] Conexão com padrões emocionais atuais

**Arquivos a criar:**
```
internal/hippocampus/memory/life_timeline.go
internal/hippocampus/memory/milestone_extractor.go
```

### US-4.4: Previsão de Estado Emocional
**Como** EVA
**Quero** prever como o paciente vai se sentir
**Para** preparar intervenções proativas

**Critérios de Aceite:**
- [ ] Modelo preditivo baseado em histórico
- [ ] Alerta para dias de risco previsto
- [ ] Sugestão de intervenções preventivas

**Arquivos a criar:**
```
internal/cortex/prediction/emotional_forecast.go
```

## Tarefas Técnicas

| ID | Tarefa | Estimativa | Responsável |
|----|--------|------------|-------------|
| T4.1 | Criar schema para padrões temporais | 2h | Dev |
| T4.2 | Implementar detector de padrões semanais | 4h | Dev |
| T4.3 | Implementar detector de padrões sazonais | 4h | Dev |
| T4.4 | Implementar inferência causal básica | 8h | Dev |
| T4.5 | Criar extrator de milestones | 4h | Dev |
| T4.6 | Implementar timeline de vida | 4h | Dev |
| T4.7 | Criar modelo de previsão emocional | 8h | Dev |
| T4.8 | Integrar com UnifiedRetrieval | 2h | Dev |

**Total Sprint 4:** ~36h (1 semana)

---

# 🟡 SPRINT 5: AUTOCONSCIÊNCIA E META-COGNIÇÃO
**Duração:** 1 semana
**Épico:** E5 - EVA sabe o que sabe (e o que não sabe)

## Objetivo
EVA tem consciência de suas limitações e estados internos.

## User Stories

### US-5.1: Detecção de Incerteza
**Como** EVA
**Quero** saber quando estou incerta
**Para** não dar respostas falsamente confiantes

**Critérios de Aceite:**
- [ ] Score de confiança em cada resposta
- [ ] Detecção de informações conflitantes
- [ ] Expressão verbal de incerteza ("não tenho certeza, mas...")
- [ ] Pedido de clarificação quando necessário

**Arquivos a criar:**
```
internal/cortex/metacognition/uncertainty_detector.go
internal/cortex/metacognition/confidence_scorer.go
```

### US-5.2: Estado Interno da EVA
**Como** desenvolvedor
**Quero** que EVA reporte seu estado interno
**Para** debuggar e melhorar o sistema

**Estados a rastrear:**
- [ ] Carga cognitiva (muitas informações simultâneas)
- [ ] Confusão (informações contraditórias)
- [ ] Preocupação (detectou risco mas não tem certeza)
- [ ] Limitação (não sabe responder algo)

**Arquivos a criar:**
```
internal/cortex/metacognition/internal_state.go
internal/cortex/metacognition/state_reporter.go
```

### US-5.3: Escalação Automática para Humano
**Como** sistema de segurança
**Quero** escalar automaticamente para humano quando necessário
**Para** não deixar situações críticas sem supervisão

**Gatilhos de escalação:**
- [ ] Risco suicida detectado
- [ ] Incerteza alta em situação médica
- [ ] Paciente pede explicitamente humano
- [ ] Detecção de abuso/violência
- [ ] 3+ tentativas sem resolução

**Arquivos a criar:**
```
internal/cortex/metacognition/human_escalation.go
```

### US-5.4: Auto-avaliação de Qualidade
**Como** EVA
**Quero** avaliar a qualidade das minhas respostas
**Para** melhorar continuamente

**Critérios de Aceite:**
- [ ] Score de relevância (respondi o que foi perguntado?)
- [ ] Score de empatia (fui acolhedora?)
- [ ] Score de segurança (não causei dano?)
- [ ] Log para análise posterior

**Arquivos a criar:**
```
internal/cortex/metacognition/self_evaluation.go
```

## Tarefas Técnicas

| ID | Tarefa | Estimativa | Responsável |
|----|--------|------------|-------------|
| T5.1 | Implementar detector de incerteza | 4h | Dev |
| T5.2 | Criar scorer de confiança | 4h | Dev |
| T5.3 | Implementar rastreador de estado interno | 4h | Dev |
| T5.4 | Criar reporter de estado | 2h | Dev |
| T5.5 | Implementar sistema de escalação | 8h | Dev |
| T5.6 | Criar gatilhos de escalação | 4h | Dev |
| T5.7 | Implementar auto-avaliação | 4h | Dev |
| T5.8 | Integrar com sistema de alertas | 4h | Dev |
| T5.9 | Criar dashboard de meta-cognição | 2h | Dev |

**Total Sprint 5:** ~36h (1 semana)

---

# 🟢 SPRINT 6: APRENDIZADO CONTÍNUO
**Duração:** 1 semana
**Épico:** E6 - EVA aprende e se adapta

## Objetivo
EVA melhora com o tempo baseado em feedback e experiência.

## User Stories

### US-6.1: Coleta de Feedback
**Como** usuário
**Quero** dar feedback sobre as respostas da EVA
**Para** ajudá-la a melhorar

**Critérios de Aceite:**
- [ ] Botões de like/dislike em respostas
- [ ] Opção de feedback textual
- [ ] Pergunta periódica "estou ajudando?"
- [ ] Armazenamento estruturado de feedback

**Arquivos a criar:**
```
internal/hippocampus/learning/feedback_collector.go
```

### US-6.2: Adaptação de Estilo
**Como** EVA
**Quero** adaptar meu estilo de comunicação
**Para** me conectar melhor com cada pessoa

**Adaptações:**
- [ ] Tom (mais formal/informal)
- [ ] Tamanho das respostas (curto/detalhado)
- [ ] Uso de metáforas (sim/não)
- [ ] Velocidade da conversa

**Arquivos a criar:**
```
internal/cortex/personality/style_adapter.go
```

### US-6.3: Aprendizado de Preferências
**Como** EVA
**Quero** lembrar o que funciona com cada pessoa
**Para** não repetir erros

**Exemplos:**
- "Maria não gosta quando falo de exercícios"
- "João prefere histórias do Nasrudin"
- "Ana se acalma com exercícios de respiração"

**Arquivos a criar:**
```
internal/hippocampus/learning/preference_learner.go
```

### US-6.4: Experimentação A/B
**Como** desenvolvedor
**Quero** testar diferentes abordagens
**Para** descobrir o que funciona melhor

**Critérios de Aceite:**
- [ ] Framework de experimentos A/B
- [ ] Métricas de sucesso por variante
- [ ] Rollout gradual de mudanças
- [ ] Dashboard de resultados

**Arquivos a criar:**
```
internal/ab/experiment_framework.go
internal/ab/metrics_collector.go
```

## Tarefas Técnicas

| ID | Tarefa | Estimativa | Responsável |
|----|--------|------------|-------------|
| T6.1 | Criar tabela de feedback | 2h | Dev |
| T6.2 | Implementar coletor de feedback | 4h | Dev |
| T6.3 | Implementar adaptador de estilo | 4h | Dev |
| T6.4 | Criar aprendiz de preferências | 8h | Dev |
| T6.5 | Expandir framework A/B existente | 4h | Dev |
| T6.6 | Criar métricas de aprendizado | 4h | Dev |
| T6.7 | Integrar com sistema de personalidade | 4h | Dev |
| T6.8 | Criar dashboard de aprendizado | 4h | Dev |
| T6.9 | Documentar sistema de aprendizado | 2h | Dev |

**Total Sprint 6:** ~36h (1 semana)

---

# 📊 RESUMO DO PROJETO

## Timeline

```
Semana 1: Sprint 1 - Testes Críticos      🔴
Semana 2: Sprint 2 - Métricas Prometheus  🟠
Semana 3: Sprint 3 - Auditoria LGPD       🟠
Semana 4: Sprint 4 - Memória Avançada     🟡
Semana 5: Sprint 5 - Autoconsciência      🟡
Semana 6: Sprint 6 - Aprendizado          🟢
```

## Esforço Total

| Sprint | Horas | Entregáveis Principais |
|--------|-------|------------------------|
| 1 | 36h | Testes C-SSRS, Alertas, Escalas |
| 2 | 36h | Prometheus + Grafana |
| 3 | 36h | Auditoria LGPD completa |
| 4 | 36h | Padrões temporais e previsões |
| 5 | 36h | Meta-cognição e escalação |
| 6 | 36h | Feedback e aprendizado |
| **Total** | **216h** | **6 semanas** |

## Definição de Pronto (DoD)

- [ ] Código revisado (self-review ou pair)
- [ ] Testes passando (quando aplicável)
- [ ] Documentação atualizada
- [ ] Métricas implementadas
- [ ] Sem erros críticos em log

## Riscos e Mitigações

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Complexidade do C-SSRS | Alta | Alto | Começar por ele |
| Integração Prometheus | Média | Médio | Usar lib padrão |
| LGPD compliance | Média | Alto | Consultar advogado |
| Padrões temporais | Alta | Médio | MVP simples primeiro |

---

# 🏁 PRÓXIMOS PASSOS

1. **Hoje:** Validar este plano com Jose
2. **Amanhã:** Iniciar Sprint 1 - Setup de testes
3. **Fim da semana:** Testes C-SSRS funcionando

---

*Documento criado por EVA + Jose R F Junior*
*Data: 2025-01-27*
*Versão: 1.0*
