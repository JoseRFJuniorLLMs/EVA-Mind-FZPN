package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"eva-mind/internal/cortex/prediction"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("🔮 Predictive Life Trajectory Simulator - Test")
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

	// Criar simulador
	simulator := prediction.NewTrajectorySimulator(db)

	// Testar simulação para um paciente
	patientID := int64(1)

	fmt.Printf("🔮 Simulando trajetória de 30 dias para paciente %d...\n\n", patientID)

	// 1. Simulação baseline
	results, err := simulator.SimulateTrajectory(patientID, 30)
	if err != nil {
		log.Fatalf("❌ Erro na simulação: %v", err)
	}

	printSimulationResults(results)

	// 2. Simular cenários de intervenção
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🎯 CENÁRIOS DE INTERVENÇÃO (What-If)")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	scenarios, err := simulator.SimulateScenarios(patientID, 30)
	if err != nil {
		log.Fatalf("❌ Erro ao simular cenários: %v", err)
	}

	for i, scenario := range scenarios {
		printScenario(i+1, scenario)
	}

	// 3. Gerar recomendações
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("💡 RECOMENDAÇÕES CLÍNICAS")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	recommendations := simulator.GenerateRecommendations(patientID, results, scenarios)

	if len(recommendations) == 0 {
		fmt.Println("✅ Nenhuma intervenção urgente necessária. Paciente em risco baixo.")
	} else {
		for i, rec := range recommendations {
			printRecommendation(i+1, rec)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("✅ Simulação completa")
	fmt.Println(strings.Repeat("=", 70))
}

func printSimulationResults(results *prediction.SimulationResults) {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📊 RESULTADOS DA SIMULAÇÃO (BASELINE)")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	// Cabeçalho
	fmt.Printf("Paciente ID: %d\n", results.PatientID)
	fmt.Printf("Simulações executadas: %d\n", results.NSimulations)
	fmt.Printf("Período simulado: %d dias\n", results.DaysAhead)
	fmt.Printf("Tempo de computação: %d ms\n", results.ComputationTimeMs)
	fmt.Println()

	// Estado inicial
	fmt.Println(strings.Repeat("─", 70))
	fmt.Println("📋 ESTADO INICIAL:")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("PHQ-9: %.1f\n", results.InitialState.PHQ9Score)
	fmt.Printf("GAD-7: %.1f\n", results.InitialState.GAD7Score)
	fmt.Printf("Adesão medicamentosa: %.1f%%\n", results.InitialState.MedicationAdherence*100)
	fmt.Printf("Sono: %.1f horas/noite\n", results.InitialState.SleepHours)
	fmt.Printf("Isolamento social: %d dias sem contato\n", results.InitialState.SocialIsolationDays)
	fmt.Printf("Carga cognitiva: %.2f\n", results.InitialState.CognitiveLoad)
	fmt.Println()

	// Probabilidades de desfechos
	fmt.Println(strings.Repeat("─", 70))
	fmt.Println("🎲 PROBABILIDADES DE DESFECHOS:")
	fmt.Println(strings.Repeat("─", 70))

	riskLevel := getRiskLevel(results.CrisisProbability30d)
	riskIcon := getRiskIcon(riskLevel)

	fmt.Printf("%s Crise em 7 dias:  %.1f%% (%s)\n",
		riskIcon,
		results.CrisisProbability7d*100,
		getRiskLevel(results.CrisisProbability7d))

	fmt.Printf("%s Crise em 30 dias: %.1f%% (%s)\n",
		riskIcon,
		results.CrisisProbability30d*100,
		riskLevel)

	fmt.Printf("🏥 Hospitalização:   %.1f%%\n", results.HospitalizationProb30d*100)
	fmt.Printf("💊 Abandono de tratamento: %.1f%%\n", results.TreatmentDropoutProb90d*100)
	fmt.Printf("🤕 Risco de queda:   %.1f%%\n", results.FallRiskProb7d*100)
	fmt.Println()

	// Projeções
	fmt.Println(strings.Repeat("─", 70))
	fmt.Println("📈 PROJEÇÕES AO FINAL DE 30 DIAS:")
	fmt.Println(strings.Repeat("─", 70))

	phq9Change := results.ProjectedPHQ9 - results.InitialState.PHQ9Score
	phq9Arrow := "→"
	if phq9Change > 0 {
		phq9Arrow = "↑"
	} else if phq9Change < 0 {
		phq9Arrow = "↓"
	}

	adherenceChange := results.ProjectedAdherence - results.InitialState.MedicationAdherence
	adherenceArrow := "→"
	if adherenceChange > 0 {
		adherenceArrow = "↑"
	} else if adherenceChange < 0 {
		adherenceArrow = "↓"
	}

	sleepChange := results.ProjectedSleepHours - results.InitialState.SleepHours
	sleepArrow := "→"
	if sleepChange > 0 {
		sleepArrow = "↑"
	} else if sleepChange < 0 {
		sleepArrow = "↓"
	}

	fmt.Printf("PHQ-9:     %.1f %s %.1f (mudança: %+.1f)\n",
		results.InitialState.PHQ9Score,
		phq9Arrow,
		results.ProjectedPHQ9,
		phq9Change)

	fmt.Printf("Adesão:    %.1f%% %s %.1f%% (mudança: %+.1f%%)\n",
		results.InitialState.MedicationAdherence*100,
		adherenceArrow,
		results.ProjectedAdherence*100,
		adherenceChange*100)

	fmt.Printf("Sono:      %.1fh %s %.1fh (mudança: %+.1fh)\n",
		results.InitialState.SleepHours,
		sleepArrow,
		results.ProjectedSleepHours,
		sleepChange)

	isolationChange := results.ProjectedIsolationDays - results.InitialState.SocialIsolationDays
	fmt.Printf("Isolamento: %d dias %s %d dias (mudança: %+d)\n",
		results.InitialState.SocialIsolationDays,
		getArrow(isolationChange),
		results.ProjectedIsolationDays,
		isolationChange)
	fmt.Println()

	// Fatores críticos
	if len(results.CriticalFactors) > 0 {
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println("⚠️ FATORES DE RISCO CRÍTICOS:")
		fmt.Println(strings.Repeat("─", 70))
		for _, factor := range results.CriticalFactors {
			fmt.Printf("• %s\n", translateFactor(factor))
		}
		fmt.Println()
	}

	// Tendências (amostra de trajetória)
	if len(results.SampleTrajectories) > 0 {
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println("📉 TENDÊNCIA PROJETADA (PHQ-9):")
		fmt.Println(strings.Repeat("─", 70))

		// Mostrar a cada 5 dias
		for i := 0; i < len(results.SampleTrajectories); i += 5 {
			sample := results.SampleTrajectories[i]
			day := int(sample["day"])
			phq9 := sample["phq9"]
			adherence := sample["adherence"]

			bar := strings.Repeat("█", int(phq9))
			fmt.Printf("Dia %2d: [%s] PHQ-9: %.1f | Adesão: %.0f%%\n",
				day, bar, phq9, adherence*100)
		}
		fmt.Println()
	}
}

func printScenario(number int, scenario prediction.InterventionScenario) {
	fmt.Printf("═══════════════════════════════════════════════════════════════════\n")
	fmt.Printf("CENÁRIO %d: %s\n", number, scenario.ScenarioName)
	fmt.Printf("═══════════════════════════════════════════════════════════════════\n")

	if scenario.ScenarioType == "baseline" {
		fmt.Println("Tipo: Linha de base (sem intervenções)")
	} else {
		fmt.Printf("Tipo: Com intervenção\n")
		fmt.Printf("Descrição: %s\n", scenario.Description)
		fmt.Println()

		if len(scenario.Interventions) > 0 {
			fmt.Println("Intervenções aplicadas:")
			for _, intervention := range scenario.Interventions {
				fmt.Printf("  • %s: %s\n", intervention.Type, intervention.Description)
				if intervention.Frequency != "" {
					fmt.Printf("    Frequência: %s\n", intervention.Frequency)
				}
			}
			fmt.Println()
		}
	}

	fmt.Printf("Risco de crise (7d):  %.1f%%\n", scenario.CrisisProbability7d*100)
	fmt.Printf("Risco de crise (30d): %.1f%%\n", scenario.CrisisProbability30d*100)
	fmt.Printf("Hospitalização (30d): %.1f%%\n", scenario.HospitalizationProb30d*100)
	fmt.Println()

	fmt.Printf("PHQ-9 projetado:  %.1f\n", scenario.ProjectedPHQ9)
	fmt.Printf("Adesão projetada: %.1f%%\n", scenario.ProjectedAdherence*100)
	fmt.Printf("Sono projetado:   %.1fh\n", scenario.ProjectedSleepHours)
	fmt.Println()

	if scenario.ScenarioType != "baseline" {
		if scenario.RiskReduction30d > 0 {
			fmt.Printf("✅ Redução de risco (7d):  %.1f%%\n", scenario.RiskReduction7d*100)
			fmt.Printf("✅ Redução de risco (30d): %.1f%%\n", scenario.RiskReduction30d*100)
		} else {
			fmt.Println("⚠️ Sem redução significativa de risco")
		}

		fmt.Printf("Efetividade: %.1f%%\n", scenario.EffectivenessScore*100)
		fmt.Printf("Custo estimado: R$ %.2f/mês\n", scenario.EstimatedCostMonthly)
		fmt.Printf("Viabilidade: %s\n", scenario.Feasibility)
	}

	fmt.Println()
}

func printRecommendation(number int, rec prediction.RecommendedIntervention) {
	priorityIcon := "📌"
	if rec.Priority == "critical" || rec.Priority == "high" {
		priorityIcon = "🚨"
	}

	fmt.Printf("%s %d. [%s] %s\n", priorityIcon, number, strings.ToUpper(rec.Priority), rec.Title)
	fmt.Printf("   Prazo: %s\n", rec.UrgencyTimeframe)
	fmt.Printf("   Descrição: %s\n", rec.Description)
	fmt.Printf("   Justificativa: %s\n", rec.Rationale)
	fmt.Printf("   Redução esperada de risco: %.1f%%\n", rec.ExpectedRiskReduction*100)

	if rec.ExpectedPHQ9Improvement != 0 {
		fmt.Printf("   Melhora esperada PHQ-9: %.1f pontos\n", rec.ExpectedPHQ9Improvement)
	}

	fmt.Printf("   Confiança: %.0f%%\n", rec.ConfidenceLevel*100)

	if len(rec.ActionSteps) > 0 {
		fmt.Println("   Passos:")
		for _, step := range rec.ActionSteps {
			fmt.Printf("     • %s\n", step)
		}
	}

	if rec.EstimatedCost > 0 {
		fmt.Printf("   Custo estimado: R$ %.2f\n", rec.EstimatedCost)
	} else {
		fmt.Println("   Custo estimado: Gratuito")
	}

	fmt.Println()
}

func getRiskLevel(probability float64) string {
	if probability >= 0.6 {
		return "CRÍTICO"
	} else if probability >= 0.4 {
		return "ALTO"
	} else if probability >= 0.2 {
		return "MODERADO"
	}
	return "BAIXO"
}

func getRiskIcon(level string) string {
	switch level {
	case "CRÍTICO":
		return "🔴"
	case "ALTO":
		return "🟠"
	case "MODERADO":
		return "🟡"
	default:
		return "🟢"
	}
}

func getArrow(change int) string {
	if change > 0 {
		return "↑"
	} else if change < 0 {
		return "↓"
	}
	return "→"
}

func translateFactor(factor string) string {
	translations := map[string]string{
		"low_medication_adherence":      "Baixa adesão medicamentosa (<50%)",
		"moderate_to_severe_depression": "Depressão moderada a severa (PHQ-9 ≥ 15)",
		"poor_sleep_quality":            "Qualidade de sono ruim (<5h/noite)",
		"social_isolation":              "Isolamento social (≥5 dias sem contato)",
		"moderate_to_severe_anxiety":    "Ansiedade moderada a severa (GAD-7 ≥ 10)",
		"high_cognitive_load":           "Carga cognitiva elevada (>0.7)",
		"low_motivation":                "Baixa motivação (<0.3)",
		"worsening_depression_trend":    "Tendência de piora na depressão",
		"declining_adherence_trend":     "Tendência de queda na adesão",
	}

	if translation, ok := translations[factor]; ok {
		return translation
	}
	return factor
}
