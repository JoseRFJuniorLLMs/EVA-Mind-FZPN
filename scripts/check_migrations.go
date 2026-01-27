// +build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:Debian23%40@104.248.219.200:5432/eva-db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping: %v", err)
	}

	fmt.Println("=== VERIFICACAO DE MIGRATIONS ===")
	fmt.Println()

	// Tabelas esperadas por sprint
	sprints := map[string][]string{
		"Sprint 1 (003 - Cognitive/Ethics)": {
			"interaction_cognitive_load",
			"cognitive_load_state",
			"cognitive_load_decisions",
			"ethical_boundary_events",
			"ethical_boundary_state",
			"ethical_redirections",
		},
		"Sprint 2 (004 - Explainability)": {
			"clinical_decision_explanations",
			"decision_factors",
			"prediction_accuracy_log",
		},
		"Sprint 3 (005 - Trajectory)": {
			"trajectory_simulations",
			"intervention_scenarios",
			"recommended_interventions",
			"trajectory_prediction_accuracy",
			"bayesian_network_parameters",
		},
		"Sprint 4 (008 - Multi-Persona)": {
			"persona_definitions",
			"persona_sessions",
			"persona_activation_rules",
			"persona_tool_permissions",
			"persona_transitions",
		},
		"Sprint 5 (007 - Research)": {
			"research_cohorts",
			"research_datapoints",
			"longitudinal_correlations",
			"statistical_analyses",
			"research_publications",
			"research_exports",
		},
	}

	for sprint, tables := range sprints {
		fmt.Printf("--- %s ---\n", sprint)
		found := 0
		for _, table := range tables {
			exists := tableExists(db, table)
			status := "❌"
			if exists {
				status = "✅"
				found++
			}
			fmt.Printf("  %s %s\n", status, table)
		}
		fmt.Printf("  Total: %d/%d\n\n", found, len(tables))
	}

	fmt.Println("=== FIM ===")
}

func tableExists(db *sql.DB, tableName string) bool {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = $1
	)`
	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
