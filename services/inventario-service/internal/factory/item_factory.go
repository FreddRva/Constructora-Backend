package factory

import (
	"errors"
	"inventory-service/internal/model"
	"strings"
)

type ItemCreator interface {
	Create(input model.Item) (*model.Item, error)
}

type MaterialCreator struct{}

func (c *MaterialCreator) Create(input model.Item) (*model.Item, error) {
	unit := strings.TrimSpace(input.Unit)
	if unit == "" {
		unit = "Unidades"
	}
	return &model.Item{
		Name:      strings.TrimSpace(input.Name),
		Category:  model.CategoryMaterial,
		Unit:      unit,
		Stock:     input.Stock,
		Condition: "Consumible",
		ProjectID: input.ProjectID,
		ImageURL:  strings.TrimSpace(input.ImageURL),
	}, nil
}

type ToolCreator struct{}

func (c *ToolCreator) Create(input model.Item) (*model.Item, error) {
	condition := strings.TrimSpace(input.Condition)
	if condition == "" {
		condition = "Operativo"
	}
	return &model.Item{
		Name:      strings.TrimSpace(input.Name),
		Category:  model.CategoryTool,
		Unit:      "Unidad",
		Stock:     input.Stock,
		Condition: condition,
		ProjectID: input.ProjectID,
		ImageURL:  strings.TrimSpace(input.ImageURL),
	}, nil
}

func CreateInventoryItem(input model.Item) (*model.Item, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("el nombre del recurso es obligatorio")
	}
	if input.Stock < 0 {
		return nil, errors.New("el stock no puede ser negativo")
	}

	var creator ItemCreator
	switch strings.ToUpper(strings.TrimSpace(input.Category)) {
	case model.CategoryMaterial:
		creator = &MaterialCreator{}
	case model.CategoryTool:
		creator = &ToolCreator{}
	default:
		return nil, errors.New("categoría inválida: use MATERIAL o HERRAMIENTA")
	}

	return creator.Create(input)
}
