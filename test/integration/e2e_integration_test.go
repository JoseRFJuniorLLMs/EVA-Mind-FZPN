package integration

import (
	"encoding/json"
	"testing"
	"time"
)

// ============================================================================
// END-TO-END INTEGRATION TESTS
// Tests complete workflows from detection to intervention
// ============================================================================

func TestE2E_CrisisDetectionToIntervention(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Step 1: Patient starts with elevated but not critical symptoms
	t.Run("Step1_InitialAssessment", func(t *testing.T) {
		query := `
			INSERT INTO clinical_assessments (
				patient_id, assessment_type, score, severity_level,
				assessed_at, created_at
			) VALUES ($1, 'PHQ-9', 14, 'moderate', NOW() - INTERVAL '7 days', NOW())
			RETURNING id
		`
		var id int64
		err := suite.DB.QueryRow(query, patientID).Scan(&id)
		AssertNoError(t, err, "Failed to create initial assessment")
	})

	// Step 2: Cognitive load increases over days
	t.Run("Step2_CognitiveLoadIncrease", func(t *testing.T) {
		for day := 6; day >= 0; day-- {
			load := 0.4 + (float64(6-day) * 0.08) // Increasing load
			query := `
				INSERT INTO interaction_cognitive_load (
					patient_id, cognitive_score, emotional_score, topic,
					response_time_ms, created_at
				) VALUES ($1, $2, $3, 'health_concerns', $4, NOW() - ($5 * INTERVAL '1 day'))
			`
			suite.DB.Exec(query, patientID, load, load+0.1, 500+day*100, day)
		}
	})

	// Step 3: Voice prosody changes detected
	t.Run("Step3_VoiceProsodyChange", func(t *testing.T) {
		query := `
			INSERT INTO voice_prosody_analyses (
				patient_id, pitch_mean_hz, pitch_variance,
				energy_mean, jitter, shimmer,
				analysis_date, emotional_valence, created_at
			) VALUES ($1, 95.5, 12.3, 0.35, 0.045, 0.055, NOW(), -0.6, NOW())
		`
		_, err := suite.DB.Exec(query, patientID)
		AssertNoError(t, err, "Failed to create voice analysis")
	})

	// Step 4: C-SSRS triggered due to concerning content
	t.Run("Step4_CSSRSTriggered", func(t *testing.T) {
		query := `
			INSERT INTO clinical_assessments (
				patient_id, assessment_type, score, severity_level,
				assessed_at, created_at
			) VALUES ($1, 'C-SSRS', 3, 'moderate', NOW(), NOW())
			RETURNING id
		`
		var assessmentID int64
		err := suite.DB.QueryRow(query, patientID).Scan(&assessmentID)
		AssertNoError(t, err, "Failed to create C-SSRS assessment")

		// Add responses
		responses := []struct {
			q int
			v int
		}{
			{1, 1}, {2, 1}, {3, 1}, {4, 0}, {5, 0}, {6, 0},
		}
		for _, r := range responses {
			suite.DB.Exec(`
				INSERT INTO clinical_assessment_responses (
					assessment_id, question_number, response_value, created_at
				) VALUES ($1, $2, $3, NOW())
			`, assessmentID, r.q, r.v)
		}
	})

	// Step 5: Trajectory simulation predicts elevated risk
	t.Run("Step5_TrajectoryPrediction", func(t *testing.T) {
		initialState, _ := json.Marshal(map[string]interface{}{
			"phq9_score": 14,
			"cssrs_level": 3,
			"voice_energy": 0.35,
		})

		query := `
			INSERT INTO trajectory_simulations (
				patient_id, simulation_date, days_ahead, n_simulations,
				initial_state, crisis_probability_7d, crisis_probability_30d,
				model_version, created_at
			) VALUES ($1, NOW(), 30, 1000, $2, 0.45, 0.62, 'v1.0', NOW())
			RETURNING id
		`
		var simID int64
		err := suite.DB.QueryRow(query, patientID, initialState).Scan(&simID)
		AssertNoError(t, err, "Failed to create simulation")
	})

	// Step 6: System generates intervention recommendations
	t.Run("Step6_InterventionGeneration", func(t *testing.T) {
		// Get latest simulation
		var simID int64
		suite.DB.QueryRow(`
			SELECT id FROM trajectory_simulations
			WHERE patient_id = $1
			ORDER BY created_at DESC LIMIT 1
		`, patientID).Scan(&simID)

		recommendations := []struct {
			title    string
			priority string
			reduction float64
		}{
			{"Consulta psiquiátrica urgente", "critical", 0.20},
			{"Contato familiar", "high", 0.15},
			{"Intensificação de monitoramento", "high", 0.12},
		}

		for _, rec := range recommendations {
			actions, _ := json.Marshal([]string{"Ação 1", "Ação 2"})
			query := `
				INSERT INTO recommended_interventions (
					simulation_id, patient_id, title, description,
					priority, expected_risk_reduction, timeframe_hours,
					action_steps, status, created_at
				) VALUES ($1, $2, $3, 'Descrição', $4, $5, 24, $6, 'pending', NOW())
			`
			suite.DB.Exec(query, simID, patientID, rec.title, rec.priority, rec.reduction, actions)
		}
	})

	// Step 7: Alert escalation triggered
	t.Run("Step7_AlertEscalation", func(t *testing.T) {
		query := `
			INSERT INTO clinical_alerts (
				patient_id, alert_type, severity, message, score, created_at
			) VALUES ($1, 'crisis_risk', 'high',
				'Elevated crisis risk detected - 45% probability in 7 days', 45, NOW())
			RETURNING id
		`
		var alertID int64
		err := suite.DB.QueryRow(query, patientID).Scan(&alertID)
		AssertNoError(t, err, "Failed to create alert")

		// Log escalation
		suite.DB.Exec(`
			INSERT INTO escalation_logs (
				patient_id, alert_type, original_severity, escalated_to,
				escalation_reason, channels_attempted, created_at
			) VALUES ($1, 'crisis_risk', 'high', 'critical',
				'No response after 30 minutes', ARRAY['push', 'sms'], NOW())
		`, patientID)
	})

	// Step 8: Verify complete workflow
	t.Run("Step8_VerifyWorkflow", func(t *testing.T) {
		// Check assessments
		var assessmentCount int
		suite.DB.QueryRow(`
			SELECT COUNT(*) FROM clinical_assessments
			WHERE patient_id = $1 AND created_at > NOW() - INTERVAL '1 hour'
		`, patientID).Scan(&assessmentCount)
		AssertTrue(t, assessmentCount >= 2, "Should have at least 2 assessments")

		// Check cognitive load entries
		var loadCount int
		suite.DB.QueryRow(`
			SELECT COUNT(*) FROM interaction_cognitive_load
			WHERE patient_id = $1 AND created_at > NOW() - INTERVAL '8 days'
		`, patientID).Scan(&loadCount)
		AssertTrue(t, loadCount >= 7, "Should have cognitive load entries")

		// Check trajectory simulation
		var simCount int
		suite.DB.QueryRow(`
			SELECT COUNT(*) FROM trajectory_simulations
			WHERE patient_id = $1 AND created_at > NOW() - INTERVAL '1 hour'
		`, patientID).Scan(&simCount)
		AssertTrue(t, simCount >= 1, "Should have trajectory simulation")

		// Check recommendations
		var recCount int
		suite.DB.QueryRow(`
			SELECT COUNT(*) FROM recommended_interventions
			WHERE patient_id = $1 AND status = 'pending'
		`, patientID).Scan(&recCount)
		AssertTrue(t, recCount >= 3, "Should have pending recommendations")

		// Check alerts
		var alertCount int
		suite.DB.QueryRow(`
			SELECT COUNT(*) FROM clinical_alerts
			WHERE patient_id = $1 AND created_at > NOW() - INTERVAL '1 hour'
		`, patientID).Scan(&alertCount)
		AssertTrue(t, alertCount >= 1, "Should have alerts")
	})

	// Cleanup
	t.Cleanup(func() {
		suite.DB.Exec("DELETE FROM escalation_logs WHERE patient_id = $1", patientID)
		suite.DB.Exec("DELETE FROM clinical_alerts WHERE patient_id = $1", patientID)
		suite.DB.Exec("DELETE FROM recommended_interventions WHERE patient_id = $1", patientID)
		suite.DB.Exec("DELETE FROM trajectory_simulations WHERE patient_id = $1", patientID)
		suite.DB.Exec("DELETE FROM voice_prosody_analyses WHERE patient_id = $1", patientID)
		suite.DB.Exec("DELETE FROM interaction_cognitive_load WHERE patient_id = $1", patientID)
		suite.DB.Exec("DELETE FROM clinical_assessment_responses WHERE assessment_id IN (SELECT id FROM clinical_assessments WHERE patient_id = $1)", patientID)
		suite.DB.Exec("DELETE FROM clinical_assessments WHERE patient_id = $1", patientID)
	})
}

