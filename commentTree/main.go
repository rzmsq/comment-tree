package main

import (
	"CommentTree/commentTree/adapters/db"
	"CommentTree/commentTree/adapters/rest"
	"CommentTree/commentTree/config"
	"CommentTree/commentTree/pkg/handler"
	httpserver "CommentTree/commentTree/pkg/http_server"
	"CommentTree/commentTree/pkg/logger"
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	log.Debug("DB init success", "db", connDB)

	h := handler.New()
	h.AddHandler("POST /comments", rest.NewCreateHandler(log))
	h.AddHandler("GET /comments", rest.NewGetHandler(log))
	h.AddHandler("DELETE /comments/{id}", rest.NewDeleteHandler(log))
	log.Debug("Handler init success", "handler", h)

	server := httpserver.New(h, cfg.ServerConfig)
	log.Debug("Server init success", "server", server)

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
