package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maxim-programmer/reimagined-journey/backend/internal/cache"
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

	redisCache := cache.NewRedisCache(cfg.RedisAddr)
	if err := redisCache.Ping(context.Background()); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Printf("connected to redis at %s", cfg.RedisAddr)

	docRepo := repository.NewDocumentRepository(db)
	chunkRepo := repository.NewChunkRepository(db)
	docSvc := service.NewDocumentService(docRepo, chunkRepo, esClient, redisCache)
	docHandler := handler.NewDocumentHandler(docSvc, cfg.UploadDir)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/documents/upload", docHandler.Upload)
	mux.HandleFunc("GET /api/v1/documents", docHandler.List)
	mux.HandleFunc("GET /api/v1/documents/{id}", docHandler.Get)
	mux.HandleFunc("GET /api/v1/search", docHandler.Search)

	mux.HandleFunc("GET /api/openapi.yaml", serveOpenAPISpec)
	mux.HandleFunc("GET /docs", serveSwaggerUI)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      middleware.CORS(mux),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", cfg.Addr)
		log.Printf("swagger UI available at http://localhost%s/docs", cfg.Addr)
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

func serveOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	http.ServeFile(w, r, "api/openapi.yaml")
}

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <title>Knowledge Base API — Swagger UI</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.17.14/swagger-ui.min.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.17.14/swagger-ui-bundle.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.17.14/swagger-ui-standalone-preset.min.js"></script>
  <script>
    window.onload = function() {
      SwaggerUIBundle({
        url: "/api/openapi.yaml",
        dom_id: "#swagger-ui",
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
        layout: "StandaloneLayout",
        deepLinking: true,
      });
    };
  </script>
</body>
</html>`))
}