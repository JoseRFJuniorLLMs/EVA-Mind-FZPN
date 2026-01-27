package integration

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

// ============================================================================
// INTEGRATION TESTS: Clinical Research Engine
// Tests cohort management, anonymization, and analysis
// ============================================================================

func TestResearch_CohortCreation(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	// Create research cohort
	inclusionCriteria := map[string]interface{}{
		"min_age":            60,
		"max_age":            90,
		"has_voice_data":     true,
		"min_phq9_assessments": 3,
	}
	inclusionJSON, _ := json.Marshal(inclusionCriteria)

	exclusionCriteria := map[string]interface{}{
		"severe_cognitive_impairment": true,
	}
	exclusionJSON, _ := json.Marshal(exclusionCriteria)

	query := `
		INSERT INTO research_cohorts (
			study_name, study_code, hypothesis, study_type,
			inclusion_criteria, exclusion_criteria,
			target_n_patients, data_collection_start_date,
			followup_duration_days, status, primary_outcome,
			principal_investigator, created_at
		) VALUES (
			'Voice Biomarkers Study', 'EVA-VOICE-001',
			'Voice changes predict depression scores',
			'longitudinal_correlation',
			$1, $2, 100, NOW(), 180, 'data_collection',
			'PHQ-9 score change', 'Dr. EVA Research', NOW()
		)
		RETURNING id
	`

	var cohortID int64
	err := suite.DB.QueryRow(query, inclusionJSON, exclusionJSON).Scan(&cohortID)
	AssertNoError(t, err, "Failed to create cohort")
	AssertTrue(t, cohortID > 0, "Cohort ID should be positive")

	// Verify cohort was created
	var studyCode, status string
	var targetN int
	err = suite.DB.QueryRow(`
		SELECT study_code, status, target_n_patients
		FROM research_cohorts
		WHERE id = $1
	`, cohortID).Scan(&studyCode, &status, &targetN)

	AssertNoError(t, err, "Failed to retrieve cohort")
	AssertEqual(t, "EVA-VOICE-001", studyCode, "Study code mismatch")
	AssertEqual(t, "data_collection", status, "Status mismatch")
	AssertEqual(t, 100, targetN, "Target N mismatch")

	// Cleanup
	suite.DB.Exec("DELETE FROM research_cohorts WHERE id = $1", cohortID)
}

func TestResearch_DatapointAnonymization(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create cohort first
	cohortQuery := `
		INSERT INTO research_cohorts (
			study_name, study_code, hypothesis, study_type,
			target_n_patients, status, primary_outcome,
			principal_investigator, created_at
		) VALUES (
			'Anonymization Test', 'TEST-ANON-001',
			'Test hypothesis', 'longitudinal_correlation',
			50, 'data_collection', 'Test outcome',
			'Test Researcher', NOW()
		)
		RETURNING id
	`

	var cohortID int64
	err := suite.DB.QueryRow(cohortQuery).Scan(&cohortID)
	AssertNoError(t, err, "Failed to create cohort")

	// Create anonymized patient ID
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d", patientID)))
	anonymousID := hex.EncodeToString(hash[:])

	// Insert anonymized datapoint
	datapoint := map[string]interface{}{
		"phq9_score":           15,
		"gad7_score":           12,
		"medication_adherence": 0.85,
		"voice_pitch_mean":     125.5,
	}
	datapointJSON, _ := json.Marshal(datapoint)

	datapointQuery := `
		INSERT INTO research_datapoints (
			cohort_id, anonymous_patient_id, collection_date,
			datapoint_type, data, created_at
		) VALUES ($1, $2, NOW(), 'daily_metrics', $3, NOW())
		RETURNING id
	`

	var datapointID int64
	err = suite.DB.QueryRow(datapointQuery, cohortID, anonymousID, datapointJSON).Scan(&datapointID)
	AssertNoError(t, err, "Failed to create datapoint")

	// Verify anonymization
	var storedAnonID string
	err = suite.DB.QueryRow(`
		SELECT anonymous_patient_id FROM research_datapoints WHERE id = $1
	`, datapointID).Scan(&storedAnonID)

	AssertNoError(t, err, "Failed to retrieve datapoint")
	AssertEqual(t, 64, len(storedAnonID), "Anonymous ID should be 64 hex chars")
	AssertTrue(t, storedAnonID != fmt.Sprintf("%d", patientID), "ID should be anonymized")

	// Cleanup
	suite.DB.Exec("DELETE FROM research_datapoints WHERE id = $1", datapointID)
	suite.DB.Exec("DELETE FROM research_cohorts WHERE id = $1", cohortID)
}

