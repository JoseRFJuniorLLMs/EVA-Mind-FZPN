package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"eva-mind/internal/cortex/explainability"
	"eva-mind/internal/cortex/prediction"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("🧪 Clinical Decision Explainer - Test")
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

	// Criar predictor
	predictor := prediction.NewCrisisPredictor(db)

	// Testar predição para um paciente
	patientID := int64(1)

	fmt.Printf("🔮 Predizendo risco de crise para paciente %d...\n\n", patientID)

	explanation, err := predictor.PredictCrisisRisk(patientID)
	if err != nil {
		log.Fatalf("❌ Erro na predição: %v", err)
	}

	// Exibir explicação
	printExplanation(explanation)
}

func printExplanation(exp *explainability.Explanation) {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📊 EXPLICAÇÃO DA DECISÃO CLÍNICA")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	// Cabeçalho
	fmt.Printf("ID: %s\n", exp.ID)
	fmt.Printf("Paciente: %d\n", exp.PatientID)
	fmt.Printf("Tipo: %s\n", exp.DecisionType)
	fmt.Printf("Probabilidade: %.1f%%\n", exp.PredictionScore*100)
	fmt.Printf("Severidade: %s\n", strings.ToUpper(exp.Severity))
	fmt.Printf("Janela temporal: %s\n", exp.Timeframe)
	fmt.Println()

	// Explicação textual
	fmt.Println(strings.Repeat("─", 70))
	fmt.Println("📝 EXPLICAÇÃO EM LINGUAGEM NATURAL:")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Println()
	fmt.Println(exp.ExplanationText)

	// Fatores primários
	if len(exp.PrimaryFactors) > 0 {
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println("🎯 FATORES PRIMÁRIOS (Top 3):")
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println()

		for i, factor := range exp.PrimaryFactors {
			fmt.Printf("%d. %s\n", i+1, factor.Factor)
			fmt.Printf("   Contribuição: %.1f%%\n", factor.Contribution*100)
			fmt.Printf("   Status: %s\n", factor.Status)
			fmt.Printf("   Detalhes: %s\n", factor.HumanReadable)
			if factor.BaselineComparison != "" {
				fmt.Printf("   Baseline: %s\n", factor.BaselineComparison)
			}
			fmt.Println()
		}
	}

	// Fatores secundários
	if len(exp.SecondaryFactors) > 0 {
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println("📋 FATORES SECUNDÁRIOS:")
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println()

		for _, factor := range exp.SecondaryFactors {
			fmt.Printf("• %s (%.1f%%): %s\n", factor.Factor, factor.Contribution*100, factor.HumanReadable)
		}
		fmt.Println()
	}

	// Recomendações
	if len(exp.Recommendations) > 0 {
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println("💡 RECOMENDAÇÕES CLÍNICAS:")
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println()

		for i, rec := range exp.Recommendations {
			urgencyIcon := "📌"
			if rec.Urgency == "high" || rec.Urgency == "critical" {
				urgencyIcon = "🚨"
			}

			fmt.Printf("%s %d. [%s] %s\n", urgencyIcon, i+1, strings.ToUpper(rec.Urgency), rec.Action)
			fmt.Printf("   Justificativa: %s\n", rec.Rationale)
			fmt.Printf("   Prazo: %s\n", rec.Timeframe)
			fmt.Println()
		}
	}

	// Evidências
	if exp.SupportingEvidence != nil && len(exp.SupportingEvidence) > 0 {
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println("📎 EVIDÊNCIAS DE SUPORTE:")
		fmt.Println(strings.Repeat("─", 70))
		fmt.Println()

		if conversations, ok := exp.SupportingEvidence["conversation_excerpts"].([]string); ok && len(conversations) > 0 {
			fmt.Println("💬 Trechos de conversas recentes:")
			for _, conv := range conversations {
				fmt.Printf("   - %s\n", conv)
			}
			fmt.Println()
		}

		if audioSamples, ok := exp.SupportingEvidence["audio_samples"].([]string); ok && len(audioSamples) > 0 {
			fmt.Println("🎙️ Amostras de áudio disponíveis:")
			for _, audio := range audioSamples {
				fmt.Printf("   - %s\n", audio)
			}
			fmt.Println()
		}
	}

	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("✅ Explicação gerada em: %s\n", exp.CreatedAt.Format("02/01/2006 15:04:05"))
	fmt.Println(strings.Repeat("=", 70))
}
