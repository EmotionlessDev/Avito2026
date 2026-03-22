package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/config"
	authHttp "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/auth/delivery/http"
	authDummyLogin "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/auth/usecases"
	userHttp "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/users/delivery/http"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/middleware"

	_ "github.com/lib/pq"
)

func main() {
	// Init config
	cfg := config.New(0, "", "")

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.DB.DSN, "dsn", os.Getenv("BOOKING_POSTGRES_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.Auth.JWTSecret, "jwt", os.Getenv("JWT_SECRET"), "JWT secret key")
	flag.Parse()

	// Init logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Init DB
	db, err := openDB(cfg)
	if err != nil {
		logger.Error("cannot connect to db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Error("error closing db", slog.String("error", err.Error()))
		}
	}()
	// JWT secret
	jwtSecret := cfg.GetJWTSecret()

	// Init usecases
	authUsecase := authDummyLogin.NewDummyLogin(jwtSecret)

	// Init handlers
	authHandler := authHttp.NewHandler(authUsecase)
	helloHandler := userHttp.NewHandler()

	// Init middlewares
	authMW := middleware.JWTMiddleware(jwtSecret)

	// Init serveMux
	mux := http.NewServeMux()
	mux.HandleFunc("/dummyLogin", authHandler.DummyLogin)
	mux.Handle("/hello", middleware.Chain(http.HandlerFunc(helloHandler.HelloUser), authMW))

	// Create http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("starting server",
		slog.Int("port", cfg.Port),
		slog.String("env", cfg.Env),
	)

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error("cannot start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func openDB(cfg config.ConfigProvider) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetDBDSN())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
