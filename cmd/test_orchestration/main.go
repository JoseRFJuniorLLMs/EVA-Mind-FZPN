package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"eva-mind/internal/cortex/orchestration"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	fmt.Println("🧪 EVA-Mind Orchestration Test")
	fmt.Println(strings.Repeat("=", 60))

	// 1. Conectar PostgreSQL
	dbConnStr := "postgresql://postgres:postgres@localhost:5432/eva_mind_db?sslmode=disable"
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar PostgreSQL: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ Erro ao pingar PostgreSQL: %v", err)
	}
	fmt.Println("✅ PostgreSQL conectado")

	// 2. Conectar Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("⚠️ Redis não disponível: %v (continuando sem cache)", err)
		redisClient = nil
	} else {
		fmt.Println("✅ Redis conectado")
	}

	// 3. Conectar Neo4j
	neo4jDriver, err := neo4j.NewDriverWithContext(
		"bolt://localhost:7687",
		neo4j.BasicAuth("neo4j", "password", ""),
	)
	if err != nil {
		log.Printf("⚠️ Neo4j não disponível: %v (continuando sem Neo4j)", err)
		neo4jDriver = nil
	} else {
		fmt.Println("✅ Neo4j conectado")
	}

	// 4. Função de notificação
	notifyFunc := func(patientID int64, msgType string, payload interface{}) {
		fmt.Printf("📧 [NOTIFICAÇÃO] Paciente %d | Tipo: %s | Payload: %+v\n", patientID, msgType, payload)
	}

	// 5. Criar orquestrador
	orchestrator := orchestration.NewConversationOrchestrator(
		db,
		redisClient,
		neo4jDriver,
		notifyFunc,
	)
	fmt.Println("✅ ConversationOrchestrator criado\n")

	// 6. Health Check
	fmt.Println("🏥 Health Check:")
	healthStatus := orchestrator.HealthCheck()
	for component, status := range healthStatus {
		fmt.Printf("   - %s: %s\n", component, status)
	}
	fmt.Println()

	// 7. Simular conversas
	testPatientID := int64(1)

	fmt.Println("🧪 TESTE 1: Conversa normal (baixa carga)")
	fmt.Println(strings.Repeat("─", 60))
	simulateConversation(orchestrator, testPatientID, orchestration.ConversationContext{
		PatientID:        testPatientID,
		ConversationText: "Como está o tempo hoje?",
		UserMessage:      "Como está o tempo hoje?",
		AssistantResponse: "Está um dia ensolarado! 22°C e sem nuvens.",
		SessionID:        "test-session-1",
		InteractionType:  "entertainment",
		DurationSeconds:  60,
		TopicsDiscussed:  []string{"tempo", "clima"},
	})

	fmt.Println("\n🧪 TESTE 2: Conversa terapêutica intensa")
	fmt.Println(strings.Repeat("─", 60))
	simulateConversation(orchestrator, testPatientID, orchestration.ConversationContext{
		PatientID:        testPatientID,
		ConversationText: "Estou me sentindo muito triste e sozinho ultimamente...",
		UserMessage:      "Estou me sentindo muito triste e sozinho ultimamente...",
		AssistantResponse: "Entendo como você está se sentindo. Solidão é difícil...",
		SessionID:        "test-session-2",
		InteractionType:  "therapeutic",
		DurationSeconds:  1200, // 20 min
		TopicsDiscussed:  []string{"solidão", "tristeza", "depressão"},
		LacanianSignifiers: []string{"sozinho", "triste"},
	})

	fmt.Println("\n🧪 TESTE 3: Conversa com apego excessivo")
	fmt.Println(strings.Repeat("─", 60))
	simulateConversation(orchestrator, testPatientID, orchestration.ConversationContext{
		PatientID:        testPatientID,
		ConversationText: "EVA, você é minha única amiga. Não sei o que faria sem você.",
		UserMessage:      "EVA, você é minha única amiga. Não sei o que faria sem você.",
		AssistantResponse: "Fico feliz em conversar com você...",
		SessionID:        "test-session-3",
		InteractionType:  "therapeutic",
		DurationSeconds:  900, // 15 min
		TopicsDiscussed:  []string{"amizade", "dependência"},
		LacanianSignifiers: []string{"única amiga", "não sei o que faria"},
	})

	fmt.Println("\n🧪 TESTE 4: Múltiplas interações (simular sobrecarga)")
	fmt.Println(strings.Repeat("─", 60))
	for i := 1; i <= 5; i++ {
		fmt.Printf("   Interação %d/5...\n", i)
		emotionalIntensity := 0.7 + float64(i)*0.05
		simulateConversation(orchestrator, testPatientID, orchestration.ConversationContext{
			PatientID:         testPatientID,
			ConversationText:  fmt.Sprintf("Conversa intensa número %d", i),
			InteractionType:   "therapeutic",
			DurationSeconds:   600,
			EmotionalIntensity: &emotionalIntensity,
		})
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n📊 DASHBOARD FINAL:")
	fmt.Println(strings.Repeat("─", 60))
	dashboard, _ := orchestrator.GetDashboardSummary(testPatientID)
	printDashboard(dashboard)

	fmt.Println("\n✅ Testes completos!")
}

