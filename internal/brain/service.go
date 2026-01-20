package brain

import (
	"database/sql"
	"eva-mind/internal/infrastructure/vector"
	"eva-mind/internal/memory"
	ps "eva-mind/internal/personality" // Explicit alias
	"eva-mind/internal/push"
)

// Service encapsulates the cognitive functions of EVA
type Service struct {
	db                 *sql.DB
	qdrantClient       *vector.QdrantClient
	fdpnEngine         *memory.FDPNEngine
	personalityService *ps.PersonalityService
	zetaRouter         *ps.ZetaRouter
	pushService        *push.FirebaseService
	embeddingService   *memory.EmbeddingService
}

// NewService creates a new Brain service
func NewService(
	db *sql.DB,
	qdrant *vector.QdrantClient,
	fdpn *memory.FDPNEngine,
	personalitySvc *ps.PersonalityService,
	zeta *ps.ZetaRouter,
	push *push.FirebaseService,
	embedder *memory.EmbeddingService,
) *Service {
	return &Service{
		db:                 db,
		qdrantClient:       qdrant,
		fdpnEngine:         fdpn,
		personalityService: personalitySvc,
		zetaRouter:         zeta,
		pushService:        push,
		embeddingService:   embedder,
	}
}
