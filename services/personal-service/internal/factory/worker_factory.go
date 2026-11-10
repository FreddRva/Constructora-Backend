package factory

import (
	"errors"
	"staff-service/internal/model"
	"strings"
)

type WorkerFactory struct{}

func NewWorkerFactory() *WorkerFactory {
	return &WorkerFactory{}
}

func (f *WorkerFactory) CreateWorker(input model.Worker) (*model.Worker, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("el nombre del trabajador es obligatorio")
	}
	if strings.TrimSpace(input.Role) == "" {
		return nil, errors.New("el cargo o especialidad es obligatorio")
	}
	if input.DailyWage < 0 {
		return nil, errors.New("el jornal diario no puede ser negativo")
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "Activo"
	}

	return &model.Worker{
		Name:      strings.TrimSpace(input.Name),
		Role:      strings.TrimSpace(input.Role),
		Phone:     strings.TrimSpace(input.Phone),
		DailyWage: input.DailyWage,
		ProjectID: input.ProjectID,
		Status:    status,
	}, nil
}
