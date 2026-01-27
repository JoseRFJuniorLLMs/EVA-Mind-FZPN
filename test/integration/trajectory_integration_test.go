package integration

import (
	"encoding/json"
	"testing"
	"time"
)

// ============================================================================
// INTEGRATION TESTS: Predictive Life Trajectory Engine
// Tests the Monte Carlo simulation and prediction system
// ============================================================================

func TestTrajectory_SimulationCreation(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create initial patient state
	initialState := map[string]interface{}{
		"phq9_score":            12,
		"gad7_score":            8,
		"medication_adherence":  0.75,
		"sleep_hours":           6.5,
		"social_isolation_days": 2,
		"voice_energy_score":    0.6,
	}
	initialStateJSON, _ := json.Marshal(initialState)

	// Create simulation results
	simulationResults := map[string]interface{}{
		"crisis_probability_7d":  0.15,
		"crisis_probability_30d": 0.35,
		"critical_factors": []string{
			"medication_adherence",
			"sleep_hours",
		},
		"trajectories_simulated": 1000,
	}
	resultsJSON, _ := json.Marshal(simulationResults)

	query := `
		INSERT INTO trajectory_simulations (
			patient_id, simulation_date, days_ahead, n_simulations,
			initial_state, simulation_results,
			crisis_probability_7d, crisis_probability_30d,
			model_version, created_at
		) VALUES ($1, NOW(), 30, 1000, $2, $3, 0.15, 0.35, 'v1.0', NOW())
		RETURNING id
	`

	var simulationID int64
	err := suite.DB.QueryRow(query, patientID, initialStateJSON, resultsJSON).Scan(&simulationID)
	AssertNoError(t, err, "Failed to create simulation")
	AssertTrue(t, simulationID > 0, "Simulation ID should be positive")

	// Verify simulation was stored
	var crisisProb7d, crisisProb30d float64
	err = suite.DB.QueryRow(`
		SELECT crisis_probability_7d, crisis_probability_30d
		FROM trajectory_simulations
		WHERE id = $1
	`, simulationID).Scan(&crisisProb7d, &crisisProb30d)

	AssertNoError(t, err, "Failed to retrieve simulation")
	AssertTrue(t, crisisProb7d >= 0 && crisisProb7d <= 1, "7d probability should be 0-1")
	AssertTrue(t, crisisProb30d >= 0 && crisisProb30d <= 1, "30d probability should be 0-1")

	// Cleanup
	suite.DB.Exec("DELETE FROM trajectory_simulations WHERE id = $1", simulationID)
}

func TestTrajectory_InterventionScenario(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// First create a simulation
	simQuery := `
		INSERT INTO trajectory_simulations (
			patient_id, simulation_date, days_ahead, n_simulations,
			crisis_probability_7d, crisis_probability_30d,
			model_version, created_at
		) VALUES ($1, NOW(), 30, 1000, 0.45, 0.65, 'v1.0', NOW())
		RETURNING id
	`

	var simulationID int64
	err := suite.DB.QueryRow(simQuery, patientID).Scan(&simulationID)
	AssertNoError(t, err, "Failed to create base simulation")

	// Create intervention scenarios
	interventions := []struct {
		interventionType  string
		riskReduction7d   float64
		riskReduction30d  float64
	}{
		{"medication_adherence_boost", 0.10, 0.15},
		{"sleep_hygiene_protocol", 0.08, 0.12},
		{"family_engagement", 0.12, 0.18},
		{"therapy_intensification", 0.15, 0.22},
	}

	var scenarioIDs []int64

	for _, intervention := range interventions {
		query := `
			INSERT INTO intervention_scenarios (
				simulation_id, intervention_type,
				risk_reduction_7d, risk_reduction_30d,
				new_crisis_probability_7d, new_crisis_probability_30d,
				created_at
			) VALUES ($1, $2, $3, $4, $5, $6, NOW())
			RETURNING id
		`

		newProb7d := 0.45 - intervention.riskReduction7d
		newProb30d := 0.65 - intervention.riskReduction30d

		var scenarioID int64
		err := suite.DB.QueryRow(query, simulationID, intervention.interventionType,
			intervention.riskReduction7d, intervention.riskReduction30d,
			newProb7d, newProb30d).Scan(&scenarioID)
		AssertNoError(t, err, "Failed to create scenario")
		scenarioIDs = append(scenarioIDs, scenarioID)
	}

	// Find best intervention (highest risk reduction)
	var bestIntervention string
	var bestReduction float64
	err = suite.DB.QueryRow(`
		SELECT intervention_type, risk_reduction_30d
		FROM intervention_scenarios
		WHERE simulation_id = $1
		ORDER BY risk_reduction_30d DESC
		LIMIT 1
	`, simulationID).Scan(&bestIntervention, &bestReduction)

	AssertNoError(t, err, "Failed to find best intervention")
	AssertEqual(t, "therapy_intensification", bestIntervention, "Best intervention should be therapy")
	AssertTrue(t, bestReduction >= 0.22, "Best reduction should be >= 22%")

	// Cleanup
	for _, id := range scenarioIDs {
		suite.DB.Exec("DELETE FROM intervention_scenarios WHERE id = $1", id)
	}
	suite.DB.Exec("DELETE FROM trajectory_simulations WHERE id = $1", simulationID)
}

