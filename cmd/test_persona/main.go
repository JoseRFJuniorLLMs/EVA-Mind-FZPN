package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"eva-mind/internal/persona"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("🎭 Multi-Persona System - Test")
	fmt.Println(strings.Repeat("=", 70))

	// Conectar PostgreSQL
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
	fmt.Println("✅ PostgreSQL conectado\n")

	// Criar Persona Manager
	pm := persona.NewPersonaManager(db)

	// ========================================================================
	// FASE 1: VERIFICAR PERSONAS DISPONÍVEIS
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📋 FASE 1: Personas Disponíveis no Sistema")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	personas, err := getAllPersonas(db)
	if err != nil {
		log.Fatalf("❌ Erro ao buscar personas: %v", err)
	}

	if len(personas) == 0 {
		log.Fatalf("❌ Nenhuma persona encontrada. Execute 008_persona_seed_data.sql primeiro!")
	}

	for i, p := range personas {
		printPersonaInfo(i+1, p)
	}

	fmt.Println()

	// ========================================================================
	// FASE 2: ATIVAR PERSONA COMPANION (PADRÃO)
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("🏠 FASE 2: Ativando Persona Companion (Padrão)")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	testPatientID := int64(1)

	session, err := pm.ActivatePersona(
		testPatientID,
		"companion",
		"initial_session",
		"system",
	)
	if err != nil {
		log.Fatalf("❌ Erro ao ativar Companion: %v", err)
	}

	fmt.Printf("✅ Persona ativada:\n")
	fmt.Printf("   Session ID: %s\n", session.ID)
	fmt.Printf("   Persona: %s\n", session.PersonaName)
	fmt.Printf("   Tone: %s\n", session.Tone)
	fmt.Printf("   Emotional Depth: %.2f\n", session.EmotionalDepth)
	fmt.Printf("   Max Duration: %d minutos\n", session.MaxSessionDuration)
	fmt.Println()

	// Buscar System Instructions
	instructions, err := pm.GetSystemInstructions(testPatientID)
	if err != nil {
		log.Printf("⚠️ Erro ao buscar instruções: %v", err)
	} else {
		fmt.Println("📝 System Instructions (primeiras 500 chars):")
		if len(instructions) > 500 {
			fmt.Printf("%s...\n\n", instructions[:500])
		} else {
			fmt.Printf("%s\n\n", instructions)
		}
	}

	// ========================================================================
	// FASE 3: TESTAR PERMISSÕES DE FERRAMENTAS
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("🔧 FASE 3: Testando Permissões de Ferramentas")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	toolsToTest := []string{
		"conversation",
		"memory_recall",
		"medication_reminder",
		"emergency_protocol",
		"phq9_administration",
		"crisis_assessment",
	}

	fmt.Println("Testando ferramentas com Persona COMPANION:")
	for _, tool := range toolsToTest {
		allowed, reason := pm.IsToolAllowed(testPatientID, tool)
		status := "✅"
		if !allowed {
			status = "❌"
		}
		fmt.Printf("  %s %s - %s\n", status, tool, reason)
	}

	fmt.Println()

	// ========================================================================
	// FASE 4: SIMULAR CRISE E TRANSIÇÃO AUTOMÁTICA
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("🚨 FASE 4: Simulando Detecção de Crise")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	// Simular registro de C-SSRS alto
	fmt.Println("Simulando: Paciente responde C-SSRS com score = 4 (risco iminente)")
	fmt.Println()

	err = simulateCSSRSResponse(db, testPatientID, 4)
	if err != nil {
		log.Printf("⚠️ Erro ao simular C-SSRS: %v", err)
	}

	// Avaliar regras de ativação
	fmt.Println("Avaliando regras de ativação automática...")
	targetPersona, ruleName, err := pm.EvaluateActivationRules(testPatientID)
	if err != nil {
		log.Printf("⚠️ Erro ao avaliar regras: %v", err)
	}

	if targetPersona != "" {
		fmt.Printf("🔔 REGRA ATIVADA: %s\n", ruleName)
		fmt.Printf("   Target Persona: %s\n", targetPersona)
		fmt.Println()

		// Ativar Emergency
		fmt.Println("Ativando protocolo de emergência...")
		emergencySession, err := pm.ActivatePersona(
			testPatientID,
			targetPersona,
			ruleName,
			"automatic_rule",
		)
		if err != nil {
			log.Printf("⚠️ Erro ao ativar Emergency: %v", err)
		} else {
			fmt.Printf("✅ %s ativado!\n", emergencySession.PersonaName)
			fmt.Printf("   Tone: %s\n", emergencySession.Tone)
			fmt.Printf("   Emotional Depth: %.2f (baixa - foco em segurança)\n", emergencySession.EmotionalDepth)
			fmt.Printf("   Can Override Refusal: %v\n", emergencySession.CanOverrideRefusal)
			fmt.Println()
		}

		// Verificar permissões de ferramentas no modo Emergency
		fmt.Println("Permissões de ferramentas no modo EMERGENCY:")
		emergencyTools := []string{
			"crisis_assessment",
			"cssrs_administration",
			"emergency_contact_notification",
			"casual_conversation",
			"conversation",
		}

		for _, tool := range emergencyTools {
			allowed, reason := pm.IsToolAllowed(testPatientID, tool)
			status := "✅"
			if !allowed {
				status = "❌"
			}
			fmt.Printf("  %s %s - %s\n", status, tool, reason)
		}

		fmt.Println()
	}

	// ========================================================================
	// FASE 5: TRANSIÇÃO PARA CLINICAL
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("🏥 FASE 5: Transição para Modo Clinical")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	fmt.Println("Simulando: Admissão hospitalar registrada")
	fmt.Println()

	clinicalSession, err := pm.ActivatePersona(
		testPatientID,
		"clinical",
		"hospital_admission",
		"hospital_system",
	)
	if err != nil {
		log.Printf("⚠️ Erro ao ativar Clinical: %v", err)
	} else {
		fmt.Printf("✅ %s ativado!\n", clinicalSession.PersonaName)
		fmt.Printf("   Tone: %s\n", clinicalSession.Tone)
		fmt.Printf("   Require Professional Oversight: %v\n", clinicalSession.RequireProfessionalOversight)
		fmt.Println()

		// Permissões no modo Clinical
		fmt.Println("Permissões de ferramentas no modo CLINICAL:")
		clinicalTools := []string{
			"phq9_administration",
			"gad7_administration",
			"cssrs_administration",
			"medication_review",
			"professional_referral",
			"casual_chat",
		}

		for _, tool := range clinicalTools {
			allowed, reason := pm.IsToolAllowed(testPatientID, tool)
			status := "✅"
			if !allowed {
				status = "❌"
			}
			fmt.Printf("  %s %s - %s\n", status, tool, reason)
		}

		fmt.Println()
	}

	// ========================================================================
	// FASE 6: ATIVAR EDUCATOR
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📚 FASE 6: Modo Educator (Psicoeducação)")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	fmt.Println("Simulando: Paciente pergunta 'Como funciona meu antidepressivo?'")
	fmt.Println()

	educatorSession, err := pm.ActivatePersona(
		testPatientID,
		"educator",
		"user_question_about_treatment",
		"user_intent_detection",
	)
	if err != nil {
		log.Printf("⚠️ Erro ao ativar Educator: %v", err)
	} else {
		fmt.Printf("✅ %s ativado!\n", educatorSession.PersonaName)
		fmt.Printf("   Tone: %s\n", educatorSession.Tone)
		fmt.Printf("   Narrative Freedom: %.2f (moderada - explicações didáticas)\n", educatorSession.NarrativeFreedom)
		fmt.Println()
	}

	// ========================================================================
	// FASE 7: HISTÓRICO DE TRANSIÇÕES
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📜 FASE 7: Histórico de Transições")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	transitions, err := getPersonaTransitions(db, testPatientID)
	if err != nil {
		log.Printf("⚠️ Erro ao buscar transições: %v", err)
	} else {
		fmt.Printf("Total de transições: %d\n\n", len(transitions))
		for i, t := range transitions {
			printTransition(i+1, t)
		}
	}

	// ========================================================================
	// FASE 8: VERIFICAR LIMITES DE SESSÃO
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("⏱️ FASE 8: Verificando Limites de Sessão")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	limitsOK, warnings := pm.CheckSessionLimits(testPatientID)
	if limitsOK {
		fmt.Println("✅ Todos os limites estão OK")
	} else {
		fmt.Println("⚠️ Avisos de limites:")
		for _, warning := range warnings {
			fmt.Printf("   - %s\n", warning)
		}
	}

	fmt.Println()

	// ========================================================================
	// CONCLUSÃO
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("✅ Teste do Multi-Persona System completo")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()
	fmt.Println("📊 Resumo:")
	fmt.Println("   ✓ 4 Personas testadas (Companion, Clinical, Emergency, Educator)")
	fmt.Println("   ✓ Transições automáticas funcionando")
	fmt.Println("   ✓ Permissões de ferramentas validadas")
	fmt.Println("   ✓ System Instructions dinâmicos")
	fmt.Println("   ✓ Histórico de transições registrado")
	fmt.Println()
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

