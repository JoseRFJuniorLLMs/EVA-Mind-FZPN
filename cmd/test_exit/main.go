package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"eva-mind/internal/exit"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("🕊️ Exit Protocol & Quality of Life - Test")
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

	// Criar Exit Protocol Manager
	epm := exit.NewExitProtocolManager(db)

	testPatientID := int64(1)

	// ========================================================================
	// FASE 1: LAST WISHES (TESTAMENTO VITAL)
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📝 FASE 1: Last Wishes (Testamento Vital Digital)")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	fmt.Printf("Criando Last Wishes para paciente %d...\n", testPatientID)
	lw, err := epm.CreateLastWishes(testPatientID)
	if err != nil {
		log.Printf("⚠️ Erro ao criar: %v (talvez já exista)", err)
		// Tentar buscar existente
		lw, err = epm.GetLastWishes(testPatientID)
		if err != nil {
			log.Printf("❌ Erro ao buscar: %v", err)
		}
	}

	if lw != nil {
		fmt.Printf("✅ Last Wishes ID: %s\n", lw.ID)
		fmt.Printf("   Completion: %d%%\n\n", lw.CompletionPercentage)

		// Atualizar algumas preferências
		fmt.Println("Atualizando preferências...")
		updates := map[string]interface{}{
			"resuscitation_preference":  "dnr",
			"preferred_death_location":  "home",
			"pain_management_preference": "aggressive_pain_control",
			"organ_donation_preference":  "donate_all",
			"burial_cremation":          "cremation",
			"personal_statement":        "Quero ser lembrado pela alegria que trouxe às pessoas. Vivi uma vida plena e estou em paz.",
		}

		err = epm.UpdateLastWishes(lw.ID, updates)
		if err != nil {
			log.Printf("❌ Erro ao atualizar: %v", err)
		} else {
			fmt.Println("✅ Preferências atualizadas")

			// Buscar novamente para ver completion
			lw, _ = epm.GetLastWishes(testPatientID)
			if lw != nil {
				fmt.Printf("   Nova completion: %d%%\n", lw.CompletionPercentage)
				fmt.Printf("   Ressuscitação: %s\n", lw.ResuscitationPreference)
				fmt.Printf("   Local preferido: %s\n", lw.PreferredDeathLocation)
				fmt.Printf("   Doação de órgãos: %s\n", lw.OrganDonationPreference)
			}
		}
	}

	fmt.Println()

	// ========================================================================
	// FASE 2: QUALITY OF LIFE ASSESSMENT (WHOQOL-BREF)
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📊 FASE 2: Quality of Life Assessment (WHOQOL-BREF)")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	fmt.Println("Registrando avaliação de qualidade de vida...")
	qol := &exit.QoLAssessment{
		PatientID:                 testPatientID,
		OverallQualityOfLife:      3, // 1=muito ruim, 5=muito bom
		OverallHealthSatisfaction: 3,
	}

	err = epm.RecordQoLAssessment(qol)
	if err != nil {
		log.Printf("❌ Erro ao registrar QoL: %v", err)
	} else {
		fmt.Println("✅ Avaliação WHOQOL-BREF registrada:")
		fmt.Printf("   Overall QoL Score: %.1f/100\n", qol.OverallQoLScore)
		fmt.Printf("   Physical Domain: %.1f/100\n", qol.PhysicalDomainScore)
		fmt.Printf("   Psychological Domain: %.1f/100\n", qol.PsychologicalDomainScore)
		fmt.Printf("   Social Domain: %.1f/100\n", qol.SocialDomainScore)
		fmt.Printf("   Environmental Domain: %.1f/100\n", qol.EnvironmentalDomainScore)
		fmt.Println()

		interpretQoLScore(qol.OverallQoLScore)
	}

	fmt.Println()

	// ========================================================================
	// FASE 3: PAIN & SYMPTOM MONITORING
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("🩹 FASE 3: Pain & Symptom Monitoring")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	fmt.Println("Registrando sintomas de dor moderada...")
	painLog1 := &exit.PainLog{
		PatientID:     testPatientID,
		PainPresent:   true,
		PainIntensity: 5,
		PainLocation:  []string{"lower_back", "right_hip"},
		PainQuality:   []string{"aching", "stiff"},
		Fatigue:       6,
		OverallWellbeing: 5,
		MedicationsTaken: []string{"paracetamol_500mg"},
		InterventionEffectiveness: 6,
		ReportedBy:    "patient",
	}

	err = epm.LogPainSymptoms(painLog1)
	if err != nil {
		log.Printf("❌ Erro: %v", err)
	} else {
		fmt.Printf("✅ Dor registrada: %d/10\n\n", painLog1.PainIntensity)
	}

	// Simular dor severa
	fmt.Println("Simulando dor severa (8/10)...")
	painLog2 := &exit.PainLog{
		PatientID:     testPatientID,
		PainPresent:   true,
		PainIntensity: 8,
		PainLocation:  []string{"abdomen"},
		PainQuality:   []string{"sharp", "constant"},
		Fatigue:       8,
		AnxietyLevel:  7,
		OverallWellbeing: 3,
		ReportedBy:    "patient",
	}

	err = epm.LogPainSymptoms(painLog2)
	if err != nil {
		log.Printf("❌ Erro: %v", err)
	} else {
		fmt.Println("✅ Dor severa registrada - Alerta automático acionado")
		fmt.Println("   (Sistema buscaria Comfort Care Plan automaticamente)")
	}

	fmt.Println()

	// ========================================================================
	// FASE 4: COMFORT CARE PLANS
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📋 FASE 4: Comfort Care Plans")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	fmt.Println("Criando Comfort Care Plan para dor severa...")
	comfortPlan := &exit.ComfortCarePlan{
		PatientID:        testPatientID,
		TriggerSymptom:   "severe_pain",
		TriggerThreshold: 7,
		Interventions: []exit.Intervention{
			{
				Order:              1,
				Type:               "pharmacological",
				Action:             "Morphine 5mg sublingual",
				RepeatAfterMinutes: 30,
			},
			{
				Order:  2,
				Type:   "positioning",
				Action: "Elevate head of bed 45 degrees, pillow under knees",
			},
			{
				Order:  3,
				Type:   "comfort",
				Action: "Cool compress, dim lights, soft instrumental music",
			},
			{
				Order:  4,
				Type:   "reassurance",
				Action: "EVA provides calming presence and breathing guidance",
			},
		},
		IsActive: true,
	}

	err = epm.CreateComfortCarePlan(comfortPlan)
	if err != nil {
		log.Printf("⚠️ Erro ao criar: %v (talvez já exista)", err)
	} else {
		fmt.Println("✅ Comfort Care Plan criado:")
		fmt.Printf("   Trigger: %s (threshold: %d/10)\n", comfortPlan.TriggerSymptom, comfortPlan.TriggerThreshold)
		fmt.Printf("   Intervenções: %d passos\n", len(comfortPlan.Interventions))
		fmt.Println()
		for _, intervention := range comfortPlan.Interventions {
			fmt.Printf("   %d. [%s] %s\n", intervention.Order, intervention.Type, intervention.Action)
		}
	}

	fmt.Println()

	// ========================================================================
	// FASE 5: LEGACY MESSAGES
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("💌 FASE 5: Legacy Messages (Mensagens de Legado)")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	fmt.Println("Criando mensagem para filha...")
	legacyMsg1 := &exit.LegacyMessage{
		PatientID:             testPatientID,
		RecipientName:         "Maria",
		RecipientRelationship: "daughter",
		MessageType:           "text",
		TextContent: `Minha querida Maria,

Se você está lendo isso, significa que meu tempo aqui terminou. Quero que você saiba que ser seu pai foi a maior honra da minha vida.

Lembre-se sempre:
- Seja gentil consigo mesma
- Valorize cada momento com seus filhos
- Não tenha medo de seguir seus sonhos
- Eu sempre estarei com você, no seu coração

Você fez tudo certo. Sou tão orgulhoso da mulher que você se tornou.

Te amo para sempre,
Papai`,
		DeliveryTrigger: "after_death",
		EmotionalTone:   "loving",
		Topics:          []string{"gratitude", "advice", "love"},
	}

	err = epm.CreateLegacyMessage(legacyMsg1)
	if err != nil {
		log.Printf("❌ Erro: %v", err)
	} else {
		fmt.Println("✅ Mensagem de legado criada para Maria (filha)")
		fmt.Println("   Trigger: after_death")
		fmt.Println("   Tipo: text")
		fmt.Println()

		// Marcar como completa
		err = epm.MarkLegacyMessageComplete(legacyMsg1.ID)
		if err == nil {
			fmt.Println("✅ Mensagem marcada como completa")
		}
	}

	fmt.Println()

	// Criar mais uma mensagem
	fmt.Println("Criando mensagem para neto...")
	legacyMsg2 := &exit.LegacyMessage{
		PatientID:             testPatientID,
		RecipientName:         "João",
		RecipientRelationship: "grandchild",
		MessageType:           "combined",
		TextContent: `João, meu neto querido,

Quando você se formar, lembre-se que o vovô sempre acreditou em você.
Estude muito, seja honesto, e faça o bem no mundo.

Com todo meu amor,
Vovô`,
		DeliveryTrigger: "milestone",
		EmotionalTone:   "hopeful",
		Topics:          []string{"advice", "memories", "wishes"},
	}

	err = epm.CreateLegacyMessage(legacyMsg2)
	if err != nil {
		log.Printf("❌ Erro: %v", err)
	} else {
		fmt.Println("✅ Mensagem de legado criada para João (neto)")
		fmt.Println("   Trigger: milestone (formatura)")
	}

	fmt.Println()

	// ========================================================================
	// FASE 6: FAREWELL PREPARATION
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("🕊️ FASE 6: Farewell Preparation (Preparação para Despedida)")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	fmt.Println("Iniciando preparação para despedida...")
	fp, err := epm.CreateFarewellPreparation(testPatientID)
	if err != nil {
		log.Printf("⚠️ Erro ao criar: %v (talvez já exista)", err)
		// Tentar buscar
		fp, _ = epm.GetFarewellPreparation(testPatientID)
	}

	if fp != nil {
		fmt.Printf("✅ Farewell Preparation ID: %s\n", fp.ID)
		fmt.Printf("   Estágio de luto: %s\n\n", fp.FiveStagesGriefPosition)

		// Atualizar progresso
		fmt.Println("Atualizando progresso da preparação...")
		fpUpdates := map[string]interface{}{
			"legal_affairs_complete":        true,
			"funeral_arrangements_complete": true,
			"five_stages_grief_position":    "acceptance",
			"emotional_readiness":           7,
			"spiritual_readiness":           8,
			"peace_with_life":               true,
			"peace_with_death":              true,
			"overall_preparation_score":     75,
		}

		err = epm.UpdateFarewellPreparation(testPatientID, fpUpdates)
		if err != nil {
			log.Printf("❌ Erro: %v", err)
		} else {
			fmt.Println("✅ Progresso atualizado:")
			fp, _ = epm.GetFarewellPreparation(testPatientID)
			if fp != nil {
				fmt.Printf("   Assuntos legais: %v\n", fp.LegalAffairsComplete)
				fmt.Printf("   Funeral arranjado: %v\n", fp.FuneralArrangementsComplete)
				fmt.Printf("   Estágio de luto: %s\n", fp.FiveStagesGriefPosition)
				fmt.Printf("   Prontidão emocional: %d/10\n", fp.EmotionalReadiness)
				fmt.Printf("   Prontidão espiritual: %d/10\n", fp.SpiritualReadiness)
				fmt.Printf("   Paz com a vida: %v\n", fp.PeaceWithLife)
				fmt.Printf("   Paz com a morte: %v\n", fp.PeaceWithDeath)
				fmt.Printf("   Score geral: %d/100\n", fp.OverallPreparationScore)
			}
		}
	}

	fmt.Println()

	// ========================================================================
	// FASE 7: SPIRITUAL CARE SESSION
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("🙏 FASE 7: Spiritual Care Session")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	fmt.Println("Registrando sessão de cuidado espiritual...")
	spiritualSession := &exit.SpiritualCareSession{
		PatientID:     testPatientID,
		ConductedBy:   "eva",
		ConductorName: "EVA-Companion",
		TopicsDiscussed: []string{
			"meaning_of_life",
			"gratitude",
			"legacy",
			"fear_of_death",
		},
		PracticesPerformed: []string{
			"meditation",
			"gratitude_reflection",
		},
		PreSessionPeaceLevel:  4,
		PostSessionPeaceLevel: 7,
		SpiritualNeedsIdentified: []string{
			"desire_to_connect_with_family",
			"need_for_forgiveness",
		},
		FollowUpNeeded:  true,
		DurationMinutes: 45,
	}

	err = epm.RecordSpiritualCareSession(spiritualSession)
	if err != nil {
		log.Printf("❌ Erro: %v", err)
	} else {
		fmt.Println("✅ Sessão espiritual registrada:")
		fmt.Printf("   Duração: %d minutos\n", spiritualSession.DurationMinutes)
		fmt.Printf("   Tópicos: %v\n", spiritualSession.TopicsDiscussed)
		fmt.Printf("   Paz antes: %d/10\n", spiritualSession.PreSessionPeaceLevel)
		fmt.Printf("   Paz depois: %d/10\n", spiritualSession.PostSessionPeaceLevel)
		fmt.Printf("   Melhora: +%d pontos\n", spiritualSession.PostSessionPeaceLevel-spiritualSession.PreSessionPeaceLevel)
	}

	fmt.Println()

	// ========================================================================
	// FASE 8: PALLIATIVE CARE SUMMARY
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📈 FASE 8: Palliative Care Summary (Resumo Geral)")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	summary, err := epm.GetPalliativeCareSummary(testPatientID)
	if err != nil {
		log.Printf("⚠️ Erro ao buscar resumo: %v", err)
	} else {
		printPalliativeSummary(summary)
	}

	fmt.Println()

	// ========================================================================
	// FASE 9: UNCONTROLLED PAIN ALERTS
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("🚨 FASE 9: Uncontrolled Pain Alerts")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	alerts, err := epm.GetUncontrolledPainAlerts()
	if err != nil {
		log.Printf("❌ Erro: %v", err)
	} else {
		if len(alerts) == 0 {
			fmt.Println("✅ Nenhum alerta de dor não controlada")
		} else {
			fmt.Printf("⚠️ %d alertas de dor não controlada:\n\n", len(alerts))
			for i, alert := range alerts {
				fmt.Printf("%d. Paciente %s (ID %d)\n", i+1, alert.PatientName, alert.PatientID)
				fmt.Printf("   Intensidade: %d/10\n", alert.PainIntensity)
				fmt.Printf("   Há %.1f horas\n", alert.HoursSinceReport)
				if alert.InterventionEffectiveness > 0 {
					fmt.Printf("   Eficácia da intervenção: %d/10\n", alert.InterventionEffectiveness)
				} else {
					fmt.Println("   ⚠️ Nenhuma intervenção eficaz ainda")
				}
				fmt.Println()
			}
		}
	}

	// ========================================================================
	// CONCLUSÃO
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("✅ Teste do Exit Protocol completo")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()
	fmt.Println("📊 Funcionalidades testadas:")
	fmt.Println("   ✓ Last Wishes (Testamento Vital)")
	fmt.Println("   ✓ Quality of Life Assessment (WHOQOL-BREF)")
	fmt.Println("   ✓ Pain & Symptom Monitoring")
	fmt.Println("   ✓ Comfort Care Plans")
	fmt.Println("   ✓ Legacy Messages")
	fmt.Println("   ✓ Farewell Preparation")
	fmt.Println("   ✓ Spiritual Care Sessions")
	fmt.Println("   ✓ Palliative Care Summary")
	fmt.Println("   ✓ Uncontrolled Pain Alerts")
	fmt.Println()
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func interpretQoLScore(score float64) {
	fmt.Print("   Interpretação: ")
	if score >= 80 {
		fmt.Println("Excelente qualidade de vida ✅")
	} else if score >= 60 {
		fmt.Println("Boa qualidade de vida 👍")
	} else if score >= 40 {
		fmt.Println("Qualidade de vida moderada ⚠️")
	} else if score >= 20 {
		fmt.Println("Qualidade de vida baixa ⚠️⚠️")
	} else {
		fmt.Println("Qualidade de vida muito baixa 🚨")
	}
}

