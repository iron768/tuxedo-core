package models

import "slices"

// Component property keys for GameObject
const (
	ComponentButton     = "Button"
	ComponentMoveTo     = "MoveTo"
	ComponentAnimation  = "Animation"
	ComponentSimpleButton = "SimpleButton"
)

// Helper methods for GameObject
func (g *GameObject) HasComponent(componentType string) bool {
	return slices.Contains(g.Components, componentType)
}

func (g *GameObject) SetVisible(visible bool) {
	if g.Properties == nil {
		g.Properties = make(map[string]interface{})
	}
	g.Properties["visible"] = visible
}

func (g *GameObject) SetOrigin(originX, originY float64) {
	if g.Properties == nil {
		g.Properties = make(map[string]interface{})
	}
	g.Properties["originX"] = originX
	g.Properties["originY"] = originY
}
