package main

import (
	"fmt"
	"net/http"

	"github.com/js-codegamer/fs-sync/config"
	"github.com/js-codegamer/fs-sync/internal/database"
	"github.com/js-codegamer/fs-sync/internal/routes"
	"github.com/js-codegamer/fs-sync/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Config load error: %v", err))
	}

	logger.InitLogger()
	defer logger.Close()

	db := database.InitDatabase(cfg.Database.Path)
	defer db.Close()

	router := http.NewServeMux()

	router.HandleFunc("/", routes.SetupRoutes(cfg))

	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Sugar.Infof("Starting server on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		logger.Sugar.Fatalf("Server start failed: %v", err)
	}
}
