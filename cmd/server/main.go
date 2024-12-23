package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/js-codegamer/fs-sync/config"
	"github.com/js-codegamer/fs-sync/internal/database"
	"github.com/js-codegamer/fs-sync/internal/routes"
	"github.com/js-codegamer/fs-sync/pkg/logger"
	"github.com/js-codegamer/fs-sync/pkg/validator"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic(fmt.Sprintf("Config load error: %v", err))
	}
	cfg := config.GetConfig()

	logger.InitLogger(
		filepath.Join(cfg.Storage.BasePath, cfg.Storage.LogsPath, fmt.Sprintf("%s.log", time.Now().Format(time.RFC3339))),
	)
	defer logger.Close()

	validator.InitValidator()

	db := database.InitDatabase()
	defer db.Close()

	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Sugar.Infof("Starting server on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, routes.SetupRoutes()); err != nil {
		logger.Sugar.Fatalf("Server start failed: %v", err)
	}
}
