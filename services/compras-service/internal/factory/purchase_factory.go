package factory

import (
	"errors"
	"procurement-service/internal/model"
	"strings"
)

type PurchaseFactory struct{}

func NewPurchaseFactory() *PurchaseFactory {
	return &PurchaseFactory{}
}

func (f *PurchaseFactory) CreatePurchase(input model.Purchase) (*model.Purchase, error) {
	if strings.TrimSpace(input.SupplierName) == "" {
		return nil, errors.New("el nombre del proveedor es obligatorio")
	}
	if strings.TrimSpace(input.ItemName) == "" {
		return nil, errors.New("el material o recurso comprado es obligatorio")
	}
	if input.Quantity <= 0 || input.UnitPrice < 0 {
		return nil, errors.New("la cantidad debe ser mayor a 0 y el precio no negativo")
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "Pendiente"
	}

	return &model.Purchase{
		SupplierName: strings.TrimSpace(input.SupplierName),
		RUC:          strings.TrimSpace(input.RUC),
		ItemName:     strings.TrimSpace(input.ItemName),
		Quantity:     input.Quantity,
		UnitPrice:    input.UnitPrice,
		TotalAmount:  input.Quantity * input.UnitPrice,
		ProjectID:    input.ProjectID,
		Status:       status,
		InvoiceURL:   strings.TrimSpace(input.InvoiceURL),
	}, nil
}
