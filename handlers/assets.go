package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

var assetsPath = "../yukon/assets"

// Media subdirectories to search for assets
var mediaSubdirectories = []string{
	"games", "rooms", "interface", "artifacts", "clothing",
	"crumbs", "flash", "furniture", "igloos", "mainmenu",
	"misc", "music", "penguin", "postcards", "preload",
	"puffles", "shared", "sounds",
}

type AssetInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

type AssetLocation struct {
	Found     bool   `json:"found"`
	Type      string `json:"type"` // "pack" or "atlas"
	Path      string `json:"path"` // Relative path from assets root
	Directory string `json:"directory,omitempty"`
}

// checkAssetFile checks if a file exists and returns its location info
func checkAssetFile(fullPath string, absAssetsPath string, isAtlas bool) (AssetLocation, bool) {
	// Check if file exists
	if _, err := os.Stat(fullPath); err != nil {
		return AssetLocation{}, false
	}

	// Convert to web path - extract path after "assets/"
	relativePath := strings.TrimPrefix(fullPath, absAssetsPath)
	relativePath = filepath.ToSlash(relativePath)
	relativePath = "/assets" + relativePath

	location := AssetLocation{
		Found: true,
		Type:  "pack",
		Path:  relativePath,
	}

	if isAtlas {
		location.Type = "atlas"
		// For atlases, also return the directory path
		dirPath := filepath.Dir(fullPath)
		relativeDirPath := strings.TrimPrefix(dirPath, absAssetsPath)
		relativeDirPath = filepath.ToSlash(relativeDirPath)
		location.Directory = "/assets" + relativeDirPath
	}

	return location, true
}

func GetAssets(w http.ResponseWriter, r *http.Request) {
	assets := []AssetInfo{}

	err := filepath.Walk(assetsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)
			// TODO: add more asset types, webm, gif, svg, etc.
			if ext == ".png" || ext == ".jpg" || ext == ".json" || ext == ".atlas" {
				relPath, _ := filepath.Rel(assetsPath, path)
				assets = append(assets, AssetInfo{
					Name: info.Name(),
					Path: relPath,
					Type: ext[1:],
					Size: info.Size(),
				})
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assets)
}

// ResolveAssetLocation finds the pack file or atlas for a given texture key
func ResolveAssetLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		http.Error(w, "Asset key is required", http.StatusBadRequest)
		return
	}

	// Convert assetsPath to absolute path for reliable path operations
	absAssetsPath, err := filepath.Abs(assetsPath)
	if err != nil {
		http.Error(w, "Failed to resolve assets path", http.StatusInternalServerError)
		return
	}

	mediaPath := filepath.Join(absAssetsPath, "media")

	// First try direct patterns in known subdirectories
	searchPatterns := []struct {
		pattern string
		isAtlas bool
	}{
		// Direct pack files
		{filepath.Join("{subdir}", key, key+"-pack.json"), false},
		// Nested game pack files
		{filepath.Join("{subdir}", "game", key, key+"-pack.json"), false},
		// Direct atlases
		{filepath.Join("{subdir}", key, key+".json"), true},
		// Nested game atlases
		{filepath.Join("{subdir}", "game", key, key+".json"), true},
	}

	// Search through subdirectories with known patterns
	for _, subdir := range mediaSubdirectories {
		for _, search := range searchPatterns {
			pattern := strings.Replace(search.pattern, "{subdir}", subdir, 1)
			fullPath := filepath.Join(mediaPath, pattern)

			if location, found := checkAssetFile(fullPath, absAssetsPath, search.isAtlas); found {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(location)
				return
			}
		}
	}

	// If not found with direct patterns, do a recursive search for pack files
	// Search for {key}-pack.json or {key}.json recursively
	packFileName := key + "-pack.json"
	atlasFileName := key + ".json"

	var foundPath string
	var isAtlas bool

	filepath.Walk(mediaPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || foundPath != "" {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		fileName := filepath.Base(path)
		if fileName == packFileName {
			foundPath = path
			isAtlas = false
			return filepath.SkipAll
		}

		// Check if it's an atlas (but not a pack file)
		if fileName == atlasFileName && !strings.HasSuffix(filepath.Dir(path), "/"+key) {
			// Make sure this is actually an atlas directory structure
			// (atlas files are in directories with the same name as the key)
			parentDir := filepath.Base(filepath.Dir(path))
			if parentDir == key {
				foundPath = path
				isAtlas = true
				return filepath.SkipAll
			}
		}

		return nil
	})

	if foundPath != "" {
		if location, found := checkAssetFile(foundPath, absAssetsPath, isAtlas); found {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(location)
			return
		}
	}

	// Not found
	location := AssetLocation{
		Found: false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)
}
