package factory

import (
	"errors"
	"project-service/internal/model"
	"strings"
)

type ProjectFactory struct{}

func NewProjectFactory() *ProjectFactory {
	return &ProjectFactory{}
}

func (f *ProjectFactory) CreateProject(input model.Project) (*model.Project, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("el nombre del proyecto es obligatorio")
	}
	if strings.TrimSpace(input.Client) == "" {
		return nil, errors.New("el cliente es obligatorio")
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "Planificación"
	}

	progress := input.Progress
	if progress < 0 {
		progress = 0
	} else if progress > 100 {
		progress = 100
	}

	return &model.Project{
		Name:      strings.TrimSpace(input.Name),
		Client:    strings.TrimSpace(input.Client),
		Location:  strings.TrimSpace(input.Location),
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		Status:    status,
		Progress:  progress,
		ImageURL:  strings.TrimSpace(input.ImageURL),
	}, nil
}
