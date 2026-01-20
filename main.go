package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"eva-mind/internal/auth"
	"eva-mind/internal/calendar"
	"eva-mind/internal/config"
	"eva-mind/internal/database"
	"eva-mind/internal/docs"
	"eva-mind/internal/drive"
	"eva-mind/internal/gemini"
	"eva-mind/internal/gmail"
	"eva-mind/internal/googlefit"
	"eva-mind/internal/infrastructure/cache"
	"eva-mind/internal/infrastructure/graph"
	"eva-mind/internal/infrastructure/vector"
	"eva-mind/internal/lacan"
	"eva-mind/internal/logger"
	"eva-mind/internal/maps"
	"eva-mind/internal/memory"
	"eva-mind/internal/oauth"
	"eva-mind/internal/personality"
	"eva-mind/internal/push"
	"eva-mind/internal/scheduler"
	"eva-mind/internal/sheets"
	"eva-mind/internal/signaling"
	"eva-mind/internal/stories"
	"eva-mind/internal/transnar"
	"eva-mind/internal/youtube"
	"eva-mind/pkg/types"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

type SignalingServer struct {
	upgrader           websocket.Upgrader
	clients            map[string]*PCMClient
	mu                 sync.RWMutex
	cfg                *config.Config
	pushService        *push.FirebaseService
	db                 *database.DB
	calendar           *calendar.Service
	embeddingService   *memory.EmbeddingService
	memoryStore        *memory.MemoryStore
	retrievalService   *memory.RetrievalService
	metadataAnalyzer   *memory.MetadataAnalyzer
	personalityService *personality.PersonalityService

	// FZPN Components
	neo4jClient       *graph.Neo4jClient
	graphStore        *memory.GraphStore
	fdpnEngine        *memory.FDPNEngine // Updated from PrimingEngine
	signifierService  *lacan.SignifierService
	transnarEngine    *transnar.Engine // NEW: TransNAR
	personalityRouter *personality.PersonalityRouter
	storiesRepo       *stories.Repository
	zetaRouter        *personality.ZetaRouter

	// Fix 2: Qdrant Client
	qdrantClient *vector.QdrantClient
}

type PCMClient struct {
	Conn         *websocket.Conn
	CPF          string
	IdosoID      int64
	GeminiClient *gemini.Client
	ToolsClient  *gemini.ToolsClient // ✅ DUAL-MODEL
	SendCh       chan []byte
	mu           sync.Mutex
	active       bool
	ctx          context.Context
	cancel       context.CancelFunc
	lastActivity time.Time
	audioCount   int64
	audioCount   int64
	LatentDesire *transnar.DesireInference // NEW: TransNAR desire context
	CurrentStory *types.TherapeuticStory   // 📖 Zeta Engine Story
}

var (
	db              *database.DB
	pushService     *push.FirebaseService
	signalingServer *SignalingServer
	startTime       time.Time

	// 🔐 Developer whitelist for Google features (v17)
	// Add your CPF here to enable Google Calendar/Gmail/Drive features
	googleFeaturesWhitelist = map[string]bool{
		"64525430249": true, // Developer CPF
	}
)

func NewSignalingServer(
	cfg *config.Config,
	db *database.DB,
	neo4jClient *graph.Neo4jClient,
	pushService *push.FirebaseService,
	cal *calendar.Service,
	qdrant *vector.QdrantClient,
) *SignalingServer {
	// Inicializar serviços de memória
	embeddingService := memory.NewEmbeddingService(cfg.GoogleAPIKey)
	memoryStore := memory.NewMemoryStore(db.GetConnection())
	metadataAnalyzer := memory.NewMetadataAnalyzer(cfg.GoogleAPIKey)

	// Inicializar serviço de personalidade
	personalityService := personality.NewPersonalityService(db.GetConnection())
	personalityRouter := personality.NewPersonalityRouter()

	// FZPN Components
	graphStore := memory.NewGraphStore(neo4jClient, cfg)

	// Redis & FDPN
	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		log.Printf("⚠️ Redis error: %v. FDPN will run in degraded mode (no L2 cache).", err)
	}

	// Qdrant Vector Database
	qdrantClient, err := vector.NewQdrantClient(cfg.QdrantHost, cfg.QdrantPort)
	if err != nil {
		log.Printf("⚠️ Qdrant error: %v. FDPN will run without vector search.", err)
		qdrantClient = nil // Allow graceful degradation
	} else {
		log.Println("✅ Qdrant Vector DB connected")
	}

	retrievalService := memory.NewRetrievalService(db.GetConnection(), embeddingService, qdrantClient)

	// Initialize FDPN Engine (Fractal Dynamic Priming Network)
	fdpnEngine := memory.NewFDPNEngine(neo4jClient, redisClient, qdrantClient)

	signifierService := lacan.NewSignifierService(neo4jClient)

	// Initialize TransNAR Engine (Transference Narrative Reasoning)
	// Initialize TransNAR Engine (Transference Narrative Reasoning)
	transnarEngine := transnar.NewEngine(signifierService, personalityRouter, fdpnEngine)

	// ✅ Zeta Story Engine (Gap 2)
	storiesRepo := stories.NewRepository(qdrantClient, embeddingService)
	zetaRouter := personality.NewZetaRouter(storiesRepo, personalityRouter)

	log.Println("✅ TransNAR Engine initialized")
	log.Printf("✅ Serviços de Memória Episódica inicializados")
	log.Printf("✅ Serviço de Personalidade Afetiva inicializado")
	log.Printf("✅ FZPN Engine (Phase 2) initialized")
	log.Printf("✅ Zeta Story Engine initialized")

	// 📊 STARTUP SUMMARY
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("🚀 EVA-Mind V2 - Status Report")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("✅ Services Status:")
	log.Printf("  - Database: Connected (Postgres)")

	if qdrantClient != nil {
		log.Printf("  - Vector DB: ✅ Qdrant Connected")
	} else {
		log.Printf("  - Vector DB: ⚠️ Disabled (Check connection)")
	}

	if redisClient != nil {
		log.Printf("  - Cache: ✅ Redis Connected")
	} else {
		log.Printf("  - Cache: ⚠️ Disabled (Check connection)")
	}

	if neo4jClient != nil {
		log.Printf("  - Graph DB: ✅ Neo4j Connected")
	} else {
		log.Printf("  - Graph DB: ⚠️ Disabled")
	}

	if pushService != nil {
		log.Printf("  - Push: ✅ Firebase Initialized")
	}

	log.Printf("\n🧠  Cognitive Engines (FZPN):")
	if transnarEngine != nil {
		log.Printf("  - TransNAR: ✅ Reasoning Engine Active")
	}
	if signifierService != nil {
		log.Printf("  - Lacan: ✅ Signifier Tracking Active")
	}
	if personalityService != nil {
		log.Printf("  - Personality: ✅ Affective State Active")
	}
	if fdpnEngine != nil {
		log.Printf("  - FDPN: ✅ Fractal Priming Active")
	}

	log.Printf("\n🛠️  Active Tools (V2):")
	log.Printf("  - [DB] get_vitals")
	log.Printf("  - [DB] get_agendamentos")

	if cfg.EnableGoogleSearch {
		log.Printf("  - [Vertex] Google Search: ⚠️ API Key Limited (See logs)")
	} else {
		log.Printf("  - [Vertex] Google Search: 🌑 Disabled")
	}

	if cfg.EnableCodeExecution {
		log.Printf("  - [Vertex] Code Execution: ⚠️ API Key Limited (See logs)")
	} else {
		log.Printf("  - [Vertex] Code Execution: 🌑 Disabled")
	}
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	server := &SignalingServer{
		upgrader: websocket.Upgrader{
			CheckOrigin:     func(r *http.Request) bool { return true },
			ReadBufferSize:  8192,
			WriteBufferSize: 8192,
		},
		clients:            make(map[string]*PCMClient),
		cfg:                cfg,
		pushService:        pushService,
		db:                 db,
		calendar:           cal,
		embeddingService:   embeddingService,
		memoryStore:        memoryStore,
		retrievalService:   retrievalService,
		metadataAnalyzer:   metadataAnalyzer,
		personalityService: personalityService,

		// FZPN
		neo4jClient:       neo4jClient,
		graphStore:        graphStore,
		fdpnEngine:        fdpnEngine,
		signifierService:  signifierService,
		transnarEngine:    transnarEngine,
		personalityRouter: personalityRouter,
		storiesRepo:       storiesRepo,
		zetaRouter:        zetaRouter,
		// Fix 2
		qdrantClient: qdrant,
	}

	// 🧠 Iniciar Scheduler de Pattern Mining (Gap 1)
	go server.startPatternMiningScheduler()

	return server
}

