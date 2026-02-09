package controller

import (
	"encoding/json"
	responsehelper "myApi/helpers/response"
	paymentform "myApi/interface/payMentForm"
	"net/http"
)

func HandlePutPayMentForms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	var pf paymentform.PayMentForms

	if err := json.NewDecoder(r.Body).Decode(&pf); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao processar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := ps.ValidatePay(); len(err) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao validar o pagamento da venda.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	//resp := responsehelper.Response(true, payMentForms, "Todas as formas de pagamento.")

	json.NewEncoder(w).Encode(resp)
}
