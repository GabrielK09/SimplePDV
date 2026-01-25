package handle

import (
	"context"
	"encoding/json"
	responsehelper "myApi/helpers/response"
	"net/http"

	_ "github.com/jackc/pgx/v5"
)

func HandleRegisterProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Method error

	//w.WriteHeader(http.StatusCreated)

	defer conn.Close(context.Background())

	resp := responsehelper.Response(true, nil, "Usuário criado com sucesso!")

	json.NewEncoder(w).Encode(resp)
}
