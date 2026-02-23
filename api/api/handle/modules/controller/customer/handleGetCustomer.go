package controller

import (
	"encoding/json"
	responsehelper "myApi/helpers/response"
	"myApi/interface/customer"
	"net/http"
)

func HandleGetCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	customers, err := customer.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao retornar todos os clientes."),
		)

		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, customers, "Todos os clientes.")

	json.NewEncoder(w).Encode(resp)
}
