package service

import (
	"procurement-service/internal/factory"
	"procurement-service/internal/model"
	"procurement-service/internal/repository"
)

type PurchaseService struct {
	repo    *repository.PurchaseRepository
	factory *factory.PurchaseFactory
}

func NewPurchaseService(repo *repository.PurchaseRepository, f *factory.PurchaseFactory) *PurchaseService {
	return &PurchaseService{
		repo:    repo,
		factory: f,
	}
}

func (s *PurchaseService) Create(input model.Purchase) (*model.Purchase, error) {
	purchase, err := s.factory.CreatePurchase(input)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(purchase); err != nil {
		return nil, err
	}
	return purchase, nil
}

func (s *PurchaseService) GetAll() ([]model.Purchase, error) {
	return s.repo.GetAll()
}

func (s *PurchaseService) GetByID(id int) (*model.Purchase, error) {
	return s.repo.GetByID(id)
}

func (s *PurchaseService) Update(id int, input model.Purchase) (*model.Purchase, error) {
	purchase, err := s.factory.CreatePurchase(input)
	if err != nil {
		return nil, err
	}
	purchase.ID = id
	if err := s.repo.Update(id, purchase); err != nil {
		return nil, err
	}
	return purchase, nil
}

func (s *PurchaseService) Delete(id int) error {
	return s.repo.Delete(id)
}
