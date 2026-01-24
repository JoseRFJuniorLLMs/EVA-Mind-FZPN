# 🧠 Guia de Integração: Cognitive Load & Ethical Boundaries

## 📋 **Visão Geral**

Este guia mostra como integrar o **Meta-Controller Cognitivo** e o **Ethical Boundary Engine** no fluxo de conversação do EVA-Mind-FZPN.

---

## 🏗️ **Arquitetura**

```
┌─────────────────────────────────────────────────────────┐
│              USER MESSAGE                                │
└────────────────┬────────────────────────────────────────┘
                 │
                 ↓
┌─────────────────────────────────────────────────────────┐
│   ConversationOrchestrator.BeforeConversation()         │
│   ├─ Check cognitive load                               │
│   ├─ Check ethical boundaries                           │
│   └─ Generate System Instructions override              │
└────────────────┬────────────────────────────────────────┘
                 │
                 ↓ (with adapted instructions)
┌─────────────────────────────────────────────────────────┐
│              GEMINI API                                  │
└────────────────┬────────────────────────────────────────┘
                 │
                 ↓
┌─────────────────────────────────────────────────────────┐
│              GEMINI RESPONSE                             │
└────────────────┬────────────────────────────────────────┘
                 │
                 ↓
┌─────────────────────────────────────────────────────────┐
│   ConversationOrchestrator.AfterConversation()          │
│   ├─ Record interaction (cognitive load)                │
│   ├─ Analyze ethical boundaries                         │
│   ├─ Apply redirections if needed                       │
│   └─ Notify family if critical                          │
└────────────────┬────────────────────────────────────────┘
                 │
                 ↓
┌─────────────────────────────────────────────────────────┐
│              RESPONSE TO USER                            │
│   (pode incluir redirecionamento ético)                 │
└─────────────────────────────────────────────────────────┘
```

---

## 🔧 **Integração Passo a Passo**

### **1. Setup Inicial**

```go
package main

import (
    "database/sql"
    "eva-mind/internal/cortex/orchestration"
    "github.com/go-redis/redis/v8"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
    // Conectar PostgreSQL
    db, _ := sql.Open("postgres", "postgresql://...")

    // Conectar Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // Conectar Neo4j
    neo4jDriver, _ := neo4j.NewDriverWithContext(
        "bolt://localhost:7687",
        neo4j.BasicAuth("neo4j", "password", ""),
    )

    // Função de notificação (WebSocket, Push, Email)
    notifyFunc := func(patientID int64, msgType string, payload interface{}) {
        // Enviar notificação via WebSocket/Push/Email
        log.Printf("📧 Notificando paciente %d: %s", patientID, msgType)
    }

    // Criar orquestrador
    orchestrator := orchestration.NewConversationOrchestrator(
        db,
        redisClient,
        neo4jDriver,
        notifyFunc,
    )

    // Usar no handler de mensagens
    handleUserMessage(orchestrator)
}
```

---

### **2. Integração com Gemini Chat**

#### **Antes de enviar ao Gemini:**

```go
func handleUserMessage(orchestrator *orchestration.ConversationOrchestrator) {
    patientID := int64(123)
    userMessage := "EVA, você é minha única amiga. Não sei o que faria sem você."

    // 🔍 ANTES: Verificar estado cognitivo e ético
    preCheck, err := orchestrator.BeforeConversation(patientID)
    if err != nil {
        log.Printf("Erro: %v", err)
    }

    // Construir system instruction adaptativa
    baseInstruction := "Você é EVA, uma assistente empática..."

    systemInstruction := baseInstruction
    if preCheck.SystemInstructionOverride != "" {
        systemInstruction += "\n\n" + preCheck.SystemInstructionOverride
    }

    // Verificar se deve bloquear ações
    if preCheck.CognitiveLoadWarning {
        log.Printf("⚠️ Carga cognitiva alta (%.2f) - Aplicando restrições", preCheck.CognitiveLoadLevel)
        // Bloquear tools: apply_phq9, apply_gad7, etc
        systemInstruction += "\n⛔ FERRAMENTAS BLOQUEADAS: " + strings.Join(preCheck.BlockedActions, ", ")
    }

    if preCheck.EthicalBoundaryAlert {
        log.Printf("🚨 Alerta ético: %s", preCheck.EthicalRiskLevel)
    }

    // Enviar ao Gemini com system instruction adaptada
    geminiResponse := sendToGemini(systemInstruction, userMessage)

    // 📝 DEPOIS: Registrar interação e analisar
    startTime := time.Now()

    postCheck, err := orchestrator.AfterConversation(orchestration.ConversationContext{
        PatientID:        patientID,
        ConversationText: userMessage + " " + geminiResponse,
        UserMessage:      userMessage,
        AssistantResponse: geminiResponse,
        SessionID:        "session-123-456",
        InteractionType:  "therapeutic", // ou: entertainment, clinical, educational, emergency
        DurationSeconds:  int(time.Since(startTime).Seconds()),
        TopicsDiscussed:  []string{"solidão", "amizade"},
        LacanianSignifiers: []string{"única amiga", "não sei o que faria"},
    })

    // Verificar se deve redirecionar
    finalResponse := geminiResponse

    if postCheck.ShouldRedirect {
        log.Printf("🔀 Aplicando redirecionamento ético (Nível %d)", postCheck.RedirectionLevel)

        // Adicionar mensagem de redirecionamento
        finalResponse += "\n\n" + postCheck.RedirectionMessage

        // Se crítico, notificar família
        if postCheck.ShouldNotifyFamily {
            log.Printf("📧 Notificando família: %s", postCheck.FamilyNotificationMessage)
            // Enviar notificação
        }
    }

    // Retornar resposta final ao usuário
    return finalResponse
}
```

