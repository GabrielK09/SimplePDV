package controller

import (
	"encoding/json"
	responsehelper "myApi/helpers/response"
	"myApi/interface/sale"

	"net/http"
)

func HandleGetSale(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	sales, err := sale.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao retornar todos as vendas.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)

	resp := responsehelper.Response(true, sales, "Todos as vendas.")

	json.NewEncoder(w).Encode(resp)

}
