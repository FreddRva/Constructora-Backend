package handler

import (
	"budget-service/internal/model"
	"budget-service/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

// BudgetHandler expone los endpoints HTTP REST del microservicio de Presupuestos.
type BudgetHandler struct {
	service *service.BudgetService
}

func NewBudgetHandler(s *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{service: s}
}

func (h *BudgetHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /budgets", h.GetAll)
	mux.HandleFunc("POST /budgets", h.Create)
	mux.HandleFunc("GET /budgets/{id}", h.GetByID)
	mux.HandleFunc("PUT /budgets/{id}", h.Update)
	mux.HandleFunc("DELETE /budgets/{id}", h.Delete)
}

func (h *BudgetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.Budget
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
		return
	}
	budget, err := h.service.Create(input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, budget)
}

func (h *BudgetHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	budgets, err := h.service.GetAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, budgets)
}

func (h *BudgetHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	budget, err := h.service.GetByID(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Presupuesto no encontrado"})
		return
	}
	writeJSON(w, http.StatusOK, budget)
}

func (h *BudgetHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	var input model.Budget
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
		return
	}
	budget, err := h.service.Update(id, input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, budget)
}

func (h *BudgetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	if err := h.service.Delete(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Presupuesto eliminado"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

