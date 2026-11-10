package repository

import (
	"database/sql"
	"errors"
	"staff-service/internal/model"
)

type WorkerRepository struct {
	db        *sql.DB
	useMemory bool
	memData   []model.Worker
	nextID    int
}

func NewWorkerRepository(db *sql.DB) *WorkerRepository {
	return &WorkerRepository{
		db:     db,
		nextID: 4,
		memData: []model.Worker{
			{ID: 1, Name: "Ing. Carlos Mendoza", Role: "Ingeniero Residente", Phone: "987-654-321", DailyWage: 180.00, ProjectID: 1, Status: "Activo"},
			{ID: 2, Name: "Jorge Quispe Huamán", Role: "Maestro de Obra", Phone: "912-345-678", DailyWage: 120.00, ProjectID: 1, Status: "Activo"},
			{ID: 3, Name: "Luis Paredes Rojas", Role: "Operario Electricista", Phone: "945-112-233", DailyWage: 95.00, ProjectID: 2, Status: "Activo"},
		},
	}
}

func (r *WorkerRepository) InitSchema() error {
	if err := r.db.Ping(); err != nil {
		r.useMemory = true
		return err
	}
	query := `
	CREATE TABLE IF NOT EXISTS workers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(150) NOT NULL,
		role VARCHAR(100) NOT NULL,
		phone VARCHAR(50),
		daily_wage NUMERIC(10,2) DEFAULT 0,
		project_id INT DEFAULT 0,
		status VARCHAR(50) NOT NULL
	);`
	if _, err := r.db.Exec(query); err != nil {
		r.useMemory = true
		return err
	}
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM workers`).Scan(&count); err == nil && count == 0 {
		for _, w := range r.memData {
			_, _ = r.db.Exec(
				`INSERT INTO workers (name, role, phone, daily_wage, project_id, status) VALUES ($1,$2,$3,$4,$5,$6)`,
				w.Name, w.Role, w.Phone, w.DailyWage, w.ProjectID, w.Status,
			)
		}
	}
	return nil
}

func (r *WorkerRepository) Create(w *model.Worker) error {
	if r.useMemory {
		w.ID = r.nextID
		r.nextID++
		r.memData = append([]model.Worker{*w}, r.memData...)
		return nil
	}
	query := `
		INSERT INTO workers (name, role, phone, daily_wage, project_id, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	return r.db.QueryRow(query, w.Name, w.Role, w.Phone, w.DailyWage, w.ProjectID, w.Status).Scan(&w.ID)
}

func (r *WorkerRepository) GetAll() ([]model.Worker, error) {
	if r.useMemory {
		return r.memData, nil
	}
	rows, err := r.db.Query(`SELECT id, name, role, phone, daily_wage, project_id, status FROM workers ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workers := []model.Worker{}
	for rows.Next() {
		var w model.Worker
		if err := rows.Scan(&w.ID, &w.Name, &w.Role, &w.Phone, &w.DailyWage, &w.ProjectID, &w.Status); err != nil {
			return nil, err
		}
		workers = append(workers, w)
	}
	return workers, nil
}

func (r *WorkerRepository) GetByID(id int) (*model.Worker, error) {
	if r.useMemory {
		for _, w := range r.memData {
			if w.ID == id {
				return &w, nil
			}
		}
		return nil, errors.New("trabajador no encontrado")
	}
	var w model.Worker
	query := `SELECT id, name, role, phone, daily_wage, project_id, status FROM workers WHERE id = $1`
	if err := r.db.QueryRow(query, id).Scan(&w.ID, &w.Name, &w.Role, &w.Phone, &w.DailyWage, &w.ProjectID, &w.Status); err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WorkerRepository) Update(id int, w *model.Worker) error {
	if r.useMemory {
		for i := range r.memData {
			if r.memData[i].ID == id {
				w.ID = id
				r.memData[i] = *w
				return nil
			}
		}
		return errors.New("trabajador no encontrado")
	}
	query := `UPDATE workers SET name=$1, role=$2, phone=$3, daily_wage=$4, project_id=$5, status=$6 WHERE id=$7`
	_, err := r.db.Exec(query, w.Name, w.Role, w.Phone, w.DailyWage, w.ProjectID, w.Status, id)
	return err
}

func (r *WorkerRepository) Delete(id int) error {
	if r.useMemory {
		for i, w := range r.memData {
			if w.ID == id {
				r.memData = append(r.memData[:i], r.memData[i+1:]...)
				return nil
			}
		}
		return nil
	}
	_, err := r.db.Exec(`DELETE FROM workers WHERE id = $1`, id)
	return err
}
