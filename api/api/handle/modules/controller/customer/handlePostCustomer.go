package controller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/customer"
	"net/http"
)

func HandlePostCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	var c customer.CustomerContract

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao validar os dados.", err)
		resp := responsehelper.Response(false, err, "Erro ao validar os dados do cliente.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := c.Create(); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Erro ao cadastrar o cliente: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Erro ao cadastrar o cliente."),
		)

		return
	}

	w.WriteHeader(http.StatusCreated)

	resp := responsehelper.Response(true, c, "Cliente cadastrado com sucesso!")

	json.NewEncoder(w).Encode(resp)
}
