package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/merev/ds-player-api/internal/config"
	"github.com/merev/ds-player-api/internal/database"
	apphttp "github.com/merev/ds-player-api/internal/http"
	"github.com/merev/ds-player-api/internal/player"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPool(cfg.DBDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := database.Migrate(ctx, db); err != nil {
		cancel()
		log.Fatalf("migration failed: %v", err)
	}
	cancel()

	repo := player.NewRepository(db)
	handler := player.NewHandler(repo)
	router := apphttp.NewRouter(handler)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("player-api running on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down player-api...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
