package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"it-asset-management/backend/internal/config"
	"it-asset-management/backend/internal/handlers"
	"it-asset-management/backend/internal/middleware"
	"it-asset-management/backend/internal/repository"
)

func main() {
	log.Println("Starting IT Asset Management API Service...")

	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Connect to Database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Fatal: Database initialization failed: %v", err)
	}
	defer db.Close()

	// 3. Initialize repository and handlers
	assetRepo := repository.NewPostgresAssetRepository(db)
	assetHandler := handlers.NewAssetHandler(assetRepo)

	// 4. Setup HTTP routes
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"UP","timestamp":"` + time.Now().UTC().Format(time.RFC3339) + `"}`))
	})

	// Assets API routes
	mux.HandleFunc("GET /api/assets/stats", assetHandler.GetStats)
	mux.HandleFunc("GET /api/assets", assetHandler.GetAssets)
	mux.HandleFunc("POST /api/assets", assetHandler.CreateAsset)
	mux.HandleFunc("GET /api/assets/{id}", assetHandler.GetAssetByID)
	mux.HandleFunc("PUT /api/assets/{id}", assetHandler.UpdateAsset)
	mux.HandleFunc("DELETE /api/assets/{id}", assetHandler.DeleteAsset)

	// Asset lifecycle actions
	mux.HandleFunc("POST /api/assets/{id}/assign", assetHandler.AssignAsset)
	mux.HandleFunc("POST /api/assets/{id}/return", assetHandler.ReturnAsset)
	mux.HandleFunc("POST /api/assets/{id}/repair", assetHandler.RepairAsset)
	mux.HandleFunc("POST /api/assets/{id}/retire", assetHandler.RetireAsset)
	mux.HandleFunc("GET /api/assets/{id}/history", assetHandler.GetAssetHistory)

	// Wrap router with CORS middleware
	handlerWithCORS := middleware.CORS(cfg.CORSOrigin)(mux)

	// 5. Start HTTP server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      handlerWithCORS,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server is running and listening at http://localhost%s", addr)
		log.Printf(" Allowed CORS Origin: %s", cfg.CORSOrigin)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen failed: %v", err)
		}
	}()

	// 6. Graceful shutdown listening to interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly.")
}