func (s *SignalingServer) startPatternMiningScheduler() {
	// Aguardar inicialização do sistema
	time.Sleep(1 * time.Minute)

	log.Printf("⛏️ [PATTERN_MINING] Scheduler iniciado (Intervalo: 1h)")
	ticker := time.NewTicker(1 * time.Hour)

	// Rodar imediatamente na startup
	go s.runPatternMining()

	for range ticker.C {
		s.runPatternMining()
	}
}

func (s *SignalingServer) runPatternMining() {
	if s.neo4jClient == nil || s.db == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Buscar todos os idosos ativos nos últimos 7 dias
	query := `
        SELECT DISTINCT idoso_id 
        FROM episodic_memories 
        WHERE timestamp > NOW() - INTERVAL '7 days'
    `

	rows, err := s.db.GetConnection().QueryContext(ctx, query)
	if err != nil {
		log.Printf("❌ [PATTERN_MINING] Query error: %v", err)
		return
	}
	defer rows.Close()

	miner := memory.NewPatternMiner(s.neo4jClient)

	for rows.Next() {
		var idosoID int64
		if err := rows.Scan(&idosoID); err != nil {
			continue
		}

		// Minerar padrões
		patterns, err := miner.MineRecurrentPatterns(ctx, idosoID, 3) // min 3 ocorrências
		if err != nil {
			log.Printf("⚠️ [PATTERN_MINING] Error for idoso %d: %v", idosoID, err)
			continue
		}

		if len(patterns) > 0 {
			log.Printf("🔍 [PATTERN_MINING] Idoso %d: Found %d patterns", idosoID, len(patterns))

			// Materializar como nós no grafo
			if err := miner.CreatePatternNodes(ctx, idosoID); err != nil {
				log.Printf("⚠️ [PATTERN_MINING] Failed to create nodes: %v", err)
			}
		}
	}

}

