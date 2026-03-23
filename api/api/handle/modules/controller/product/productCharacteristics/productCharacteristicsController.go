package productcharacteristicscontroller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	productcharacteristics "myApi/interface/product/productCharacteristics"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func HandlePostCreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))

		return
	}

	var productCharacteristics productcharacteristics.ProductCharacteristicsContract

	if err := json.NewDecoder(r.Body).Decode(&productCharacteristics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.ErrorLogger.Println("Erro ao processar os dados: ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar os dados."))
		return
	}

	if err := productCharacteristics.Create(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.ErrorLogger.Println("Erro ao criar a característica do produto: ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar os dados."))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responsehelper.Response(true, productCharacteristics, "Característica do produto cadastrada com sucesso!"))
}

func HandlePutUpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))

		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Id inválido."))

		return
	}

	productData, err := productcharacteristics.Show(id)

	if err := json.NewDecoder(r.Body).Decode(&productData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.ErrorLogger.Println("Erro ao processar os dados: ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar os dados."))
		return
	}

	productData.Id = id

	if err := productData.Update(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao alterar os dados da característica. ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao alterar os dados da característica."))
		return
	}

	json.NewEncoder(w).Encode(responsehelper.Response(true, nil, "Característica do produto alterada com sucesso!"))
}
