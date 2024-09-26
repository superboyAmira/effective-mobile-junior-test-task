package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"online-song-library/internal/controller"
	"online-song-library/internal/model"
	"online-song-library/internal/repository"
	"online-song-library/internal/router"
	"online-song-library/internal/service"
	"online-song-library/pkg/logger"
	"online-song-library/pkg/storage/postgresql"
	test "online-song-library/test/external_api_mock"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const (
	envPath = "../../config/.env"
)

func main() {
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("WARN:didn't load .env file")
	}

	log := logger.SetupLogger()

	// mock external ip
	if os.Getenv("EXTERNAL_API_HTTPTEST_SERVER") == "true" {
		externalAPI := test.CreateMockExternalAPIServer(log)
		defer externalAPI.Close()
		log.Info("Mock server is starting", slog.String("port", externalAPI.URL))
		os.Setenv("PATH_EXTERNAL_API_HTTPTEST_SERVER", strings.Split(externalAPI.URL, ":")[2])
	}

	// db
	db, err := postgresql.Connect(log)
	if err != nil {
		log.Error("unable to connect db", slog.String("err", err.Error()))
		return
	}
	log.Info("db connection successfully", slog.String("port", os.Getenv("DB_PORT")), slog.String("db_name", os.Getenv("DB_NAME")))

	err = postgresql.Migrate(db, model.Song{})
	if err != nil {
		log.Error("unable to migrate entity", slog.String("err", err.Error()))
		return
	}
	log.Info("db migration successfully", slog.String("port", os.Getenv("DB_PORT")), slog.String("db_name", os.Getenv("DB_NAME")))

	// logic
	repo := repository.NewSongRepository(db)
	serv := service.NewSongService(repo)
	cntrler := controller.NewSongController(serv, log)
	ginRouter := router.SetupRouter(cntrler, log)

	server := http.Server{
		Addr:    os.Getenv("API_PORT"),
		Handler: ginRouter,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTSTP)

	// start server
	go func() {
		log.Info("Server is starting", slog.String("port", os.Getenv("API_PORT")))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("ListenAndServe error: %v", slog.String("err", err.Error()))
		}
	}()

	// gracefull stop
	<-quit

	log.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown: %v", slog.String("err", err.Error()))
	}

	log.Info("Server exited gracefully")
}
