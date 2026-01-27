// Package integration provides integration tests for EVA-Mind-FZPN
package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// TestSuite holds shared resources for integration tests
type TestSuite struct {
	DB      *sql.DB
	Ctx     context.Context
	Cancel  context.CancelFunc
	TestIDs struct {
		PatientID int64
	}
}

// Global test suite
var suite *TestSuite

// TestMain sets up and tears down the test suite
func TestMain(m *testing.M) {
	suite = &TestSuite{}

	// Setup
	if err := suite.Setup(); err != nil {
		panic("Failed to setup test suite: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Teardown
	suite.Teardown()

	os.Exit(code)
}

// Setup initializes the test suite
func (s *TestSuite) Setup() error {
	// Get database URL
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/eva_test?sslmode=disable"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	s.DB = db
	s.Ctx, s.Cancel = context.WithTimeout(context.Background(), 5*time.Minute)

	// Create test patient
	s.createTestPatient()

	return nil
}

// Teardown cleans up the test suite
func (s *TestSuite) Teardown() {
	if s.Cancel != nil {
		s.Cancel()
	}

	// Clean up test data
	if s.DB != nil {
		s.cleanupTestData()
		s.DB.Close()
	}
}

// createTestPatient creates a test patient for integration tests
func (s *TestSuite) createTestPatient() {
	query := `
		INSERT INTO idosos (nome, cpf, email, telefone, data_nascimento)
		VALUES ('Teste Integration', '99999999999', 'test@integration.com', '11999999999', '1950-01-01')
		ON CONFLICT (cpf) DO UPDATE SET nome = 'Teste Integration'
		RETURNING id
	`
	s.DB.QueryRow(query).Scan(&s.TestIDs.PatientID)
}

// cleanupTestData removes test data
func (s *TestSuite) cleanupTestData() {
	if s.TestIDs.PatientID > 0 {
		// Clean up in order (respecting foreign keys)
		tables := []string{
			"assessment_responses",
			"clinical_assessments",
			"clinical_alerts",
			"interaction_cognitive_load",
			"cognitive_load_decisions",
			"ethical_boundary_events",
			"persona_sessions",
			"trajectory_simulations",
		}

		for _, table := range tables {
			s.DB.Exec("DELETE FROM "+table+" WHERE patient_id = $1", s.TestIDs.PatientID)
		}

		// Don't delete the patient - it might be needed for other tests
	}
}

// Helper functions for tests

// AssertNoError fails if there's an error
func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

// AssertError fails if there's NO error
func AssertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got nil", msg)
	}
}

// AssertEqual compares two values
func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if expected != actual {
		t.Fatalf("%s: expected %v, got %v", msg, expected, actual)
	}
}

// AssertTrue fails if condition is false
func AssertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Fatalf("%s: expected true", msg)
	}
}

// AssertNotNil fails if value is nil
func AssertNotNil(t *testing.T, value interface{}, msg string) {
	t.Helper()
	if value == nil {
		t.Fatalf("%s: expected non-nil value", msg)
	}
}
