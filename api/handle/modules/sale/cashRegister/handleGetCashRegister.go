package cashregister

import (
	"encoding/json"
	responsehelper "myApi/helpers/response"
	"myApi/interface/cashRegister"

	"net/http"
)

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
		resp := responsehelper.Response(false, err, "Erro ao retornar todos os produtos.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, cashRegisters, "Todo movimento do caixa.")

	json.NewEncoder(w).Encode(resp)

}