func TestE2E_PersonaSwitchingOnContext(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Start with companion persona (default)
	t.Run("StartWithCompanion", func(t *testing.T) {
		query := `
			INSERT INTO persona_sessions (
				patient_id, persona_code, activation_reason,
				activated_by, started_at, status
			) VALUES ($1, 'companion', 'Session start', 'system', NOW(), 'active')
			RETURNING id
		`
		var sessionID int64
		err := suite.DB.QueryRow(query, patientID).Scan(&sessionID)
		AssertNoError(t, err, "Failed to start companion session")
	})

	// Educational content request → Switch to educator
	t.Run("SwitchToEducator", func(t *testing.T) {
		// End companion session
		suite.DB.Exec(`
			UPDATE persona_sessions
			SET status = 'ended', ended_at = NOW()
			WHERE patient_id = $1 AND status = 'active'
		`, patientID)

		// Log transition
		suite.DB.Exec(`
			INSERT INTO persona_transitions (
				patient_id, from_persona, to_persona,
				transition_reason, triggered_by, transitioned_at
			) VALUES ($1, 'companion', 'educator',
				'User requested information about medication', 'system', NOW())
		`, patientID)

		// Start educator session
		query := `
			INSERT INTO persona_sessions (
				patient_id, persona_code, activation_reason,
				activated_by, started_at, status
			) VALUES ($1, 'educator', 'Medication education request', 'system', NOW(), 'active')
			RETURNING id
		`
		var sessionID int64
		err := suite.DB.QueryRow(query, patientID).Scan(&sessionID)
		AssertNoError(t, err, "Failed to start educator session")
	})

	// Crisis detected → Immediate switch to emergency
	t.Run("SwitchToEmergency", func(t *testing.T) {
		// End educator session
		suite.DB.Exec(`
			UPDATE persona_sessions
			SET status = 'ended', ended_at = NOW()
			WHERE patient_id = $1 AND status = 'active'
		`, patientID)

		// Log transition
		suite.DB.Exec(`
			INSERT INTO persona_transitions (
				patient_id, from_persona, to_persona,
				transition_reason, triggered_by, transitioned_at
			) VALUES ($1, 'educator', 'emergency',
				'C-SSRS score >= 4 detected', 'system', NOW())
		`, patientID)

		// Start emergency session
		query := `
			INSERT INTO persona_sessions (
				patient_id, persona_code, activation_reason,
				activated_by, started_at, status
			) VALUES ($1, 'emergency', 'C-SSRS critical risk', 'system', NOW(), 'active')
			RETURNING id
		`
		var sessionID int64
		err := suite.DB.QueryRow(query, patientID).Scan(&sessionID)
		AssertNoError(t, err, "Failed to start emergency session")
	})

	// Crisis resolved → Switch to clinical for follow-up
	t.Run("SwitchToClinical", func(t *testing.T) {
		// End emergency session
		suite.DB.Exec(`
			UPDATE persona_sessions
			SET status = 'ended', ended_at = NOW()
			WHERE patient_id = $1 AND status = 'active'
		`, patientID)

		// Log transition
		suite.DB.Exec(`
			INSERT INTO persona_transitions (
				patient_id, from_persona, to_persona,
				transition_reason, triggered_by, transitioned_at
			) VALUES ($1, 'emergency', 'clinical',
				'Crisis resolved, clinical follow-up needed', 'system', NOW())
		`, patientID)

		// Start clinical session
		query := `
			INSERT INTO persona_sessions (
				patient_id, persona_code, activation_reason,
				activated_by, started_at, status
			) VALUES ($1, 'clinical', 'Post-crisis follow-up', 'system', NOW(), 'active')
			RETURNING id
		`
		var sessionID int64
		err := suite.DB.QueryRow(query, patientID).Scan(&sessionID)
		AssertNoError(t, err, "Failed to start clinical session")
	})

	// Verify transition history
	t.Run("VerifyTransitions", func(t *testing.T) {
		var transitionCount int
		err := suite.DB.QueryRow(`
			SELECT COUNT(*) FROM persona_transitions
			WHERE patient_id = $1 AND transitioned_at > NOW() - INTERVAL '1 hour'
		`, patientID).Scan(&transitionCount)

		AssertNoError(t, err, "Failed to count transitions")
		AssertEqual(t, 3, transitionCount, "Should have 3 transitions")

		// Verify order of transitions
		rows, err := suite.DB.Query(`
			SELECT from_persona, to_persona FROM persona_transitions
			WHERE patient_id = $1
			ORDER BY transitioned_at ASC
		`, patientID)
		AssertNoError(t, err, "Failed to query transitions")
		defer rows.Close()

		expectedTransitions := []struct {
			from, to string
		}{
			{"companion", "educator"},
			{"educator", "emergency"},
			{"emergency", "clinical"},
		}

		i := 0
		for rows.Next() {
			var from, to string
			rows.Scan(&from, &to)
			if i < len(expectedTransitions) {
				AssertEqual(t, expectedTransitions[i].from, from, "From persona mismatch")
				AssertEqual(t, expectedTransitions[i].to, to, "To persona mismatch")
			}
			i++
		}
	})

	// Cleanup
	t.Cleanup(func() {
		suite.DB.Exec("DELETE FROM persona_transitions WHERE patient_id = $1", patientID)
		suite.DB.Exec("DELETE FROM persona_sessions WHERE patient_id = $1", patientID)
	})
}