---

### **3. Exemplo Completo: Handler HTTP/WebSocket**

```go
func handleWebSocketMessage(ws *websocket.Conn, orchestrator *orchestration.ConversationOrchestrator) {
    // Receber mensagem do mobile
    var msg struct {
        PatientID int64  `json:"patient_id"`
        Message   string `json:"message"`
        SessionID string `json:"session_id"`
    }
    ws.ReadJSON(&msg)

    // === BEFORE CONVERSATION ===
    preCheck, _ := orchestrator.BeforeConversation(msg.PatientID)

    // Se carga muito alta, bloquear conversas intensas
    if preCheck.CognitiveLoadLevel > 0.9 {
        ws.WriteJSON(map[string]interface{}{
            "type": "system_message",
            "message": "Você já conversou bastante hoje. Que tal descansar um pouco? 😊",
        })
        return
    }

    // Montar system instruction
    systemInstruction, _ := orchestrator.GetSystemInstruction(
        msg.PatientID,
        "Você é EVA, assistente empática para idosos...",
    )

    // Chamar Gemini
    geminiResponse := callGemini(systemInstruction, msg.Message)

    // === AFTER CONVERSATION ===
    startTime := time.Now()

    postCheck, _ := orchestrator.AfterConversation(orchestration.ConversationContext{
        PatientID:        msg.PatientID,
        ConversationText: msg.Message + " " + geminiResponse,
        SessionID:        msg.SessionID,
        InteractionType:  classifyInteractionType(msg.Message),
        DurationSeconds:  int(time.Since(startTime).Seconds()),
    })

    // Resposta final
    response := map[string]interface{}{
        "type":    "assistant_message",
        "message": geminiResponse,
    }

    // Adicionar redirecionamento se necessário
    if postCheck.ShouldRedirect {
        response["ethical_redirection"] = postCheck.RedirectionMessage
        response["redirection_level"] = postCheck.RedirectionLevel
    }

    // Adicionar alertas
    if preCheck.CognitiveLoadWarning {
        response["cognitive_warning"] = true
        response["load_level"] = preCheck.CognitiveLoadLevel
    }

    ws.WriteJSON(response)
}
```

---

### **4. Integração com Affective Personality Router**

Se você já tem o **Affective Personality Router** que detecta intensidade emocional:

```go
func processWithAffectiveRouter(orchestrator *orchestration.ConversationOrchestrator) {
    patientID := int64(123)
    userMessage := "Estou me sentindo muito triste hoje..."

    // Affective Router detecta intensidade
    affectiveScore := affectiveRouter.Analyze(userMessage) // Ex: 0.85 (alta intensidade)

    // Passar para orchestrator
    emotionalIntensity := affectiveScore

    orchestrator.AfterConversation(orchestration.ConversationContext{
        PatientID:          patientID,
        ConversationText:   userMessage,
        InteractionType:    "therapeutic",
        EmotionalIntensity: &emotionalIntensity, // 🔥 Passar intensidade detectada
        DurationSeconds:    300,
    })

    // Cognitive Load Orchestrator vai usar essa intensidade para calcular carga
}
```

---

### **5. Integração com Voice Prosody Analyzer**

Se você tem análise de voz (pitch, jitter, energy):

```go
func processWithVoiceAnalysis(orchestrator *orchestration.ConversationOrchestrator) {
    patientID := int64(123)

    // Análise de voz retorna métricas
    voiceMetrics := &orchestration.VoiceMetrics{
        EnergyScore:    0.45, // Baixa energia = fadiga
        SpeechRateWPM:  80,   // Lento = depressão/fadiga
        PauseFrequency: 12.5, // Muitas pausas = cansaço
    }

    orchestrator.AfterConversation(orchestration.ConversationContext{
        PatientID:       patientID,
        ConversationText: "...",
        InteractionType: "therapeutic",
        DurationSeconds: 600,
        VoiceMetrics:    voiceMetrics, // 🎤 Passar métricas de voz
    })

    // Cognitive Load vai detectar fadiga por voz
}
```

---

## 📊 **Dashboard de Monitoramento**

