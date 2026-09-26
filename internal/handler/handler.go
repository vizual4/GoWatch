package handler

import (
	"GoWatch/internal/domain"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/monitors", createMonHandler)
}

func createMonHandler(w http.ResponseWriter, r *http.Request) {
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
}

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