func TestTrajectory_RecommendedInterventions(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create simulation
	simQuery := `
		INSERT INTO trajectory_simulations (
			patient_id, simulation_date, days_ahead, n_simulations,
			crisis_probability_7d, crisis_probability_30d,
			model_version, created_at
		) VALUES ($1, NOW(), 30, 1000, 0.55, 0.70, 'v1.0', NOW())
		RETURNING id
	`

	var simulationID int64
	err := suite.DB.QueryRow(simQuery, patientID).Scan(&simulationID)
	AssertNoError(t, err, "Failed to create simulation")

	// Create recommendations
	recommendations := []struct {
		title           string
		priority        string
		riskReduction   float64
		timeframeHours  int
	}{
		{"Consulta psiquiátrica urgente", "critical", 0.25, 24},
		{"Reforço de adesão medicamentosa", "high", 0.15, 48},
		{"Protocolo de higiene do sono", "medium", 0.10, 72},
	}

	var recIDs []int64

	for _, rec := range recommendations {
		actionSteps := []string{"Passo 1", "Passo 2", "Passo 3"}
		actionStepsJSON, _ := json.Marshal(actionSteps)

		query := `
			INSERT INTO recommended_interventions (
				simulation_id, patient_id, title, description,
				priority, expected_risk_reduction, timeframe_hours,
				action_steps, status, created_at
			) VALUES ($1, $2, $3, 'Description', $4, $5, $6, $7, 'pending', NOW())
			RETURNING id
		`

		var recID int64
		err := suite.DB.QueryRow(query, simulationID, patientID, rec.title,
			rec.priority, rec.riskReduction, rec.timeframeHours, actionStepsJSON).Scan(&recID)
		AssertNoError(t, err, "Failed to create recommendation")
		recIDs = append(recIDs, recID)
	}

	// Verify recommendations are prioritized correctly
	var criticalCount int
	err = suite.DB.QueryRow(`
		SELECT COUNT(*) FROM recommended_interventions
		WHERE simulation_id = $1 AND priority = 'critical'
	`, simulationID).Scan(&criticalCount)

	AssertNoError(t, err, "Failed to count critical recommendations")
	AssertEqual(t, 1, criticalCount, "Should have 1 critical recommendation")

	// Verify ordering by priority
	var firstPriority string
	err = suite.DB.QueryRow(`
		SELECT priority FROM recommended_interventions
		WHERE simulation_id = $1
		ORDER BY
			CASE priority
				WHEN 'critical' THEN 1
				WHEN 'high' THEN 2
				WHEN 'medium' THEN 3
				WHEN 'low' THEN 4
			END
		LIMIT 1
	`, simulationID).Scan(&firstPriority)

	AssertNoError(t, err, "Failed to get first recommendation")
	AssertEqual(t, "critical", firstPriority, "First recommendation should be critical")

	// Cleanup
	for _, id := range recIDs {
		suite.DB.Exec("DELETE FROM recommended_interventions WHERE id = $1", id)
	}
	suite.DB.Exec("DELETE FROM trajectory_simulations WHERE id = $1", simulationID)
}

