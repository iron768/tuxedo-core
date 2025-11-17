package models

type Scene struct {
	ID          string        `json:"id"`
	SceneType   string        `json:"sceneType"`
	Settings    SceneSettings `json:"settings"`
	DisplayList []GameObject  `json:"displayList"`
	Lists       []ObjectList  `json:"lists,omitempty"`
}

type SceneSettings struct {
	SceneKey     string   `json:"sceneKey"`
	BorderWidth  int      `json:"borderWidth"`
	BorderHeight int      `json:"borderHeight"`
	PreloadPacks []string `json:"preloadPackFiles"`
}

type GameObject struct {
	Type       string       `json:"type,omitempty"`
	ID         string       `json:"id"`
	Label      string       `json:"label"`
	X          float64      `json:"x,omitempty"`
	Y          float64      `json:"y,omitempty"`
	OriginX    *float64     `json:"originX,omitempty"`
	OriginY    *float64     `json:"originY,omitempty"`
	ScaleX     *float64     `json:"scaleX,omitempty"`
	ScaleY     *float64     `json:"scaleY,omitempty"`
	Angle      *float64     `json:"angle,omitempty"`
	Visible    *bool        `json:"visible,omitempty"`
	Width      *float64     `json:"width,omitempty"`
	Height     *float64     `json:"height,omitempty"`
	Texture    *Texture     `json:"texture,omitempty"`
	List       []GameObject `json:"list,omitempty"` // Container children
	Components []string     `json:"components,omitempty"`
	PrefabId   string       `json:"prefabId,omitempty"` // Reference to prefab definition
	Unlock     []string     `json:"unlock,omitempty"`   // Properties that override prefab

	// Text-specific properties
	Text            string   `json:"text,omitempty"`
	FontFamily      string   `json:"fontFamily,omitempty"`
	FontSize        string   `json:"fontSize,omitempty"`
	FontStyle       string   `json:"fontStyle,omitempty"`
	Color           string   `json:"color,omitempty"`
	Stroke          string   `json:"stroke,omitempty"`
	StrokeThickness *float64 `json:"strokeThickness,omitempty"`
	Align           string   `json:"align,omitempty"`
	PaddingLeft     *float64 `json:"paddingLeft,omitempty"`
	PaddingTop      *float64 `json:"paddingTop,omitempty"`
	PaddingRight    *float64 `json:"paddingRight,omitempty"`
	PaddingBottom   *float64 `json:"paddingBottom,omitempty"`

	Properties map[string]any `json:"-"` // For component properties
}

type Texture struct {
	Key   string `json:"key"`
	Frame string `json:"frame,omitempty"`
}

type ObjectList struct {
	ID        string   `json:"id"`
	Label     string   `json:"label"`
	ObjectIDs []string `json:"objectIds"`
}