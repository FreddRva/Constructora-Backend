package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"project-service/internal/factory"
	"project-service/internal/handler"
	"project-service/internal/repository"
	"project-service/internal/service"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := getEnv("DB_URL", "postgres://postgres:Constructora2026@constructora-db.cab8ecee8zdw.us-east-1.rds.amazonaws.com:5432/projects_db?sslmode=require")
	port := getEnv("PORT", "8081")

	ensureDatabase(dbURL, "projects_db")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error abriendo conexión a BD: %v", err)
	}
	defer db.Close()

	repo := repository.NewProjectRepository(db)
	if err := repo.InitSchema(); err != nil {
		log.Printf("Aviso: no se pudo inicializar tabla en Postgres (%v); usando respaldo en memoria", err)
	}

	projFactory := factory.NewProjectFactory()
	projService := service.NewProjectService(repo, projFactory)
	projHandler := handler.NewProjectHandler(projService)

	mux := http.NewServeMux()
	projHandler.RegisterRoutes(mux)

	log.Printf("Servicio de Proyectos corriendo en el puerto %s", port)
	if err := http.ListenAndServe(":"+port, enableCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func ensureDatabase(dbURL, dbName string) {
	adminURL := strings.Replace(dbURL, "/"+dbName, "/postgres", 1)
	if adminDB, err := sql.Open("postgres", adminURL); err == nil {
		defer adminDB.Close()
		_, _ = adminDB.Exec("CREATE DATABASE " + dbName)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
