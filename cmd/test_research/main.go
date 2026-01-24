package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"eva-mind/internal/research"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("🧬 Clinical Research Engine - Test")
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

	// Criar Research Engine
	engine := research.NewResearchEngine(db)

	// ========================================================================
	// FASE 1: CRIAR ESTUDOS PRÉ-CONFIGURADOS
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📚 FASE 1: Criando Estudos Pré-configurados")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	err = engine.CreatePreconfiguredStudies()
	if err != nil {
		log.Printf("⚠️ Erro ao criar estudos: %v", err)
	}

	fmt.Println()

	// ========================================================================
	// FASE 2: COLETAR DADOS PARA ESTUDO 1 (Voice → PHQ-9)
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📊 FASE 2: Coletando Dados para Estudo 1 (Voice → PHQ-9)")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	// Buscar ID do estudo
	cohort1ID, err := getCohortIDByCode(db, "EVA-VOICE-PHQ9-001")
	if err != nil {
		log.Fatalf("❌ Erro ao buscar coorte: %v", err)
	}

	fmt.Printf("Coorte ID: %s\n\n", cohort1ID)

	// Coletar dados (AVISO: pode demorar!)
	fmt.Println("⏳ Coletando e anonimizando dados longitudinais...")
	fmt.Println("   (Isso pode levar alguns minutos dependendo do volume de dados)")
	fmt.Println()

	err = engine.CollectDataForCohort(cohort1ID)
	if err != nil {
		log.Printf("⚠️ Erro ao coletar dados: %v", err)
	}

	fmt.Println()

	// ========================================================================
	// FASE 3: EXECUTAR LAG CORRELATION ANALYSIS
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("🔬 FASE 3: Executando Lag Correlation Analysis")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	fmt.Println("Analisando: Voice Pitch → PHQ-9 (lag 0-14 dias)")
	fmt.Println()

	err = engine.RunLagCorrelationAnalysis(
		cohort1ID,
		"voice_pitch_mean",
		"phq9",
		14, // Testar lags de 0 a 14 dias
	)

	if err != nil {
		log.Printf("⚠️ Erro na análise: %v", err)
	}

	fmt.Println()

	// ========================================================================
	// FASE 4: VISUALIZAR RESULTADOS
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📈 FASE 4: Resultados da Análise")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	// Buscar correlações significativas
	significantCorrelations, err := getSignificantCorrelations(db, cohort1ID)
	if err != nil {
		log.Printf("⚠️ Erro ao buscar correlações: %v", err)
	}

	if len(significantCorrelations) == 0 {
		fmt.Println("⚠️ Nenhuma correlação significativa encontrada.")
		fmt.Println("   Isso pode acontecer se:")
		fmt.Println("   - Não há dados suficientes no banco")
		fmt.Println("   - Os dados não têm correlação temporal real")
		fmt.Println("   - O período de followup é muito curto")
	} else {
		fmt.Printf("✅ %d correlações significativas encontradas:\n\n", len(significantCorrelations))

		for i, corr := range significantCorrelations {
			printCorrelation(i+1, corr)
		}
	}

	fmt.Println()

	// ========================================================================
	// FASE 5: GERAR RELATÓRIO DO ESTUDO
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📄 FASE 5: Relatório do Estudo")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	report, err := engine.GenerateStudyReport(cohort1ID)
	if err != nil {
		log.Printf("⚠️ Erro ao gerar relatório: %v", err)
	} else {
		printReport(report)
	}

	fmt.Println()

	// ========================================================================
	// FASE 6: VERIFICAR STATUS DE TODOS OS ESTUDOS
	// ========================================================================
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("📊 FASE 6: Status de Todos os Estudos")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	studies, err := getAllStudies(db)
	if err != nil {
		log.Printf("⚠️ Erro ao buscar estudos: %v", err)
	} else {
		for i, study := range studies {
			printStudyStatus(i+1, study)
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("✅ Teste do Research Engine completo")
	fmt.Println(strings.Repeat("=", 70))
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func getCohortIDByCode(db *sql.DB, code string) (string, error) {
	var id string
	query := `SELECT id FROM research_cohorts WHERE study_code = $1`
	err := db.QueryRow(query, code).Scan(&id)
	return id, err
}

type Correlation struct {
	Predictor    string
	Outcome      string
	LagDays      int
	Coefficient  float64
	PValue       float64
	NObservations int
	NPatients    int
}

func getSignificantCorrelations(db *sql.DB, cohortID string) ([]Correlation, error) {
	query := `
		SELECT
			predictor_variable,
			outcome_variable,
			lag_days,
			correlation_coefficient,
			p_value,
			n_observations,
			n_patients
		FROM longitudinal_correlations
		WHERE cohort_id = $1
		  AND is_significant = TRUE
		ORDER BY ABS(correlation_coefficient) DESC
	`

	rows, err := db.Query(query, cohortID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	correlations := []Correlation{}
	for rows.Next() {
		var c Correlation
		err := rows.Scan(
			&c.Predictor, &c.Outcome, &c.LagDays,
			&c.Coefficient, &c.PValue,
			&c.NObservations, &c.NPatients,
		)
		if err != nil {
			continue
		}
		correlations = append(correlations, c)
	}

	return correlations, nil
}

func printCorrelation(num int, corr Correlation) {
	direction := "↑"
	if corr.Coefficient < 0 {
		direction = "↓"
	}

	effectSize := "pequeno"
	if corr.Coefficient > 0.5 || corr.Coefficient < -0.5 {
		effectSize = "grande"
	} else if corr.Coefficient > 0.3 || corr.Coefficient < -0.3 {
		effectSize = "médio"
	}

	fmt.Printf("%d. %s → %s (lag %d dias)\n", num, corr.Predictor, corr.Outcome, corr.LagDays)
	fmt.Printf("   %s Correlação: r = %.3f (efeito %s)\n", direction, corr.Coefficient, effectSize)
	fmt.Printf("   Significância: p = %.6f\n", corr.PValue)
	fmt.Printf("   Dados: %d observações, %d pacientes\n", corr.NObservations, corr.NPatients)

	// Interpretação clínica
	if corr.Predictor == "voice_pitch_mean" && corr.Outcome == "phq9" {
		if corr.Coefficient < 0 {
			fmt.Printf("   💡 Interpretação: Queda no pitch vocal PREDIZ piora no PHQ-9 após %d dias\n", corr.LagDays)
		} else {
			fmt.Printf("   💡 Interpretação: Aumento no pitch vocal PREDIZ melhora no PHQ-9 após %d dias\n", corr.LagDays)
		}
	}

	fmt.Println()
}

type Study struct {
	Name     string
	Code     string
	Status   string
	NPatients int
	Target   int
}

func getAllStudies(db *sql.DB) ([]Study, error) {
	query := `
		SELECT study_name, study_code, status, current_n_patients, target_n_patients
		FROM research_cohorts
		ORDER BY created_at
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	studies := []Study{}
	for rows.Next() {
		var s Study
		err := rows.Scan(&s.Name, &s.Code, &s.Status, &s.NPatients, &s.Target)
		if err != nil {
			continue
		}
		studies = append(studies, s)
	}

	return studies, nil
}

func printStudyStatus(num int, study Study) {
	statusIcon := "🔄"
	switch study.Status {
	case "completed":
		statusIcon = "✅"
	case "published":
		statusIcon = "📄"
	case "recruiting":
		statusIcon = "🔍"
	case "analyzing":
		statusIcon = "🔬"
	}

	progress := float64(study.NPatients) / float64(study.Target) * 100

	fmt.Printf("%s %d. [%s] %s\n", statusIcon, num, study.Code, study.Name)
	fmt.Printf("   Status: %s\n", study.Status)
	fmt.Printf("   Pacientes: %d/%d (%.1f%%)\n", study.NPatients, study.Target, progress)
	fmt.Println()
}

func printReport(report map[string]interface{}) {
	if study, ok := report["study"].(map[string]interface{}); ok {
		fmt.Printf("Estudo: %s\n", study["name"])
		fmt.Printf("Código: %s\n", study["code"])
		fmt.Printf("Hipótese: %s\n", study["hypothesis"])
		fmt.Printf("Status: %s\n", study["status"])
		fmt.Printf("N Pacientes: %.0f\n", study["n_patients"])
		fmt.Println()
	}

	if correlations, ok := report["significant_correlations"].([]interface{}); ok && len(correlations) > 0 {
		fmt.Println("Correlações Significativas:")
		for _, corrInterface := range correlations {
			if corr, ok := corrInterface.(map[string]interface{}); ok {
				fmt.Printf("  • %s → %s (lag %.0f dias): r=%.3f, p=%.6f\n",
					corr["predictor"], corr["outcome"],
					corr["lag_days"], corr["r"], corr["p"])
			}
		}
		fmt.Println()
	}

	if analyses, ok := report["analyses"].([]interface{}); ok && len(analyses) > 0 {
		fmt.Println("Análises Realizadas:")
		for _, analysisInterface := range analyses {
			if analysis, ok := analysisInterface.(map[string]interface{}); ok {
				fmt.Printf("  • %s: %s\n", analysis["type"], analysis["name"])
			}
		}
		fmt.Println()
	}
}
