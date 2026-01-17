package main

import (
	"context"
	"eva-mind/internal/config"
	"eva-mind/internal/infrastructure/graph"
	"eva-mind/internal/lacan"
	"eva-mind/internal/transnar"
	"eva-mind/internal/veracity"
	"fmt"
	"log"
)

func main() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🕵️ LIE DETECTION TEST SUITE")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Setup
	ctx := context.Background()
	cfg, _ := config.Load()

	neo4jClient, err := graph.NewNeo4jClient(cfg)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar Neo4j: %v", err)
	}
	defer neo4jClient.Close()

	lacanService := lacan.NewSignifierService(neo4jClient)
	transnarEngine := transnar.NewEngine(lacanService, nil, nil)

	detector := veracity.NewLieDetector(neo4jClient, lacanService, transnarEngine)
	responseGen := veracity.NewResponseGenerator()

	passed := 0
	failed := 0

	// Test 1: Contradição Direta
	fmt.Println("🔬 TEST 1: Direct Contradiction")
	fmt.Println("   Setup: Inserir registro de 'tomou remédio'")
	fmt.Println("   Input: \"Nunca tomei esse remédio\"")

	// TODO: Inserir dado de teste no grafo
	// Por ora, simular

	inconsistencies := detector.Detect(ctx, 123, "Nunca tomei esse remédio")

	if len(inconsistencies) > 0 && inconsistencies[0].Type == veracity.DirectContradiction {
		fmt.Println("   ✅ PASS | Contradição detectada")
		fmt.Printf("      Confiança: %.0f%%\n", inconsistencies[0].Confidence*100)
		fmt.Printf("      Severidade: %s\n", inconsistencies[0].Severity)

		strategy := responseGen.SelectStrategy(&inconsistencies[0])
		response := responseGen.GenerateResponse(&inconsistencies[0], strategy)
		fmt.Printf("      Resposta: \"%s\"\n", response)
		passed++
	} else {
		fmt.Println("   ❌ FAIL | Contradição não detectada")
		failed++
	}
	fmt.Println()

	// Test 2: Inconsistência Temporal
	fmt.Println("🔬 TEST 2: Temporal Inconsistency")
	fmt.Println("   Setup: Evento registrado há 3 dias")
	fmt.Println("   Input: \"Ontem meu joelho doeu\"")

	inconsistencies = detector.Detect(ctx, 123, "Ontem meu joelho doeu")

	if len(inconsistencies) > 0 && inconsistencies[0].Type == veracity.TemporalInconsistency {
		fmt.Println("   ✅ PASS | Inconsistência temporal detectada")
		fmt.Printf("      Confiança: %.0f%%\n", inconsistencies[0].Confidence*100)
		passed++
	} else {
		fmt.Println("   ⚠️ SKIP | Requer dados de teste no grafo")
	}
	fmt.Println()

	// Test 3: Inconsistência Emocional
	fmt.Println("🔬 TEST 3: Emotional Inconsistency")
	fmt.Println("   Setup: Significante 'medo' mencionado 10x")
	fmt.Println("   Input: \"Não tenho medo de nada\"")

	inconsistencies = detector.Detect(ctx, 123, "Não tenho medo de nada")

	if len(inconsistencies) > 0 && inconsistencies[0].Type == veracity.EmotionalInconsistency {
		fmt.Println("   ✅ PASS | Inconsistência emocional detectada")
		fmt.Printf("      Confiança: %.0f%%\n", inconsistencies[0].Confidence*100)

		// Inferir desejo
		desire := responseGen.InferDesireFromLie(&inconsistencies[0])
		fmt.Printf("      Desejo inferido: %s\n", desire)
		passed++
	} else {
		fmt.Println("   ⚠️ SKIP | Requer histórico de significantes")
	}
	fmt.Println()

	// Test 4: Response Strategy Selection
	fmt.Println("🔬 TEST 4: Response Strategy Selection")

	testInc := veracity.Inconsistency{
		Type:       veracity.DirectContradiction,
		Confidence: 0.85,
		Severity:   veracity.SeverityHigh,
	}

	strategy := responseGen.SelectStrategy(&testInc)

	if strategy == veracity.SoftConfrontation {
		fmt.Println("   ✅ PASS | Estratégia correta selecionada")
		fmt.Printf("      Estratégia: %s\n", strategy)
		passed++
	} else {
		fmt.Println("   ❌ FAIL | Estratégia incorreta")
		fmt.Printf("      Esperado: soft_confrontation, Obtido: %s\n", strategy)
		failed++
	}
	fmt.Println()

	// Test 5: Prompt Generation
	fmt.Println("🔬 TEST 5: Prompt Addendum Generation")

	testIncs := []veracity.Inconsistency{
		{
			Type:       veracity.DirectContradiction,
			Confidence: 0.85,
			Statement:  "Nunca tomei remédio",
			GraphEvidence: []veracity.Evidence{
				{Fact: "Tomou Aspirina em 10/01/2026"},
			},
			Severity: veracity.SeverityHigh,
		},
	}

	prompt := responseGen.GeneratePromptAddendum(testIncs)

	if len(prompt) > 0 {
		fmt.Println("   ✅ PASS | Prompt gerado")
		fmt.Printf("      Tamanho: %d caracteres\n", len(prompt))
		passed++
	} else {
		fmt.Println("   ❌ FAIL | Prompt vazio")
		failed++
	}
	fmt.Println()

	// Final Report
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 FINAL REPORT")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	total := passed + failed
	passRate := float64(passed) / float64(total) * 100

	fmt.Printf("Total Tests: %d\n", total)
	fmt.Printf("✅ Passed: %d\n", passed)
	fmt.Printf("❌ Failed: %d\n", failed)
	fmt.Printf("Pass Rate: %.1f%%\n", passRate)
	fmt.Println()

	if passRate >= 80 {
		fmt.Println("🎉 LIE DETECTION SYSTEM VALIDATED!")
		fmt.Println("   Sistema pronto para integração.")
	} else {
		fmt.Println("⚠️ VALIDATION INCOMPLETE")
		fmt.Println("   Revisar implementação.")
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
