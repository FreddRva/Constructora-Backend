package model

type Project struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Client    string  `json:"client"`
	Location  string  `json:"location"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	Status    string  `json:"status"`
	Progress  float64 `json:"progress"`
	ImageURL  string  `json:"image_url"`
}
