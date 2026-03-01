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

func HandleGetSaleWithProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
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

	saleCommissionDetails, err := sale.ShowTotalCommission(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao procurar a comissão venda: ", err)
		resp := responsehelper.Response(false, err, "Erro ao procurar a comissão venda.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	responseData := map[string]any{
		"sale":       saleDetail,
		"commission": saleCommissionDetails,
	}

	resp := responsehelper.Response(true, responseData, "Detalhes das vendas!")

	json.NewEncoder(w).Encode(resp)
}
