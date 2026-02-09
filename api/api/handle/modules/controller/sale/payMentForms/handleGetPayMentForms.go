package controller

import (
	"encoding/json"
	responsehelper "myApi/helpers/response"
	paymentform "myApi/interface/payMentForm"
	"net/http"
)

func HandleGetPayMentForms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	payMentForms, err := paymentform.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao retornar todos as formas de pagamento.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, payMentForms, "Todas as formas de pagamento.")

	json.NewEncoder(w).Encode(resp)
}
