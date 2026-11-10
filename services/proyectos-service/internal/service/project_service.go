package service

import (
	"project-service/internal/factory"
	"project-service/internal/model"
	"project-service/internal/repository"
)

type ProjectService struct {
	repo    *repository.ProjectRepository
	factory *factory.ProjectFactory
}

func NewProjectService(repo *repository.ProjectRepository, f *factory.ProjectFactory) *ProjectService {
	return &ProjectService{
		repo:    repo,
		factory: f,
	}
}

func (s *ProjectService) Create(input model.Project) (*model.Project, error) {
	project, err := s.factory.CreateProject(input)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) GetAll() ([]model.Project, error) {
	return s.repo.GetAll()
}

func (s *ProjectService) GetByID(id int) (*model.Project, error) {
	return s.repo.GetByID(id)
}

func (s *ProjectService) Update(id int, input model.Project) (*model.Project, error) {
	project, err := s.factory.CreateProject(input)
	if err != nil {
		return nil, err
	}
	project.ID = id
	if err := s.repo.Update(id, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) Delete(id int) error {
	return s.repo.Delete(id)
}
