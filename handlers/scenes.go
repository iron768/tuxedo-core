package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"tuxedo-core/models"

	"github.com/gorilla/mux"
)

var projectPath = "../yukon/src/scenes/"

func GetScenes(w http.ResponseWriter, r *http.Request) {
    scenes := []string{}
    
    err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
        if filepath.Ext(path) == ".scene" {
            // Get relative path from project root
            relPath, _ := filepath.Rel(projectPath, path)
            // Remove .scene extension and normalize to forward slashes
            sceneName := filepath.ToSlash(relPath[:len(relPath)-6])
            scenes = append(scenes, sceneName)
        }
        return nil
    })
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(scenes)
}

func GetScene(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    
    scenePath := filepath.Join(projectPath, name+".scene")
    
    // Check if file exists
    if _, err := os.Stat(scenePath); os.IsNotExist(err) {
        http.Error(w, "Scene not found", http.StatusNotFound)
        return
    }
    
    data, err := os.ReadFile(scenePath)
    if err != nil {
        http.Error(w, "Scene not found", http.StatusNotFound)
        return
    }
    
    var scene models.Scene
    if err := json.Unmarshal(data, &scene); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(scene)
}

func UpdateScene(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    var scene models.Scene
    if err := json.Unmarshal(body, &scene); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    scenePath := filepath.Join(projectPath, name+".scene")
    
    // Pretty print JSON
    prettyJSON, _ := json.MarshalIndent(scene, "", "    ")
    
    if err := os.WriteFile(scenePath, prettyJSON, 0644); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func CreateScene(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    var scene models.Scene
    if err := json.Unmarshal(body, &scene); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if scene.Settings.SceneKey == "" {
        http.Error(w, "Scene key is required", http.StatusBadRequest)
        return
    }
    
    scenePath := filepath.Join(projectPath, scene.Settings.SceneKey+".scene")
    
    // Check if scene already exists
    if _, err := os.Stat(scenePath); err == nil {
        http.Error(w, "Scene already exists", http.StatusConflict)
        return
    }
    
    // Pretty print JSON
    prettyJSON, _ := json.MarshalIndent(scene, "", "    ")
    
    if err := os.WriteFile(scenePath, prettyJSON, 0644); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "created", "path": scenePath})
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    // WebSocket implementation for hot reload
    // This will be implemented when needed for file watching
    w.WriteHeader(http.StatusNotImplemented)
    json.NewEncoder(w).Encode(map[string]string{"status": "not implemented yet"})
}