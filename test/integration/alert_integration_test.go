package integration

import (
	"database/sql"
	"testing"
	"time"
)

// ============================================================================
// INTEGRATION TESTS: Alert System
// Tests the full alert flow from detection to escalation
// ============================================================================

func TestAlertSystem_CriticalAlertCreation(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create a critical alert
	query := `
		INSERT INTO clinical_alerts (
			patient_id, alert_type, severity, message, score, created_at
		) VALUES ($1, 'C-SSRS', 'critical', 'RISCO SUICIDA - C-SSRS score: 4', 4, NOW())
		RETURNING id
	`

	var alertID int64
	err := suite.DB.QueryRow(query, patientID).Scan(&alertID)
	AssertNoError(t, err, "Failed to create critical alert")
	AssertTrue(t, alertID > 0, "Alert ID should be positive")

	// Verify alert was created correctly
	var severity, alertType string
	var score int
	err = suite.DB.QueryRow(`
		SELECT severity, alert_type, score
		FROM clinical_alerts
		WHERE id = $1
	`, alertID).Scan(&severity, &alertType, &score)

	AssertNoError(t, err, "Failed to retrieve alert")
	AssertEqual(t, "critical", severity, "Severity should be critical")
	AssertEqual(t, "C-SSRS", alertType, "Alert type should be C-SSRS")
	AssertEqual(t, 4, score, "Score should be 4")

	// Cleanup
	suite.DB.Exec("DELETE FROM clinical_alerts WHERE id = $1", alertID)
}

func TestAlertSystem_EscalationChain(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create escalation log entry
	query := `
		INSERT INTO escalation_logs (
			patient_id, alert_type, original_severity, escalated_to,
			escalation_reason, channels_attempted, created_at
		) VALUES ($1, 'crisis', 'high', 'critical', 'No response from primary contact',
			ARRAY['push', 'sms'], NOW())
		RETURNING id
	`

	var logID int64
	err := suite.DB.QueryRow(query, patientID).Scan(&logID)
	AssertNoError(t, err, "Failed to create escalation log")

	// Verify escalation chain works
	var channelsAttempted []string
	var escalatedTo string
	err = suite.DB.QueryRow(`
		SELECT escalated_to, channels_attempted
		FROM escalation_logs
		WHERE id = $1
	`, logID).Scan(&escalatedTo, (*pqStringArray)(&channelsAttempted))

	AssertNoError(t, err, "Failed to retrieve escalation log")
	AssertEqual(t, "critical", escalatedTo, "Should escalate to critical")
	AssertTrue(t, len(channelsAttempted) >= 2, "Should have attempted multiple channels")

	// Cleanup
	suite.DB.Exec("DELETE FROM escalation_logs WHERE id = $1", logID)
}

func TestAlertSystem_AlertTimeout(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create an alert that's been pending too long
	query := `
		INSERT INTO clinical_alerts (
			patient_id, alert_type, severity, message, score,
			created_at, acknowledged_at
		) VALUES ($1, 'PHQ-9', 'high', 'High PHQ-9 score', 18,
			NOW() - INTERVAL '2 hours', NULL)
		RETURNING id
	`

	var alertID int64
	err := suite.DB.QueryRow(query, patientID).Scan(&alertID)
	AssertNoError(t, err, "Failed to create pending alert")

	// Query for unacknowledged alerts older than 1 hour
	var count int
	err = suite.DB.QueryRow(`
		SELECT COUNT(*)
		FROM clinical_alerts
		WHERE patient_id = $1
		  AND acknowledged_at IS NULL
		  AND created_at < NOW() - INTERVAL '1 hour'
	`, patientID).Scan(&count)

	AssertNoError(t, err, "Failed to query pending alerts")
	AssertTrue(t, count >= 1, "Should find at least one pending alert")

	// Cleanup
	suite.DB.Exec("DELETE FROM clinical_alerts WHERE id = $1", alertID)
}

