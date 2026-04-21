package paymentcontroller

import (
	"encoding/json"
	"fmt"
	processpayment "myApi/api/services/payMent"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	_ "myApi/interface/sale"
	_ "myApi/interface/shopping"
	"net/http"
	"strings"
)

func HandlePutPaySaleOrShopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var label string

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))
		return
	} // Erro de método da rota

	var payMents processpayment.PayContract

	if err := json.NewDecoder(r.Body).Decode(&payMents); err != nil {
		u.ErrorLogger.Println("Erro ao processar os dados:", err)
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar os dados."))
		return
	}

	if err := payMents.Validate(); len(err) > 0 {
		u.ErrorLogger.Println("Erro ao validar o pagamento da compra.")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao validar o pagamento da compra."))
		return
	}

	var totalPaide float64

	for _, p := range payMents.Species {
		totalPaide += p.AmountPaid
	}

	if err := processpayment.PayMentShoppingOrSale(payMents, totalPaide); err != nil {
		u.ErrorLogger.Println("Erro ao pagar a compra/venda: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar o pagamento da venda."))
		return
	}

	w.WriteHeader(http.StatusOK)

	if payMents.SaleId > 0 {
		label = "Venda"
	}

	if payMents.ShoppingId > 0 {
		label = "Compra"
	}

	json.NewEncoder(w).Encode(responsehelper.Response(true, payMents, fmt.Sprintf("%s concluída com sucesso!", label)))
}

func HandlePutCancelSaleOrShopping(w http.ResponseWriter, r *http.Request) {
	u.InfoLogger.Println("Called HandlePutCancelSaleOrShopping")

	var label string
	var cancelBody processpayment.CancelContract

	if err := json.NewDecoder(r.Body).Decode(&cancelBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar os dados."))
		return
	}

	if cancelBody.SaleId > 0 {
		label = "Venda"
	}

	if cancelBody.ShoppingId > 0 {
		label = "Compra"
	}

	if err := processpayment.CancelSaleOrShopping(cancelBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, fmt.Sprintf("Erro ao cancelar %s", strings.ToLower(label))))
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsehelper.Response(true, nil, fmt.Sprintf("%s cancelada com sucesso!", label)))
}
