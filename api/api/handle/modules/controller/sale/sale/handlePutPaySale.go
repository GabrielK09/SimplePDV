package controller

import (
	"encoding/json"
	responsehelper "myApi/helpers/response"
	"myApi/interface/sale"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool

func HandlePutPaySale(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	var payMents sale.PaySaleContract

	if err := json.NewDecoder(r.Body).Decode(&payMents); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao processar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := payMents.ValidatePay(); len(err) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao validar o pagamento da venda.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := sale.PaySale(payMents.SaleId, payMents); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao processar o pagamento da venda.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, payMents, "Venda concluída com sucesso!")

	json.NewEncoder(w).Encode(resp)
}
