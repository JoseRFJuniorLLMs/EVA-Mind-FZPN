package main

import (
	"context"
	"eva-mind/internal/lacan"
	"eva-mind/internal/personality"
	"eva-mind/internal/transnar"
	"fmt"
	"log"
)

func main() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🧠 TRANSNAR VALIDATION TEST SUITE")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	ctx := context.Background()
	analyzer := transnar.NewAnalyzer()
	detector := transnar.NewDesireDetector()

	passed := 0
	failed := 0

	// Test 1: Negation Pattern
	fmt.Println("🔬 TEST 1: Negation Pattern")
	fmt.Println("   Input: \"Não quero tomar o remédio\"")

	chain1 := analyzer.Analyze("Não quero tomar o remédio")
	desire1 := detector.Detect(ctx, chain1, []lacan.Signifier{}, personality.Type9)

	if desire1.Desire == transnar.DesireSecurity && desire1.Confidence > 0.6 {
		fmt.Printf("   ✅ PASS | Desejo: %s (%.0f%%)\n", desire1.Desire, desire1.Confidence*100)
		passed++
	} else {
		fmt.Printf("   ❌ FAIL | Esperado: security, Obtido: %s (%.0f%%)\n", desire1.Desire, desire1.Confidence*100)
		failed++
	}
	fmt.Println()

	// Test 2: Repetition (Loneliness)
	fmt.Println("🔬 TEST 2: Repetition Detection")
	fmt.Println("   Input: \"Estou sozinho\" (com histórico de 'solidão')")

	history := []lacan.Signifier{
		{Word: "solidão", Frequency: 5},
	}
	chain2 := analyzer.Analyze("Estou sozinho novamente")
	desire2 := detector.Detect(ctx, chain2, history, personality.Type9)

	if desire2.Desire == transnar.DesireConnection && desire2.Confidence > 0.7 {
		fmt.Printf("   ✅ PASS | Desejo: %s (%.0f%%)\n", desire2.Desire, desire2.Confidence*100)
		passed++
	} else {
		fmt.Printf("   ❌ FAIL | Esperado: connection, Obtido: %s (%.0f%%)\n", desire2.Desire, desire2.Confidence*100)
		failed++
	}
	fmt.Println()

	// Test 3: Type 6 + Fear
	fmt.Println("🔬 TEST 3: Type 6 + Negative Emotion")
	fmt.Println("   Input: \"Tenho medo\" (Tipo 6)")

	chain3 := analyzer.Analyze("Tenho medo de cair")
	desire3 := detector.Detect(ctx, chain3, []lacan.Signifier{}, personality.Type6)

	if desire3.Desire == transnar.DesireSecurity && desire3.Confidence > 0.8 {
		fmt.Printf("   ✅ PASS | Desejo: %s (%.0f%%)\n", desire3.Desire, desire3.Confidence*100)
		passed++
	} else {
		fmt.Printf("   ❌ FAIL | Esperado: security (>80%%), Obtido: %s (%.0f%%)\n", desire3.Desire, desire3.Confidence*100)
		failed++
	}
	fmt.Println()

	// Test 4: Signifier Chain Analysis
	fmt.Println("🔬 TEST 4: Signifier Chain Extraction")
	fmt.Println("   Input: \"Não quero esse remédio horrível\"")

	chain4 := analyzer.Analyze("Não quero esse remédio horrível")

	hasNegation := len(chain4.Negations) > 0
	hasModal := len(chain4.Modals) > 0
	highIntensity := chain4.Intensity > 0.7

	if hasNegation && hasModal && highIntensity {
		fmt.Printf("   ✅ PASS | Negações: %v, Modais: %v, Intensidade: %.2f\n",
			chain4.Negations, chain4.Modals, chain4.Intensity)
		passed++
	} else {
		fmt.Printf("   ❌ FAIL | Negação: %v, Modal: %v, Intensidade: %.2f\n",
			hasNegation, hasModal, chain4.Intensity)
		failed++
	}
	fmt.Println()

	// Test 5: Loneliness Keyword
	fmt.Println("🔬 TEST 5: Loneliness Signifier")
	fmt.Println("   Input: \"A solidão é difícil\"")

	chain5 := analyzer.Analyze("A solidão é difícil")
	desire5 := detector.Detect(ctx, chain5, []lacan.Signifier{}, personality.Type9)

	if desire5.Desire == transnar.DesireConnection && desire5.Confidence > 0.8 {
		fmt.Printf("   ✅ PASS | Desejo: %s (%.0f%%)\n", desire5.Desire, desire5.Confidence*100)
		passed++
	} else {
		fmt.Printf("   ❌ FAIL | Esperado: connection (>80%%), Obtido: %s (%.0f%%)\n", desire5.Desire, desire5.Confidence*100)
		failed++
	}
	fmt.Println()

	// Test 6: Response Strategy Selection
	fmt.Println("🔬 TEST 6: Response Strategy Selection")

	generator := transnar.NewResponseGenerator()
	strategy := generator.SelectStrategy(desire1, chain1)

	if strategy == transnar.Punctuation || strategy == transnar.Reflection {
		fmt.Printf("   ✅ PASS | Estratégia: %s\n", strategy)
		passed++
	} else {
		fmt.Printf("   ❌ FAIL | Estratégia inesperada: %s\n", strategy)
		failed++
	}
	fmt.Println()

	// Final Report
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 FINAL REPORT")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	total := passed + failed
	passRate := float64(passed) / float64(total) * 100

	fmt.Printf("Total Tests: %d\n", total)
	fmt.Printf("✅ Passed: %d\n", passed)
	fmt.Printf("❌ Failed: %d\n", failed)
	fmt.Printf("Pass Rate: %.1f%%\n", passRate)
	fmt.Println()

	if passRate >= 80 {
		fmt.Println("🎉 TRANSNAR VALIDATED!")
		fmt.Println("   Sistema pronto para uso.")
	} else {
		fmt.Println("⚠️ VALIDATION INCOMPLETE")
		fmt.Println("   Revisar regras de inferência.")
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if passRate < 80 {
		log.Fatal("Tests failed")
	}
}
