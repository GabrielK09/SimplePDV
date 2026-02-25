package controller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/customer"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func HandlePutCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	if id <= 1 {
		//O cliente padrão não pode ser desativado.
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("O Id precisa ser maior que 1")
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	if id == 1 {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("O cliente padrão não pode ser alterado.")

		json.NewEncoder(w).Encode(
			responsehelper.Response(false, nil, "O cliente padrão não pode ser desativado."),
		)

		return
	}

	customer, err := customer.Show(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Cliente não localizado.")

		json.NewEncoder(w).Encode(
			responsehelper.Response(false, nil, "Cliente não localizado."),
		)

		return
	}

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("O cliente padrão não pode ser alterado.")

		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao processar dados"),
		)

		return
	}

	customer.Id = id

	updated, err := customer.Update()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		u.ErrorLogger.Println("O cliente padrão não pode ser alterado.")

		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao processar dados"),
		)

		return
	}

	resp := responsehelper.Response(true, updated, "Cliente alterado com sucesso!")

	json.NewEncoder(w).Encode(resp)

}
