package main

import (
	"CommentTree/comment_tree/adapters/db"
	"CommentTree/comment_tree/adapters/rest"
	"CommentTree/comment_tree/config"
	"CommentTree/comment_tree/pkg/handler"
	httpserver "CommentTree/comment_tree/pkg/http_server"
	"CommentTree/comment_tree/pkg/logger"
	"CommentTree/comment_tree/usecase"
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "config file path")
	flag.Parse()

	cfg := config.MustLoadConfig(configPath)

	log := logger.New(cfg.AppConfig.Env)

	log.Info("Load config success")
	log.Debug("Debug mode enabled")

	if err := run(cfg, log); err != nil {
		log.Error("Run err:", "error", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config, log logger.Interface) error {
	connDB, err := db.NewDB(log, cfg.DBConfig)
	if err != nil {
		return err
	}
	log.Debug("DB init success", "db", *connDB)

	validate := validator.New()

	uc := usecase.NewUseCase(connDB, log, validate)

	h := handler.New()
	h.AddHandler("/", http.FileServer(http.Dir("./web")))

	h.AddHandlerFunc("POST /comments", rest.NewCreateHandler(uc))
	h.AddHandlerFunc("GET /comments", rest.NewGetHandler(uc))
	h.AddHandlerFunc("DELETE /comments/{id}", rest.NewDeleteHandler(uc))
	log.Debug("Handler init success", "handler", *h)

	server := httpserver.New(h, cfg.ServerConfig)
	log.Debug("Server init success", "server", *server)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Info("Shutting down...")
		ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = server.Stop(ctxShutDown); err != nil {
			log.Error("Failed to shutdown server server", "error", err)
		}
	}()

	if err = server.Start(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
