package database

import (
	"database/sql"
	"fmt"
	"time"
)

type VideoSession struct {
	ID        string
	SessionID string
	IdosoID   int64
	Status    string
	SdpOffer  string
	SdpAnswer sql.NullString
	CreatedAt time.Time
}

type SignalingMessage struct {
	ID        int64
	SessionID string
	Sender    string
	Type      string
	Payload   string // JSON
	CreatedAt time.Time
}

func (db *DB) CreateVideoSession(sessionID string, idosoID int64, sdpOffer string) error {
	query := `
		INSERT INTO video_sessions (session_id, idoso_id, status, sdp_offer, created_em)
		VALUES ($1, $2, 'waiting_operator', $3, CURRENT_TIMESTAMP)
	`
	// Usamos ExecContext para boas práticas, mas aqui com context.Background() se não vier de cima
	_, err := db.conn.Exec(query, sessionID, idosoID, sdpOffer)
	if err != nil {
		return fmt.Errorf("failed to create video session: %w", err)
	}
	return nil
}

func (db *DB) CreateSignalingMessage(sessionID string, sender string, msgType string, payload string) error {
	query := `
		INSERT INTO signaling_messages (session_id, sender, type, payload)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.conn.Exec(query, sessionID, sender, msgType, payload)
	if err != nil {
		return fmt.Errorf("failed to insert signaling message: %w", err)
	}
	return nil
}

func (db *DB) GetVideoSessionAnswer(sessionID string) (string, error) {
	query := `SELECT sdp_answer FROM video_sessions WHERE session_id = $1`

	var sdpAnswer sql.NullString
	err := db.conn.QueryRow(query, sessionID).Scan(&sdpAnswer)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Não encontrou a sessão ou não tem answer ainda
		}
		return "", fmt.Errorf("failed to get session answer: %w", err)
	}

	if sdpAnswer.Valid {
		return sdpAnswer.String, nil
	}
	return "", nil
}

// Opcional: Pegar candidatos do Operador para o Mobile
func (db *DB) GetOperatorCandidates(sessionID string, sinceID int64) ([]SignalingMessage, error) {
	query := `
		SELECT id, session_id, sender, type, payload 
		FROM signaling_messages 
		WHERE session_id = $1 AND sender = 'operator' AND id > $2
		ORDER BY id ASC
	`

	rows, err := db.conn.Query(query, sessionID, sinceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []SignalingMessage
	for rows.Next() {
		var m SignalingMessage
		if err := rows.Scan(&m.ID, &m.SessionID, &m.Sender, &m.Type, &m.Payload); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}