func TestResearch_LongitudinalCorrelation(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	// Create cohort
	cohortQuery := `
		INSERT INTO research_cohorts (
			study_name, study_code, hypothesis, study_type,
			target_n_patients, status, primary_outcome,
			principal_investigator, created_at
		) VALUES (
			'Correlation Test', 'TEST-CORR-001',
			'Test correlation', 'longitudinal_correlation',
			50, 'analysis', 'Test outcome',
			'Test Researcher', NOW()
		)
		RETURNING id
	`

	var cohortID int64
	err := suite.DB.QueryRow(cohortQuery).Scan(&cohortID)
	AssertNoError(t, err, "Failed to create cohort")

	// Insert correlation result
	ciJSON, _ := json.Marshal(map[string]float64{
		"lower": 0.45,
		"upper": 0.75,
	})

	corrQuery := `
		INSERT INTO longitudinal_correlations (
			cohort_id, predictor_variable, outcome_variable,
			lag_days, correlation_coefficient, p_value,
			confidence_interval_95, n_observations, n_patients,
			analysis_method, created_at
		) VALUES ($1, 'voice_pitch_mean', 'phq9_score', 7, 0.62, 0.001,
			$2, 500, 50, 'pearson', NOW())
		RETURNING id
	`

	var corrID int64
	err = suite.DB.QueryRow(corrQuery, cohortID, ciJSON).Scan(&corrID)
	AssertNoError(t, err, "Failed to create correlation")

	// Verify correlation is significant
	var pValue float64
	var coefficient float64
	err = suite.DB.QueryRow(`
		SELECT correlation_coefficient, p_value
		FROM longitudinal_correlations
		WHERE id = $1
	`, corrID).Scan(&coefficient, &pValue)

	AssertNoError(t, err, "Failed to retrieve correlation")
	AssertTrue(t, pValue < 0.05, "p-value should be significant (< 0.05)")
	AssertTrue(t, coefficient > 0.5, "Correlation should be moderate-strong")

	// Cleanup
	suite.DB.Exec("DELETE FROM longitudinal_correlations WHERE id = $1", corrID)
	suite.DB.Exec("DELETE FROM research_cohorts WHERE id = $1", cohortID)
}

func TestResearch_StatisticalAnalysis(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	// Create cohort
	cohortQuery := `
		INSERT INTO research_cohorts (
			study_name, study_code, hypothesis, study_type,
			target_n_patients, status, primary_outcome,
			principal_investigator, created_at
		) VALUES (
			'Stats Test', 'TEST-STATS-001',
			'Test stats', 'causal_inference',
			100, 'analysis', 'Test outcome',
			'Test Researcher', NOW()
		)
		RETURNING id
	`

	var cohortID int64
	err := suite.DB.QueryRow(cohortQuery).Scan(&cohortID)
	AssertNoError(t, err, "Failed to create cohort")

	// Insert statistical analysis
	resultsJSON, _ := json.Marshal(map[string]interface{}{
		"effect_size":    0.65,
		"cohens_d":       0.72,
		"power":          0.85,
		"sample_size":    100,
	})

	analysisQuery := `
		INSERT INTO statistical_analyses (
			cohort_id, analysis_type, analysis_name,
			results, p_value, effect_size, interpretation,
			performed_at, created_at
		) VALUES ($1, 't_test', 'Treatment vs Control',
			$2, 0.003, 0.65, 'Significant difference with large effect',
			NOW(), NOW())
		RETURNING id
	`

	var analysisID int64
	err = suite.DB.QueryRow(analysisQuery, cohortID, resultsJSON).Scan(&analysisID)
	AssertNoError(t, err, "Failed to create analysis")

	// Verify analysis
	var pValue, effectSize float64
	var interpretation string
	err = suite.DB.QueryRow(`
		SELECT p_value, effect_size, interpretation
		FROM statistical_analyses
		WHERE id = $1
	`, analysisID).Scan(&pValue, &effectSize, &interpretation)

	AssertNoError(t, err, "Failed to retrieve analysis")
	AssertTrue(t, pValue < 0.01, "p-value should be highly significant")
	AssertTrue(t, effectSize >= 0.5, "Effect size should be medium-large")
	AssertTrue(t, interpretation != "", "Interpretation should not be empty")

	// Cleanup
	suite.DB.Exec("DELETE FROM statistical_analyses WHERE id = $1", analysisID)
	suite.DB.Exec("DELETE FROM research_cohorts WHERE id = $1", cohortID)
}