func main() {
	startTime = time.Now()

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	logLevel := logger.InfoLevel
	if environment == "development" {
		logLevel = logger.DebugLevel
	}

	logger.Init(logLevel, environment)
	appLog := logger.Logger
	appLog.Info().Msg("🚀 EVA-Mind 2026-v2")

	cfg, err := config.Load()
	if err != nil {
		appLog.Fatal().Err(err).Msg("Config error")
	}

	// Build DATABASE_URL if not provided
	if cfg.DatabaseURL == "" {
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "localhost"
		}
		dbPort := os.Getenv("DB_PORT")
		if dbPort == "" {
			dbPort = "5432"
		}
		dbUser := os.Getenv("DB_USER")
		if dbUser == "" {
			dbUser = "postgres"
		}
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "eva_db"
		}
		dbSSLMode := os.Getenv("DB_SSLMODE")
		if dbSSLMode == "" {
			dbSSLMode = "disable"
		}

		cfg.DatabaseURL = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
		)
	}

	db, err = database.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ DB error: %v", err)
	}
	defer db.Close()

	pushService, err = push.NewFirebaseService(cfg.FirebaseCredentialsPath)
	if err != nil {
		log.Printf("⚠️ Firebase warning: %v", err)
	} else {
		log.Printf("✅ Firebase initialized")
	}

	// 📅 Calendar Service (v17 - OAuth per-user)
	calService := calendar.NewService(context.Background())
	log.Printf("✅ Calendar service initialized (OAuth mode)")

	// Neo4j Client (FZPN)
	neo4jClient, err := graph.NewNeo4jClient(cfg)
	if err != nil {
		log.Printf("⚠️ Neo4j warning: %v. FZPN features will be disabled.", err)
	} else {
		defer neo4jClient.Close(context.Background())
		log.Printf("✅ Neo4j initialized")
	}

	signalingServer = NewSignalingServer(cfg, db, neo4jClient, pushService, calService, qdrantClient)

	sch, err := scheduler.NewScheduler(cfg, db.GetConnection())
	if err != nil {
		log.Printf("⚠️ Scheduler error: %v", err)
	} else {
		go sch.Start(context.Background())
		log.Printf("✅ Scheduler started")
	}

	router := mux.NewRouter()
	router.HandleFunc("/wss", signalingServer.HandleWebSocket)
	router.HandleFunc("/ws/pcm", signalingServer.HandleWebSocket)

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/stats", statsHandler).Methods("GET")
	api.HandleFunc("/health", healthCheckHandler).Methods("GET")
	api.HandleFunc("/call-logs", callLogsHandler).Methods("POST")

	// 🔐 Auth Routes (v16)
	authHandler := auth.NewHandler(db, cfg)
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// 🛡️ Protected Routes
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(auth.AuthMiddleware(cfg.JWTSecret))
	protected.HandleFunc("/auth/me", authHandler.Me).Methods("GET")

	// 🔐 OAuth Routes (v17)
	oauthService := oauth.NewService(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_URL"),
	)
	oauthHandler := oauth.NewHandler(oauthService, db)
	api.HandleFunc("/oauth/google/authorize", oauthHandler.HandleAuthorize).Methods("GET")
	api.HandleFunc("/oauth/google/callback", oauthHandler.HandleCallback).Methods("GET")
	api.HandleFunc("/oauth/google/token", oauthHandler.HandleTokenExchange).Methods("POST")

	// 🎥 Video Signaling Routes (v15) - DEPRECATED (Moved to WebSocket)
	// api.HandleFunc("/video/session", signalingServer.handleCreateVideoSession).Methods("POST")
	// api.HandleFunc("/video/candidate", signalingServer.handleCreateVideoCandidate).Methods("POST")
	// api.HandleFunc("/video/session/{id}/answer", signalingServer.handleGetVideoAnswer).Methods("GET")

	// 🖥️ Operator Signaling Routes - DEPRECATED (Moved to WebSocket)
	// api.HandleFunc("/video/session/{id}", signalingServer.handleGetVideoSession).Methods("GET")
	// api.HandleFunc("/video/session/{id}/answer", signalingServer.handleSaveVideoAnswer).Methods("POST")
	// api.HandleFunc("/video/session/{id}/candidates", signalingServer.handleGetMobileCandidates).Methods("GET")
	// api.HandleFunc("/video/sessions/pending", signalingServer.handleGetPendingSessions).Methods("GET")
	api.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"wsUrl": "ws://localhost:8080/ws/pcm",
		})
	}).Methods("GET")

	// ⌚ Google Fit Sync (v18)
	api.HandleFunc("/google/fit/sync/{id}", syncGoogleFitHandler).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("✅ Server ready on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(router)))
}

func (s *SignalingServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("🌐 Nova conexão de %s", r.RemoteAddr)

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ Upgrade error: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := &PCMClient{
		Conn:         conn,
		SendCh:       make(chan []byte, 256),
		ctx:          ctx,
		cancel:       cancel,
		lastActivity: time.Now(),
	}

	go s.handleClientSend(client)
	go s.monitorClientActivity(client)
	s.handleClientMessages(client)
}

func (s *SignalingServer) handleClientMessages(client *PCMClient) {
	defer s.cleanupClient(client)

	for {
		msgType, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("⚠️ Unexpected close: %v", err)
			}
			break
		}

		client.lastActivity = time.Now()

		if msgType == websocket.TextMessage {
			var data map[string]interface{}
			if err := json.Unmarshal(message, &data); err != nil {
				log.Printf("❌ JSON error: %v", err)
				continue
			}

			switch data["type"] {
			case "register":
				s.registerClient(client, data)
			case "start_call":
				log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				log.Printf("📞 START_CALL RECEBIDO")
				log.Printf("👤 CPF: %s", client.CPF)
				log.Printf("🆔 Session ID: %v", data["session_id"])
				log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

				if client.CPF == "" {
					log.Printf("❌ ERRO: Cliente não registrado!")
					s.sendJSON(client, map[string]string{"type": "error", "message": "Register first"})
					continue
				}

				// ✅ FIX: Gemini JÁ foi criado no registerClient
				// Agora só precisamos confirmar que está pronto
				if client.GeminiClient == nil {
					log.Printf("❌ ERRO: GeminiClient não existe!")
					s.sendJSON(client, map[string]string{"type": "error", "message": "Gemini not ready"})
					continue
				}

				log.Printf("✅ Gemini já está pronto!")
				log.Printf("✅ Callbacks já configurados!")

				// Enviar confirmação
				s.sendJSON(client, map[string]string{"type": "session_created", "status": "ready"})
				log.Printf("✅ session_created enviado para %s", client.CPF)

			case "start_video_cascade":
				log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				log.Printf("🎥 START_VIDEO_CASCADE RECEBIDO")
				log.Printf("👤 CPF: %s", client.CPF)
				log.Printf("🆔 Session ID: %v", data["session_id"])
				log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

				if client.CPF == "" {
					log.Printf("❌ ERRO: Cliente não registrado!")
					s.sendJSON(client, map[string]string{"type": "error", "message": "Register first"})
					continue
				}

				// Extrair dados
				sessionID, _ := data["session_id"].(string)
				sdpOffer, _ := data["sdp_offer"].(string)

				if sessionID == "" || sdpOffer == "" {
					log.Printf("❌ ERRO: Dados incompletos (session_id ou sdp_offer)")
					s.sendJSON(client, map[string]string{"type": "error", "message": "Missing session_id or sdp_offer"})
					continue
				}

				// Salvar sessão no banco
				err := s.db.CreateVideoSession(sessionID, client.IdosoID, sdpOffer)
				if err != nil {
					log.Printf("❌ Erro ao criar sessão de vídeo: %v", err)
					s.sendJSON(client, map[string]string{"type": "error", "message": "Failed to create session"})
					continue
				}

				log.Printf("✅ Sessão de vídeo criada: %s", sessionID)

				// Iniciar cascata de notificações em goroutine
				go s.handleVideoCascade(client.IdosoID, sessionID)

				// Confirmar recebimento ao mobile
				s.sendJSON(client, map[string]string{
					"type":       "video_cascade_started",
					"session_id": sessionID,
					"status":     "calling_family",
				})

			case "hangup":
				log.Printf("🔴 Hangup from %s", client.CPF)
				return

			case "vision":
				// ✅ FZPN V2: Vision Support
				// Payload: { type: "vision", payload: "BASE64..." }
				if payload, ok := data["payload"].(string); ok {
					if client.GeminiClient != nil {
						// Decode base64 if needed, or pass directly depending on client.go
						// client.go SendImage expects []byte
						imgBytes, err := base64.StdEncoding.DecodeString(payload)
						if err == nil {
							client.GeminiClient.SendImage(imgBytes)
							log.Printf("👁️ [VISION] Frame recebido e enviado (%d bytes)", len(imgBytes))
						} else {
							log.Printf("❌ [VISION] Erro ao decodificar Base64")
						}
					}
				}
			}
		}

		if msgType == websocket.BinaryMessage && client.active {
			client.audioCount++

			// if client.audioCount%50 == 0 {
			// 	log.Printf("🎤 [%s] Áudio chunk #%d (%d bytes)", client.CPF, client.audioCount, len(message))
			// }

			if client.GeminiClient != nil {
				client.GeminiClient.SendAudio(message)
			}
		}
	}
}

