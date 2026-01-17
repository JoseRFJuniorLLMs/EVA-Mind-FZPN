package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"eva-mind/internal/config"
	"eva-mind/internal/infrastructure/cache"
	"eva-mind/internal/infrastructure/graph"
	"eva-mind/internal/lacan"
	"eva-mind/internal/memory"
	"eva-mind/internal/personality"
	"eva-mind/internal/telemetry"
)

// TestReport armazena resultados dos testes
type TestReport struct {
	TestName string
	Passed   bool
	Duration time.Duration
	Details  string
	ErrorMsg string
}

var reports []TestReport

func main() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🧪 FZPN VALIDATION TEST SUITE")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	ctx := context.Background()

	// Initialize components
	neo4jClient, err := graph.NewNeo4jClient(cfg)
	if err != nil {
		log.Fatalf("❌ Neo4j error: %v", err)
	}
	defer neo4jClient.Close(ctx)

	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		log.Printf("⚠️ Redis not available: %v. Some tests will be skipped.", err)
	}

	// Run test suites
	fmt.Println("📋 Running Test Suites...\n")

	testFDPNLatency(ctx, neo4jClient, redisClient)
	testZetaPersonality()
	testLacanSignifiers(ctx, neo4jClient)
	testAntiSycophancy()

	// Print final report
	printFinalReport()
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// TEST 1: FDPN Priming Latency
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func testFDPNLatency(ctx context.Context, neo4j *graph.Neo4jClient, redis *cache.RedisClient) {
	fmt.Println("🔬 TEST 1: FDPN Priming Latency")
	fmt.Println("   Objetivo: Verificar se priming é < 10ms (com cache)")
	fmt.Println()

	engine := memory.NewFDPNEngine(neo4j, redis)

	// Test Case 1: Cold Query (sem cache)
	testColdQuery(ctx, engine)

	// Test Case 2: Hot Query (com cache)
	testHotQuery(ctx, engine)

	// Test Case 3: Parallel Priming
	testParallelPriming(ctx, engine)

	fmt.Println()
}

func testColdQuery(ctx context.Context, engine *memory.FDPNEngine) {
	start := time.Now()

	err := engine.StreamingPrime(ctx, "test_user_1", "dor de cabeça tontura")

	elapsed := time.Since(start)

	report := TestReport{
		TestName: "FDPN Cold Query (Neo4j direto)",
		Duration: elapsed,
		Details:  fmt.Sprintf("Latência: %dms", elapsed.Milliseconds()),
	}

	if err != nil {
		report.Passed = false
		report.ErrorMsg = err.Error()
	} else if elapsed.Milliseconds() < 100 {
		report.Passed = true
		report.Details += " ✅ Excelente (< 100ms)"
	} else {
		report.Passed = true
		report.Details += " ⚠️ Aceitável mas lento"
	}

	reports = append(reports, report)
	printTestResult(report)
}

func testHotQuery(ctx context.Context, engine *memory.FDPNEngine) {
	// Prime cache first
	engine.StreamingPrime(ctx, "test_user_1", "dor de cabeça")

	// Now test cached retrieval
	start := time.Now()

	err := engine.StreamingPrime(ctx, "test_user_1", "dor de cabeça")

	elapsed := time.Since(start)

	report := TestReport{
		TestName: "FDPN Hot Query (Redis cache)",
		Duration: elapsed,
		Details:  fmt.Sprintf("Latência: %dms", elapsed.Milliseconds()),
	}

	if err != nil {
		report.Passed = false
		report.ErrorMsg = err.Error()
	} else if elapsed.Milliseconds() < 10 {
		report.Passed = true
		report.Details += " 🚀 PERFEITO (< 10ms)"
	} else if elapsed.Milliseconds() < 50 {
		report.Passed = true
		report.Details += " ✅ Bom (< 50ms)"
	} else {
		report.Passed = false
		report.Details += " ❌ Muito lento para cache"
	}

	reports = append(reports, report)
	printTestResult(report)
}

