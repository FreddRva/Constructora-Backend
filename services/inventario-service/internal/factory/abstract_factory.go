package factory

import (
	"errors"
	"inventory-service/internal/model"
)

type ConstructoraKitFactory interface {
	CreateBaseMaterial(projectID int) *model.Item
	CreateBaseTool(projectID int) *model.Item
}

type ObraGruesaFactory struct{}

func (f *ObraGruesaFactory) CreateBaseMaterial(projectID int) *model.Item {
	return &model.Item{
		Name:      "Cemento Portland Tipo I",
		Category:  model.CategoryMaterial,
		Unit:      "Bolsas",
		Stock:     100,
		Condition: "Consumible",
		ProjectID: projectID,
	}
}

func (f *ObraGruesaFactory) CreateBaseTool(projectID int) *model.Item {
	return &model.Item{
		Name:      "Mezcladora de Concreto 11p3",
		Category:  model.CategoryTool,
		Unit:      "Unidad",
		Stock:     2,
		Condition: "Operativo",
		ProjectID: projectID,
	}
}

type AcabadosFactory struct{}

func (f *AcabadosFactory) CreateBaseMaterial(projectID int) *model.Item {
	return &model.Item{
		Name:      "Ladrillo Pastelero / Cerámico",
		Category:  model.CategoryMaterial,
		Unit:      "Millar",
		Stock:     10,
		Condition: "Consumible",
		ProjectID: projectID,
	}
}

func (f *AcabadosFactory) CreateBaseTool(projectID int) *model.Item {
	return &model.Item{
		Name:      "Taladro Percutor Industrial",
		Category:  model.CategoryTool,
		Unit:      "Unidad",
		Stock:     4,
		Condition: "Operativo",
		ProjectID: projectID,
	}
}

func GetKitFactory(kitType string) (ConstructoraKitFactory, error) {
	switch kitType {
	case "obra_gruesa":
		return &ObraGruesaFactory{}, nil
	case "acabados":
		return &AcabadosFactory{}, nil
	default:
		return nil, errors.New("tipo de kit desconocido: use 'obra_gruesa' o 'acabados'")
	}
}
