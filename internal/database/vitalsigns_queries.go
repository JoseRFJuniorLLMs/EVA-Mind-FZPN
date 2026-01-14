package database

import (
	"fmt"
	"time"
)

// SaveVitalSign saves a vital sign measurement to the database
func (db *DB) SaveVitalSign(idosoID int64, tipo, valor, unidade, metodo, observacao string) error {
	query := `
		INSERT INTO sinais_vitais (idoso_id, tipo, valor, unidade, metodo, data_medicao, observacao)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.conn.Exec(query, idosoID, tipo, valor, unidade, metodo, time.Now(), observacao)
	if err != nil {
		return fmt.Errorf("failed to save vital sign: %w", err)
	}
	return nil
}

// GetRecentVitalSigns gets recent vital signs for an idoso
func (db *DB) GetRecentVitalSigns(idosoID int64, tipo string, limit int) ([]VitalSign, error) {
	query := `
		SELECT id, tipo, valor, unidade, metodo, data_medicao, observacao
		FROM sinais_vitais
		WHERE idoso_id = $1 AND tipo = $2
		ORDER BY data_medicao DESC
		LIMIT $3
	`
	rows, err := db.conn.Query(query, idosoID, tipo, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get vital signs: %w", err)
	}
	defer rows.Close()

	var signs []VitalSign
	for rows.Next() {
		var sign VitalSign
		err := rows.Scan(&sign.ID, &sign.Tipo, &sign.Valor, &sign.Unidade, &sign.Metodo, &sign.DataMedicao, &sign.Observacao)
		if err != nil {
			continue
		}
		signs = append(signs, sign)
	}
	return signs, nil
}

type VitalSign struct {
	ID          int64
	Tipo        string
	Valor       string
	Unidade     string
	Metodo      string
	DataMedicao time.Time
	Observacao  string
}