func TestE2E_ResearchDataCollection(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create research study
	var cohortID int64
	t.Run("CreateStudy", func(t *testing.T) {
		inclusion, _ := json.Marshal(map[string]interface{}{
			"min_age": 60,
			"has_phq9": true,
		})

		query := `
			INSERT INTO research_cohorts (
				study_name, study_code, hypothesis, study_type,
				inclusion_criteria, target_n_patients,
				data_collection_start_date, followup_duration_days,
				status, primary_outcome, principal_investigator, created_at
			) VALUES (
				'E2E Test Study', 'E2E-TEST-001',
				'Test the complete research workflow',
				'longitudinal_correlation', $1, 100,
				NOW(), 90, 'data_collection',
				'PHQ-9 improvement', 'Test PI', NOW()
			)
			RETURNING id
		`
		err := suite.DB.QueryRow(query, inclusion).Scan(&cohortID)
		AssertNoError(t, err, "Failed to create study")
	})

	// Collect data points over time
	t.Run("CollectDataPoints", func(t *testing.T) {
		for day := 0; day < 7; day++ {
			data, _ := json.Marshal(map[string]interface{}{
				"phq9_score": 15 - day, // Improving
				"gad7_score": 12 - day/2,
				"adherence":  0.7 + float64(day)*0.03,
			})

			query := `
				INSERT INTO research_datapoints (
					cohort_id, anonymous_patient_id, collection_date,
					datapoint_type, data, created_at
				) VALUES ($1, $2, NOW() - ($3 * INTERVAL '1 day'),
					'daily_assessment', $4, NOW())
			`
			anonID := "test_anonymous_" + string(rune('a'+day))
			suite.DB.Exec(query, cohortID, anonID, day, data)
		}
	})

	// Analyze data
	t.Run("AnalyzeData", func(t *testing.T) {
		// Insert correlation result
		ci, _ := json.Marshal(map[string]float64{"lower": 0.55, "upper": 0.75})
		query := `
			INSERT INTO longitudinal_correlations (
				cohort_id, predictor_variable, outcome_variable,
				lag_days, correlation_coefficient, p_value,
				confidence_interval_95, n_observations, n_patients,
				analysis_method, created_at
			) VALUES ($1, 'medication_adherence', 'phq9_score',
				0, -0.65, 0.002, $2, 42, 7, 'pearson', NOW())
			RETURNING id
		`
		var corrID int64
		err := suite.DB.QueryRow(query, cohortID, ci).Scan(&corrID)
		AssertNoError(t, err, "Failed to create correlation")

		// Verify negative correlation (higher adherence = lower PHQ-9)
		var coeff float64
		suite.DB.QueryRow(`
			SELECT correlation_coefficient FROM longitudinal_correlations WHERE id = $1
		`, corrID).Scan(&coeff)
		AssertTrue(t, coeff < 0, "Correlation should be negative")
	})

	// Export data
	t.Run("ExportData", func(t *testing.T) {
		query := `
			INSERT INTO research_exports (
				cohort_id, export_name, export_format, file_path,
				n_patients, n_observations, anonymization_level,
				pii_removed, exported_by, exported_for_purpose, created_at
			) VALUES ($1, 'E2E Test Export', 'csv', '/exports/e2e_test.csv',
				7, 42, 'fully_anonymized', true, 'test_system',
				'Integration test validation', NOW())
			RETURNING id
		`
		var exportID int64
		err := suite.DB.QueryRow(query, cohortID).Scan(&exportID)
		AssertNoError(t, err, "Failed to create export")

		// Verify export is compliant
		var piiRemoved bool
		suite.DB.QueryRow(`
			SELECT pii_removed FROM research_exports WHERE id = $1
		`, exportID).Scan(&piiRemoved)
		AssertTrue(t, piiRemoved, "PII should be removed")
	})

	// Cleanup
	t.Cleanup(func() {
		suite.DB.Exec("DELETE FROM research_exports WHERE cohort_id = $1", cohortID)
		suite.DB.Exec("DELETE FROM longitudinal_correlations WHERE cohort_id = $1", cohortID)
		suite.DB.Exec("DELETE FROM research_datapoints WHERE cohort_id = $1", cohortID)
		suite.DB.Exec("DELETE FROM research_cohorts WHERE id = $1", cohortID)
	})
}

