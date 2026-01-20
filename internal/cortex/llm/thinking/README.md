# 🧠 Gemini Thinking Mode - Health Triage

## Visão Geral

O módulo **Gemini Thinking** adiciona capacidades avançadas de raciocínio médico ao EVA-Mind, permitindo análise passo-a-passo de preocupações de saúde e triagem inteligente de riscos.

## Componentes

### 1. **Client** (`client.go`)
Cliente principal do Gemini 2.0 Flash Thinking Mode.

**Principais Funções**:
- `NewThinkingClient(apiKey)`: Cria cliente configurado
- `AnalyzeHealthConcern(ctx, concern, patientContext)`: Analisa preocupação de saúde
- Retorna `ThinkingResponse` com processo de raciocínio, nível de risco, e ações recomendadas

### 2. **Detector** (`detector.go`)
Detecta preocupações de saúde em mensagens.

**Funções**:
- `IsHealthConcern(message)`: Verifica se mensagem contém tópico de saúde
- `IsCriticalConcern(message)`: Identifica emergências médicas
- Usa keywords: dor, febre, medicamento, etc.

### 3. **Audit Logger** (`audit.go`)
Gerencia auditoria de análises no banco de dados.

**Funções**:
- `LogHealthAnalysis()`: Salva análise completa
- `MarkCaregiverNotified()`: Marca notificação enviada
- `GetPendingCriticalAlerts()`: Busca alertas não notificados
- `GetHealthSummary()`: Resumo de saúde do idoso

### 4. **Notification Service** (`notification.go`)
Envia notificações push para cuidadores.

**Funções**:
- `NotifyCaregiver()`: Envia alerta para cuidador
- `CheckPendingAlerts()`: Verifica alertas pendentes (cronjob)
- Mensagens customizadas por nível de risco

### 5. **Integration Service** (`integration.go`)
Orquestra o fluxo completo de triagem.

**Função Principal**:
```go
ProcessHealthConcern(ctx, idosoID, message, patientContext) (string, error)
```

**Fluxo**:
1. Detecta se é preocupação de saúde
2. Analisa com Thinking Mode
3. Salva auditoria
4. Notifica cuidador (se risco alto/crítico)
5. Retorna resposta com disclaimer

## Níveis de Risco

| Nível | Descrição | Ação |
|-------|-----------|------|
| **CRÍTICO** | Emergência médica | Notificação imediata + sugestão pronto-socorro |
| **ALTO** | Requer atenção médica urgente | Notificação + consulta em 24h |
| **MÉDIO** | Sintomas que precisam monitoramento | Orientação + sugestão consulta |
| **BAIXO** | Sintomas leves | Orientação geral |

## Uso no EVA-Mind

### Integração no Websocket Handler

```go
// No SignalingServer, adicionar:
type SignalingServer struct {
    // ... campos existentes ...
    healthTriage *thinking.HealthTriageService
}

// No NewSignalingServer:
healthTriage, err := thinking.NewHealthTriageService(cfg.GoogleAPIKey, db, pushService)
if err != nil {
    log.Fatalf("Erro ao criar health triage: %v", err)
}
server.healthTriage = healthTriage

// No handleGeminiResponse ou saveTranscription:
if server.healthTriage.ShouldUseThinkingMode(userMessage) {
    thinkingResponse, err := server.healthTriage.ProcessHealthConcern(
        ctx,
        session.IdosoID,
        userMessage,
        patientContext,
    )
    
    if err == nil && thinkingResponse != "" {
        // Usar resposta do Thinking Mode ao invés da resposta normal
        // Enviar para o usuário
    }
}
```

## Banco de Dados

### Tabela: `health_thinking_audit`

```sql
CREATE TABLE health_thinking_audit (
    id BIGSERIAL PRIMARY KEY,
    idoso_id BIGINT NOT NULL,
    concern TEXT NOT NULL,
    thought_process JSONB,
    risk_level VARCHAR(20),
    recommended_actions JSONB,
    seek_medical_care BOOLEAN,
    urgency_level VARCHAR(20),
    caregiver_notified BOOLEAN,
    notified_at TIMESTAMP,
    final_answer TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Views Úteis

- `v_health_concerns_summary`: Resumo por idoso (30 dias)
- `v_critical_alerts_pending`: Alertas críticos não notificados

## Testes

```bash
# Executar testes unitários
go test ./internal/llm/thinking/... -v

# Testes específicos
go test ./internal/llm/thinking -run TestIsHealthConcern
go test ./internal/llm/thinking -run TestIsCriticalConcern
```

## Segurança e Compliance

### Disclaimers Obrigatórios

Todas as respostas incluem disclaimer apropriado:
- ✅ "Sou uma assistente virtual e não substituo um profissional de saúde"
- ✅ Sempre recomenda consultar médico para sintomas preocupantes
- ✅ Não fornece diagnósticos

### Auditoria Completa

- ✅ Todas as análises são registradas
- ✅ Processo de raciocínio é salvo (transparência)
- ✅ Notificações são rastreadas
- ✅ Logs sanitizados (sem PII desnecessário)

## Monitoramento

### Queries Úteis

```sql
-- Alertas críticos das últimas 24h
SELECT * FROM v_critical_alerts_pending;

-- Resumo de saúde de um idoso
SELECT * FROM v_health_concerns_summary WHERE idoso_id = 123;

-- Taxa de notificação
SELECT 
    risk_level,
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE caregiver_notified) as notified,
    ROUND(100.0 * COUNT(*) FILTER (WHERE caregiver_notified) / COUNT(*), 2) as notification_rate
FROM health_thinking_audit
WHERE created_at >= NOW() - INTERVAL '7 days'
GROUP BY risk_level;
```

## Próximos Passos

1. ✅ Integrar no websocket handler principal
2. ✅ Testar com casos reais
3. ✅ Configurar cronjob para `CheckPendingAlerts()`
4. ✅ Adicionar métricas de performance
5. ✅ Criar dashboard de monitoramento

## Troubleshooting

### Thinking Mode não está sendo ativado
- Verificar se keywords de saúde estão na mensagem
- Aumentar sensibilidade em `IsHealthConcern()`

### Notificações não estão sendo enviadas
- Verificar se FCM token está configurado
- Checar logs de `NotifyCaregiver()`
- Verificar `v_critical_alerts_pending`

### Erro ao parsear JSON da resposta
- Gemini pode retornar formato diferente
- Fallback automático está implementado em `createFallbackResponse()`

---

**Criado em**: 15 de janeiro de 2026  
**Versão**: 1.0
