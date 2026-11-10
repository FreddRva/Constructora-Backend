package service

import (
	"staff-service/internal/factory"
	"staff-service/internal/model"
	"staff-service/internal/repository"
)

type WorkerService struct {
	repo    *repository.WorkerRepository
	factory *factory.WorkerFactory
}

func NewWorkerService(repo *repository.WorkerRepository, f *factory.WorkerFactory) *WorkerService {
	return &WorkerService{
		repo:    repo,
		factory: f,
	}
}

func (s *WorkerService) Create(input model.Worker) (*model.Worker, error) {
	worker, err := s.factory.CreateWorker(input)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(worker); err != nil {
		return nil, err
	}
	return worker, nil
}

func (s *WorkerService) GetAll() ([]model.Worker, error) {
	return s.repo.GetAll()
}

func (s *WorkerService) GetByID(id int) (*model.Worker, error) {
	return s.repo.GetByID(id)
}

func (s *WorkerService) Update(id int, input model.Worker) (*model.Worker, error) {
	worker, err := s.factory.CreateWorker(input)
	if err != nil {
		return nil, err
	}
	worker.ID = id
	if err := s.repo.Update(id, worker); err != nil {
		return nil, err
	}
	return worker, nil
}

func (s *WorkerService) Delete(id int) error {
	return s.repo.Delete(id)
}