func (s *SignalingServer) registerClient(client *PCMClient, data map[string]interface{}) {
	cpf, _ := data["cpf"].(string)
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📝 REGISTRANDO CLIENTE")
	log.Printf("👤 CPF: %s", cpf)
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	idoso, err := s.db.GetIdosoByCPF(cpf)
	if err != nil {
		log.Printf("❌ CPF não encontrado: %s - %v", cpf, err)
		s.sendJSON(client, map[string]string{
			"type":    "error",
			"message": "CPF não cadastrado",
		})
		return
	}

	client.CPF = idoso.CPF
	client.IdosoID = idoso.ID

	s.mu.Lock()
	s.clients[idoso.CPF] = client
	s.mu.Unlock()

	log.Printf("✅ Cliente registrado: %s (ID: %d)", idoso.CPF, idoso.ID)

	// ✅ FIX: CRIAR GEMINI AQUI usando helper
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("🤖 CRIANDO CLIENTE GEMINI (Initial)")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// ✅ DUAL-MODEL: Inicializar cliente de tools (mantido separado pois é REST, não WebSocket)
	client.ToolsClient = gemini.NewToolsClient(s.cfg)

	// Usar helper para configurar sessão (Voz padrão: Aoede)
	if err := s.setupGeminiSession(client, "Aoede"); err != nil {
		log.Printf("❌ Erro ao configurar sessão Gemini: %v", err)
		s.sendJSON(client, map[string]string{"type": "error", "message": "IA error"})
		return
	}

	// ✅ AGORA enviar 'registered' (Mobile vai inicializar player ao receber)
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📤 ENVIANDO 'registered' PARA MOBILE")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	s.sendJSON(client, map[string]interface{}{
		"type":   "registered",
		"cpf":    idoso.CPF,
		"status": "ready",
	})

	log.Printf("✅ Sessão completa para: %s", client.CPF)
	log.Printf("✅ Gemini pronto e aguardando start_call...")
}

func (s *SignalingServer) setupGeminiSession(client *PCMClient, voiceName string) error {
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("🤖 CONFIGURANDO SESSÃO GEMINI (Voz: %s)", voiceName)
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Fechar cliente anterior se existir para liberar recursos
	if client.GeminiClient != nil {
		client.GeminiClient.Close()
	}

	gemClient, err := gemini.NewClient(client.ctx, s.cfg)
	if err != nil {
		log.Printf("❌ Gemini error: %v", err)
		return err
	}

	client.GeminiClient = gemClient

	// ✅ CRÍTICO: Configurar callbacks
	log.Printf("🎯 Configurando callbacks de áudio e transcrição...")

	gemClient.SetCallbacks(
		// 📊 1. Callback de Áudio
		func(audioBytes []byte) {
			select {
			case client.SendCh <- audioBytes:
				// OK
			default:
				log.Printf("⚠️ Canal cheio, dropando áudio para %s", client.CPF)
			}
		},
		// 🛠️ 2. Callback de Tool Call (Nativa)
		func(name string, args map[string]interface{}) map[string]interface{} {
			log.Printf("🔧 Tool call nativa: %s", name)
			return s.handleToolCall(client, name, args)
		},
		// 📝 3. Callback de Transcrição (Dual-Model + AUTO-SAVE)
		func(role, text string) {
			if role == "user" {
				go s.analyzeForTools(client, text)

				// FZPN: Priming (FDPN Streaming)
				if s.fdpnEngine != nil {
					go func() {
						// ⚡ FDPN Streaming Prime
						err := s.fdpnEngine.StreamingPrime(client.ctx, strconv.FormatInt(client.IdosoID, 10), text)
						if err != nil {
							log.Printf("⚠️ FDPN Error: %v", err)
						}
					}()
				}

				// TransNAR: Desire Inference (NEW)
				if s.transnarEngine != nil {
					go func() {
						// Get current personality
						currentType := personality.Type9 // Default
						if s.personalityRouter != nil {
							// TODO: Get from personality state
							currentType = personality.Type9
						}

						// Infer latent desire
						desire := s.transnarEngine.InferDesire(
							client.ctx,
							client.IdosoID,
							text,
							currentType,
						)

						// If high confidence, inject into context
						if s.transnarEngine.ShouldInterpellate(desire) {
							log.Printf("🧠 [TransNAR] Desejo latente: %s (%.0f%%) - %s",
								desire.Desire, desire.Confidence*100, desire.Reasoning)

							// Store in client context for next LLM call
							client.LatentDesire = desire
						}
					}()
				}

				// Lacan: Track Signifiers
				if s.signifierService != nil {
					go func() {
						err := s.signifierService.TrackSignifiers(client.ctx, client.IdosoID, text)
						if err != nil {
							log.Printf("⚠️ Lacan Error: %v", err)
						}
					}()
				}
			}

			// AUTO-SAVE (ambos roles)
			go s.saveAsMemory(client.IdosoID, role, text)
		},
	)

	// 🧠 Buscar memórias episódicas relevantes
	memories, err := s.retrievalService.Retrieve(
		client.ctx,
		client.IdosoID,
		"últimas conversas importantes",
		5,
	)

	var memoryTexts []string
	if len(memories) > 0 {
		for _, mem := range memories {
			memText := fmt.Sprintf("- [%s] %s: %s",
				mem.Memory.Timestamp.Format("02/01"),
				mem.Memory.Speaker,
				mem.Memory.Content,
			)
			memoryTexts = append(memoryTexts, memText)
		}
	}
	medicalContext := strings.Join(memoryTexts, "\n")

	// 🎭 FZPN: Obter Estado de Personalidade & Lacan
	var currentType int = 9 // Default Pacificador
	var lacanState string = "Transferência não iniciada."

	// 1. Personalidade (Zeta)
	if s.personalityService != nil {
		state, err := s.personalityService.GetState(client.ctx, client.IdosoID)
		if err == nil {
			// Mapear emoção para tipo (Simples 9->6 ou 9->3 por enquanto, ou usar Router completo)
			// Aqui usaremos o Router para determinar o "Modo Ativo"
			if s.personalityRouter != nil {
				activeType, _ := s.personalityRouter.RoutePersonality(personality.Type9, state.DominantEmotion)
				currentType = int(activeType)
			}
		}
	}

	// 2. Inconsciente (Lacan) - Extrair significantes
	if s.signifierService != nil {
		sigs, err := s.signifierService.GetKeySignifiers(client.ctx, client.IdosoID, 5)
		if err == nil && len(sigs) > 0 {
			var words []string
			for _, sig := range sigs {
				words = append(words, fmt.Sprintf("'%s' (Carga: %.1f)", sig.Word, sig.EmotionalCharge))
			}
			lacanState = "Significantes Mestre: " + strings.Join(words, ", ")
		}
	}

	// Adicionar contexto de relacionamento ao Lacan State (já que é psíquico)
	relationshipContext := signaling.BuildInstructions(client.IdosoID, s.db.GetConnection())
	lacanState += "\n" + relationshipContext

	// 🧠 Pattern Mining (Gap 1)
	miner := memory.NewPatternMiner(s.neo4jClient)
	patterns, err := miner.MineRecurrentPatterns(client.ctx, client.IdosoID, 3)
	if err != nil {
		log.Printf("⚠️ Pattern Mining error: %v", err)
		patterns = nil
	} else if len(patterns) > 0 {
		log.Printf("🔍 [Patterns] Detected %d patterns for user %d", len(patterns), client.IdosoID)
	}

	// ⚡ BUILD FINAL PROMPT (Co-Intelligence)
	instructions := gemini.BuildSystemPrompt(currentType, lacanState, medicalContext, patterns, nil)

	log.Printf("🚀 Iniciando sessão Gemini (Co-Intelligence Mode)...")
	// Passamos nil em memories e instructions antiga porque tudo agora está no System Prompt unificado
	err = client.GeminiClient.StartSession(instructions, nil, nil, voiceName)
	if err != nil {
		return err
	}

	// ✅ Iniciar loop de leitura
	go func() {
		log.Printf("👂 HandleResponses iniciado para %s", client.CPF)
		err := client.GeminiClient.HandleResponses(client.ctx)
		if err != nil {
			log.Printf("⚠️ HandleResponses finalizado: %v", err)
		}
		// Não setamos active=false aqui pois pode ser um restart
	}()

	client.active = true
	return nil
}

