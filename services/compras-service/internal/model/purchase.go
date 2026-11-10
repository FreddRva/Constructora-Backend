package model

type Purchase struct {
	ID           int     `json:"id"`
	SupplierName string  `json:"supplier_name"`
	RUC          string  `json:"ruc"`
	ItemName     string  `json:"item_name"`
	Quantity     float64 `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	TotalAmount  float64 `json:"total_amount"`
	ProjectID    int     `json:"project_id"`
	Status       string  `json:"status"`
	InvoiceURL   string  `json:"invoice_url"`
}
