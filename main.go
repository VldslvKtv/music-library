package main

import (
	"log/slog"
	"music_library/config"
	"music_library/internal/http_server/handlers/add_song"
	"music_library/internal/http_server/handlers/delete_song"
	"music_library/internal/http_server/handlers/get_all_data"
	"music_library/internal/http_server/handlers/get_song"
	"music_library/internal/http_server/handlers/update_song"
	"music_library/internal/http_server/lib/logger"
	"music_library/internal/http_server/storage/pg"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "music_library/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Music Library API
// @version 1.0
// @description This is a sample server for a music library.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8002
// @BasePath /

func main() {
	config := config.MustLoad()
	log := setupLogger(config.Env)
	log.Info("start app", slog.String("env", config.Env))

	_ = log

	storage, err := pg.New(&config)
	if err != nil {
		log.Error("failed to init storage", logger.Err(err))
		os.Exit(1)
	}
	_ = storage
	log.Info("migration run is completed")
	defer storage.Close()

	router := chi.NewRouter()
	router.Use(middleware.Recoverer) // воостановление полсе паники (чтобы не падало приложение после 1 ошибки в хендлере)
	router.Use(middleware.URLFormat)

	// Swagger UI
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Route("/get_data", func(r chi.Router) {
		r.Get("/songs", get_all_data.New(log, storage))
		r.Get("/text", get_song.New(log, storage))
	})
	router.Post("/add", add_song.New(log, config.ExtAPIUrl, storage))
	router.Delete("/delete/{id}", delete_song.New(log, storage))
	router.Patch("/update/{id}", update_song.New(log, storage))

	log.Info("starting server", slog.String("address", config.Address))

	srv := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT) // graceful shutdown
	check := <-stop

	log.Debug("server stopped", slog.String("signal", check.String()))
}

// Настройка уровня логирования
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
