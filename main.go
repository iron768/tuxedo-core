package main

import (
	"fmt"
	"log"
	"net/http"

	"tuxedo-core/config"
	"tuxedo-core/handlers"
	"tuxedo-core/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Printf("Warning: Failed to load config, using defaults: %v", err)
		cfg, _ = config.Load("") // Get default config
	}

	log.Printf("Starting Tuxedo Core server...")
	log.Printf("Configuration:")
	log.Printf("  - Port: %s", cfg.Server.Port)
	log.Printf("  - Assets Path: %s", cfg.GetAssetsPath())
	log.Printf("  - Scenes Path: %s", cfg.GetScenesPath())

	r := mux.NewRouter()

	// Serve yukon assets FIRST (for loading textures in editor)
	// This must come before the catch-all static file handler
	assetsPath := cfg.GetAssetsPath()
	assetsFileServer := http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsPath)))
	r.PathPrefix("/assets/").Handler(middleware.CORS(assetsFileServer))
	log.Printf("Serving assets from: %s", assetsPath)

	// API routes with better pattern matching
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/scenes", handlers.GetScenes).Methods("GET")
	api.HandleFunc("/scenes/{name:.+}", handlers.GetScene).Methods("GET")
	api.HandleFunc("/scenes/{name:.+}", handlers.UpdateScene).Methods("PUT")
	api.HandleFunc("/scenes", handlers.CreateScene).Methods("POST")
	api.HandleFunc("/assets", handlers.GetAssets).Methods("GET")
	api.HandleFunc("/assets/resolve/{key}", handlers.ResolveAssetLocation).Methods("GET")
	api.HandleFunc("/project", handlers.GetProjectInfo).Methods("GET")
	api.HandleFunc("/prefab/{id}", handlers.GetPrefab).Methods("GET")

	// File watching endpoint for hot reload
	api.HandleFunc("/ws", handlers.WebSocketHandler)

	// Serve static files
	// Vite serves the frontend for now, uncomment later
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("../tuxedo/dist")))

	// Wrap with middleware
	handler := middleware.CORS(middleware.Logger(r))

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Tuxedo Core server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}