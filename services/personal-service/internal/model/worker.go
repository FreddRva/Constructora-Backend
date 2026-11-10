package model

type Worker struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Role      string  `json:"role"`
	Phone     string  `json:"phone"`
	DailyWage float64 `json:"daily_wage"`
	ProjectID int     `json:"project_id"`
	Status    string  `json:"status"`
}
