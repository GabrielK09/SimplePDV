package shoppingcontroller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/shopping"
	"net/http"
)

func HandleGetSaleWithProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

}

func HandlePostCreateShopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	var payload shopping.ShoppingContract

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.ErrorLogger.Println("Erro ao processar os dados: ", err)
		resp := responsehelper.Response(false, err, "Erro ao processar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := payload.Validate(); len(err) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		u.ErrorLogger.Println("Campos obrigatórios ausentes.", err)
		resp := responsehelper.Response(false, err, "Campos obrigatórios ausentes.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	shoppingId, err := payload.Create()

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		u.ErrorLogger.Println("Erro ao cadastrar a compra.", err)
		resp := responsehelper.Response(false, err, "Erro ao cadastrar a compra.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := responsehelper.Response(true, shoppingId, "Compra cadastrada com sucesso!")

	json.NewEncoder(w).Encode(resp)
}