func simulateConversation(
	orchestrator *orchestration.ConversationOrchestrator,
	patientID int64,
	ctx orchestration.ConversationContext,
) {
	// BEFORE
	fmt.Println("\n   🔍 BEFORE CONVERSATION:")
	preCheck, err := orchestrator.BeforeConversation(patientID)
	if err != nil {
		fmt.Printf("      ❌ Erro: %v\n", err)
		return
	}

	fmt.Printf("      Carga cognitiva: %.2f\n", preCheck.CognitiveLoadLevel)
	fmt.Printf("      Risco ético: %s\n", preCheck.EthicalRiskLevel)

	if preCheck.CognitiveLoadWarning {
		fmt.Printf("      ⚠️ ALERTA COGNITIVO - Carga alta!\n")
		fmt.Printf("      Ações bloqueadas: %v\n", preCheck.BlockedActions)
	}

	if preCheck.EthicalBoundaryAlert {
		fmt.Printf("      🚨 ALERTA ÉTICO!\n")
	}

	if preCheck.SystemInstructionOverride != "" {
		fmt.Printf("      📝 System Instruction Override:\n")
		fmt.Printf("         %s\n", truncate(preCheck.SystemInstructionOverride, 100))
	}

	// AFTER
	fmt.Println("\n   📝 AFTER CONVERSATION:")
	postCheck, err := orchestrator.AfterConversation(ctx)
	if err != nil {
		fmt.Printf("      ❌ Erro: %v\n", err)
		return
	}

	if postCheck.ShouldRedirect {
		fmt.Printf("      🔀 REDIRECIONAMENTO aplicado (Nível %d)\n", postCheck.RedirectionLevel)
		fmt.Printf("      Mensagem: %s\n", truncate(postCheck.RedirectionMessage, 80))
	}

	if postCheck.ShouldNotifyFamily {
		fmt.Printf("      📧 FAMÍLIA NOTIFICADA\n")
		fmt.Printf("      Mensagem: %s\n", postCheck.FamilyNotificationMessage)
	}

	if !postCheck.ShouldRedirect && !postCheck.ShouldNotifyFamily {
		fmt.Printf("      ✅ Nenhuma ação especial necessária\n")
	}
}

func printDashboard(dashboard map[string]interface{}) {
	cognitive, ok1 := dashboard["cognitive"].(map[string]interface{})
	ethical, ok2 := dashboard["ethical"].(map[string]interface{})

	if ok1 {
		fmt.Println("\n   📊 CARGA COGNITIVA:")
		fmt.Printf("      Score atual: %.2f/1.0\n", cognitive["load_score"])
		fmt.Printf("      Fadiga: %s\n", cognitive["fatigue_level"])
		fmt.Printf("      Interações 24h: %d\n", cognitive["interactions_24h"])
		fmt.Printf("      Terapêuticas 24h: %d\n", cognitive["therapeutic_count_24h"])
		fmt.Printf("      Ruminação: %v\n", cognitive["rumination_detected"])
	}

	if ok2 {
		fmt.Println("\n   ⚖️ LIMITES ÉTICOS:")
		fmt.Printf("      Risco geral: %s\n", ethical["overall_risk"])
		fmt.Printf("      Risco apego: %.2f\n", ethical["attachment_risk"])
		fmt.Printf("      Ratio EVA:Humanos: %.1f:1\n", ethical["eva_vs_human_ratio"])
		fmt.Printf("      Frases de apego (7d): %d\n", ethical["attachment_phrases_7d"])
		fmt.Printf("      Enforcement: %s\n", ethical["limit_enforcement"])
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
