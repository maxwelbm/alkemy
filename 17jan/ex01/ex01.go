package ex01

import (
	"encoding/json"
	"net/http"
)

// Handler simples que retorna uma mensagem em JSON
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
}
