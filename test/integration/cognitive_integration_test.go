package integration

import (
	"testing"
	"time"
)

// ============================================================================
// INTEGRATION TESTS: Cognitive Load Orchestrator
// Tests the cognitive load management system
// ============================================================================

func TestCognitiveLoad_StateTracking(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create cognitive load state
	query := `
		INSERT INTO cognitive_load_state (
			patient_id, current_load, emotional_intensity, fatigue_level,
			last_interaction, session_duration_minutes, updated_at
		) VALUES ($1, 0.5, 0.6, 0.3, NOW(), 30, NOW())
		ON CONFLICT (patient_id) DO UPDATE SET
			current_load = 0.5,
			emotional_intensity = 0.6,
			fatigue_level = 0.3,
			updated_at = NOW()
		RETURNING patient_id
	`

	var returnedID int64
	err := suite.DB.QueryRow(query, patientID).Scan(&returnedID)
	AssertNoError(t, err, "Failed to create cognitive load state")
	AssertEqual(t, patientID, returnedID, "Patient ID mismatch")

	// Verify state was saved
	var currentLoad, emotionalIntensity, fatigueLevel float64
	err = suite.DB.QueryRow(`
		SELECT current_load, emotional_intensity, fatigue_level
		FROM cognitive_load_state
		WHERE patient_id = $1
	`, patientID).Scan(&currentLoad, &emotionalIntensity, &fatigueLevel)

	AssertNoError(t, err, "Failed to retrieve cognitive load state")
	AssertTrue(t, currentLoad == 0.5, "Current load should be 0.5")
	AssertTrue(t, emotionalIntensity == 0.6, "Emotional intensity should be 0.6")
	AssertTrue(t, fatigueLevel == 0.3, "Fatigue level should be 0.3")
}

func TestCognitiveLoad_InteractionLogging(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Log multiple interactions
	interactions := []struct {
		cognitiveScore    float64
		emotionalScore    float64
		topic             string
		responseTimeMs    int
	}{
		{0.3, 0.4, "weather", 500},
		{0.5, 0.6, "health", 800},
		{0.7, 0.8, "family_death", 1200}, // High cognitive load topic
	}

	var interactionIDs []int64

	for _, interaction := range interactions {
		query := `
			INSERT INTO interaction_cognitive_load (
				patient_id, cognitive_score, emotional_score, topic,
				response_time_ms, created_at
			) VALUES ($1, $2, $3, $4, $5, NOW())
			RETURNING id
		`

		var id int64
		err := suite.DB.QueryRow(query, patientID, interaction.cognitiveScore,
			interaction.emotionalScore, interaction.topic, interaction.responseTimeMs).Scan(&id)
		AssertNoError(t, err, "Failed to log interaction")
		interactionIDs = append(interactionIDs, id)
	}

	// Verify interactions were logged
	var count int
	err := suite.DB.QueryRow(`
		SELECT COUNT(*) FROM interaction_cognitive_load
		WHERE patient_id = $1 AND created_at > NOW() - INTERVAL '1 hour'
	`, patientID).Scan(&count)

	AssertNoError(t, err, "Failed to count interactions")
	AssertTrue(t, count >= 3, "Should have at least 3 interactions")

	// Calculate average cognitive load
	var avgLoad float64
	err = suite.DB.QueryRow(`
		SELECT AVG(cognitive_score) FROM interaction_cognitive_load
		WHERE patient_id = $1 AND created_at > NOW() - INTERVAL '1 hour'
	`, patientID).Scan(&avgLoad)

	AssertNoError(t, err, "Failed to calculate average load")
	AssertTrue(t, avgLoad > 0.4, "Average load should be > 0.4")

	// Cleanup
	for _, id := range interactionIDs {
		suite.DB.Exec("DELETE FROM interaction_cognitive_load WHERE id = $1", id)
	}
}

func TestCognitiveLoad_DecisionLogging(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Log a cognitive load decision
	query := `
		INSERT INTO cognitive_load_decisions (
			patient_id, decision_type, action_taken, reason,
			cognitive_load_at_decision, emotional_state, created_at
		) VALUES ($1, 'redirect', 'topic_change', 'Rumination detected',
			0.85, 'anxious', NOW())
		RETURNING id
	`

	var decisionID int64
	err := suite.DB.QueryRow(query, patientID).Scan(&decisionID)
	AssertNoError(t, err, "Failed to log decision")

	// Verify decision was logged
	var decisionType, actionTaken, reason string
	err = suite.DB.QueryRow(`
		SELECT decision_type, action_taken, reason
		FROM cognitive_load_decisions
		WHERE id = $1
	`, decisionID).Scan(&decisionType, &actionTaken, &reason)

	AssertNoError(t, err, "Failed to retrieve decision")
	AssertEqual(t, "redirect", decisionType, "Decision type mismatch")
	AssertEqual(t, "topic_change", actionTaken, "Action taken mismatch")
	AssertTrue(t, reason != "", "Reason should not be empty")

	// Cleanup
	suite.DB.Exec("DELETE FROM cognitive_load_decisions WHERE id = $1", decisionID)
}