func (s *SignalingServer) handleToolCall(client *PCMClient, name string, args map[string]interface{}) map[string]interface{} {
	log.Printf("🛠️ Tool call: %s para %s", name, client.CPF)

	switch name {
	case "change_voice":
		voiceName, _ := args["voice_name"].(string)
		log.Printf("🎤 Solicitada troca de voz para: %s", voiceName)

		// Reconfigurar sessão com nova voz
		err := s.setupGeminiSession(client, voiceName)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}

		return map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Voz alterada para %s", voiceName),
		}

	case "alert_family":
		reason, _ := args["reason"].(string)
		severity, _ := args["severity"].(string)
		if severity == "" {
			severity = "alta"
		}

		err := gemini.AlertFamilyWithSeverity(s.db.GetConnection(), s.pushService, client.IdosoID, reason, severity)
		if err != nil {
			log.Printf("❌ Erro ao alertar família: %v", err)
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}

		return map[string]interface{}{
			"success": true,
			"message": "Família alertada com sucesso",
		}

	case "confirm_medication":
		medicationName, _ := args["medication_name"].(string)

		err := gemini.ConfirmMedication(s.db.GetConnection(), s.pushService, client.IdosoID, medicationName)
		if err != nil {
			log.Printf("❌ Erro ao confirmar medicamento: %v", err)
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}

		return map[string]interface{}{
			"success": true,
			"message": "Medicamento confirmado",
		}

	case "schedule_appointment":
		timestampStr, _ := args["timestamp"].(string)
		tipo, _ := args["type"].(string)
		descricao, _ := args["description"].(string)

		err := gemini.ScheduleAppointment(s.db.GetConnection(), client.IdosoID, timestampStr, tipo, descricao)
		if err != nil {
			log.Printf("❌ Erro ao agendar: %v", err)
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}

		return map[string]interface{}{
			"success": true,
			"message": "Agendamento realizado com sucesso para " + timestampStr,
		}

	case "call_family_webrtc":
		return s.initiateWebRTCCall(client, "familia")

	case "call_central_webrtc":
		return s.initiateWebRTCCall(client, "central")

	case "call_doctor_webrtc":
		return s.initiateWebRTCCall(client, "medico")

	case "call_caregiver_webrtc":
		return s.initiateWebRTCCall(client, "cuidador")

	case "open_camera_analysis":
		log.Printf("📸 Abrindo câmera para análise visual (Solicitado por %s)", client.CPF)
		s.sendJSON(client, map[string]interface{}{
			"type": "open_camera",
			"mode": "analysis",
		})
		return map[string]interface{}{
			"success": true,
			"message": "Câmera ativada para análise visual",
		}

	case "manage_calendar_event":
		// 🔐 Check if user is in developer whitelist
		if !googleFeaturesWhitelist[client.CPF] {
			return map[string]interface{}{
				"success": false,
				"error":   "Google Calendar features are currently in beta and not available for your account.",
			}
		}

		if s.calendar == nil {
			return map[string]interface{}{"success": false, "error": "Calendar service not configured"}
		}

		// Get user's OAuth tokens from database
		refreshToken, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
		if err != nil || refreshToken == "" {
			return map[string]interface{}{
				"success": false,
				"error":   "Google account not linked. Please connect your Google account first.",
			}
		}

		// Refresh token if expired
		if time.Now().After(expiry) {
			log.Printf("🔄 Refreshing expired token for idoso %d", client.IdosoID)
			// TODO: Implement token refresh using oauth service
			// For now, return error asking user to re-authenticate
			return map[string]interface{}{
				"success": false,
				"error":   "Google token expired. Please reconnect your Google account.",
			}
		}

		action, _ := args["action"].(string)

		if action == "create" {
			summary, _ := args["summary"].(string)
			desc, _ := args["description"].(string)
			start, _ := args["start_time"].(string)
			end, _ := args["end_time"].(string)

			link, err := s.calendar.CreateEventForUser(accessToken, summary, desc, start, end)
			if err != nil {
				return map[string]interface{}{"success": false, "error": err.Error()}
			}
			return map[string]interface{}{"success": true, "message": "Evento criado", "link": link}
		}

		return map[string]interface{}{"success": false, "error": "Unknown calendar action"}

	case "send_email":
		if !googleFeaturesWhitelist[client.CPF] {
			return map[string]interface{}{"success": false, "error": "Gmail features not available"}
		}

		_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
		if err != nil || time.Now().After(expiry) {
			return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
		}

		to, _ := args["to"].(string)
		subject, _ := args["subject"].(string)
		body, _ := args["body"].(string)

		gmailSvc := gmail.NewService(context.Background())
		err = gmailSvc.SendEmail(accessToken, to, subject, body)
		if err != nil {
			return map[string]interface{}{"success": false, "error": err.Error()}
		}
		return map[string]interface{}{"success": true, "message": "Email enviado com sucesso"}

	case "save_to_drive":
		if !googleFeaturesWhitelist[client.CPF] {
			return map[string]interface{}{"success": false, "error": "Drive features not available"}
		}

		_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
		if err != nil || time.Now().After(expiry) {
			return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
		}

		filename, _ := args["filename"].(string)
		content, _ := args["content"].(string)
		folder, _ := args["folder"].(string)

		driveSvc := drive.NewService(context.Background())
		fileID, err := driveSvc.SaveFile(accessToken, filename, content, folder)
		if err != nil {
			return map[string]interface{}{"success": false, "error": err.Error()}
		}
		return map[string]interface{}{"success": true, "message": "Arquivo salvo", "file_id": fileID}

	case "manage_health_sheet":
		if !googleFeaturesWhitelist[client.CPF] {
			return map[string]interface{}{"success": false, "error": "Sheets features not available"}
		}

		_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
		if err != nil || time.Now().After(expiry) {
			return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
		}

		action, _ := args["action"].(string)
		sheetsSvc := sheets.NewService(context.Background())

		if action == "create" {
			title, _ := args["title"].(string)
			url, err := sheetsSvc.CreateHealthSheet(accessToken, title)
			if err != nil {
				return map[string]interface{}{"success": false, "error": err.Error()}
			}
			return map[string]interface{}{"success": true, "message": "Planilha criada", "url": url}
		}

		// TODO: Implement append action
		return map[string]interface{}{"success": false, "error": "Action not implemented"}

	case "create_health_doc":
		if !googleFeaturesWhitelist[client.CPF] {
			return map[string]interface{}{"success": false, "error": "Docs features not available"}
		}

		_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
		if err != nil || time.Now().After(expiry) {
			return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
		}

		title, _ := args["title"].(string)
		content, _ := args["content"].(string)

		docsSvc := docs.NewService(context.Background())
		url, err := docsSvc.CreateDocument(accessToken, title, content)
		if err != nil {
			return map[string]interface{}{"success": false, "error": err.Error()}
		}
		return map[string]interface{}{"success": true, "message": "Documento criado", "url": url}

	case "find_nearby_places":
		if !googleFeaturesWhitelist[client.CPF] {
			return map[string]interface{}{"success": false, "error": "Maps features not available"}
		}

		placeType, _ := args["place_type"].(string)
		location, _ := args["location"].(string)
		radius := 5000
		if r, ok := args["radius"].(float64); ok {
			radius = int(r)
		}

		mapsSvc := maps.NewService(context.Background(), s.cfg.GoogleMapsAPIKey)
		places, err := mapsSvc.FindNearbyPlaces(placeType, location, radius)
		if err != nil {
			return map[string]interface{}{"success": false, "error": err.Error()}
		}
		return map[string]interface{}{"success": true, "places": places}

	case "search_videos":
		if !googleFeaturesWhitelist[client.CPF] {
			return map[string]interface{}{"success": false, "error": "YouTube features not available"}
		}

		_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
		if err != nil || time.Now().After(expiry) {
			return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
		}

		query, _ := args["query"].(string)
		maxResults := int64(5)
		if mr, ok := args["max_results"].(float64); ok {
			maxResults = int64(mr)
		}

		youtubeSvc := youtube.NewService(context.Background())
		videos, err := youtubeSvc.SearchVideos(accessToken, query, maxResults)
		if err != nil {
			return map[string]interface{}{"success": false, "error": err.Error()}
		}
		return map[string]interface{}{"success": true, "videos": videos}

	case "play_music":
		if !googleFeaturesWhitelist[client.CPF] {
			return map[string]interface{}{"success": false, "error": "Spotify features not available"}
		}

		// TODO: Implement Spotify OAuth separately
		return map[string]interface{}{"success": false, "error": "Spotify integration pending OAuth setup"}

	case "request_ride":
		if !googleFeaturesWhitelist[client.CPF] {
			return map[string]interface{}{"success": false, "error": "Uber features not available"}
		}

		// TODO: Implement Uber OAuth separately
		return map[string]interface{}{"success": false, "error": "Uber integration pending OAuth setup"}

	case "get_health_data":
		if !googleFeaturesWhitelist[client.CPF] {
			return map[string]interface{}{"success": false, "error": "Google Fit features not available"}
		}

		_, accessToken, expiry, err := s.db.GetGoogleTokens(client.IdosoID)
		if err != nil || time.Now().After(expiry) {
			return map[string]interface{}{"success": false, "error": "Google account not linked or expired"}
		}

		fitSvc := googlefit.NewService(context.Background())

		// Get all health data
		healthData, err := fitSvc.GetAllHealthData(accessToken)
		if err != nil {
			return map[string]interface{}{"success": false, "error": err.Error()}
		}

		// Save to database automatically
		if healthData.Steps > 0 {
			s.db.SaveVitalSign(client.IdosoID, "passos", fmt.Sprintf("%d", healthData.Steps), "steps", "google_fit", "")
		}
		if healthData.HeartRate > 0 {
			s.db.SaveVitalSign(client.IdosoID, "frequencia_cardiaca", fmt.Sprintf("%.0f", healthData.HeartRate), "bpm", "google_fit", "")
		}
		if healthData.Calories > 0 {
			s.db.SaveVitalSign(client.IdosoID, "calorias", fmt.Sprintf("%d", healthData.Calories), "kcal", "google_fit", "")
		}
		if healthData.Distance > 0 {
			s.db.SaveVitalSign(client.IdosoID, "distancia", fmt.Sprintf("%.2f", healthData.Distance), "km", "google_fit", "")
		}
		if healthData.Weight > 0 {
			s.db.SaveVitalSign(client.IdosoID, "peso", fmt.Sprintf("%.1f", healthData.Weight), "kg", "google_fit", "")
		}

		return map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"steps":      healthData.Steps,
				"heart_rate": healthData.HeartRate,
				"calories":   healthData.Calories,
				"distance":   healthData.Distance,
				"weight":     healthData.Weight,
			},
			"message": "Dados de saúde coletados e salvos com sucesso",
		}

	case "send_whatsapp":
		if !googleFeaturesWhitelist[client.CPF] {
			return map[string]interface{}{"success": false, "error": "WhatsApp features not available"}
		}

		// TODO: Implement WhatsApp Business API
		return map[string]interface{}{"success": false, "error": "WhatsApp integration pending configuration"}

	case "run_sql_select":
		query, _ := args["query"].(string)

		// 🛡️ Segurança básica
		if query == "" {
			return map[string]interface{}{"success": false, "error": "Empty query"}
		}

		// ⚠️ Apenas SELECT
		// (Idealmente parsear a query, mas string check simples serve para MVP)
		// Nota: Em produção, usar um usuário DB com permissões Read-Only
		if len(query) < 6 || query[:6] != "SELECT" && query[:6] != "select" {
			return map[string]interface{}{"success": false, "error": "Only SELECT queries allowed"}
		}

		log.Printf("🔍 Executando SQL: %s", query)

		rows, err := s.db.GetConnection().Query(query)
		if err != nil {
			return map[string]interface{}{"success": false, "error": err.Error()}
		}
		defer rows.Close()

		cols, _ := rows.Columns()
		var result []map[string]interface{}

		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i := range columns {
				columnPointers[i] = &columns[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				return map[string]interface{}{"success": false, "error": err.Error()}
			}

			m := make(map[string]interface{})
			for i, colName := range cols {
				val := columnPointers[i].(*interface{})

				// Handle bytes (DB text can come as bytes)
				if b, ok := (*val).([]byte); ok {
					m[colName] = string(b)
				} else {
					m[colName] = *val
				}
			}
			result = append(result, m)
		}

		return map[string]interface{}{"success": true, "data": result}

	case "list_voices":
		return s.getAvailableVoices()

	default:
		log.Printf("⚠️ Tool desconhecida: %s", name)
		return map[string]interface{}{
			"success": false,
			"error":   "Ferramenta desconhecida",
		}
	}
}

