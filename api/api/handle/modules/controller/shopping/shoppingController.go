package shoppingcontroller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/shopping"
	"net/http"
)

func HandleGetAllShopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	shoppings, err := shopping.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, err, "Erro ao retornar todos os produtos.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsehelper.Response(true, shoppings, "Todas as compras cadastradas."))
}

func HandleGetLastShoppingId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	shoppingId, err := shopping.ReturnLastShoppingId()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao retornar o ID: ", err)
		resp := responsehelper.Response(false, err, "Falha na operação.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsehelper.Response(true, shoppingId, ""))
}

func HandlePostCreateShopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))
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
		if err["load"] != "" {
			w.WriteHeader(http.StatusInternalServerError)
			u.ErrorLogger.Println("Carga duplicada.", err)

			json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Dados duplicados"))
			return
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		u.ErrorLogger.Println("Campos obrigatórios ausentes.", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Campos obrigatórios ausentes."))
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
	json.NewEncoder(w).Encode(responsehelper.Response(true, shoppingId, "Compra cadastrada com sucesso!"))
}
