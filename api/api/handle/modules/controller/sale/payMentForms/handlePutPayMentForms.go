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

	payMent, err := paymentform.Show()

	if err := json.NewDecoder(r.Body).Decode(&payMent); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao processar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := payMent.Validate(); len(err) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao validar o pagamento da venda.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	updated, err := payMent.Update()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao atualizar o produto não localizado."),
		)

		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, updated, "Forma de pagamento alterada com sucesso!")

	json.NewEncoder(w).Encode(resp)
}
