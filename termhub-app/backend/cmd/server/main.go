// Package main is the entry point for the TermHub backend server.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/termhub/backend/internal/config"
	"github.com/termhub/backend/internal/handler"
	"github.com/termhub/backend/internal/middleware"
	"github.com/termhub/backend/internal/policy"
	ptymanager "github.com/termhub/backend/internal/pty"
	"github.com/termhub/backend/internal/recording"
	"github.com/termhub/backend/internal/repository"
	"github.com/termhub/backend/internal/service"
	"github.com/termhub/backend/internal/ws"
)

func main() {
	// Configure zerolog with human-friendly console output for development.
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Timestamp().
		Caller().
		Logger()

	// Load configuration from environment variables.
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	// Set the global log level.
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	log.Info().Str("port", cfg.Port).Msg("starting TermHub backend")

	// Connect to PostgreSQL.
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer dbpool.Close()

	if err := dbpool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}
	log.Info().Msg("connected to PostgreSQL")

	// Initialize PTY Manager.
	ptyMgr := ptymanager.NewManager()

	// Initialize repositories.
	userRepo := repository.NewUserRepository(dbpool)
	terminalRepo := repository.NewTerminalRepository(dbpool)

	// Initialize services.
	authService := service.NewAuthService(userRepo, cfg)
	terminalService := service.NewTerminalService(terminalRepo, ptyMgr)
	auditService := service.NewAuditService(dbpool)

	// Seed default admin user if no users exist.
	seedAdminUser(ctx, authService, userRepo)

	// Initialize WebSocket hub.
	hub := ws.NewHub()
	go hub.Run()

	// Initialize policy engine.
	policyEngine := policy.NewEngine()

	// Initialize recording manager.
	recordingMgr := recording.NewManager("./recordings")

	// Initialize handlers.
	authHandler := handler.NewAuthHandler(authService, auditService, cfg.JWTSecret)
	terminalHandler := handler.NewTerminalHandler(terminalService, auditService, cfg.JWTSecret)
	adminHandler := handler.NewAdminHandler(auditService, terminalService, userRepo, ptyMgr, hub, cfg.JWTSecret)
	wsHandler := handler.NewWSHandler(hub, ptyMgr, terminalService, authService, auditService, policyEngine, recordingMgr)
	healthHandler := handler.NewHealthHandler(dbpool)
	recordingHandler := handler.NewRecordingHandler(recordingMgr, authService, cfg.JWTSecret)

	// Set up the chi router with global middleware.
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logging(log.Logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS(cfg.CORSOrigins).Handler)

	// Mount route groups.
	r.Route("/api", func(r chi.Router) {
		r.Mount("/auth", authHandler.Routes())
		r.Mount("/terminals", terminalHandler.Routes())
		r.Mount("/admin", adminHandler.Routes())
		r.Mount("/ws", wsHandler.Routes())
		r.Mount("/health", healthHandler.Routes())
		r.Mount("/recordings", recordingHandler.Routes())

		// Terminal-scoped recording routes: POST /api/terminals/{id}/recording/start|stop
		r.Route("/terminals/{id}/recording", func(r chi.Router) {
			r.Mount("/", recordingHandler.TerminalRecordingRoutes())
		})
	})

	// Create the HTTP server.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine.
	go func() {
		log.Info().Str("addr", srv.Addr).Msg("HTTP server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	// Wait for interrupt signal for graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("shutting down server")

	// Clean up PTY sessions.
	ptyMgr.Cleanup()

	// Clean up active recordings.
	recordingMgr.Cleanup()

	// Give outstanding requests 30 seconds to complete.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server stopped gracefully")
}

// seedAdminUser creates a default admin user if no users exist in the database.
func seedAdminUser(ctx context.Context, authService *service.AuthService, userRepo *repository.UserRepository) {
	count, err := userRepo.Count(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to check user count for seeding; skipping admin seed")
		return
	}

	if count > 0 {
		log.Debug().Int64("user_count", count).Msg("users already exist, skipping admin seed")
		return
	}

	user, err := authService.Register(ctx, "admin", "admin123", "admin")
	if err != nil {
		log.Error().Err(err).Msg("failed to seed admin user")
		return
	}

	log.Info().
		Str("username", user.Username).
		Str("role", string(user.Role)).
		Str("id", user.ID.String()).
		Msg("seeded default admin user (username: admin, password: admin123)")
}
