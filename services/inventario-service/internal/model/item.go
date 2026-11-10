package model

const (
	CategoryMaterial = "MATERIAL"
	CategoryTool     = "HERRAMIENTA"
)

type Item struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	Unit      string  `json:"unit"`
	Stock     float64 `json:"stock"`
	Condition string  `json:"condition"`
	ProjectID int     `json:"project_id"`
	ImageURL  string  `json:"image_url"`
}