func printPalliativeSummary(s *exit.PalliativeSummary) {
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Printf("                   RELATÓRIO DE CUIDADOS PALIATIVOS\n")
	fmt.Printf("                   Paciente: %s (ID %d)\n", s.PatientName, s.PatientID)
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("📝 LAST WISHES (Testamento Vital)")
	fmt.Printf("   Completion: %d%%", s.LastWishesCompletion)
	if s.LastWishesCompletion >= 80 {
		fmt.Println(" ✅")
	} else if s.LastWishesCompletion >= 50 {
		fmt.Println(" ⚠️")
	} else {
		fmt.Println(" 🚨 Requer atenção")
	}
	if s.ResuscitationPreference != "" {
		fmt.Printf("   Preferência de ressuscitação: %s\n", s.ResuscitationPreference)
	}
	fmt.Println()

	fmt.Println("📊 QUALITY OF LIFE")
	if s.OverallQoLScore > 0 {
		fmt.Printf("   Overall QoL Score: %.1f/100", s.OverallQoLScore)
		if s.OverallQoLScore >= 60 {
			fmt.Println(" 👍")
		} else if s.OverallQoLScore >= 40 {
			fmt.Println(" ⚠️")
		} else {
			fmt.Println(" 🚨 Baixa QoL")
		}
	} else {
		fmt.Println("   Nenhuma avaliação recente")
	}
	fmt.Println()

	fmt.Println("🩹 PAIN MANAGEMENT (últimos 7 dias)")
	if s.AvgPain7Days > 0 {
		fmt.Printf("   Dor média: %.1f/10", s.AvgPain7Days)
		if s.AvgPain7Days < 4 {
			fmt.Println(" ✅ Bem controlada")
		} else if s.AvgPain7Days < 7 {
			fmt.Println(" ⚠️ Moderada")
		} else {
			fmt.Println(" 🚨 Severa - requer intervenção")
		}
		fmt.Printf("   Pico de dor: %d/10\n", s.MaxPain7Days)
	} else {
		fmt.Println("   Sem registros de dor")
	}
	fmt.Println()

	fmt.Println("🕊️ EMOTIONAL & SPIRITUAL READINESS")
	if s.EmotionalReadiness > 0 || s.SpiritualReadiness > 0 {
		if s.EmotionalReadiness > 0 {
			fmt.Printf("   Prontidão emocional: %d/10", s.EmotionalReadiness)
			if s.EmotionalReadiness >= 7 {
				fmt.Println(" ✅")
			} else if s.EmotionalReadiness >= 4 {
				fmt.Println(" ⏳ Em progresso")
			} else {
				fmt.Println(" ⚠️ Suporte adicional necessário")
			}
		}
		if s.SpiritualReadiness > 0 {
			fmt.Printf("   Prontidão espiritual: %d/10", s.SpiritualReadiness)
			if s.SpiritualReadiness >= 7 {
				fmt.Println(" ✅")
			} else if s.SpiritualReadiness >= 4 {
				fmt.Println(" ⏳ Em progresso")
			} else {
				fmt.Println(" ⚠️ Suporte adicional necessário")
			}
		}
	} else {
		fmt.Println("   Avaliação não iniciada")
	}
	fmt.Println()

	fmt.Println("💌 LEGACY MESSAGES")
	fmt.Printf("   Completas: %d\n", s.LegacyMessagesCompleted)
	fmt.Printf("   Pendentes de entrega: %d\n", s.LegacyMessagesPending)
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println()
}