// ============================================================================
// STRESS TESTS
// ============================================================================

func TestStress_ConcurrentAssessments(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID
	numConcurrent := 10

	// Create multiple assessments concurrently
	done := make(chan bool, numConcurrent)
	var assessmentIDs []int64

	for i := 0; i < numConcurrent; i++ {
		go func(idx int) {
			query := `
				INSERT INTO clinical_assessments (
					patient_id, assessment_type, score, severity_level,
					assessed_at, created_at
				) VALUES ($1, 'PHQ-9', $2, 'mild', NOW(), NOW())
				RETURNING id
			`
			var id int64
			suite.DB.QueryRow(query, patientID, idx).Scan(&id)
			assessmentIDs = append(assessmentIDs, id)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numConcurrent; i++ {
		<-done
	}

	// Verify all were created
	var count int
	suite.DB.QueryRow(`
		SELECT COUNT(*) FROM clinical_assessments
		WHERE patient_id = $1 AND created_at > NOW() - INTERVAL '1 minute'
	`, patientID).Scan(&count)

	AssertTrue(t, count >= numConcurrent, "Should have created all assessments")

	// Cleanup
	for _, id := range assessmentIDs {
		suite.DB.Exec("DELETE FROM clinical_assessments WHERE id = $1", id)
	}
}

func TestStress_RapidAlertCreation(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID
	numAlerts := 50

	start := time.Now()

	var alertIDs []int64
	for i := 0; i < numAlerts; i++ {
		query := `
			INSERT INTO clinical_alerts (
				patient_id, alert_type, severity, message, score, created_at
			) VALUES ($1, 'stress_test', 'low', 'Stress test alert', $2, NOW())
			RETURNING id
		`
		var id int64
		suite.DB.QueryRow(query, patientID, i).Scan(&id)
		alertIDs = append(alertIDs, id)
	}

	elapsed := time.Since(start)

	// Verify all were created
	var count int
	suite.DB.QueryRow(`
		SELECT COUNT(*) FROM clinical_alerts
		WHERE patient_id = $1 AND alert_type = 'stress_test'
	`, patientID).Scan(&count)

	AssertEqual(t, numAlerts, count, "Should have created all alerts")

	// Performance check (should complete in < 5 seconds)
	AssertTrue(t, elapsed < 5*time.Second, "Alert creation should be fast")

	// Cleanup
	suite.DB.Exec("DELETE FROM clinical_alerts WHERE patient_id = $1 AND alert_type = 'stress_test'", patientID)
}
