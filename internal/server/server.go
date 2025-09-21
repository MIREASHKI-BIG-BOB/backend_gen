package server

import (
	"backend_gen/config"
	"backend_gen/internal/handlers/health"
	healthUC "backend_gen/internal/usecase/health"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"time"
)

type Server struct {
	cfg *config.Config

	router   *chi.Mux
	server   *http.Server
	healthUC healthUC.HealthUseCase
}

func New(cfg *config.Config) (*Server, error) {
	s := &Server{cfg: cfg}
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) init() error {
	s.initUseCases()
	s.initRouter()
	s.initHTTPServer()
	return nil
}

func (s *Server) initUseCases() {
	s.healthUC = healthUC.NewHealthUseCase()
}

func (s *Server) initHTTPServer() {
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.cfg.Server.Addr, s.cfg.Server.Port),
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
	})
}

func (s *Server) Run() {
	log.Println("Server started")
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
