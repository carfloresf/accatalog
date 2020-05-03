package service

import (
	"github.com/hellerox/AcCatalog/model"
	"github.com/hellerox/AcCatalog/storage"
	log "github.com/sirupsen/logrus"
)

// AcCatalogService controller
type AcCatalogService struct {
	Storage storage.Storage
}

// GetFullCostume full costume
func (s *AcCatalogService) GetFullCostume(cID int) (c *model.Costume, err error) {
	costume, err := s.Storage.GetCostume(cID)
	if err != nil {
		log.Errorf("error getting costume: %+v", err)
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
func (s *AcCatalogService) GetAllCostumes() (cs []model.Costume, err error) {
	costumes, err := s.Storage.GetAllCostumes()
	if err != nil {
		log.Errorf("error getting costume: %+v", err)
		return cs, err
	}

	for i, costume := range costumes {
		fullCm, _ := s.GetFullCostume(costume.CostumeID)
		costumes[i] = *fullCm
	}

	return costumes, err
}

// GetAllMaterials all the materials
func (s *AcCatalogService) GetAllMaterials() (ms []model.Material, err error) {
	materials, err := s.Storage.GetAllMaterials()
	if err != nil {
		log.Errorf("error getting costume: %+v", err)
		return materials, err
	}

	return materials, err
}

// CreateMaterial create material
func (s *AcCatalogService) CreateMaterial(m model.Material) (mID int64, err error) {
	mID, err = s.Storage.InsertMaterial(m)
	if err != nil {
		log.Errorf("error creating material: %+v", err)
		return mID, err
	}

	return mID, err
}

// CreateCostume create costume
func (s *AcCatalogService) CreateCostume(c model.Costume) (cID int, err error) {
	cID, err = s.Storage.InsertCostume(c)
	if err != nil {
		log.Errorf("error creating costume: %+v", err)
		return
	}

	return cID, err
}
