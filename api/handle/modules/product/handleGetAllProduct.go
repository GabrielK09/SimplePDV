package handle

import (
	"encoding/json"
	responsehelper "myApi/helpers/response"
	"net/http"
)

func HandleGetAllProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Method error

	//w.WriteHeader(http.StatusCreated)

	resp := responsehelper.Response(true, nil, "Usuário criado com sucesso!")

	json.NewEncoder(w).Encode(resp)
}
