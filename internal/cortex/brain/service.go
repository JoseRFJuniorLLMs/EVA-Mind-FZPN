package brain

import (
	"database/sql"
	"eva-mind/internal/brainstem/infrastructure/vector"
	"eva-mind/internal/brainstem/push"
	"eva-mind/internal/cortex/lacan"
	ps "eva-mind/internal/cortex/personality"
	"eva-mind/internal/hippocampus/knowledge"
	"eva-mind/internal/hippocampus/memory"
)

// Service encapsulates the cognitive functions of EVA
type Service struct {
	db                 *sql.DB
	qdrantClient       *vector.QdrantClient
	fdpnEngine         *lacan.FDPNEngine // Updated type
	personalityService *ps.PersonalityService
	zetaRouter         *ps.ZetaRouter
	pushService        *push.FirebaseService
	embeddingService   *memory.EmbeddingService // This one keeps using memory package for now as I created embedding_service in hippocampus/knowledge but service.go uses memory.EmbeddingService?
	// Wait, I created internal/hippocampus/knowledge/embedding_service.go
	// But existing service.go has `*memory.EmbeddingService`.
	// I should update it to `*knowledge.EmbeddingService` if I want to use the new one.
	// The previous refactoring moved `internal/memory` to `internal/hippocampus/memory`.
	// So `embeddingService` in service.go is likely referring to `internal/hippocampus/memory`.
	// The new one is in `internal/hippocampus/knowledge`.
	// I will add unifiedRetrieval and potentially update existing fields if they are redundant.

	knowledgeEmbedder *knowledge.EmbeddingService // New one
	unifiedRetrieval  *lacan.UnifiedRetrieval
}

// NewService creates a new Brain service
func NewService(
	db *sql.DB,
	qdrant *vector.QdrantClient,
	// fdpn *memory.FDPNEngine, // Replaced
	unified *lacan.UnifiedRetrieval,
	personalitySvc *ps.PersonalityService,
	zeta *ps.ZetaRouter,
	push *push.FirebaseService,
	embedder *memory.EmbeddingService, // Restored
) *Service {
	return &Service{
		db:                 db,
		qdrantClient:       qdrant,
		personalityService: personalitySvc,
		zetaRouter:         zeta,
		pushService:        push,
		embeddingService:   embedder, // Restored
		unifiedRetrieval:   unified,
	}
}