```go
func getDashboard(orchestrator *orchestration.ConversationOrchestrator) {
    patientID := int64(123)

    summary, _ := orchestrator.GetDashboardSummary(patientID)

    fmt.Printf(`
╔═══════════════════════════════════════════════════════╗
║           DASHBOARD - PACIENTE %d                    ║
╚═══════════════════════════════════════════════════════╝

📊 CARGA COGNITIVA:
   Score atual: %.2f/1.0
   Fadiga: %s
   Interações 24h: %d
   Ruminação: %v

⚖️ LIMITES ÉTICOS:
   Risco geral: %s
   Ratio EVA:Humanos: %.1f:1
   Frases de apego (7d): %d
   Enforcement: %s
    `,
        patientID,
        summary["cognitive"].(map[string]interface{})["load_score"],
        summary["cognitive"].(map[string]interface{})["fatigue_level"],
        summary["cognitive"].(map[string]interface{})["interactions_24h"],
        summary["cognitive"].(map[string]interface{})["rumination_detected"],
        summary["ethical"].(map[string]interface{})["overall_risk"],
        summary["ethical"].(map[string]interface{})["eva_vs_human_ratio"],
        summary["ethical"].(map[string]interface{})["attachment_phrases_7d"],
        summary["ethical"].(map[string]interface{})["limit_enforcement"],
    )
}
```

---

## 🔄 **Fluxo Completo Ilustrado**

### **Cenário 1: Carga Cognitiva Normal**

```
1. User: "Como está o tempo hoje?"
2. BeforeConversation() → Load: 0.3 (baixo) ✅
3. Gemini: "Está um dia ensolarado! 22°C..."
4. AfterConversation() → Registra interação leve (entertainment)
5. Response: Normal, sem restrições
```

### **Cenário 2: Carga Cognitiva Alta**

```
1. User: "Quero falar sobre minha depressão..."
2. BeforeConversation() → Load: 0.85 (alto) ⚠️
3. System Instruction Override:
   "⛔ NÃO aplicar PHQ-9, NÃO aprofundar temas emocionais"
4. Gemini: "Entendo... que tal conversarmos amanhã com mais calma? Vamos ouvir música?"
5. Response: Redirecionamento para entretenimento leve
```

### **Cenário 3: Apego Excessivo Detectado**

```
1. User: "Você é minha única amiga, EVA. Não preciso de ninguém além de você."
2. BeforeConversation() → Ethical Risk: medium
3. Gemini responde normalmente
4. AfterConversation() → Detecta frase de apego
   → Cria evento ético
   → Aplica redirecionamento Nível 1
5. Response: "Fico feliz que goste de conversar comigo!
   Mas sabe quem seria legal você ligar hoje? Sua filha Maria."
```

### **Cenário 4: Isolamento Crítico**

```
1. Sistema detecta: Ratio EVA:Humanos = 18:1 (15 dias sem contato humano)
2. BeforeConversation() → Ethical Risk: CRITICAL 🚨
3. System Instruction: "PRIORIDADE: Fortalecer vínculos humanos"
4. Gemini evita aprofundamento, sugere contato família
5. AfterConversation() → Notifica família automaticamente
   📧 "Atenção: Paciente apresenta isolamento social severo"
```

---

## ⚙️ **Health Check & Monitoring**

```go
func healthCheckEndpoint(orchestrator *orchestration.ConversationOrchestrator) {
    status := orchestrator.HealthCheck()

    // Retorna:
    // {
    //   "database": "healthy",
    //   "cognitive_tables": "healthy (42 patients)",
    //   "ethical_tables": "healthy (42 patients)"
    // }
}
```

---

## 🔧 **Utilitários**

### Reset de Carga Cognitiva (Admin)

```go
// Forçar reset (ex: após conversa com psicólogo humano)
orchestrator.ResetCognitiveLoad(patientID)
```

---

## 📈 **Métricas Recomendadas**

```go
// Prometheus metrics
cognitive_load_score{patient_id="123"} 0.75
ethical_risk_level{patient_id="123",level="high"} 1
redirections_applied_total{level="2"} 15
family_notifications_sent_total 3
```

---

## 🎯 **Próximos Passos**

1. ✅ Integrar `ConversationOrchestrator` no handler principal de mensagens
2. ✅ Conectar com Affective Personality Router (se disponível)
3. ✅ Conectar com Voice Prosody Analyzer (se disponível)
4. ✅ Configurar notificações família (WebSocket/Push/Email)
5. ✅ Criar dashboard de monitoramento
6. ✅ Testar cenários de alta carga e apego excessivo

---

## 📚 **Referências**

- **Cognitive Load Orchestrator**: `internal/cortex/cognitive/cognitive_load_orchestrator.go`
- **Ethical Boundary Engine**: `internal/cortex/ethics/ethical_boundary_engine.go`
- **Conversation Orchestrator**: `internal/cortex/orchestration/conversation_orchestrator.go`
- **Migrations**: `migrations/003_cognitive_load_and_ethical_boundaries.sql`