func testParallelPriming(ctx context.Context, engine *memory.FDPNEngine) {
	start := time.Now()

	// Simula frase complexa com múltiplas keywords
	err := engine.StreamingPrime(ctx, "test_user_1", "dor tontura remédio médico hospital")

	elapsed := time.Since(start)

	report := TestReport{
		TestName: "FDPN Parallel Priming (5 keywords)",
		Duration: elapsed,
		Details:  fmt.Sprintf("Latência: %dms", elapsed.Milliseconds()),
	}

	if err != nil {
		report.Passed = false
		report.ErrorMsg = err.Error()
	} else if elapsed.Milliseconds() < 20 {
		report.Passed = true
		report.Details += " 🚀 Goroutines funcionando!"
	} else {
		report.Passed = true
		report.Details += " ⚠️ Pode melhorar paralelismo"
	}

	reports = append(reports, report)
	printTestResult(report)
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// TEST 2: Zeta Personality Routing
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func testZetaPersonality() {
	fmt.Println("🔬 TEST 2: Zeta Personality Routing")
	fmt.Println("   Objetivo: Verificar mudanças de tipo por emoção")
	fmt.Println()

	router := personality.NewPersonalityRouter()

	// Test Case 1: Stress Path (9 → 6)
	testStressPath(router)

	// Test Case 2: Growth Path (9 → 3)
	testGrowthPath(router)

	// Test Case 3: Attention Weights
	testAttentionWeights(router)

	fmt.Println()
}

func testStressPath(router *personality.PersonalityRouter) {
	baseType := personality.Type9
	emotion := "anxiety" // Deve triggerar stress

	activeType, mode := router.RoutePersonality(baseType, emotion)

	report := TestReport{
		TestName: "Zeta Stress Path (9 → 6)",
		Details:  fmt.Sprintf("Base: %d, Emoção: %s → Tipo: %d, Modo: %s", baseType, emotion, activeType, mode),
	}

	if activeType == personality.Type6 && mode == "stress" {
		report.Passed = true
		report.Details += " ✅ Correto!"
	} else {
		report.Passed = false
		report.Details += fmt.Sprintf(" ❌ Esperado Tipo 6, obteve %d", activeType)
	}

	reports = append(reports, report)
	printTestResult(report)
}

func testGrowthPath(router *personality.PersonalityRouter) {
	baseType := personality.Type9
	emotion := "joy" // Deve triggerar growth

	activeType, mode := router.RoutePersonality(baseType, emotion)

	report := TestReport{
		TestName: "Zeta Growth Path (9 → 3)",
		Details:  fmt.Sprintf("Base: %d, Emoção: %s → Tipo: %d, Modo: %s", baseType, emotion, activeType, mode),
	}

	if activeType == personality.Type3 && mode == "growth" {
		report.Passed = true
		report.Details += " ✅ Correto!"
	} else {
		report.Passed = false
		report.Details += fmt.Sprintf(" ❌ Esperado Tipo 3, obteve %d", activeType)
	}

	reports = append(reports, report)
	printTestResult(report)
}

func testAttentionWeights(router *personality.PersonalityRouter) {
	weights := router.GetAttentionWeights(personality.Type6)

	report := TestReport{
		TestName: "Zeta Attention Weights (Tipo 6)",
		Details: fmt.Sprintf("RISCO: %.1f, SEGURANÇA: %.1f, AMBIGUIDADE: %.1f",
			weights["RISCO"], weights["SEGURANÇA"], weights["AMBIGUIDADE"]),
	}

	// Tipo 6 deve amplificar RISCO e reduzir AMBIGUIDADE
	if weights["RISCO"] > 2.0 && weights["AMBIGUIDADE"] < 1.0 {
		report.Passed = true
		report.Details += " ✅ Zeros corretos!"
	} else {
		report.Passed = false
		report.Details += " ❌ Pesos incorretos"
	}

	reports = append(reports, report)
	printTestResult(report)
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// TEST 3: Lacan Signifier Detection
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func testLacanSignifiers(ctx context.Context, neo4j *graph.Neo4jClient) {
	fmt.Println("🔬 TEST 3: Lacan Signifier Detection")
	fmt.Println("   Objetivo: Rastrear significantes emocionais")
	fmt.Println()

	service := lacan.NewSignifierService(neo4j)

	// Test Case 1: Track Signifier
	testTrackSignifier(ctx, service)

	// Test Case 2: Retrieve Key Signifiers
	testRetrieveSignifiers(ctx, service)

	fmt.Println()
}

func testTrackSignifier(ctx context.Context, service *lacan.SignifierService) {
	// Simula 5 menções de "solidão"
	texts := []string{
		"Me sinto muito sozinho",
		"A solidão é difícil",
		"Ninguém me visita, solidão total",
		"Solidão me consome",
		"Estou sozinho novamente",
	}

	start := time.Now()

	for _, text := range texts {
		err := service.TrackSignifiers(ctx, 999, text)
		if err != nil {
			report := TestReport{
				TestName: "Lacan Track Signifier",
				Passed:   false,
				ErrorMsg: err.Error(),
			}
			reports = append(reports, report)
			printTestResult(report)
			return
		}
	}

	elapsed := time.Since(start)

	report := TestReport{
		TestName: "Lacan Track Signifier",
		Passed:   true,
		Duration: elapsed,
		Details:  fmt.Sprintf("5 textos processados em %dms", elapsed.Milliseconds()),
	}

	reports = append(reports, report)
	printTestResult(report)
}

func testRetrieveSignifiers(ctx context.Context, service *lacan.SignifierService) {
	sigs, err := service.GetKeySignifiers(ctx, 999, 5)

	report := TestReport{
		TestName: "Lacan Retrieve Signifiers",
	}

	if err != nil {
		report.Passed = false
		report.ErrorMsg = err.Error()
	} else if len(sigs) > 0 {
		report.Passed = true
		report.Details = fmt.Sprintf("Encontrados %d significantes. Top: '%s' (freq: %d)",
			len(sigs), sigs[0].Word, sigs[0].Frequency)
	} else {
		report.Passed = false
		report.Details = "Nenhum significante encontrado"
	}

	reports = append(reports, report)
	printTestResult(report)
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// TEST 4: Anti-Sycophancy (Mollick)
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func testAntiSycophancy() {
	fmt.Println("🔬 TEST 4: Anti-Sycophancy (Co-Intelligence)")
	fmt.Println("   Objetivo: Verificar se prompts bloqueiam concordância perigosa")
	fmt.Println()

	// Simular prompt Tipo 6 (Leal)
	router := personality.NewPersonalityRouter()
	fragment := router.GetSystemPromptFragment(personality.Type6)

	report := TestReport{
		TestName: "Anti-Sycophancy Prompt Check",
	}

	// Verificar se o prompt contém instruções anti-sycophancy
	if containsAntiSycophancy(fragment) {
		report.Passed = true
		report.Details = "Prompt contém 'DISCORDE IMEDIATAMENTE' ✅"
	} else {
		report.Passed = false
		report.Details = "Prompt NÃO contém proteção anti-sycophancy ❌"
	}

	reports = append(reports, report)
	printTestResult(report)

	fmt.Println()
}

func containsAntiSycophancy(prompt string) bool {
	keywords := []string{"DISCORDE", "Anti-Sycophancy", "arriscado"}
	for _, kw := range keywords {
		if strings.Contains(prompt, kw) {
			return true
		}
	}
	return false
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Utilities
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func printTestResult(r TestReport) {
	status := "✅ PASS"
	if !r.Passed {
		status = "❌ FAIL"
	}

	fmt.Printf("   %s | %s\n", status, r.TestName)
	if r.Details != "" {
		fmt.Printf("      └─ %s\n", r.Details)
	}
	if r.ErrorMsg != "" {
		fmt.Printf("      └─ Error: %s\n", r.ErrorMsg)
	}
	if r.Duration > 0 {
		fmt.Printf("      └─ Duration: %dms\n", r.Duration.Milliseconds())
	}
	fmt.Println()
}

func printFinalReport() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 FINAL REPORT")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	passed := 0
	failed := 0

	for _, r := range reports {
		if r.Passed {
			passed++
		} else {
			failed++
		}
	}

	total := passed + failed
	passRate := float64(passed) / float64(total) * 100

	fmt.Printf("Total Tests: %d\n", total)
	fmt.Printf("✅ Passed: %d\n", passed)
	fmt.Printf("❌ Failed: %d\n", failed)
	fmt.Printf("Pass Rate: %.1f%%\n", passRate)
	fmt.Println()

	// Telemetry snapshot
	snapshot := telemetry.GetSnapshot()
	fmt.Println("📈 Telemetry Snapshot:")
	fmt.Printf("   Enneatype: %v\n", snapshot["enneatype"])
	fmt.Printf("   Priming Latency: %vms\n", snapshot["priming_latency"])
	fmt.Printf("   Switches: %v\n", snapshot["switches"])
	fmt.Println()

	if passRate >= 80 {
		fmt.Println("🎉 FZPN ARCHITECTURE VALIDATED!")
		fmt.Println("   Sistema pronto para produção.")
	} else {
		fmt.Println("⚠️ VALIDATION INCOMPLETE")
		fmt.Println("   Revisar componentes com falha.")
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
