package main

import (
	"ditdah/internal/adapter/config"
	"ditdah/internal/adapter/postgres/storage"
	router "ditdah/internal/controller/http"

	"log/slog"
	"os"

	auth "ditdah/internal/features/auth"
	user "ditdah/internal/features/user"
	lesson "ditdah/internal/features/lesson"

	"github.com/joho/godotenv"
)

func main()  {
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err := os.MkdirAll("data", 0o755); err != nil {
		slog.Error("create data dir", "error", err)
		return
	}

	store, err := storage.New(cfg.DBPath)
	if err != nil {
		slog.Error("open storage", "error", err)
		return
	}
	defer store.Close()

	uRepo := user.New(store.DB)
	lRepo := lesson.New(store.DB)

	uUseCase := user.NewUserUseCase(uRepo)
	aUseCase := auth.NewAuthUseCase(uRepo, cfg.JwtSecret)
	lessonUseCase := lesson.NewLessonUseCase(lRepo, uRepo)

	r := router.NewRouter(router.UseCases{
		Auth:      aUseCase,
		User:      uUseCase,
		Lesson:    lessonUseCase,

		JWTSecret: cfg.JwtSecret,
		DB:        store.DB,
	})

	slog.Info("server listening", "addr", cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		slog.Error("server failed", "error", err)
	}
}