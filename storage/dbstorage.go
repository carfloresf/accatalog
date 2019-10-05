package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/hellerox/AcCatalog/model"
)

// Connect function to start db connection
func connect(connectionString string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

// InsertCostume for ACCAT
func (ds *DatabaseStorage) InsertCostume(c model.Costume) (int, error) {
	err := ds.db.QueryRow(
		`INSERT INTO costume (name, color, costume_code, genre, costume_category_id) 
			VALUES ($1,$2,$3,$4,$5)
			RETURNING costume_id`,
		c.Name, c.Color, c.CostumeCode, c.Genre, c.CostumeCategory.CategoryID).Scan(&c.CostumeID)
	if err != nil {
		log.Errorf("error while inserting costumes: %s", err.Error())
		return 0, err
	}
	return c.CostumeID, nil
}

// InsertMaterial for ACCAT
func (ds *DatabaseStorage) InsertMaterial(m model.Material) (int, error) {
	err := ds.db.QueryRow(
		`INSERT INTO material (description, cost, measure, material_type_id, brand_id)
			SELECT CAST($1 AS VARCHAR), $2, $3, $4, $5
			WHERE NOT EXISTS(SELECT material_id FROM material WHERE description = $1) RETURNING material_id`,
		m.Description, m.Cost, m.Measure, m.MaterialType.MaterialTypeID, m.BrandID).Scan(&m.MaterialID)
	fmt.Printf("%+v", m)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, fmt.Errorf("didn't insert data, it already exists")
		}
		log.Errorf("error while inserting material: %s", err.Error())
		return 0, err
	}
	return m.MaterialID, nil
}

// InsertCostumeMaterialRelation for ACCAT
func (ds *DatabaseStorage) InsertCostumeMaterialRelation(cm model.CostumeMaterialRelation) error {
	var rows int
	err := ds.db.QueryRow(
		`INSERT INTO costume_material_relation(costume_id, material_id, quantity) 
		VALUES ($1,$2,$3) ON CONFLICT ON CONSTRAINT costume_material_relation_pk DO UPDATE SET quantity = $3
		RETURNING 1`,
		&cm.CostumeID, &cm.MaterialID, &cm.Quantity).Scan(&rows)
	if err != nil {
		log.Errorf("error while inserting costume material: %s", err.Error())
		return fmt.Errorf("error inserting Costume Material")
	}
	return nil
}

// InsertMaterialType for ACCAT
func (ds *DatabaseStorage) InsertMaterialType(m model.MaterialType) (int, error) {
	err := ds.db.QueryRow(
		`INSERT INTO material_type (name) 
		SELECT CAST($1 AS VARCHAR) WHERE NOT EXISTS(
			SELECT material_type_id FROM material_type 
			WHERE name = $1) 
			RETURNING material_type_id`,
		m.Name).Scan(&m.MaterialTypeID)
	if err != nil {
		log.Errorf("error while inserting material type: %s", err.Error())
		return 0, err
	}
	return m.MaterialTypeID, nil
}

// GetCostume for ACCAT
func (ds *DatabaseStorage) GetCostume(cID int) (c model.Costume, err error) {
	var cc model.Category
	if err = ds.db.QueryRow(
		"SELECT c.costume_id, c.name, c.color, c.costume_code, c.genre, c.created_at, c.costume_category_id, cc.costume_category_name,cc.created_at FROM Costume c JOIN costume_category cc USING (costume_category_id) where costume_id=$1", cID).
		Scan(&c.CostumeID, &c.Name, &c.Color, &c.CostumeCode, &c.Genre, &c.CreatedAt, &cc.CategoryID, &cc.Name, &cc.CreatedAt); err != nil {
		log.Errorf("error while gettings costumes: %s", err.Error())
		return c, fmt.Errorf("error while getting costumes ")
	}
	c.CostumeCategory = cc
	log.Infof("costume found %+v", c)
	return c, err
}

// GetAllCostumes for ACCAT
func (ds *DatabaseStorage) GetAllCostumes() (cs []model.Costume, err error) {
	rows, err := ds.db.Query(
		"SELECT c.costume_id FROM Costume c")
	if err != nil {
		log.Errorf("error while gettings costumes: %s", err.Error())
		return cs, fmt.Errorf("error while getting costumes ")
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Warnf("The closing of the rows failed, the error is %s", err)
		}
	}()

	for rows.Next() {
		var c model.Costume
		err = rows.Scan(
			&c.CostumeID)
		if err != nil {
			return cs, err
		}
		cs = append(cs, c)
	}
	return cs, err
}

// GetCostumeMaterial for ACCAT
func (ds *DatabaseStorage) GetCostumeMaterial(cID int) (cm []model.CostumeMaterialRelation, err error) {
	rows, errq := ds.db.Query(
		"SELECT costume_id, material_id, quantity FROM costume_material_relation where costume_id = $1", cID)
	if errq != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var costumeMaterial model.CostumeMaterialRelation
		err := rows.Scan(&costumeMaterial.CostumeID, &costumeMaterial.MaterialID, &costumeMaterial.Quantity)
		if err != nil {
			log.Fatal(err)
		}
		cm = append(cm, costumeMaterial)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	log.Infof("costume found %+v", cm)
	return cm, err
}

// GetPermission for ACCAT
func (ds *DatabaseStorage) GetPermission(apikey string) (permission string, user string, err error) {
	if err = ds.db.QueryRow(
		`SELECT "user",permission FROM api_key WHERE api_key = $1`, apikey).Scan(&user, &permission); err != nil {
		log.Errorf("error getting permissions: %s", err.Error())
		return permission, user, fmt.Errorf("error while getting costumes ")
	}
	log.Debugf("permission for %s found %s", apikey, permission)
	return permission, user, err
}

// GetMaterial for ACCAT
func (ds *DatabaseStorage) GetMaterial(mID int) (m model.Material, err error) {
	var mt model.MaterialType
	if err = ds.db.QueryRow(
		`SELECT m.material_id, m.description, m.cost, m.material_type_id, m.brand_id, m.created_at, m.active, mt.name 
		FROM material m JOIN material_type mt USING (material_type_id) 
		WHERE m.material_id = $1`, mID).
		Scan(&m.MaterialID, &m.Description, &m.Cost, &mt.MaterialTypeID, &m.BrandID, &m.CreatedAt, &m.Active, &mt.Name); err != nil {
		log.Errorf("error while gettings costumes: %s", err.Error())
		return m, fmt.Errorf("error while getting costumes ")
	}
	m.MaterialType = mt
	log.Infof("material found %+v", m)
	return m, err
}
