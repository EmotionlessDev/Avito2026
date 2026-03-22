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
	roomHttp "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms/delivery/http"
	roomStorage "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms/storage"
	roomUsecase "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/rooms/usecases"
	scheduleHttp "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules/delivery/http"
	scheduleStorage "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules/storage"
	scheduleUsecase "github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules/usecases"
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

	// Init storages
	roomStorage := roomStorage.NewStorage()
	scheduleStorage := scheduleStorage.NewStorage()

	// Init usecases
	authUsecase := authDummyLogin.NewDummyLogin(jwtSecret)

	createRoomUsecase := roomUsecase.NewCreateRoom(roomStorage, db)
	getRoomsUsecase := roomUsecase.NewGetRooms(roomStorage, db)

	createScheduleUsecase := scheduleUsecase.NewCreateSchedule(scheduleStorage, db)

	// Init handlers
	authHandler := authHttp.NewHandler(authUsecase)

	createRoomHandler := roomHttp.NewCreateHandler(createRoomUsecase)
	getRoomsHandler := roomHttp.NewGetHandler(getRoomsUsecase)

	createScheduleHandler := scheduleHttp.NewScheduleHandler(createScheduleUsecase)

	// Init serveMux
	mux := http.NewServeMux()
	mux.HandleFunc("/dummyLogin", authHandler.DummyLogin)
	// rooms
	mux.Handle("/rooms/create", middleware.Chain(
		http.HandlerFunc(createRoomHandler.CreateRoom),
		middleware.JWTMiddleware(jwtSecret),
		middleware.RoleBased("admin")),
	)
	mux.Handle("/rooms/list", middleware.Chain(
		http.HandlerFunc(getRoomsHandler.GetRooms),
		middleware.JWTMiddleware(jwtSecret)),
	)
	// schedules
	mux.Handle("/rooms/{roomId}/schedule/create", middleware.Chain(
		http.HandlerFunc(createScheduleHandler.CreateSchedule),
		middleware.JWTMiddleware(jwtSecret),
		middleware.RoleBased("admin")),
	)

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
