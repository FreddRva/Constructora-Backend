package service

import (
	"budget-service/internal/factory"
	"budget-service/internal/model"
	"budget-service/internal/repository"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type BudgetService struct {
	repo              *repository.BudgetRepository
	factory           *factory.BudgetFactory
	projectServiceURL string
}

func NewBudgetService(repo *repository.BudgetRepository, f *factory.BudgetFactory, projectServiceURL string) *BudgetService {
	return &BudgetService{
		repo:              repo,
		factory:           f,
		projectServiceURL: projectServiceURL,
	}
}

func (s *BudgetService) Create(input model.Budget) (*model.Budget, error) {
	budget, err := s.factory.CreateBudget(input)
	if err != nil {
		return nil, err
	}

	if err := s.verifyProjectExists(budget.ProjectID); err != nil {
		return nil, err
	}

	if err := s.repo.Create(budget); err != nil {
		return nil, err
	}
	return budget, nil
}

func (s *BudgetService) GetAll() ([]model.Budget, error) {
	return s.repo.GetAll()
}

func (s *BudgetService) GetByID(id int) (*model.Budget, error) {
	return s.repo.GetByID(id)
}

func (s *BudgetService) Update(id int, input model.Budget) (*model.Budget, error) {
	budget, err := s.factory.CreateBudget(input)
	if err != nil {
		return nil, err
	}
	budget.ID = id
	if err := s.repo.Update(id, budget); err != nil {
		return nil, err
	}
	return budget, nil
}

func (s *BudgetService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *BudgetService) verifyProjectExists(projectID int) error {
	client := &http.Client{Timeout: 2 * time.Second}
	url := fmt.Sprintf("%s/projects/%d", s.projectServiceURL, projectID)

	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("el proyecto asociado no existe en el servicio de proyectos")
	}
	return nil
}
