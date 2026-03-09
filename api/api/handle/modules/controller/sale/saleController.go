package salecontroller

import (
	"encoding/json"
	"fmt"
	calchelper "myApi/helpers/calc"
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

func HandleGetSale(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	sales, err := sale.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao retornar todos as vendas.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)

	resp := responsehelper.Response(true, sales, "Todos as vendas.")

	json.NewEncoder(w).Encode(resp)

}

func HandlePostSale(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	var payload sale.SaleContract

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao validar os dados.", err)
		resp := responsehelper.Response(false, err, "Erro ao validar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if payload.Customer == "" {
		payload.Customer = "Consumidor padrão"
	}

	if err := payload.Validate(); len(err) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		u.ErrorLogger.Println("Campos obrigatórios ausentes:", err)
		resp := responsehelper.Response(false, err, "Campos obrigatórios ausentes.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	var subTotal float64

	for _, p := range payload.Products {
		u.GeneralLogger.Println("Produto aqui: ", p)

		subTotal += calchelper.CalculateTotalSale(p.SaleValue, p.Qtde)
	}

	payload.SaleValue = subTotal

	saleId, err := payload.Create()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro no create:", err)

		resp := responsehelper.Response(false, err, "Erro ao salvar a venda.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)

	payload.Id = saleId

	resp := responsehelper.Response(true, payload, "Dados da venda cadastrado com sucesso.")

	json.NewEncoder(w).Encode(resp)

}

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

	if err := sale.PaySale(payMents); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao processar o pagamento da venda.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, payMents, "Venda concluída com sucesso!")

	json.NewEncoder(w).Encode(resp)
}

func HandleInsertNewProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	var sale sale.SaleContract

	if err := json.NewDecoder(r.Body).Decode(&sale); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao processar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := sale.InsertNewItens(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao fazer o insert dos novos produtos: ", err.Error())

		resp := responsehelper.Response(false, nil, fmt.Sprintf("Erro ao fazer o insert dos novos produtos: %s", err.Error()))

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, nil, "Dados salvos com sucesso!")

	json.NewEncoder(w).Encode(resp)
}
