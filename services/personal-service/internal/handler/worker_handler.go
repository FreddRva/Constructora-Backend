package handler

import (
	"encoding/json"
	"net/http"
	"staff-service/internal/model"
	"staff-service/internal/service"
	"strconv"
)

type WorkerHandler struct {
	service *service.WorkerService
}

func NewWorkerHandler(s *service.WorkerService) *WorkerHandler {
	return &WorkerHandler{service: s}
}

func (h *WorkerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /staff", h.GetAll)
	mux.HandleFunc("POST /staff", h.Create)
	mux.HandleFunc("GET /staff/{id}", h.GetByID)
	mux.HandleFunc("PUT /staff/{id}", h.Update)
	mux.HandleFunc("DELETE /staff/{id}", h.Delete)
}

func (h *WorkerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.Worker
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
		return
	}
	worker, err := h.service.Create(input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, worker)
}

func (h *WorkerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	workers, err := h.service.GetAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, workers)
}

func (h *WorkerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	worker, err := h.service.GetByID(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Trabajador no encontrado"})
		return
	}
	writeJSON(w, http.StatusOK, worker)
}

func (h *WorkerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	var input model.Worker
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
		return
	}
	worker, err := h.service.Update(id, input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, worker)
}

func (h *WorkerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	if err := h.service.Delete(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Trabajador eliminado"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