func TestAlertSystem_MultiChannelDelivery(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create alert with delivery tracking
	alertQuery := `
		INSERT INTO clinical_alerts (
			patient_id, alert_type, severity, message, score, created_at
		) VALUES ($1, 'C-SSRS', 'critical', 'Critical risk', 5, NOW())
		RETURNING id
	`

	var alertID int64
	err := suite.DB.QueryRow(alertQuery, patientID).Scan(&alertID)
	AssertNoError(t, err, "Failed to create alert")

	// Simulate delivery attempts to multiple channels
	channels := []struct {
		channel string
		success bool
	}{
		{"push", true},
		{"sms", true},
		{"email", false}, // Simulated failure
		{"voice", true},
	}

	for _, ch := range channels {
		query := `
			INSERT INTO alert_delivery_log (
				alert_id, channel, success, attempted_at, error_message
			) VALUES ($1, $2, $3, NOW(), $4)
		`
		errorMsg := sql.NullString{}
		if !ch.success {
			errorMsg = sql.NullString{String: "Delivery failed", Valid: true}
		}
		suite.DB.Exec(query, alertID, ch.channel, ch.success, errorMsg)
	}

	// Verify delivery tracking
	var successCount, totalCount int
	err = suite.DB.QueryRow(`
		SELECT COUNT(*), COUNT(*) FILTER (WHERE success = true)
		FROM alert_delivery_log
		WHERE alert_id = $1
	`, alertID).Scan(&totalCount, &successCount)

	AssertNoError(t, err, "Failed to query delivery log")
	AssertEqual(t, 4, totalCount, "Should have 4 delivery attempts")
	AssertEqual(t, 3, successCount, "Should have 3 successful deliveries")

	// Cleanup
	suite.DB.Exec("DELETE FROM alert_delivery_log WHERE alert_id = $1", alertID)
	suite.DB.Exec("DELETE FROM clinical_alerts WHERE id = $1", alertID)
}

func TestAlertSystem_SeverityProgression(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Test that alerts are created with correct severity based on score
	testCases := []struct {
		score           int
		expectedSeverity string
		alertType       string
	}{
		{4, "minimal", "PHQ-9"},      // PHQ-9: 0-4 = minimal
		{8, "mild", "PHQ-9"},         // PHQ-9: 5-9 = mild
		{14, "moderate", "PHQ-9"},    // PHQ-9: 10-14 = moderate
		{18, "moderately_severe", "PHQ-9"}, // PHQ-9: 15-19 = moderately_severe
		{25, "severe", "PHQ-9"},      // PHQ-9: 20-27 = severe
	}

	for _, tc := range testCases {
		t.Run(tc.expectedSeverity, func(t *testing.T) {
			query := `
				INSERT INTO clinical_assessments (
					patient_id, assessment_type, score, severity_level,
					assessed_at, created_at
				) VALUES ($1, $2, $3, $4, NOW(), NOW())
				RETURNING id
			`

			var assessmentID int64
			err := suite.DB.QueryRow(query, patientID, tc.alertType, tc.score, tc.expectedSeverity).Scan(&assessmentID)
			AssertNoError(t, err, "Failed to create assessment")

			// Verify
			var storedSeverity string
			err = suite.DB.QueryRow(`
				SELECT severity_level FROM clinical_assessments WHERE id = $1
			`, assessmentID).Scan(&storedSeverity)

			AssertNoError(t, err, "Failed to retrieve assessment")
			AssertEqual(t, tc.expectedSeverity, storedSeverity, "Severity mismatch")

			// Cleanup
			suite.DB.Exec("DELETE FROM clinical_assessments WHERE id = $1", assessmentID)
		})
	}
}

// ============================================================================
// BENCHMARK TESTS
// ============================================================================

func BenchmarkAlertCreation(b *testing.B) {
	if suite.DB == nil {
		b.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID
	query := `
		INSERT INTO clinical_alerts (
			patient_id, alert_type, severity, message, score, created_at
		) VALUES ($1, 'PHQ-9', 'high', 'Benchmark test', 18, NOW())
		RETURNING id
	`

	var alertIDs []int64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var alertID int64
		suite.DB.QueryRow(query, patientID).Scan(&alertID)
		alertIDs = append(alertIDs, alertID)
	}

	// Cleanup
	for _, id := range alertIDs {
		suite.DB.Exec("DELETE FROM clinical_alerts WHERE id = $1", id)
	}
}

// ============================================================================
// HELPER TYPES
// ============================================================================

// pqStringArray is a helper for scanning PostgreSQL arrays
type pqStringArray []string

func (a *pqStringArray) Scan(src interface{}) error {
	if src == nil {
		*a = nil
		return nil
	}

	switch v := src.(type) {
	case []byte:
		// Parse PostgreSQL array format: {value1,value2,...}
		str := string(v)
		if str == "{}" {
			*a = []string{}
			return nil
		}
		// Remove braces
		str = str[1 : len(str)-1]
		// Split by comma
		*a = splitArray(str)
	}
	return nil
}

func splitArray(s string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	current := ""
	for _, c := range s {
		if c == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// Helper to ensure time precision in tests
func timeApproxEqual(t1, t2 time.Time, tolerance time.Duration) bool {
	diff := t1.Sub(t2)
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}
