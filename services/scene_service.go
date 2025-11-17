package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"tuxedo-core/models"
)

type SceneService struct {
	projectPath string
}

func NewSceneService(projectPath string) *SceneService {
	return &SceneService{projectPath: projectPath}
}

func (s *SceneService) LoadScene(name string) (*models.Scene, error) {
	scenePath := filepath.Join(s.projectPath, name+".scene")

	data, err := os.ReadFile(scenePath)
	if err != nil {
		return nil, err
	}

	var scene models.Scene
	if err := json.Unmarshal(data, &scene); err != nil {
		return nil, err
	}

	return &scene, nil
}

func (s *SceneService) SaveScene(name string, scene *models.Scene) error {
	scenePath := filepath.Join(s.projectPath, name+".scene")

	prettyJSON, err := json.MarshalIndent(scene, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(scenePath, prettyJSON, 0644)
}

func (s *SceneService) ListScenes() ([]string, error) {
	scenes := []string{}

	err := filepath.Walk(s.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".scene" {
			relPath, _ := filepath.Rel(s.projectPath, path)
			scenes = append(scenes, relPath)
		}
		return nil
	})

	return scenes, err
}
