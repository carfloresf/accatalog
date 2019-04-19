package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/hellerox/AcCatalog/model"
)

// DatabaseStorage with config data
type DatabaseStorage struct {
	db *sqlx.DB
}

// Storage executes functions on storage resources
type Storage interface {
	InsertCostume(c model.Costume) (int, error)
	InsertMaterial(m model.Material) (int, error)
	InsertCostumeMaterialRelation(cm model.CostumeMaterialRelation) (err error)
	InsertMaterialType(cm model.MaterialType) (int, error)
	GetCostume(cID int) (c model.Costume, err error)
	GetAllCostumes() (cs []model.Costume, err error)
	GetMaterial(mID int) (c model.Material, err error)
	GetCostumeMaterial(cID int) (cm []model.CostumeMaterialRelation, err error)
	GetPermission(apikey string) (permission string, user string, err error)
}

// NewStorage returns a new DatabaseOperator
func NewStorage(connectionString string) Storage {
	storage := DatabaseStorage{
		db: connect(connectionString),
	}
	return &storage
}
