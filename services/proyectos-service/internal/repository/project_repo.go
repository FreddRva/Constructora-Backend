package repository

import (
	"database/sql"
	"errors"
	"project-service/internal/model"
)

type ProjectRepository struct {
	db        *sql.DB
	useMemory bool
	memData   []model.Project
	nextID    int
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		db:     db,
		nextID: 3,
		memData: []model.Project{
			{
				ID:        1,
				Name:      "Residencial Los Pinos Torre A",
				Client:    "Inmobiliaria Horizonte SAC",
				Location:  "Av. Primavera 1250, Surco",
				StartDate: "2026-03-01",
				EndDate:   "2026-12-15",
				Status:    "En Ejecución",
				Progress:  45,
			},
			{
				ID:        2,
				Name:      "Centro Logístico Industrial Sur",
				Client:    "Grupo Logístico Andino",
				Location:  "Panamericana Sur Km 38, Lurín",
				StartDate: "2026-05-10",
				EndDate:   "2027-02-28",
				Status:    "Planificación",
				Progress:  10,
			},
		},
	}
}

func (r *ProjectRepository) InitSchema() error {
	if err := r.db.Ping(); err != nil {
		r.useMemory = true
		return err
	}
	query := `
	CREATE TABLE IF NOT EXISTS projects (
		id SERIAL PRIMARY KEY,
		name VARCHAR(150) NOT NULL,
		client VARCHAR(150) NOT NULL,
		location VARCHAR(200),
		start_date VARCHAR(50),
		end_date VARCHAR(50),
		status VARCHAR(50) NOT NULL,
		progress NUMERIC(5,2) DEFAULT 0,
		image_url TEXT
	);`
	if _, err := r.db.Exec(query); err != nil {
		r.useMemory = true
		return err
	}
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM projects`).Scan(&count); err == nil && count == 0 {
		for _, p := range r.memData {
			_, _ = r.db.Exec(
				`INSERT INTO projects (name, client, location, start_date, end_date, status, progress, image_url) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
				p.Name, p.Client, p.Location, p.StartDate, p.EndDate, p.Status, p.Progress, p.ImageURL,
			)
		}
	}
	return nil
}

func (r *ProjectRepository) Create(p *model.Project) error {
	if r.useMemory {
		p.ID = r.nextID
		r.nextID++
		r.memData = append([]model.Project{*p}, r.memData...)
		return nil
	}
	query := `
		INSERT INTO projects (name, client, location, start_date, end_date, status, progress, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`
	return r.db.QueryRow(
		query, p.Name, p.Client, p.Location, p.StartDate, p.EndDate, p.Status, p.Progress, p.ImageURL,
	).Scan(&p.ID)
}

func (r *ProjectRepository) GetAll() ([]model.Project, error) {
	if r.useMemory {
		return r.memData, nil
	}
	rows, err := r.db.Query(`SELECT id, name, client, location, start_date, end_date, status, progress, image_url FROM projects ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []model.Project{}
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Client, &p.Location, &p.StartDate, &p.EndDate, &p.Status, &p.Progress, &p.ImageURL); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepository) GetByID(id int) (*model.Project, error) {
	if r.useMemory {
		for _, p := range r.memData {
			if p.ID == id {
				return &p, nil
			}
		}
		return nil, errors.New("proyecto no encontrado")
	}
	var p model.Project
	query := `SELECT id, name, client, location, start_date, end_date, status, progress, image_url FROM projects WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Client, &p.Location, &p.StartDate, &p.EndDate, &p.Status, &p.Progress, &p.ImageURL)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepository) Update(id int, p *model.Project) error {
	if r.useMemory {
		for i := range r.memData {
			if r.memData[i].ID == id {
				p.ID = id
				r.memData[i] = *p
				return nil
			}
		}
		return errors.New("proyecto no encontrado")
	}
	query := `
		UPDATE projects
		SET name=$1, client=$2, location=$3, start_date=$4, end_date=$5, status=$6, progress=$7, image_url=$8
		WHERE id=$9`
	_, err := r.db.Exec(query, p.Name, p.Client, p.Location, p.StartDate, p.EndDate, p.Status, p.Progress, p.ImageURL, id)
	return err
}

func (r *ProjectRepository) Delete(id int) error {
	if r.useMemory {
		for i, p := range r.memData {
			if p.ID == id {
				r.memData = append(r.memData[:i], r.memData[i+1:]...)
				return nil
			}
		}
		return nil
	}
	_, err := r.db.Exec(`DELETE FROM projects WHERE id = $1`, id)
	return err
}
