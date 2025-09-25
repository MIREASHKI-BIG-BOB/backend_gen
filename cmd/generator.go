package main

import (
	"backend_gen/config"
	"backend_gen/internal/server"
	"flag"
	"log"
	"log/slog"
	"os"
)

func main() {
	cfgPath := flag.String("c", "config/config.yaml", "path to config file")
	flag.Parse()

	// Настройка логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting application", "config_path", *cfgPath)

	cfg, err := config.ReadConfig(*cfgPath)
	if err != nil {
		slog.Error("Failed to read config", "error", err, "config_path", *cfgPath)
		log.Fatal(err)
	}

	slog.Info("Config loaded successfully", "server_addr", cfg.Server.Addr, "server_port", cfg.Server.Port)

	s, err := server.New(cfg)
	if err != nil {
		slog.Error("Failed to create server", "error", err)
		log.Fatal(err)
	}

	slog.Info("Server created successfully")
	s.Run()
}
