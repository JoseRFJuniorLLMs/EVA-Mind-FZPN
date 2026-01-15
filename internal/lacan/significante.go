package lacan

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

// SignifierService rastreia significantes recorrentes (palavras-chave que se repetem)
type SignifierService struct {
	db *sql.DB
}

// NewSignifierService cria novo serviço
func NewSignifierService(db *sql.DB) *SignifierService {
	return &SignifierService{db: db}
}

// Signifier representa um significante rastreado
type Signifier struct {
	Word            string
	Frequency       int
	Contexts        []string
	FirstOccurrence time.Time
	LastOccurrence  time.Time
	EmotionalCharge float64 // 0.0 (neutro) a 1.0 (altamente carregado)
}

// TrackSignifiers extrai e rastreia significantes emocionalmente carregados
func (s *SignifierService) TrackSignifiers(ctx context.Context, idosoID int64, text string) error {
	keywords := extractEmotionalKeywords(text)

	for _, word := range keywords {
		err := s.incrementSignifier(ctx, idosoID, word, text)
		if err != nil {
			return err
		}
	}

	return nil
}

// incrementSignifier incrementa frequência de um significante
func (s *SignifierService) incrementSignifier(ctx context.Context, idosoID int64, word, context string) error {
	// Verificar se já existe
	var exists bool
	err := s.db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM significantes_recorrentes WHERE idoso_id = $1 AND palavra = $2)
	`, idosoID, word).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		// Atualizar
		query := `
			UPDATE significantes_recorrentes
			SET frequencia = frequencia + 1,
			    contextos = array_append(contextos, $1),
			    ultima_ocorrencia = NOW()
			WHERE idoso_id = $2 AND palavra = $3
		`
		_, err = s.db.ExecContext(ctx, query, context, idosoID, word)
	} else {
		// Inserir
		query := `
			INSERT INTO significantes_recorrentes 
			(idoso_id, palavra, frequencia, contextos, primeira_ocorrencia, ultima_ocorrencia)
			VALUES ($1, $2, 1, ARRAY[$3], NOW(), NOW())
		`
		_, err = s.db.ExecContext(ctx, query, idosoID, word, context)
	}

	return err
}

// GetKeySignifiers retorna os N significantes mais frequentes
func (s *SignifierService) GetKeySignifiers(ctx context.Context, idosoID int64, topN int) ([]Signifier, error) {
	query := `
		SELECT palavra, frequencia, contextos, primeira_ocorrencia, ultima_ocorrencia
		FROM significantes_recorrentes
		WHERE idoso_id = $1
		  AND frequencia >= 3
		ORDER BY frequencia DESC
		LIMIT $2
	`

	rows, err := s.db.QueryContext(ctx, query, idosoID, topN)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signifiers []Signifier
	for rows.Next() {
		var s Signifier
		var contextsStr string

		err := rows.Scan(
			&s.Word,
			&s.Frequency,
			&contextsStr,
			&s.FirstOccurrence,
			&s.LastOccurrence,
		)
		if err != nil {
			continue
		}

		// Parse contexts (PostgreSQL array)
		s.Contexts = parsePostgresTextArray(contextsStr)
		s.EmotionalCharge = calculateEmotionalCharge(s.Word)

		signifiers = append(signifiers, s)
	}

	return signifiers, rows.Err()
}

// ShouldInterpelSignifier decide se é momento de interpelar o significante
func (s *SignifierService) ShouldInterpelSignifier(ctx context.Context, idosoID int64, word string) (bool, error) {
	var frequency int
	var lastInterpellation sql.NullTime

	query := `
		SELECT frequencia, ultima_interpelacao
		FROM significantes_recorrentes
		WHERE idoso_id = $1 AND palavra = $2
	`

	err := s.db.QueryRowContext(ctx, query, idosoID, word).Scan(&frequency, &lastInterpellation)
	if err != nil {
		return false, err
	}

	// Interpelar se:
	// 1. Frequência >= 5
	// 2. Não foi interpelado recentemente (último > 7 dias)
	if frequency >= 5 {
		if !lastInterpellation.Valid {
			return true, nil
		}

		daysSince := int(time.Since(lastInterpellation.Time).Hours() / 24)
		if daysSince > 7 {
			return true, nil
		}
	}

	return false, nil
}

// MarkAsInterpelled marca que o significante foi interpelado
func (s *SignifierService) MarkAsInterpelled(ctx context.Context, idosoID int64, word string) error {
	query := `
		UPDATE significantes_recorrentes
		SET ultima_interpelacao = NOW()
		WHERE idoso_id = $1 AND palavra = $2
	`
	_, err := s.db.ExecContext(ctx, query, idosoID, word)
	return err
}

// GenerateInterpellation gera frase para interpelar o significante
func GenerateInterpellation(word string, frequency int) string {
	return "Percebi que você frequentemente menciona a palavra '" + word + "'. " +
		"Ela apareceu " + string(rune(frequency)) + " vezes em nossas conversas. " +
		"O que essa palavra representa para você?"
}

// Helper functions

func extractEmotionalKeywords(text string) []string {
	// Palavras com carga emocional (extração simples)
	emotionalWords := map[string]bool{
		"solidão":    true,
		"tristeza":   true,
		"medo":       true,
		"saudade":    true,
		"abandono":   true,
		"dor":        true,
		"sofrimento": true,
		"angústia":   true,
		"ansiedade":  true,
		"depressão":  true,
		"alegria":    true,
		"felicidade": true,
		"amor":       true,
		"morte":      true,
		"vida":       true,
		"família":    true,
		"filho":      true,
		"filha":      true,
		"esposa":     true,
		"marido":     true,
		"vazio":      true,
		"falta":      true,
		"perda":      true,
		"culpa":      true,
		"raiva":      true,
		"ódio":       true,
		"perdão":     true,
		"esperança":  true,
		"desespero":  true,
	}

	words := strings.Fields(strings.ToLower(text))
	var keywords []string

	for _, word := range words {
		// Remove pontuação
		cleaned := strings.Trim(word, ".,!?;:")
		if emotionalWords[cleaned] {
			keywords = append(keywords, cleaned)
		}
	}

	return keywords
}

func calculateEmotionalCharge(word string) float64 {
	// Palavras de alta carga emocional
	highCharge := map[string]bool{
		"morte": true, "abandono": true, "solidão": true,
		"desespero": true, "ódio": true, "culpa": true,
		"vazio": true, "perda": true,
	}

	if highCharge[word] {
		return 1.0
	}

	return 0.5 // Carga média por padrão
}

func parsePostgresTextArray(s string) []string {
	if s == "{}" || s == "" {
		return []string{}
	}

	// Remove {}
	s = strings.Trim(s, "{}")
	if s == "" {
		return []string{}
	}

	// Split respeitando aspas
	var result []string
	var current strings.Builder
	inQuotes := false

	for _, c := range s {
		switch c {
		case '"':
			inQuotes = !inQuotes
		case ',':
			if !inQuotes && current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			} else if inQuotes {
				current.WriteRune(c)
			}
		default:
			current.WriteRune(c)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
