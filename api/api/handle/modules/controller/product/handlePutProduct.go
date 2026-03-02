package controller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/product"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5"
)

func HandlePutProduct(w http.ResponseWriter, r *http.Request) {
	var productData *product.ProductContract

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

	productData, err = product.Show(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		u.ErrorLogger.Println("Produto não localizado:", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Produto não localizado."),
		)

		return
	}

	u.InfoLogger.Println("Produto:", productData)

	if err := json.NewDecoder(r.Body).Decode(&productData); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Erro ao processar dados: ", err)

		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao processar dados"),
		)

		return
	}

	productData.Id = id

	updated, err := productData.Update()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		u.ErrorLogger.Println("Erro ao alterar o produto: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao atualizar"),
		)

		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		responsehelper.Response(true, updated, "Produto alterado com sucesso!"),
	)
}
