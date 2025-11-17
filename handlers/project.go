package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type ProjectInfo struct {
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	SceneCount int      `json:"sceneCount"`
	Folders    []string `json:"folders"`
}

func GetProjectInfo(w http.ResponseWriter, r *http.Request) {
	info := ProjectInfo{
		Name:    "Club Penguin",
		Path:    projectPath,
		Folders: []string{},
	}

	// Count scenes
	sceneCount := 0
	err := filepath.Walk(projectPath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fileInfo.IsDir() && path != projectPath {
			relPath, _ := filepath.Rel(projectPath, path)
			info.Folders = append(info.Folders, relPath)
		}

		if filepath.Ext(path) == ".scene" {
			sceneCount++
		}
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info.SceneCount = sceneCount

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
