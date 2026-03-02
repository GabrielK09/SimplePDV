package controller

import (
	"encoding/json"
	calchelper "myApi/helpers/calc"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/sale"

	"net/http"
)

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
