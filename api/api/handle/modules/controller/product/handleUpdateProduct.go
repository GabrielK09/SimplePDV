package controller

import (
	"encoding/json"
	"log"
	responsehelper "myApi/helpers/response"
	"myApi/interface/product"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5"
)

func HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
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

		log.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	product, err := product.Show(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		log.Println("Produto não localizado:", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Produto não localizado."),
		)

		return
	}

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao processar dados"),
		)

		return
	}

	product.Id = id

	updated, err := product.Update()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao atualizar o produto não localizado."),
		)

		return
	}

	json.NewEncoder(w).Encode(
		responsehelper.Response(true, updated, "Produto alterado com sucesso!"),
	)
}
