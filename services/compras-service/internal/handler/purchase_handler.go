package handler

import (
	"encoding/json"
	"net/http"
	"procurement-service/internal/model"
	"procurement-service/internal/service"
	"strconv"
)

type PurchaseHandler struct {
	service *service.PurchaseService
}

func NewPurchaseHandler(s *service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: s}
}

func (h *PurchaseHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /purchases", h.GetAll)
	mux.HandleFunc("POST /purchases", h.Create)
	mux.HandleFunc("GET /purchases/{id}", h.GetByID)
	mux.HandleFunc("PUT /purchases/{id}", h.Update)
	mux.HandleFunc("DELETE /purchases/{id}", h.Delete)
}

func (h *PurchaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.Purchase
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
		return
	}
	purchase, err := h.service.Create(input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, purchase)
}

func (h *PurchaseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.GetAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *PurchaseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	purchase, err := h.service.GetByID(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Compra no encontrada"})
		return
	}
	writeJSON(w, http.StatusOK, purchase)
}

func (h *PurchaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	var input model.Purchase
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
		return
	}
	purchase, err := h.service.Update(id, input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, purchase)
}

func (h *PurchaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	if err := h.service.Delete(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Orden de compra eliminada"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
