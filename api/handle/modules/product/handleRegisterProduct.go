package handle

import (
	"encoding/json"
	responsehelper "myApi/helpers/response"
	"myApi/interface/product"
	"net/http"
)

func HandleRegisterProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	var payload product.ProductContract

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := responsehelper.Response(false, err, "Erro ao processar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := payload.Validate(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp := responsehelper.Response(false, err, "Campos obrigatórios ausentes.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := payload.Create(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao gravar o produto.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := responsehelper.Response(true, payload, "Produto criado com sucesso!")

	json.NewEncoder(w).Encode(resp)
}
