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

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))
		return
	} // Erro de método da rota

	var payMents processpayment.PayContract

	if err := json.NewDecoder(r.Body).Decode(&payMents); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar os dados."))
		return
	}

	if payMents.ShoppingId <= 0 && payMents.SaleId <= 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		u.ErrorLogger.Printf("IDs ausentes, shopping_id: %d, sale_id: %d", payMents.ShoppingId, payMents.SaleId)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Identificadores de compra ou venda ausentes ausentes."))
		return
	}

	if payMents.ShoppingId == 0 && payMents.SaleId > 0 {
		u.InfoLogger.Println("Vai finalizar uma venda.")

		if err := payMents.ValidatePay(payMents.SaleId); len(err) > 0 {
			u.ErrorLogger.Println("Erro ao validar o pagamento da venda.")
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao validar o pagamento da venda."))
			return
		}
	}

	if payMents.SaleId == 0 && payMents.ShoppingId > 0 {
		u.InfoLogger.Println("Vai finalizar uma compra.")

		if err := payMents.ValidatePay(payMents.ShoppingId); len(err) > 0 {
			u.ErrorLogger.Println("Erro ao validar o pagamento da compra.")
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao validar o pagamento da compra."))
			return
		}
	}

	if err := processpayment.PayMentShoppingOrSale(payMents); err != nil {
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