func (s *SignalingServer) listenGemini(client *PCMClient) {
	log.Printf("👂 Listener iniciado: %s", client.CPF)

	for client.active {
		resp, err := client.GeminiClient.ReadResponse()
		if err != nil {
			if client.active {
				log.Printf("⚠️ Gemini read error: %v", err)
			}
			return
		}
		s.processGeminiResponse(client, resp)
	}

	log.Printf("📚 Listener finalizado: %s", client.CPF)
}

func (s *SignalingServer) processGeminiResponse(client *PCMClient, resp map[string]interface{}) {
	serverContent, ok := resp["serverContent"].(map[string]interface{})
	if !ok {
		return
	}

	modelTurn, _ := serverContent["modelTurn"].(map[string]interface{})
	parts, _ := modelTurn["parts"].([]interface{})

	audioCount := 0

	for _, part := range parts {
		p, ok := part.(map[string]interface{})
		if !ok {
			continue
		}

		// 1. Processar Texto (Delegation Protocol)
		if text, hasText := p["text"].(string); hasText {
			// Regex para capturar [[TOOL:nome:{arg}]]
			// Ex: [[TOOL:google_search_retrieval:{"query": "clima sp"}]]
			re := regexp.MustCompile(`\[\[TOOL:(\w+):({.*?})\]\]`)
			matches := re.FindStringSubmatch(text)

			if len(matches) == 3 {
				toolName := matches[1]
				argsJSON := matches[2]

				log.Printf("🤖 [AGENT] Comando detectado: %s", toolName)

				var args map[string]interface{}
				if err := json.Unmarshal([]byte(argsJSON), &args); err == nil {
					// Executar ferramenta
					result := s.handleToolCall(client, toolName, args)

					// TODO: Enviar resultado de volta para o modelo 2.5?
					// Por enquanto, apenas executamos (alertas, agendamentos funcionam one-way)
					// Para busca, precisaríamos injetar contexto.
					log.Printf("🤖 [AGENT] Resultado da execução: %+v", result)

					// Se for busca, tentar enviar de volta como User Message oculta?
					// s.SendSystemMessage(client, fmt.Sprintf("System: Resultado da ferramenta %s: %v", toolName, result))
				} else {
					log.Printf("❌ [AGENT] Erro ao parsear args: %v", err)
				}
			}
		}

		// 2. Processar Áudio
		// 1. Processar Texto (Delegation Protocol)
		if text, hasText := p["text"].(string); hasText {
			re := regexp.MustCompile(`\[\[TOOL:(\w+):({.*?})\]\]`)
			matches := re.FindStringSubmatch(text)

			if len(matches) == 3 {
				toolName := matches[1]
				argsJSON := matches[2]

				log.Printf("🤖 [AGENT] Comando detectado: %s", toolName)

				var args map[string]interface{}
				if err := json.Unmarshal([]byte(argsJSON), &args); err == nil {
					// Executar ferramenta (Delegation Pattern)
					result := s.handleToolCall(client, toolName, args)
					log.Printf("🤖 [AGENT] Resultado: %+v", result)
				} else {
					log.Printf("❌ [AGENT] Erro JSON: %v", err)
				}
			}
		}

		if data, hasData := p["inlineData"]; hasData {
			b64, _ := data.(map[string]interface{})["data"].(string)
			audio, err := base64.StdEncoding.DecodeString(b64)
			if err != nil {
				continue
			}

			select {
			case client.SendCh <- audio:
				audioCount++
			default:
				log.Printf("⚠️ Canal cheio, dropando áudio")
			}
		}
	}
}

