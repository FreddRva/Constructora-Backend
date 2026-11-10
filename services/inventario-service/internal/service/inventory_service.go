package service

import (
	"inventory-service/internal/factory"
	"inventory-service/internal/model"
	"inventory-service/internal/repository"
)

type InventoryService struct {
	repo *repository.InventoryRepository
}

func NewInventoryService(repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) Create(input model.Item) (*model.Item, error) {
	item, err := factory.CreateInventoryItem(input)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *InventoryService) CreateStarterKit(kitType string, projectID int) ([]model.Item, error) {
	kitFactory, err := factory.GetKitFactory(kitType)
	if err != nil {
		return nil, err
	}

	material := kitFactory.CreateBaseMaterial(projectID)
	tool := kitFactory.CreateBaseTool(projectID)

	if err := s.repo.Create(material); err != nil {
		return nil, err
	}
	if err := s.repo.Create(tool); err != nil {
		return nil, err
	}

	return []model.Item{*material, *tool}, nil
}

func (s *InventoryService) GetAll() ([]model.Item, error) {
	return s.repo.GetAll()
}

func (s *InventoryService) GetByID(id int) (*model.Item, error) {
	return s.repo.GetByID(id)
}

func (s *InventoryService) Update(id int, input model.Item) (*model.Item, error) {
	item, err := factory.CreateInventoryItem(input)
	if err != nil {
		return nil, err
	}
	item.ID = id
	if err := s.repo.Update(id, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *InventoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
