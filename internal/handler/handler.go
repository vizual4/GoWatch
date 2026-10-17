package handler

import (
	"GoWatch/internal/domain"
	"GoWatch/internal/repository"
	"GoWatch/internal/service"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

type MonitorsResp struct {
	Monitors []*domain.Monitor `json:"monitors"`
}

type MonitorHandler struct {
	service *service.MonitorService
}

func NewMonitorHandler(service *service.MonitorService) *MonitorHandler {
	return &MonitorHandler{service: service}
}

func (h *MonitorHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/monitors", h.CreateMonHandler)
	mux.HandleFunc("GET /api/v1/monitors", h.MonsHandler)
	mux.HandleFunc("GET /api/v1/monitors/{id}", h.GetMonHandler)
	mux.HandleFunc("DELETE /api/v1/monitors/{id}", h.DelMonHandler)
}

// Создает объект, возвращает код 201 и сам объект с присвоенным id и временем создания
//
// В случае ошибки возвращает ошибку и соответствующий статус
func (h *MonitorHandler) CreateMonHandler(w http.ResponseWriter, r *http.Request) {
	var info domain.Monitor

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error:": err.Error()}, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &info)
	if err != nil {
		writeJson(w, map[string]string{"error:": err.Error()}, http.StatusBadRequest)
		return
	}

	//вызвать метод Create из service
	err = h.service.CreateMon(&info)
	if err != nil {
		writeJson(w, map[string]string{"error:": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, info, http.StatusCreated)
}

// Выводит список всех созданных объектов
func (h *MonitorHandler) MonsHandler(w http.ResponseWriter, r *http.Request) {
	var result []*domain.Monitor

	result, err := h.service.AllMons()
	if err != nil {
		log.Println(err)
		writeJson(w, map[string]string{"error": "internal error"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, MonitorsResp{Monitors: result}, http.StatusOK)
}

func (h *MonitorHandler) GetMonHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJson(w, map[string]string{"error:": "invalid id format"}, http.StatusBadRequest)
		return
	}

	res, err := h.service.GetMon(id)
	if err != nil {
		if errors.Is(err, repository.NF) {
			writeJson(w, map[string]string{"error": "monitor not found"}, http.StatusNotFound)
			return
		}

		log.Println(err)
		writeJson(w, map[string]string{"error:": "internal error"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, res, http.StatusOK)
}

func (h *MonitorHandler) DelMonHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJson(w, map[string]string{"error:": "invalid id format"}, http.StatusBadRequest)
		return
	}

	err = h.service.DelMon(id)
	if err != nil {
		if errors.Is(err, repository.NF) {
			writeJson(w, map[string]string{"error:": "monitor not found"}, http.StatusNotFound)
			return
		}

		log.Println(err)
		writeJson(w, map[string]string{"error:": "internal error"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK)
}

// Сериализует json объект
func writeJson(w http.ResponseWriter, data any, status int) {
	resp, err := json.Marshal(data)
	if err != nil {
		log.Println("ERROR: failed to marshal json:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(resp)
}
