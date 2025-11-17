package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"tuxedo-core/models"

	"github.com/gorilla/mux"
)

// GetPrefab retrieves a prefab scene by its ID
// Prefabs are stored in shared_prefabs directory and identified by the "id" field in the scene file
func GetPrefab(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prefabId := vars["id"]

	if prefabId == "" {
		http.Error(w, "Prefab ID is required", http.StatusBadRequest)
		return
	}

	// Search for the prefab file by ID
	prefabPath, err := findPrefabById(prefabId)
	if err != nil {
		http.Error(w, "Prefab not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Read the prefab file
	data, err := os.ReadFile(prefabPath)
	if err != nil {
		http.Error(w, "Error reading prefab file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse and validate it's a prefab
	var scene models.Scene
	if err := json.Unmarshal(data, &scene); err != nil {
		http.Error(w, "Error parsing prefab file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if scene.SceneType != "PREFAB" {
		http.Error(w, "Scene is not a prefab", http.StatusBadRequest)
		return
	}

	// Return the prefab scene
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// findPrefabById searches for a prefab file with the given ID
// It searches in all scene directories, not just shared_prefabs
// TODO: better logging and error handling for file access issues
// TODO: consider caching prefab paths for faster lookups
func findPrefabById(prefabId string) (string, error) {
	// Start from the root scenes directory to search everywhere
	scenesPath := projectPath
	var foundPath string

	err := filepath.Walk(scenesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a .scene file
		if !strings.HasSuffix(path, ".scene") {
			return nil
		}

		// Read and check if ID matches
		data, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}

		var scene models.Scene
		if err := json.Unmarshal(data, &scene); err != nil {
			return nil // Skip invalid JSON files
		}

		// Check if this is the prefab we're looking for
		if scene.ID == prefabId && scene.SceneType == "PREFAB" {
			foundPath = path
			return filepath.SkipAll // Stop walking once found
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", os.ErrNotExist
	}

	return foundPath, nil
}
