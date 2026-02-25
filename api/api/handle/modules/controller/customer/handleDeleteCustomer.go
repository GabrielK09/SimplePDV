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

func HandleDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
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
		//O cliente padrão não pode ser desativado.
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("O cliente padrão não pode ser desativado.")

		json.NewEncoder(w).Encode(
			responsehelper.Response(false, nil, "O cliente padrão não pode ser desativado."),
		)

		return
	}

	if err := customer.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao deletar o cliente: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao deletar o Cliente."),
		)

		return
	}

	resp := responsehelper.Response(true, nil, "Cliente deletado com sucesso!")

	json.NewEncoder(w).Encode(resp)
}