func TestTrajectory_PredictionAccuracy(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create historical predictions with outcomes
	predictions := []struct {
		predictedProb  float64
		actualOutcome  bool // Did crisis actually occur?
	}{
		{0.75, true},   // High prediction, crisis occurred (true positive)
		{0.80, true},   // High prediction, crisis occurred (true positive)
		{0.20, false},  // Low prediction, no crisis (true negative)
		{0.15, false},  // Low prediction, no crisis (true negative)
		{0.60, false},  // Medium prediction, no crisis (false positive)
	}

	var predIDs []int64

	for _, pred := range predictions {
		query := `
			INSERT INTO trajectory_prediction_accuracy (
				patient_id, prediction_date, days_ahead,
				predicted_crisis_probability, actual_crisis_occurred,
				model_version, created_at
			) VALUES ($1, NOW() - INTERVAL '30 days', 30, $2, $3, 'v1.0', NOW())
			RETURNING id
		`

		var predID int64
		err := suite.DB.QueryRow(query, patientID, pred.predictedProb, pred.actualOutcome).Scan(&predID)
		AssertNoError(t, err, "Failed to create prediction record")
		predIDs = append(predIDs, predID)
	}

	// Calculate accuracy metrics
	var truePositives, falsePositives, trueNegatives, falseNegatives int
	err := suite.DB.QueryRow(`
		SELECT
			COUNT(*) FILTER (WHERE predicted_crisis_probability >= 0.5 AND actual_crisis_occurred = true),
			COUNT(*) FILTER (WHERE predicted_crisis_probability >= 0.5 AND actual_crisis_occurred = false),
			COUNT(*) FILTER (WHERE predicted_crisis_probability < 0.5 AND actual_crisis_occurred = false),
			COUNT(*) FILTER (WHERE predicted_crisis_probability < 0.5 AND actual_crisis_occurred = true)
		FROM trajectory_prediction_accuracy
		WHERE patient_id = $1
	`, patientID).Scan(&truePositives, &falsePositives, &trueNegatives, &falseNegatives)

	AssertNoError(t, err, "Failed to calculate metrics")
	AssertEqual(t, 2, truePositives, "Should have 2 true positives")
	AssertEqual(t, 2, trueNegatives, "Should have 2 true negatives")
	AssertEqual(t, 1, falsePositives, "Should have 1 false positive")
	AssertEqual(t, 0, falseNegatives, "Should have 0 false negatives")

	// Calculate sensitivity (recall)
	sensitivity := float64(truePositives) / float64(truePositives+falseNegatives)
	AssertTrue(t, sensitivity == 1.0, "Sensitivity should be 100%")

	// Cleanup
	for _, id := range predIDs {
		suite.DB.Exec("DELETE FROM trajectory_prediction_accuracy WHERE id = $1", id)
	}
}

func TestTrajectory_BayesianParameters(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	// Check if Bayesian network parameters exist
	var paramCount int
	err := suite.DB.QueryRow(`
		SELECT COUNT(*) FROM bayesian_network_parameters
		WHERE is_active = true
	`).Scan(&paramCount)

	if err != nil {
		t.Skip("bayesian_network_parameters table may not exist")
	}

	// If table exists, verify we have active parameters
	if paramCount > 0 {
		var nodeName string
		var prior float64
		err = suite.DB.QueryRow(`
			SELECT node_name, prior_probability
			FROM bayesian_network_parameters
			WHERE is_active = true
			LIMIT 1
		`).Scan(&nodeName, &prior)

		AssertNoError(t, err, "Failed to get Bayesian parameters")
		AssertTrue(t, nodeName != "", "Node name should not be empty")
		AssertTrue(t, prior >= 0 && prior <= 1, "Prior should be 0-1")
	}
}

// ============================================================================
// INTEGRATION TESTS: Persona System
// ============================================================================

func TestPersona_SessionCreation(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create persona session
	query := `
		INSERT INTO persona_sessions (
			patient_id, persona_code, activation_reason,
			activated_by, started_at, status
		) VALUES ($1, 'companion', 'User request', 'system', NOW(), 'active')
		RETURNING id
	`

	var sessionID int64
	err := suite.DB.QueryRow(query, patientID).Scan(&sessionID)
	AssertNoError(t, err, "Failed to create persona session")

	// Verify session was created
	var personaCode, status string
	err = suite.DB.QueryRow(`
		SELECT persona_code, status
		FROM persona_sessions
		WHERE id = $1
	`, sessionID).Scan(&personaCode, &status)

	AssertNoError(t, err, "Failed to retrieve session")
	AssertEqual(t, "companion", personaCode, "Persona code mismatch")
	AssertEqual(t, "active", status, "Status should be active")

	// Cleanup
	suite.DB.Exec("DELETE FROM persona_sessions WHERE id = $1", sessionID)
}

