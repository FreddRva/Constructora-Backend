package repository

import (
	"database/sql"
	"errors"
	"inventory-service/internal/model"
)

type InventoryRepository struct {
	db        *sql.DB
	useMemory bool
	memData   []model.Item
	nextID    int
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{
		db:     db,
		nextID: 5,
		memData: []model.Item{
			{ID: 1, Name: "Cemento Portland Tipo I", Category: "MATERIAL", Unit: "Bolsas", Stock: 250, Condition: "Consumible", ProjectID: 1},
			{ID: 2, Name: "Fierro Corrugado 1/2\"", Category: "MATERIAL", Unit: "Varillas", Stock: 180, Condition: "Consumible", ProjectID: 1},
			{ID: 3, Name: "Mezcladora de Concreto 11p3", Category: "HERRAMIENTA", Unit: "Unidad", Stock: 3, Condition: "Operativo", ProjectID: 1},
			{ID: 4, Name: "Taladro Percutor Industrial", Category: "HERRAMIENTA", Unit: "Unidad", Stock: 5, Condition: "Operativo", ProjectID: 0},
		},
	}
}

func (r *InventoryRepository) InitSchema() error {
	if err := r.db.Ping(); err != nil {
		r.useMemory = true
		return err
	}
	query := `
	CREATE TABLE IF NOT EXISTS inventory_items (
		id SERIAL PRIMARY KEY,
		name VARCHAR(150) NOT NULL,
		category VARCHAR(50) NOT NULL,
		unit VARCHAR(50) NOT NULL,
		stock NUMERIC(10,2) DEFAULT 0,
		condition VARCHAR(50) NOT NULL,
		project_id INT DEFAULT 0,
		image_url TEXT
	);`
	if _, err := r.db.Exec(query); err != nil {
		r.useMemory = true
		return err
	}
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM inventory_items`).Scan(&count); err == nil && count == 0 {
		for _, it := range r.memData {
			_, _ = r.db.Exec(
				`INSERT INTO inventory_items (name, category, unit, stock, condition, project_id, image_url) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
				it.Name, it.Category, it.Unit, it.Stock, it.Condition, it.ProjectID, it.ImageURL,
			)
		}
	}
	return nil
}

func (r *InventoryRepository) Create(item *model.Item) error {
	if r.useMemory {
		item.ID = r.nextID
		r.nextID++
		r.memData = append([]model.Item{*item}, r.memData...)
		return nil
	}
	query := `
		INSERT INTO inventory_items (name, category, unit, stock, condition, project_id, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	return r.db.QueryRow(
		query, item.Name, item.Category, item.Unit, item.Stock, item.Condition, item.ProjectID, item.ImageURL,
	).Scan(&item.ID)
}

func (r *InventoryRepository) GetAll() ([]model.Item, error) {
	if r.useMemory {
		return r.memData, nil
	}
	rows, err := r.db.Query(`SELECT id, name, category, unit, stock, condition, project_id, image_url FROM inventory_items ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Item{}
	for rows.Next() {
		var it model.Item
		if err := rows.Scan(&it.ID, &it.Name, &it.Category, &it.Unit, &it.Stock, &it.Condition, &it.ProjectID, &it.ImageURL); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}

func (r *InventoryRepository) GetByID(id int) (*model.Item, error) {
	if r.useMemory {
		for _, it := range r.memData {
			if it.ID == id {
				return &it, nil
			}
		}
		return nil, errors.New("recurso no encontrado")
	}
	var it model.Item
	query := `SELECT id, name, category, unit, stock, condition, project_id, image_url FROM inventory_items WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&it.ID, &it.Name, &it.Category, &it.Unit, &it.Stock, &it.Condition, &it.ProjectID, &it.ImageURL)
	if err != nil {
		return nil, err
	}
	return &it, nil
}

func (r *InventoryRepository) Update(id int, item *model.Item) error {
	if r.useMemory {
		for i := range r.memData {
			if r.memData[i].ID == id {
				item.ID = id
				r.memData[i] = *item
				return nil
			}
		}
		return errors.New("recurso no encontrado")
	}
	query := `
		UPDATE inventory_items
		SET name=$1, category=$2, unit=$3, stock=$4, condition=$5, project_id=$6, image_url=$7
		WHERE id=$8`
	_, err := r.db.Exec(query, item.Name, item.Category, item.Unit, item.Stock, item.Condition, item.ProjectID, item.ImageURL, id)
	return err
}

func (r *InventoryRepository) Delete(id int) error {
	if r.useMemory {
		for i, it := range r.memData {
			if it.ID == id {
				r.memData = append(r.memData[:i], r.memData[i+1:]...)
				return nil
			}
		}
		return nil
	}
	_, err := r.db.Exec(`DELETE FROM inventory_items WHERE id = $1`, id)
	return err
}