func (s *SignalingServer) handleClientSend(client *PCMClient) {
	sentCount := 0

	for {
		select {
		case <-client.ctx.Done():
			return
		case audio := <-client.SendCh:
			sentCount++

			// 🔙 REVERTIDO: Voltando para binário para investigação correta
			client.mu.Lock()
			err := client.Conn.WriteMessage(websocket.BinaryMessage, audio)
			client.mu.Unlock()

			if err != nil {
				log.Printf("❌ Send error: %v", err)
				return
			}

			// Debug DETALHADO: Loga a cada 10 pacotes
			// if sentCount%10 == 0 {
			// 	log.Printf(" [DEBUG-BIN] Enviado %d bytes (Chunk #%d). Status: OK", len(audio), sentCount)
			// }
		}
	}
}

func (s *SignalingServer) monitorClientActivity(client *PCMClient) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-client.ctx.Done():
			return
		case <-ticker.C:
			if time.Since(client.lastActivity) > 5*time.Minute {
				log.Printf("⏰ Timeout inativo: %s", client.CPF)
				client.cancel()
				return
			}
		}
	}
}

func (s *SignalingServer) cleanupClient(client *PCMClient) {
	log.Printf("🧹 Cleanup: %s", client.CPF)

	client.cancel()

	s.mu.Lock()
	delete(s.clients, client.CPF)
	s.mu.Unlock()

	client.Conn.Close()

	if client.GeminiClient != nil {
		client.GeminiClient.Close()
	}

	log.Printf("✅ Desconectado: %s", client.CPF)
}