func TestCognitiveLoad_RuminationDetection(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create pattern of rumination (same topic repeated)
	topic := "death_of_spouse"
	for i := 0; i < 4; i++ {
		query := `
			INSERT INTO interaction_cognitive_load (
				patient_id, cognitive_score, emotional_score, topic,
				response_time_ms, created_at
			) VALUES ($1, 0.8, 0.9, $2, 1000, NOW() - ($3 * INTERVAL '10 minutes'))
		`
		suite.DB.Exec(query, patientID, topic, i)
	}

	// Detect rumination (same topic 3+ times in 2 hours)
	var ruminationCount int
	err := suite.DB.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT topic, COUNT(*) as cnt
			FROM interaction_cognitive_load
			WHERE patient_id = $1
			  AND created_at > NOW() - INTERVAL '2 hours'
			GROUP BY topic
			HAVING COUNT(*) >= 3
		) rumination
	`, patientID).Scan(&ruminationCount)

	AssertNoError(t, err, "Failed to detect rumination")
	AssertTrue(t, ruminationCount >= 1, "Should detect rumination pattern")

	// Cleanup
	suite.DB.Exec(`
		DELETE FROM interaction_cognitive_load
		WHERE patient_id = $1 AND topic = $2
	`, patientID, topic)
}

func TestCognitiveLoad_FatigueProgression(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Simulate session with increasing fatigue
	sessionStart := time.Now().Add(-60 * time.Minute)

	fatigueProgression := []struct {
		minutesIn   int
		fatigueLevel float64
	}{
		{10, 0.1},
		{20, 0.2},
		{30, 0.35},
		{40, 0.5},
		{50, 0.7},
		{60, 0.85},
	}

	for _, fp := range fatigueProgression {
		query := `
			INSERT INTO interaction_cognitive_load (
				patient_id, cognitive_score, emotional_score, topic,
				fatigue_indicator, response_time_ms, created_at
			) VALUES ($1, 0.5, 0.5, 'general', $2, $3, $4)
		`
		createdAt := sessionStart.Add(time.Duration(fp.minutesIn) * time.Minute)
		// Response time increases with fatigue
		responseTime := 500 + int(fp.fatigueLevel*1000)
		suite.DB.Exec(query, patientID, fp.fatigueLevel, responseTime, createdAt)
	}

	// Verify fatigue progression
	var latestFatigue float64
	err := suite.DB.QueryRow(`
		SELECT fatigue_indicator
		FROM interaction_cognitive_load
		WHERE patient_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, patientID).Scan(&latestFatigue)

	AssertNoError(t, err, "Failed to get latest fatigue")
	AssertTrue(t, latestFatigue >= 0.8, "Final fatigue should be high")

	// Cleanup
	suite.DB.Exec(`
		DELETE FROM interaction_cognitive_load
		WHERE patient_id = $1 AND created_at > NOW() - INTERVAL '2 hours'
	`, patientID)
}

// ============================================================================
// INTEGRATION TESTS: Ethical Boundary Engine
// ============================================================================

func TestEthicalBoundary_AttachmentDetection(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Simulate pathological attachment phrase
	attachmentPhrases := []string{
		"Você é minha única amiga",
		"Só confio em você",
		"Não preciso de mais ninguém",
		"Você me entende melhor que minha família",
	}

	var eventIDs []int64

	for _, phrase := range attachmentPhrases {
		query := `
			INSERT INTO ethical_boundary_events (
				patient_id, event_type, severity, description,
				phrase_detected, intervention_level, created_at
			) VALUES ($1, 'pathological_attachment', 'warning', $2, $2, 1, NOW())
			RETURNING id
		`

		var eventID int64
		err := suite.DB.QueryRow(query, patientID, phrase).Scan(&eventID)
		AssertNoError(t, err, "Failed to log attachment event")
		eventIDs = append(eventIDs, eventID)
	}

	// Verify events were logged
	var attachmentCount int
	err := suite.DB.QueryRow(`
		SELECT COUNT(*) FROM ethical_boundary_events
		WHERE patient_id = $1
		  AND event_type = 'pathological_attachment'
		  AND created_at > NOW() - INTERVAL '1 hour'
	`, patientID).Scan(&attachmentCount)

	AssertNoError(t, err, "Failed to count attachment events")
	AssertTrue(t, attachmentCount >= 4, "Should have logged all attachment phrases")

	// Cleanup
	for _, id := range eventIDs {
		suite.DB.Exec("DELETE FROM ethical_boundary_events WHERE id = $1", id)
	}
}

