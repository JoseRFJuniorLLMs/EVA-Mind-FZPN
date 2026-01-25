package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Database connection
var db *sql.DB

// DTOs (simplified - em produção usar internal/integration)
type PatientDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	DateOfBirth string    `json:"date_of_birth"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Email       *string   `json:"email,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AssessmentDTO struct {
	ID             string    `json:"id"`
	PatientID      int64     `json:"patient_id"`
	AssessmentType string    `json:"assessment_type"`
	TotalScore     *int      `json:"total_score,omitempty"`
	Severity       *string   `json:"severity,omitempty"`
	CompletedAt    time.Time `json:"completed_at"`
}

type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
	Signature string                 `json:"signature,omitempty"`
}

func main() {
	// Conectar ao database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:Debian23@@104.248.219.200:5432/eva-db?sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Erro ao conectar database:", err)
	}
	defer db.Close()

	// Testar conexão
	err = db.Ping()
	if err != nil {
		log.Fatal("Database não acessível:", err)
	}

	log.Println("✓ Conectado ao PostgreSQL remoto")

	// Setup router
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// Serialize endpoints
	r.HandleFunc("/serialize/patient/{id}", serializePatientHandler).Methods("GET")
	r.HandleFunc("/serialize/assessment/{id}", serializeAssessmentHandler).Methods("GET")

	// FHIR endpoints
	r.HandleFunc("/fhir/patient/{id}", convertPatientToFHIRHandler).Methods("GET")
	r.HandleFunc("/fhir/bundle/{patient_id}", createFHIRBundleHandler).Methods("GET")

	// Webhook endpoints
	r.HandleFunc("/webhook/build", buildWebhookHandler).Methods("POST")

	// Export endpoints
	r.HandleFunc("/export/lgpd/{patient_id}", exportLGPDHandler).Methods("GET")

	// CORS middleware
	r.Use(corsMiddleware)

	// Start server
	port := ":8081"
	log.Printf("🚀 Eva Integration Service rodando em http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"service": "eva-integration-service",
		"version": "1.0.0",
		"time":    time.Now(),
	})
}

func serializePatientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	// Buscar patient do database
	var patient PatientDTO
	var email, phone sql.NullString

	err = db.QueryRow(`
		SELECT id, name, date_of_birth,
		       EXTRACT(YEAR FROM AGE(date_of_birth::date))::int as age,
		       gender, email, phone, created_at, updated_at
		FROM patients
		WHERE id = $1
	`, patientID).Scan(
		&patient.ID,
		&patient.Name,
		&patient.DateOfBirth,
		&patient.Age,
		&patient.Gender,
		&email,
		&phone,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if email.Valid {
		patient.Email = &email.String
	}
	if phone.Valid {
		patient.Phone = &phone.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

func serializeAssessmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assessmentID := vars["id"]

	var assessment AssessmentDTO
	var totalScore sql.NullInt64
	var severity sql.NullString

	err := db.QueryRow(`
		SELECT id, patient_id, assessment_type, total_score, severity, completed_at
		FROM mental_health_assessments
		WHERE id = $1
	`, assessmentID).Scan(
		&assessment.ID,
		&assessment.PatientID,
		&assessment.AssessmentType,
		&totalScore,
		&severity,
		&assessment.CompletedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Assessment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if totalScore.Valid {
		score := int(totalScore.Int64)
		assessment.TotalScore = &score
	}
	if severity.Valid {
		assessment.Severity = &severity.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assessment)
}

func convertPatientToFHIRHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID := vars["id"]

	// Simplified FHIR Patient
	fhirPatient := map[string]interface{}{
		"resourceType": "Patient",
		"id":           patientID,
		"meta": map[string]interface{}{
			"lastUpdated": time.Now(),
			"source":      "EVA-Mind",
		},
		"identifier": []map[string]interface{}{
			{
				"use":    "official",
				"system": "https://eva-mind.com/patient-id",
				"value":  patientID,
			},
		},
		"active": true,
	}

	w.Header().Set("Content-Type", "application/fhir+json")
	json.NewEncoder(w).Encode(fhirPatient)
}

func createFHIRBundleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID := vars["patient_id"]

	bundle := map[string]interface{}{
		"resourceType": "Bundle",
		"type":         "collection",
		"entry": []map[string]interface{}{
			{
				"fullUrl":  fmt.Sprintf("Patient/%s", patientID),
				"resource": map[string]string{"resourceType": "Patient", "id": patientID},
			},
		},
	}

	w.Header().Set("Content-Type", "application/fhir+json")
	json.NewEncoder(w).Encode(bundle)
}

func buildWebhookHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EventType string                 `json:"event_type"`
		Data      map[string]interface{} `json:"data"`
		Secret    string                 `json:"secret"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event := WebhookEvent{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Type:      req.EventType,
		Timestamp: time.Now(),
		Source:    "EVA-Mind",
		Data:      req.Data,
	}

	// TODO: Add HMAC-SHA256 signature
	event.Signature = "signature_placeholder"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func exportLGPDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID := vars["patient_id"]

	export := map[string]interface{}{
		"export_date": time.Now(),
		"patient_id":  patientID,
		"patient":     map[string]string{},
		"assessments": []map[string]interface{}{},
		"message":     "LGPD export - dados seriam buscados do database",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=lgpd_export_%s.json", patientID))
	json.NewEncoder(w).Encode(export)
}
