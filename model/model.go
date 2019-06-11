package model

import "time"

// Material contains info about the materials used in AC
type Material struct {
	MaterialID   int          `json:"materialID,omitempty"`
	Description  string       `json:"description,omitempty"`
	Cost         float64      `json:"cost,omitempty"`
	Measure      string       `json:"measure,omitempty"`
	BrandID      int          `json:"brandID,omitempty"`
	CreatedAt    time.Time    `json:"createdAt,omitempty"`
	Active       bool         `json:"active,omitempty"`
	MaterialType MaterialType `json:"materialType,omitempty"`
}

// Costume contains info about the costumes made in AC
type Costume struct {
	CostumeID       int                       `json:"costumeID,omitempty"`
	Name            string                    `json:"name,omitempty"`
	Color           string                    `json:"color,omitempty"`
	ActualCost      float64                   `json:"actualCost,omitempty"`
	CostumeCode     string                    `json:"costumeCode,omitempty"`
	Genre           string                    `json:"genre,omitempty"`
	CreatedAt       time.Time                 `json:"createdAt,omitempty"`
	CostumeCategory Category                  `json:"costumeCategory,omitempty"`
	CostumeMaterial []CostumeMaterialResponse `json:"materials,omitempty"`
}

// CostumeMaterialRelation struct
type CostumeMaterialRelation struct {
	CostumeID  int `json:"costumeID,omitempty"`
	MaterialID int `json:"material,omitempty"`
	Quantity   int `json:"quantity,omitempty"`
}

// CostumeMaterialResponse struct
type CostumeMaterialResponse struct {
	Material Material `json:"materialID,omitempty"`
	Quantity int      `json:"quantity,omitempty"`
}

// MaterialType struct
type MaterialType struct {
	MaterialTypeID int    `json:"materialTypeID,omitempty"`
	Name           string `json:"name,omitempty"`
}

// Category category
type Category struct {
	CategoryID int       `json:"categoryID,omitempty"`
	Name       string    `json:"name,omitempty"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
}
