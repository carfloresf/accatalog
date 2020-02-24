package service

import (
	"log"

	"github.com/hellerox/AcCatalog/model"
	"github.com/hellerox/AcCatalog/storage"
)

// NewService returns a new DatabaseOperator
func NewService(storage storage.Storage) Service {
	service := Service{
		Storage: storage,
	}

	return service
}

// Service controller
type Service struct {
	Storage storage.Storage
}

// GetFullCostume full costume
func (s *Service) GetFullCostume(cID int) (c *model.Costume, err error) {
	costume, err := s.Storage.GetCostume(cID)
	if err != nil {
		log.Printf("error getting costume: %s", err.Error())
		return nil, err
	}

	cms, err := s.Storage.GetCostumeMaterial(cID)
	if err != nil {
		return costume, err
	}

	var cmrs []model.CostumeMaterialResponse

	for _, material := range cms {
		var cmr model.CostumeMaterialResponse

		cms, err := s.Storage.GetMaterial(material.MaterialID)
		if err != nil {
			return costume, err
		}

		cmr.Material = cms
		cmr.Quantity = material.Quantity
		cmrs = append(cmrs, cmr)
	}

	costume.CostumeMaterial = cmrs

	return costume, err
}

// GetAllCostumes all the costumes
func (s *Service) GetAllCostumes() (cs []model.Costume, err error) {
	costumes, err := s.Storage.GetAllCostumes()
	if err != nil {
		log.Printf("error getting costume: %s", err.Error())
		return cs, err
	}

	for i, costume := range costumes {
		fullCm, _ := s.GetFullCostume(costume.CostumeID)
		costumes[i] = *fullCm
	}

	return costumes, err
}

// GetAllCostumes all the costumes
func (s *Service) GetAllMaterials() (ms []model.Material, err error) {
	materials, err := s.Storage.GetAllMaterials()
	if err != nil {
		log.Printf("error getting costume: %s", err.Error())
		return materials, err
	}

	return materials, err
}