func TestEthicalBoundary_InterventionLevels(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Test intervention level progression
	interventionLevels := []struct {
		level       int
		description string
		expectedAction string
	}{
		{1, "Soft redirect", "gentle_redirect"},
		{2, "Explicit human recommendation", "recommend_human"},
		{3, "Block + family notification", "block_and_notify"},
	}

	var eventIDs []int64

	for _, intervention := range interventionLevels {
		query := `
			INSERT INTO ethical_boundary_events (
				patient_id, event_type, severity, description,
				intervention_level, action_taken, created_at
			) VALUES ($1, 'dependency_pattern', 'warning', $2, $3, $4, NOW())
			RETURNING id
		`

		var eventID int64
		err := suite.DB.QueryRow(query, patientID, intervention.description,
			intervention.level, intervention.expectedAction).Scan(&eventID)
		AssertNoError(t, err, "Failed to log intervention")
		eventIDs = append(eventIDs, eventID)
	}

	// Verify max intervention level
	var maxLevel int
	err := suite.DB.QueryRow(`
		SELECT MAX(intervention_level) FROM ethical_boundary_events
		WHERE patient_id = $1 AND created_at > NOW() - INTERVAL '1 hour'
	`, patientID).Scan(&maxLevel)

	AssertNoError(t, err, "Failed to get max intervention level")
	AssertEqual(t, 3, maxLevel, "Max intervention level should be 3")

	// Cleanup
	for _, id := range eventIDs {
		suite.DB.Exec("DELETE FROM ethical_boundary_events WHERE id = $1", id)
	}
}

func TestEthicalBoundary_HumanInteractionRatio(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Update or create ethical boundary state
	query := `
		INSERT INTO ethical_boundary_state (
			patient_id, eva_interaction_count, human_interaction_count,
			interaction_ratio, last_human_contact, updated_at
		) VALUES ($1, 50, 5, 10.0, NOW() - INTERVAL '3 days', NOW())
		ON CONFLICT (patient_id) DO UPDATE SET
			eva_interaction_count = 50,
			human_interaction_count = 5,
			interaction_ratio = 10.0,
			last_human_contact = NOW() - INTERVAL '3 days',
			updated_at = NOW()
	`

	_, err := suite.DB.Exec(query, patientID)
	AssertNoError(t, err, "Failed to create ethical boundary state")

	// Check for high ratio (>10:1 is concerning)
	var ratio float64
	err = suite.DB.QueryRow(`
		SELECT interaction_ratio FROM ethical_boundary_state
		WHERE patient_id = $1
	`, patientID).Scan(&ratio)

	AssertNoError(t, err, "Failed to get interaction ratio")
	AssertTrue(t, ratio >= 10.0, "Ratio should indicate concerning dependency")

	// Check days since human contact
	var daysSinceHuman float64
	err = suite.DB.QueryRow(`
		SELECT EXTRACT(EPOCH FROM (NOW() - last_human_contact))/86400
		FROM ethical_boundary_state
		WHERE patient_id = $1
	`, patientID).Scan(&daysSinceHuman)

	AssertNoError(t, err, "Failed to get days since human contact")
	AssertTrue(t, daysSinceHuman >= 3, "Should be at least 3 days since human contact")
}

func TestEthicalBoundary_FamilyNotification(t *testing.T) {
	if suite.DB == nil {
		t.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID

	// Create family notification record
	query := `
		INSERT INTO family_notifications (
			patient_id, notification_type, reason, sent_at,
			recipient_email, recipient_phone, created_at
		) VALUES ($1, 'dependency_warning', 'High EVA:Human ratio detected',
			NOW(), 'family@test.com', '11999999999', NOW())
		RETURNING id
	`

	var notificationID int64
	err := suite.DB.QueryRow(query, patientID).Scan(&notificationID)
	AssertNoError(t, err, "Failed to create notification")

	// Verify notification was created
	var notificationType, reason string
	err = suite.DB.QueryRow(`
		SELECT notification_type, reason
		FROM family_notifications
		WHERE id = $1
	`, notificationID).Scan(&notificationType, &reason)

	AssertNoError(t, err, "Failed to retrieve notification")
	AssertEqual(t, "dependency_warning", notificationType, "Notification type mismatch")
	AssertTrue(t, reason != "", "Reason should not be empty")

	// Cleanup
	suite.DB.Exec("DELETE FROM family_notifications WHERE id = $1", notificationID)
}

// ============================================================================
// BENCHMARK TESTS
// ============================================================================

func BenchmarkCognitiveLoadInsertion(b *testing.B) {
	if suite.DB == nil {
		b.Skip("Database not available")
	}

	patientID := suite.TestIDs.PatientID
	query := `
		INSERT INTO interaction_cognitive_load (
			patient_id, cognitive_score, emotional_score, topic,
			response_time_ms, created_at
		) VALUES ($1, 0.5, 0.5, 'benchmark', 500, NOW())
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
		suite.DB.Exec("DELETE FROM interaction_cognitive_load WHERE id = $1", id)
	}
}
