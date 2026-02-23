package controller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/sale"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func HandlePutCancelSale(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

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

	saleDetail, err := sale.Show(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao procurar a venda: ", err)
		resp := responsehelper.Response(false, err, "Erro ao procurar a venda.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	canceledSale, err := saleDetail.CancelSale()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao cancelar a venda: ", err)
		resp := responsehelper.Response(false, err, "Erro ao cancelar a venda.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, canceledSale, "Venda cancelada com sucesso!")

	json.NewEncoder(w).Encode(resp)
}
