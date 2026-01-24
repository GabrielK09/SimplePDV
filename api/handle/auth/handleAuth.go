package handle

import (
	"encoding/json"
	responsehelper "myApi/helpers/response"
	"net/http"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, nil, "Rota para criar usuário")

	json.NewEncoder(w).Encode(resp)
}
