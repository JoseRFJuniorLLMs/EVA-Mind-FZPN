package main

import (
	"context"
	"eva-mind/internal/transnar"
	"fmt"
	"time"
)

func main() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🧪 A/B TESTING DEMO - TransNAR")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	manager := transnar.NewABTestManager()

	// Simular 150 usuários
	fmt.Println("Simulando 150 sessões de usuários...")
	fmt.Println()

	for userID := int64(1); userID <= 150; userID++ {
		variant := manager.AssignVariant(userID)

		// Simular métricas (em produção, viriam de dados reais)
		desire := &transnar.DesireInference{
			Desire:     transnar.DesireSecurity,
			Confidence: 0.75,
		}

		// Simular engajamento (varia por variante)
		engagement := 0.7
		switch variant {
		case transnar.VariantAggressive:
			engagement = 0.65 // Menos engajamento
		case transnar.VariantEmpathetic:
			engagement = 0.82 // Mais engajamento
		case transnar.VariantDirective:
			engagement = 0.71
		case transnar.VariantControl:
			engagement = 0.75
		}

		manager.RecordIntervention(variant, desire, 8)
		manager.RecordSession(variant, engagement)

		if userID%30 == 0 {
			fmt.Printf("  Processados %d usuários...\n", userID)
		}
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Gerar relatório
	report := manager.GetReport()
	fmt.Println(report)

	// Demonstrar logging contínuo
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("CONTINUOUS MONITORING")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("Em produção, métricas seriam logadas a cada 1 hora.")
	fmt.Println("Use: go manager.LogMetrics(ctx, 1*time.Hour)")
	fmt.Println()

	// Exemplo de uso em produção
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go manager.LogMetrics(ctx, 2*time.Second)

	// Aguardar alguns logs
	time.Sleep(5 * time.Second)

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ A/B Testing Framework Ready!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
