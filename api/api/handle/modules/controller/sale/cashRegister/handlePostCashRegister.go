package controller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/cashRegister"

	"net/http"
)

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
		resp := responsehelper.Response(false, err, "Erro ao processaro os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Valida o Json

	if err := c.Validate(); len(err) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp := responsehelper.Response(false, err, "Campos incorretos.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Valida os dados

	if err := c.Create(nil, c.InputValue, c.OutputValue, 0, 0, c.CustomerId, c.Customer); len(err) > 0 {
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
