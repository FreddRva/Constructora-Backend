package factory

import (
	"budget-service/internal/model"
	"errors"
	"strings"
)

type BudgetFactory struct{}

func NewBudgetFactory() *BudgetFactory {
	return &BudgetFactory{}
}

func (f *BudgetFactory) CreateBudget(input model.Budget) (*model.Budget, error) {
	if input.ProjectID <= 0 {
		return nil, errors.New("debe seleccionar un proyecto válido (project_id > 0)")
	}
	if input.MaterialCost < 0 || input.LaborCost < 0 || input.EquipmentCost < 0 {
		return nil, errors.New("los costos no pueden ser negativos")
	}

	total := CalculateTotalCost(input.MaterialCost, input.LaborCost, input.EquipmentCost)

	return &model.Budget{
		ProjectID:     input.ProjectID,
		MaterialCost:  input.MaterialCost,
		LaborCost:     input.LaborCost,
		EquipmentCost: input.EquipmentCost,
		TotalCost:     total,
		Notes:         strings.TrimSpace(input.Notes),
		DocumentURL:   strings.TrimSpace(input.DocumentURL),
	}, nil
}

func CalculateTotalCost(material, labor, equipment float64) float64 {
	return material + labor + equipment
}
