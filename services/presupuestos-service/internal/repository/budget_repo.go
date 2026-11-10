package repository

import (
	"budget-service/internal/model"
	"database/sql"
	"errors"
)

type BudgetRepository struct {
	db        *sql.DB
	useMemory bool
	memData   []model.Budget
	nextID    int
}

func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{
		db:     db,
		nextID: 3,
		memData: []model.Budget{
			{
				ID:            1,
				ProjectID:     1,
				MaterialCost:  85000.00,
				LaborCost:     62000.00,
				EquipmentCost: 18500.00,
				TotalCost:     165500.00,
				Notes:         "Presupuesto Casco Estructural y Cimentación",
			},
			{
				ID:            2,
				ProjectID:     2,
				MaterialCost:  120000.00,
				LaborCost:     95000.00,
				EquipmentCost: 34000.00,
				TotalCost:     249000.00,
				Notes:         "Presupuesto Base Nave Industrial",
			},
		},
	}
}

func (r *BudgetRepository) InitSchema() error {
	if err := r.db.Ping(); err != nil {
		r.useMemory = true
		return err
	}
	query := `
	CREATE TABLE IF NOT EXISTS budgets (
		id SERIAL PRIMARY KEY,
		project_id INT NOT NULL,
		material_cost NUMERIC(12,2) NOT NULL DEFAULT 0,
		labor_cost NUMERIC(12,2) NOT NULL DEFAULT 0,
		equipment_cost NUMERIC(12,2) NOT NULL DEFAULT 0,
		total_cost NUMERIC(12,2) NOT NULL DEFAULT 0,
		notes TEXT,
		document_url TEXT
	);`
	if _, err := r.db.Exec(query); err != nil {
		r.useMemory = true
		return err
	}
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM budgets`).Scan(&count); err == nil && count == 0 {
		for _, b := range r.memData {
			_, _ = r.db.Exec(
				`INSERT INTO budgets (project_id, material_cost, labor_cost, equipment_cost, total_cost, notes, document_url) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
				b.ProjectID, b.MaterialCost, b.LaborCost, b.EquipmentCost, b.TotalCost, b.Notes, b.DocumentURL,
			)
		}
	}
	return nil
}

func (r *BudgetRepository) Create(b *model.Budget) error {
	if r.useMemory {
		b.ID = r.nextID
		r.nextID++
		r.memData = append([]model.Budget{*b}, r.memData...)
		return nil
	}
	query := `
		INSERT INTO budgets (project_id, material_cost, labor_cost, equipment_cost, total_cost, notes, document_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	return r.db.QueryRow(
		query, b.ProjectID, b.MaterialCost, b.LaborCost, b.EquipmentCost, b.TotalCost, b.Notes, b.DocumentURL,
	).Scan(&b.ID)
}

func (r *BudgetRepository) GetAll() ([]model.Budget, error) {
	if r.useMemory {
		return r.memData, nil
	}
	rows, err := r.db.Query(`SELECT id, project_id, material_cost, labor_cost, equipment_cost, total_cost, notes, document_url FROM budgets ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	budgets := []model.Budget{}
	for rows.Next() {
		var b model.Budget
		if err := rows.Scan(&b.ID, &b.ProjectID, &b.MaterialCost, &b.LaborCost, &b.EquipmentCost, &b.TotalCost, &b.Notes, &b.DocumentURL); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, nil
}

func (r *BudgetRepository) GetByID(id int) (*model.Budget, error) {
	if r.useMemory {
		for _, b := range r.memData {
			if b.ID == id {
				return &b, nil
			}
		}
		return nil, errors.New("presupuesto no encontrado")
	}
	var b model.Budget
	query := `SELECT id, project_id, material_cost, labor_cost, equipment_cost, total_cost, notes, document_url FROM budgets WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&b.ID, &b.ProjectID, &b.MaterialCost, &b.LaborCost, &b.EquipmentCost, &b.TotalCost, &b.Notes, &b.DocumentURL)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BudgetRepository) Update(id int, b *model.Budget) error {
	if r.useMemory {
		for i := range r.memData {
			if r.memData[i].ID == id {
				b.ID = id
				r.memData[i] = *b
				return nil
			}
		}
		return errors.New("presupuesto no encontrado")
	}
	query := `
		UPDATE budgets
		SET project_id=$1, material_cost=$2, labor_cost=$3, equipment_cost=$4, total_cost=$5, notes=$6, document_url=$7
		WHERE id=$8`
	_, err := r.db.Exec(query, b.ProjectID, b.MaterialCost, b.LaborCost, b.EquipmentCost, b.TotalCost, b.Notes, b.DocumentURL, id)
	return err
}

func (r *BudgetRepository) Delete(id int) error {
	if r.useMemory {
		for i, b := range r.memData {
			if b.ID == id {
				r.memData = append(r.memData[:i], r.memData[i+1:]...)
				return nil
			}
		}
		return nil
	}
	_, err := r.db.Exec(`DELETE FROM budgets WHERE id = $1`, id)
	return err
}
