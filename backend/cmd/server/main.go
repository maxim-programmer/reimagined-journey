package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maxim-programmer/reimagined-journey/backend/internal/config"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/elastic"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/handler"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/middleware"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/repository"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := repository.RunMigrations(context.Background(), db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	esClient := elastic.NewClient(cfg.ElasticsearchURL)
	if err := esClient.EnsureIndex(context.Background()); err != nil {
		log.Fatalf("failed to ensure elasticsearch index: %v", err)
	}

	docRepo := repository.NewDocumentRepository(db)
	chunkRepo := repository.NewChunkRepository(db)
	docSvc := service.NewDocumentService(docRepo, chunkRepo, esClient)
	docHandler := handler.NewDocumentHandler(docSvc, cfg.UploadDir)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/documents/upload", docHandler.Upload)
	mux.HandleFunc("GET /api/v1/documents", docHandler.List)
	mux.HandleFunc("GET /api/v1/search", docHandler.Search)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      middleware.CORS(mux),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
	log.Println("server stopped")
}