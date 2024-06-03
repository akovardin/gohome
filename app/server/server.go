package server

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"go.uber.org/fx"

	"gohome.4gophers.ru/getapp/gohome/app/handlers/home"
)

type Config struct {
	Addr string
}

type Server struct {
	server *http.Server

	home *home.Handler
}

func New(lc fx.Lifecycle, cfg Config, home *home.Handler) *Server {
	s := Server{
		home: home,
		server: &http.Server{
			Addr: cfg.Addr,
		},
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.server.Handler = s.routing()
			s.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.Stop(ctx)
			return nil
		},
	})

	return &Server{}
}

func (s *Server) Start() {
	go func() {
		log.Printf("start server on %s\n", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	s.server.Shutdown(ctx)
}

func (s *Server) routing() http.Handler {
	logger := httplog.NewLogger("gohome", httplog.Options{
		JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(logger))
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler, middleware.Recoverer, middleware.NoCache)

	r.Get("/*", s.home.Home)

	return r
}