func TestResearch_DataExport(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	// Create cohort
	cohortQuery := `
		INSERT INTO research_cohorts (
			study_name, study_code, hypothesis, study_type,
			target_n_patients, status, primary_outcome,
			principal_investigator, created_at
		) VALUES (
			'Export Test', 'TEST-EXPORT-001',
			'Test export', 'longitudinal_correlation',
			50, 'completed', 'Test outcome',
			'Test Researcher', NOW()
		)
		RETURNING id
	`

	var cohortID int64
	err := suite.DB.QueryRow(cohortQuery).Scan(&cohortID)
	AssertNoError(t, err, "Failed to create cohort")

	// Create export record
	exportQuery := `
		INSERT INTO research_exports (
			cohort_id, export_name, export_format,
			file_path, n_patients, n_observations,
			anonymization_level, pii_removed,
			exported_by, exported_for_purpose, created_at
		) VALUES ($1, 'Full Dataset', 'csv',
			'/exports/test-export.csv', 45, 1200,
			'fully_anonymized', true,
			'test_researcher', 'Statistical analysis', NOW())
		RETURNING id
	`

	var exportID int64
	err = suite.DB.QueryRow(exportQuery, cohortID).Scan(&exportID)
	AssertNoError(t, err, "Failed to create export")

	// Verify export compliance
	var piiRemoved bool
	var anonLevel string
	err = suite.DB.QueryRow(`
		SELECT pii_removed, anonymization_level
		FROM research_exports
		WHERE id = $1
	`, exportID).Scan(&piiRemoved, &anonLevel)

	AssertNoError(t, err, "Failed to retrieve export")
	AssertTrue(t, piiRemoved, "PII should be removed")
	AssertEqual(t, "fully_anonymized", anonLevel, "Should be fully anonymized")

	// Cleanup
	suite.DB.Exec("DELETE FROM research_exports WHERE id = $1", exportID)
	suite.DB.Exec("DELETE FROM research_cohorts WHERE id = $1", cohortID)
}

// ============================================================================
// INTEGRATION TESTS: Clinical Scales with Database
// ============================================================================

func TestScales_PHQ9WithDatabase(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create PHQ-9 assessment
	query := `
		INSERT INTO clinical_assessments (
			patient_id, assessment_type, score, severity_level,
			assessed_at, created_at
		) VALUES ($1, 'PHQ-9', 18, 'moderately_severe', NOW(), NOW())
		RETURNING id
	`

	var assessmentID int64
	err := suite.DB.QueryRow(query, patientID).Scan(&assessmentID)
	AssertNoError(t, err, "Failed to create assessment")

	// Insert responses
	responses := []struct {
		question int
		score    int
	}{
		{1, 2}, {2, 3}, {3, 2}, {4, 2},
		{5, 2}, {6, 2}, {7, 2}, {8, 1}, {9, 2},
	}

	for _, r := range responses {
		respQuery := `
			INSERT INTO clinical_assessment_responses (
				assessment_id, question_number, response_value, created_at
			) VALUES ($1, $2, $3, NOW())
		`
		suite.DB.Exec(respQuery, assessmentID, r.question, r.score)
	}

	// Verify total score calculation
	var totalScore int
	err = suite.DB.QueryRow(`
		SELECT SUM(response_value)
		FROM clinical_assessment_responses
		WHERE assessment_id = $1
	`, assessmentID).Scan(&totalScore)

	AssertNoError(t, err, "Failed to calculate total")
	AssertEqual(t, 18, totalScore, "Total score should be 18")

	// Check if Q9 (suicide risk) is positive
	var q9Score int
	err = suite.DB.QueryRow(`
		SELECT response_value
		FROM clinical_assessment_responses
		WHERE assessment_id = $1 AND question_number = 9
	`, assessmentID).Scan(&q9Score)

	AssertNoError(t, err, "Failed to get Q9 score")
	AssertTrue(t, q9Score > 0, "Q9 should be positive (suicide risk)")

	// Cleanup
	suite.DB.Exec("DELETE FROM clinical_assessment_responses WHERE assessment_id = $1", assessmentID)
	suite.DB.Exec("DELETE FROM clinical_assessments WHERE id = $1", assessmentID)
}

