package paymentcontroller

import (
	"encoding/json"
	processpayment "myApi/api/services/payMent"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	_ "myApi/interface/sale"
	_ "myApi/interface/shopping"
	"net/http"
)

func HandlePutPaySaleOrShopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		u.ErrorLogger.Println("IDs ausentes.")

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Identificadores de compra ou venda ausentes ausentes."))
		return
	}

	if payMents.ShoppingId == 0 && payMents.SaleId > 0 {
		if err := payMents.ValidatePay(payMents.SaleId); len(err) > 0 {
			u.ErrorLogger.Println("Erro ao validar o pagamento da venda.")
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao validar o pagamento da venda."))
			return
		}

		if err := processpayment.PaySale(payMents); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar o pagamento da venda."))
			return
		}
	}

	if payMents.SaleId == 0 && payMents.ShoppingId > 0 {
		if err := payMents.ValidatePay(payMents.ShoppingId); len(err) > 0 {
			u.ErrorLogger.Println("Erro ao validar o pagamento da compra.")
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao validar o pagamento da compra."))
			return
		}

		if err := processpayment.PaySale(payMents); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar o pagamento da venda."))
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsehelper.Response(true, payMents, "Venda concluída com sucesso!"))
}

func HandlePutCancelSale(w http.ResponseWriter, r *http.Request) {}
