package controller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))
		return
	}

	u.ErrorLogger.Println("Ping")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsehelper.Response(true, "pong", "pong"))
}
