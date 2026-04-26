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

	if payMents.SaleId > 0 {
		label = "Venda"
	}

	if payMents.ShoppingId > 0 {
		label = "Compra"
	}

	if err := payMents.Validate(); len(err) > 0 {
		u.ErrorLogger.Println("Erro ao validar o pagamento da compra.")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao validar o pagamento da compra."))
		return
	}

	if err := processpayment.PayMentShoppingOrSale(r.Context(), payMents); err != nil {
		u.ErrorLogger.Printf("Erro ao pagar a %s: %s", label, err)
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, fmt.Sprintf("Erro ao processar o pagamento da %s.", label)))
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsehelper.Response(true, payMents, fmt.Sprintf("%s concluída com sucesso!", label)))
}

func HandlePutCancelSaleOrShopping(w http.ResponseWriter, r *http.Request) {
	u.InfoLogger.Println("Called HandlePutCancelSaleOrShopping")

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))
		return
	} // Erro de método da rota

	var cancelBody processpayment.CancelContract
	if err := json.NewDecoder(r.Body).Decode(&cancelBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar os dados."))
		return
	}

	if cancelBody.SaleId > 0 && cancelBody.ShoppingId > 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Uma venda e uma compra não podem ser canceladas ao mesmo tempo."))
		return
	}

	var label string

	if cancelBody.SaleId > 0 {
		label = "Venda"
	}

	if cancelBody.ShoppingId > 0 {
		label = "Compra"
	}

	if err := processpayment.CancelSaleOrShopping(r.Context(), cancelBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro no CancelSaleOrShopping:", err)

		errorsField := make(map[string]string)
		errorsField["cancel"] = fmt.Sprintf("%s", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, errorsField, fmt.Sprintf("Erro ao cancelar %s", strings.ToLower(label))))
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsehelper.Response(true, nil, fmt.Sprintf("%s cancelada com sucesso!", label)))
}