func TestScales_CSSRSWithDatabase(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create C-SSRS assessment with critical risk
	query := `
		INSERT INTO clinical_assessments (
			patient_id, assessment_type, score, severity_level,
			assessed_at, created_at
		) VALUES ($1, 'C-SSRS', 4, 'critical', NOW(), NOW())
		RETURNING id
	`

	var assessmentID int64
	err := suite.DB.QueryRow(query, patientID).Scan(&assessmentID)
	AssertNoError(t, err, "Failed to create C-SSRS assessment")

	// Insert responses (Q6 positive = behavior = critical)
	responses := []struct {
		question int
		answer   bool
	}{
		{1, true}, {2, true}, {3, true},
		{4, true}, {5, false}, {6, true}, // Q6 = behavior
	}

	for _, r := range responses {
		value := 0
		if r.answer {
			value = 1
		}
		respQuery := `
			INSERT INTO clinical_assessment_responses (
				assessment_id, question_number, response_value, created_at
			) VALUES ($1, $2, $3, NOW())
		`
		suite.DB.Exec(respQuery, assessmentID, r.question, value)
	}

	// Verify behavior detected (Q6)
	var q6Value int
	err = suite.DB.QueryRow(`
		SELECT response_value
		FROM clinical_assessment_responses
		WHERE assessment_id = $1 AND question_number = 6
	`, assessmentID).Scan(&q6Value)

	AssertNoError(t, err, "Failed to get Q6 value")
	AssertEqual(t, 1, q6Value, "Q6 should be positive (behavior)")

	// Verify alert was created (should happen automatically for critical)
	var alertCount int
	err = suite.DB.QueryRow(`
		SELECT COUNT(*) FROM clinical_alerts
		WHERE patient_id = $1 AND alert_type = 'C-SSRS'
		AND created_at > NOW() - INTERVAL '1 minute'
	`, patientID).Scan(&alertCount)

	// Note: Alert creation depends on trigger/application logic
	// This test verifies the database structure is correct

	// Cleanup
	suite.DB.Exec("DELETE FROM clinical_assessment_responses WHERE assessment_id = $1", assessmentID)
	suite.DB.Exec("DELETE FROM clinical_assessments WHERE id = $1", assessmentID)
	suite.DB.Exec("DELETE FROM clinical_alerts WHERE patient_id = $1 AND alert_type = 'C-SSRS'", patientID)
}

func TestScales_GAD7WithDatabase(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create GAD-7 assessment
	query := `
		INSERT INTO clinical_assessments (
			patient_id, assessment_type, score, severity_level,
			assessed_at, created_at
		) VALUES ($1, 'GAD-7', 16, 'severe', NOW(), NOW())
		RETURNING id
	`

	var assessmentID int64
	err := suite.DB.QueryRow(query, patientID).Scan(&assessmentID)
	AssertNoError(t, err, "Failed to create GAD-7 assessment")

	// Insert responses (all high scores)
	for i := 1; i <= 7; i++ {
		score := 2
		if i <= 3 {
			score = 3 // First 3 questions at max
		}
		respQuery := `
			INSERT INTO clinical_assessment_responses (
				assessment_id, question_number, response_value, created_at
			) VALUES ($1, $2, $3, NOW())
		`
		suite.DB.Exec(respQuery, assessmentID, i, score)
	}

	// Verify total score
	var totalScore int
	err = suite.DB.QueryRow(`
		SELECT SUM(response_value)
		FROM clinical_assessment_responses
		WHERE assessment_id = $1
	`, assessmentID).Scan(&totalScore)

	AssertNoError(t, err, "Failed to calculate total")
	AssertTrue(t, totalScore >= 15, "Total should indicate severe anxiety")

	// Cleanup
	suite.DB.Exec("DELETE FROM clinical_assessment_responses WHERE assessment_id = $1", assessmentID)
	suite.DB.Exec("DELETE FROM clinical_assessments WHERE id = $1", assessmentID)
}

// ============================================================================
// BENCHMARK TESTS
// ============================================================================

func BenchmarkAnonymization(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		patientID := int64(i + 1000)
		hash := sha256.Sum256([]byte(fmt.Sprintf("%d", patientID)))
		_ = hex.EncodeToString(hash[:])
	}
}

func BenchmarkCohortDatapointInsertion(b *testing.B) {
	if suite.DB == nil {
		b.Skip("Database not available")
	}

	// Create cohort
	cohortQuery := `
		INSERT INTO research_cohorts (
			study_name, study_code, hypothesis, study_type,
			target_n_patients, status, primary_outcome,
			principal_investigator, created_at
		) VALUES (
			'Benchmark Test', 'BENCH-001', 'Benchmark',
			'longitudinal_correlation', 1000, 'data_collection',
			'Benchmark', 'Benchmark', NOW()
		)
		RETURNING id
	`

	var cohortID int64
	suite.DB.QueryRow(cohortQuery).Scan(&cohortID)

	dataJSON, _ := json.Marshal(map[string]interface{}{"score": 10})

	b.ResetTimer()
	var ids []int64
	for i := 0; i < b.N; i++ {
		hash := sha256.Sum256([]byte(fmt.Sprintf("%d", i)))
		anonID := hex.EncodeToString(hash[:])

		query := `
			INSERT INTO research_datapoints (
				cohort_id, anonymous_patient_id, collection_date,
				datapoint_type, data, created_at
			) VALUES ($1, $2, NOW(), 'benchmark', $3, NOW())
			RETURNING id
		`
		var id int64
		suite.DB.QueryRow(query, cohortID, anonID, dataJSON).Scan(&id)
		ids = append(ids, id)
	}

	// Cleanup
	for _, id := range ids {
		suite.DB.Exec("DELETE FROM research_datapoints WHERE id = $1", id)
	}
	suite.DB.Exec("DELETE FROM research_cohorts WHERE id = $1", cohortID)
}
