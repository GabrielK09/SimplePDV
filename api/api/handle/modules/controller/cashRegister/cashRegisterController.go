package cashregistercontroller

import (
	"context"
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/cashRegister"
	"myApi/interface/customer"
	paymentform "myApi/interface/payMentForm"

	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func HandlePostCashRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	var c cashRegister.CashRegisterContract

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao processaro os dados:", err)
		resp := responsehelper.Response(false, err, "Erro ao processaro os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Valida o Json

	if err := c.Validate(); len(err) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		u.ErrorLogger.Println("Erro ao validar o registro no caixa:", err)
		resp := responsehelper.Response(false, err, "Campos incorretos.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Valida os dados

	tx, err := conn.Begin(ctx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao iniciar a transação:", err)
		resp := responsehelper.Response(false, err, "Erro ao iniciar a transação.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	customerData, err := customer.Show(c.CustomerId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao localizar o cliente da venda:", err)
		resp := responsehelper.Response(false, err, "Erro ao localizar o cliente da venda.")

		json.NewEncoder(w).Encode(resp)
		return

	}

	paymentformData, err := paymentform.ShowById(c.SpecieId)

	c.CustomerId = customerData.Id
	c.Customer = customerData.Name

	c.SpecieId = paymentformData.Id
	c.Specie = paymentformData.Specie

	if err := c.Create(tx, c.InputValue, c.OutputValue); len(err) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao salvar o registro no caixa:", err)
		resp := responsehelper.Response(false, err, "Erro ao salvar o registro no caixa.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Cria

	w.WriteHeader(http.StatusCreated)

	resp := responsehelper.Response(true, c, "Todo movimento do caixa.")

	json.NewEncoder(w).Encode(resp)

}

func HandleGetCashRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	cashRegisters, err := cashRegister.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao retornar todo o caixa.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, cashRegisters, "Todo movimento do caixa.")

	json.NewEncoder(w).Encode(resp)

}