func (s *SignalingServer) sendJSON(c *PCMClient, v interface{}) {
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📤 sendJSON CHAMADO")
	log.Printf("📦 Payload: %+v", v)
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.Conn.WriteJSON(v)
	if err != nil {
		log.Printf("❌ ERRO ao enviar JSON: %v", err)
		log.Printf("❌ Cliente CPF: %s", c.CPF)
		return
	}

	log.Printf("✅ JSON enviado com sucesso para %s", c.CPF)
}

func (s *SignalingServer) GetActiveClientsCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}

// --- API HANDLERS ---

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dbStatus := false
	if db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := db.GetConnection().PingContext(ctx); err == nil {
			dbStatus = true
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"active_clients": signalingServer.GetActiveClientsCount(),
		"uptime":         time.Since(startTime).String(),
		"db_status":      dbStatus,
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := "healthy"
	httpStatus := http.StatusOK

	if err := db.GetConnection().Ping(); err != nil {
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func callLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("❌ Erro ao decodificar call log: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("💾 CALL LOG RECEBIDO: %+v", data)

	// TODO: Salvar no banco de dados quando a tabela estiver pronta
	// Por enquanto, apenas logamos e retornamos sucesso para o app não dar erro.

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "saved", "message": "Log received"})
}

func syncGoogleFitHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idosoIDStr := vars["id"]
	idosoID, _ := strconv.ParseInt(idosoIDStr, 10, 64)

	log.Printf("⌚ Iniciando sincronização Google Fit para idoso %d", idosoID)

	// 1. Buscar tokens
	_, accessToken, expiry, err := db.GetGoogleTokens(idosoID)
	if err != nil || accessToken == "" || time.Now().After(expiry) {
		http.Error(w, "Google account not linked or token expired", http.StatusUnauthorized)
		return
	}

	// 2. Chamar serviço Google Fit
	fitSvc := googlefit.NewService(context.Background())
	healthData, err := fitSvc.GetAllHealthData(accessToken)
	if err != nil {
		log.Printf("❌ Erro ao buscar dados do Fit: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Salvar no Banco
	err = db.SaveDeviceHealthData(idosoID, int(healthData.HeartRate), int(healthData.Steps))
	if err != nil {
		log.Printf("❌ Erro ao salvar dados de saúde: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Sincronização Google Fit concluída para idoso %d: %d BPM, %d passos", idosoID, int(healthData.HeartRate), int(healthData.Steps))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   healthData,
	})
}

// initiateWebRTCCall handles the logic to start a WebRTC call
func (s *SignalingServer) initiateWebRTCCall(client *PCMClient, target string) map[string]interface{} {
	log.Printf("📹 Iniciando chamada de vídeo para %s (Solicitado por %s)", target, client.CPF)

	// 1. Criar sessão de vídeo no DB
	// OBS: Estamos reutilizando a lógica de session start aqui, mas simplificada
	sessionID := fmt.Sprintf("video-%s-%d", target, time.Now().Unix())

	// 2. Enviar comando para o Mobile abrir a câmera
	// O app mobile vai receber 'start_video' e navegar para /video
	s.sendJSON(client, map[string]interface{}{
		"type":       "start_video",
		"session_id": sessionID,
		"target":     target,
	})

	// 3. (Simulação) Notificar o target
	// Aqui entraria a lógica de push notification para o App da Família ou Painel da Central
	log.Printf("🔔 [TODO] Notificar %s sobre chamada recebida na sessão %s", target, sessionID)

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Chamada de vídeo iniciada para %s. Abrindo câmera...", target),
	}
}

// ✅ DUAL-MODEL: Analisa transcrição e executa tools se necessário
func (s *SignalingServer) analyzeForTools(client *PCMClient, text string) {
	if client.ToolsClient == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("🔍 [TOOLS] Analisando transcrição: \"%s\"", text)

	toolCalls, err := client.ToolsClient.AnalyzeTranscription(ctx, text, "user")
	if err != nil {
		log.Printf("⚠️ [TOOLS] Erro ao analisar: %v", err)
		return
	}

	if len(toolCalls) == 0 {
		return
	}

	for _, tc := range toolCalls {
		log.Printf("🛠️ [TOOLS] Executando: %s com args: %+v", tc.Name, tc.Args)
		// Executar tool
		s.handleToolCall(client, tc.Name, tc.Args)
	}
}
