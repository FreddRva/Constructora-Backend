package repository

import (
	"database/sql"
	"errors"
	"procurement-service/internal/model"
)

type PurchaseRepository struct {
	db        *sql.DB
	useMemory bool
	memData   []model.Purchase
	nextID    int
}

func NewPurchaseRepository(db *sql.DB) *PurchaseRepository {
	return &PurchaseRepository{
		db:     db,
		nextID: 3,
		memData: []model.Purchase{
			{
				ID:           1,
				SupplierName: "UNACEM Distribuidora Oficial",
				RUC:          "20100137390",
				ItemName:     "Cemento Portland Tipo I (Lote 200 bls)",
				Quantity:     200,
				UnitPrice:    28.50,
				TotalAmount:  5700.00,
				ProjectID:    1,
				Status:       "Entregado",
			},
			{
				ID:           2,
				SupplierName: "Corporación Aceros Arequipa S.A.",
				RUC:          "20370146994",
				ItemName:     "Fierro Corrugado 1/2\" (Paquete)",
				Quantity:     150,
				UnitPrice:    39.00,
				TotalAmount:  5850.00,
				ProjectID:    1,
				Status:       "Pendiente",
			},
		},
	}
}

func (r *PurchaseRepository) InitSchema() error {
	if err := r.db.Ping(); err != nil {
		r.useMemory = true
		return err
	}
	query := `
	CREATE TABLE IF NOT EXISTS purchases (
		id SERIAL PRIMARY KEY,
		supplier_name VARCHAR(150) NOT NULL,
		ruc VARCHAR(30),
		item_name VARCHAR(150) NOT NULL,
		quantity NUMERIC(10,2) NOT NULL DEFAULT 1,
		unit_price NUMERIC(12,2) NOT NULL DEFAULT 0,
		total_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
		project_id INT DEFAULT 0,
		status VARCHAR(50) NOT NULL,
		invoice_url TEXT
	);`
	if _, err := r.db.Exec(query); err != nil {
		r.useMemory = true
		return err
	}
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM purchases`).Scan(&count); err == nil && count == 0 {
		for _, p := range r.memData {
			_, _ = r.db.Exec(
				`INSERT INTO purchases (supplier_name, ruc, item_name, quantity, unit_price, total_amount, project_id, status, invoice_url) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
				p.SupplierName, p.RUC, p.ItemName, p.Quantity, p.UnitPrice, p.TotalAmount, p.ProjectID, p.Status, p.InvoiceURL,
			)
		}
	}
	return nil
}

func (r *PurchaseRepository) Create(p *model.Purchase) error {
	if r.useMemory {
		p.ID = r.nextID
		r.nextID++
		r.memData = append([]model.Purchase{*p}, r.memData...)
		return nil
	}
	query := `
		INSERT INTO purchases (supplier_name, ruc, item_name, quantity, unit_price, total_amount, project_id, status, invoice_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`
	return r.db.QueryRow(
		query, p.SupplierName, p.RUC, p.ItemName, p.Quantity, p.UnitPrice, p.TotalAmount, p.ProjectID, p.Status, p.InvoiceURL,
	).Scan(&p.ID)
}

func (r *PurchaseRepository) GetAll() ([]model.Purchase, error) {
	if r.useMemory {
		return r.memData, nil
	}
	rows, err := r.db.Query(`SELECT id, supplier_name, ruc, item_name, quantity, unit_price, total_amount, project_id, status, invoice_url FROM purchases ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []model.Purchase{}
	for rows.Next() {
		var p model.Purchase
		if err := rows.Scan(&p.ID, &p.SupplierName, &p.RUC, &p.ItemName, &p.Quantity, &p.UnitPrice, &p.TotalAmount, &p.ProjectID, &p.Status, &p.InvoiceURL); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PurchaseRepository) GetByID(id int) (*model.Purchase, error) {
	if r.useMemory {
		for _, p := range r.memData {
			if p.ID == id {
				return &p, nil
			}
		}
		return nil, errors.New("compra no encontrada")
	}
	var p model.Purchase
	query := `SELECT id, supplier_name, ruc, item_name, quantity, unit_price, total_amount, project_id, status, invoice_url FROM purchases WHERE id = $1`
	if err := r.db.QueryRow(query, id).Scan(&p.ID, &p.SupplierName, &p.RUC, &p.ItemName, &p.Quantity, &p.UnitPrice, &p.TotalAmount, &p.ProjectID, &p.Status, &p.InvoiceURL); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PurchaseRepository) Update(id int, p *model.Purchase) error {
	if r.useMemory {
		for i := range r.memData {
			if r.memData[i].ID == id {
				p.ID = id
				r.memData[i] = *p
				return nil
			}
		}
		return errors.New("compra no encontrada")
	}
	query := `
		UPDATE purchases
		SET supplier_name=$1, ruc=$2, item_name=$3, quantity=$4, unit_price=$5, total_amount=$6, project_id=$7, status=$8, invoice_url=$9
		WHERE id=$10`
	_, err := r.db.Exec(query, p.SupplierName, p.RUC, p.ItemName, p.Quantity, p.UnitPrice, p.TotalAmount, p.ProjectID, p.Status, p.InvoiceURL, id)
	return err
}

func (r *PurchaseRepository) Delete(id int) error {
	if r.useMemory {
		for i, p := range r.memData {
			if p.ID == id {
				r.memData = append(r.memData[:i], r.memData[i+1:]...)
				return nil
			}
		}
		return nil
	}
	_, err := r.db.Exec(`DELETE FROM purchases WHERE id = $1`, id)
	return err
}