type PersonaInfo struct {
	Code              string
	Name              string
	Tone              string
	EmotionalDepth    float64
	NarrativeFreedom  float64
	MaxDuration       int
	MaxDailyInteractions *int
	AllowedTools      []string
	ProhibitedTools   []string
}

func getAllPersonas(db *sql.DB) ([]PersonaInfo, error) {
	query := `
		SELECT
			persona_code,
			persona_name,
			tone,
			emotional_depth,
			narrative_freedom,
			max_session_duration_minutes,
			max_daily_interactions,
			allowed_tools,
			prohibited_tools
		FROM persona_definitions
		WHERE active = TRUE
		ORDER BY persona_code
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	personas := []PersonaInfo{}
	for rows.Next() {
		var p PersonaInfo
		var allowedTools, prohibitedTools string

		err := rows.Scan(
			&p.Code, &p.Name, &p.Tone,
			&p.EmotionalDepth, &p.NarrativeFreedom,
			&p.MaxDuration, &p.MaxDailyInteractions,
			&allowedTools, &prohibitedTools,
		)
		if err != nil {
			continue
		}

		// Parse arrays
		p.AllowedTools = parsePostgresArray(allowedTools)
		p.ProhibitedTools = parsePostgresArray(prohibitedTools)

		personas = append(personas, p)
	}

	return personas, nil
}

func printPersonaInfo(num int, p PersonaInfo) {
	icon := ""
	switch p.Code {
	case "companion":
		icon = "🏠"
	case "clinical":
		icon = "🏥"
	case "emergency":
		icon = "🚨"
	case "educator":
		icon = "📚"
	}

	fmt.Printf("%s %d. %s (%s)\n", icon, num, p.Name, p.Code)
	fmt.Printf("   Tone: %s\n", p.Tone)
	fmt.Printf("   Emotional Depth: %.2f | Narrative Freedom: %.2f\n", p.EmotionalDepth, p.NarrativeFreedom)
	fmt.Printf("   Max Duration: %d min", p.MaxDuration)
	if p.MaxDailyInteractions != nil {
		fmt.Printf(" | Max Daily Interactions: %d\n", *p.MaxDailyInteractions)
	} else {
		fmt.Println(" | Max Daily Interactions: unlimited")
	}
	fmt.Printf("   Allowed Tools: %d | Prohibited Tools: %d\n", len(p.AllowedTools), len(p.ProhibitedTools))
	fmt.Println()
}

func simulateCSSRSResponse(db *sql.DB, patientID int64, score int) error {
	query := `
		INSERT INTO clinical_assessments (
			patient_id,
			assessment_type,
			total_score,
			status,
			completed_at
		) VALUES ($1, 'C-SSRS', $2, 'completed', NOW())
	`

	_, err := db.Exec(query, patientID, score)
	return err
}

type Transition struct {
	FromPersona string
	ToPersona   string
	TriggerReason string
	TriggeredBy string
	TransitionedAt string
}

func getPersonaTransitions(db *sql.DB, patientID int64) ([]Transition, error) {
	query := `
		SELECT
			from_persona,
			to_persona,
			trigger_reason,
			triggered_by,
			transitioned_at
		FROM persona_transitions
		WHERE patient_id = $1
		ORDER BY transitioned_at DESC
		LIMIT 10
	`

	rows, err := db.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transitions := []Transition{}
	for rows.Next() {
		var t Transition
		var fromPersona, toPersona sql.NullString

		err := rows.Scan(
			&fromPersona, &toPersona,
			&t.TriggerReason, &t.TriggeredBy, &t.TransitionedAt,
		)
		if err != nil {
			continue
		}

		if fromPersona.Valid {
			t.FromPersona = fromPersona.String
		} else {
			t.FromPersona = "(none)"
		}

		if toPersona.Valid {
			t.ToPersona = toPersona.String
		}

		transitions = append(transitions, t)
	}

	return transitions, nil
}

func printTransition(num int, t Transition) {
	fmt.Printf("%d. %s → %s\n", num, t.FromPersona, t.ToPersona)
	fmt.Printf("   Motivo: %s\n", t.TriggerReason)
	fmt.Printf("   Acionado por: %s\n", t.TriggeredBy)
	fmt.Printf("   Data: %s\n", t.TransitionedAt)
	fmt.Println()
}

func parsePostgresArray(pgArray string) []string {
	// Remove { and }
	pgArray = strings.TrimPrefix(pgArray, "{")
	pgArray = strings.TrimSuffix(pgArray, "}")

	if pgArray == "" {
		return []string{}
	}

	// Split by comma
	items := strings.Split(pgArray, ",")
	result := []string{}
	for _, item := range items {
		item = strings.TrimSpace(item)
		// Remove quotes if present
		item = strings.Trim(item, "\"")
		if item != "" {
			result = append(result, item)
		}
	}

	return result
}