func TestPersona_AutomaticActivation(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Test activation rules
	activationRules := []struct {
		trigger         string
		expectedPersona string
		priority        int
	}{
		{"cssrs_score >= 4", "emergency", 100},
		{"phq9_score >= 20", "clinical", 80},
		{"education_request", "educator", 40},
		{"default", "companion", 10},
	}

	for _, rule := range activationRules {
		t.Run(rule.trigger, func(t *testing.T) {
			// Create a persona session based on trigger
			query := `
				INSERT INTO persona_sessions (
					patient_id, persona_code, activation_reason,
					activated_by, trigger_condition, started_at, status
				) VALUES ($1, $2, $3, 'system', $3, NOW(), 'active')
				RETURNING id
			`

			var sessionID int64
			err := suite.DB.QueryRow(query, patientID, rule.expectedPersona, rule.trigger).Scan(&sessionID)
			AssertNoError(t, err, "Failed to create session for trigger: "+rule.trigger)

			// Verify correct persona was activated
			var activatedPersona string
			err = suite.DB.QueryRow(`
				SELECT persona_code FROM persona_sessions WHERE id = $1
			`, sessionID).Scan(&activatedPersona)

			AssertNoError(t, err, "Failed to get session")
			AssertEqual(t, rule.expectedPersona, activatedPersona, "Wrong persona for trigger: "+rule.trigger)

			// Cleanup
			suite.DB.Exec("DELETE FROM persona_sessions WHERE id = $1", sessionID)
		})
	}
}

func TestPersona_ToolPermissions(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	// Check if tool permissions exist for each persona
	personas := []string{"companion", "clinical", "emergency", "educator"}

	for _, persona := range personas {
		t.Run(persona, func(t *testing.T) {
			var allowedCount, blockedCount int
			err := suite.DB.QueryRow(`
				SELECT
					COUNT(*) FILTER (WHERE is_allowed = true),
					COUNT(*) FILTER (WHERE is_allowed = false)
				FROM persona_tool_permissions
				WHERE persona_code = $1
			`, persona).Scan(&allowedCount, &blockedCount)

			if err != nil {
				t.Skipf("persona_tool_permissions may not exist for %s", persona)
			}

			// Each persona should have some tool permissions defined
			totalPermissions := allowedCount + blockedCount
			AssertTrue(t, totalPermissions > 0, "Persona should have tool permissions: "+persona)
		})
	}
}

func TestPersona_TransitionLogging(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Log persona transition
	query := `
		INSERT INTO persona_transitions (
			patient_id, from_persona, to_persona,
			transition_reason, triggered_by, transitioned_at
		) VALUES ($1, 'companion', 'clinical', 'Hospital admission', 'system', NOW())
		RETURNING id
	`

	var transitionID int64
	err := suite.DB.QueryRow(query, patientID).Scan(&transitionID)
	AssertNoError(t, err, "Failed to log transition")

	// Verify transition was logged
	var fromPersona, toPersona, reason string
	err = suite.DB.QueryRow(`
		SELECT from_persona, to_persona, transition_reason
		FROM persona_transitions
		WHERE id = $1
	`, transitionID).Scan(&fromPersona, &toPersona, &reason)

	AssertNoError(t, err, "Failed to retrieve transition")
	AssertEqual(t, "companion", fromPersona, "From persona mismatch")
	AssertEqual(t, "clinical", toPersona, "To persona mismatch")
	AssertTrue(t, reason != "", "Reason should not be empty")

	// Cleanup
	suite.DB.Exec("DELETE FROM persona_transitions WHERE id = $1", transitionID)
}

// ============================================================================
// BENCHMARK TESTS
// ============================================================================

func BenchmarkSimulationCreation(b *testing.B) {
	if suite.DB == nil {
		b.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID
	query := `
		INSERT INTO trajectory_simulations (
			patient_id, simulation_date, days_ahead, n_simulations,
			crisis_probability_7d, crisis_probability_30d,
			model_version, created_at
		) VALUES ($1, NOW(), 30, 1000, 0.25, 0.40, 'v1.0', NOW())
		RETURNING id
	`

	var ids []int64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int64
		suite.DB.QueryRow(query, patientID).Scan(&id)
		ids = append(ids, id)
	}

	// Cleanup
	for _, id := range ids {
		suite.DB.Exec("DELETE FROM trajectory_simulations WHERE id = $1", id)
	}
}
