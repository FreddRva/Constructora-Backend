package model

type Budget struct {
	ID            int     `json:"id"`
	ProjectID     int     `json:"project_id"`
	MaterialCost  float64 `json:"material_cost"`
	LaborCost     float64 `json:"labor_cost"`
	EquipmentCost float64 `json:"equipment_cost"`
	TotalCost     float64 `json:"total_cost"`
	Notes         string  `json:"notes"`
	DocumentURL   string  `json:"document_url"`
}
