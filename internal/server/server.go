package server

import (
	"backend_gen/config"
	generatorAdapter "backend_gen/internal/adapter/generator"
	wsAdapter "backend_gen/internal/adapter/websocket"
	"backend_gen/internal/handlers/health"
	wsHandler "backend_gen/internal/handlers/websocket"
	"backend_gen/internal/ports/generator"
	"backend_gen/internal/ports/websocket"
	"backend_gen/internal/usecase"
	healthUC "backend_gen/internal/usecase/health"
	wsUC "backend_gen/internal/usecase/websocket"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	cfg *config.Config

	router *chi.Mux
	server *http.Server

	// adapters
	wsClient      websocket.Client
	dataGenerator generator.DataGenerator

	// usecases
	healthUC         usecase.HealthUseCase
	websocketUseCase usecase.WebSocketUseCase
}

func New(cfg *config.Config) (*Server, error) {
	s := &Server{cfg: cfg}
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) init() error {
	s.initAdapters()
	s.initUseCases()
	s.initRouter()
	s.initHTTPServer()
	return nil
}

func (s *Server) initAdapters() {
	s.wsClient = wsAdapter.NewClient()
	s.dataGenerator = generatorAdapter.NewSinusoidalGenerator()
}

func (s *Server) initUseCases() {
	s.healthUC = healthUC.NewHealthUseCase()
	s.websocketUseCase = wsUC.NewWebSocketUseCase(s.wsClient, s.dataGenerator)
}

func (s *Server) initHTTPServer() {
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Addr, s.cfg.Server.Port)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) initRouter() {
	s.router = chi.NewRouter()
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Content-Length"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.router.Route("/api", func(r chi.Router) {
		r.Get("/health", health.NewHealthHandler(s.healthUC))
		r.Get("/on", wsHandler.OnSocket(s.websocketUseCase))
		r.Get("/off", wsHandler.OffSocket(s.websocketUseCase))
	})
}

func (s *Server) Run() {
	slog.Info("Starting HTTP server", "addr", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
