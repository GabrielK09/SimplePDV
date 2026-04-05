package productcharacteristicscontroller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	productcharacteristics "myApi/interface/product/productCharacteristics"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type FindByIds struct {
	GridId        int `json:"grid_id"`
	ProductGridId int `json:"product_grid_id"`
}

func HandlePostCreateProductCharacteristics(w http.ResponseWriter, r *http.Request) {
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

func HandleDeleteProductCharacteristics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))

		return
	}

	params := mux.Vars(r)
	productId, err := strconv.Atoi(params["product_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id do produto inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id do produto inválido."),
		)

		return
	}

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	if err := productcharacteristics.Delete(id, productId, time.Now()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao deletar o produto."),
		)

		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(responsehelper.Response(true, nil, "Produto deletado com sucesso!"))
}

func HandlePutUpdateProductCharacteristics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))

		return
	}

	var productCharacteristicsPayLoad productcharacteristics.ProductCharacteristicsContract

	if err := json.NewDecoder(r.Body).Decode(&productCharacteristicsPayLoad); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.ErrorLogger.Println("Erro ao processar os dados: ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar os dados."))
		return
	}

	if productCharacteristicsPayLoad.Id > 0 && productCharacteristicsPayLoad.ProductId > 0 {
		if err := productCharacteristicsPayLoad.Update(productCharacteristicsPayLoad.Id, productCharacteristicsPayLoad.ProductId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			u.ErrorLogger.Println("Erro ao alterar os dados da característica. ", err)

			json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao alterar os dados da característica."))
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		if err := productCharacteristicsPayLoad.Create(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			u.ErrorLogger.Println("Erro ao criar a característica do produto: ", err)

			json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar os dados."))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}

	json.NewEncoder(w).Encode(responsehelper.Response(true, nil, "Característica do produto alterada com sucesso!"))
}

func HandleGetProductCharacteristicsByGridId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))

		return
	}

	params := mux.Vars(r)
	productId, err := strconv.Atoi(params["product_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id do produto inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id do produto inválido."),
		)

		return
	}

	gridId, err := strconv.Atoi(params["grid_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id da grade inválida: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id da grade inválido."),
		)

		return
	}

	if productId <= 0 || gridId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		u.ErrorLogger.Println("Ids inválidos.")

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Ids inválidos."))
		return
	}

	productcharacteristicData, err := productcharacteristics.ShowById(gridId, productId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao alterar os dados da característica. ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao localizar os dados da característica."))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsehelper.Response(true, productcharacteristicData, "Característica do produto localizada com sucesso!"))
}

func HandleGetShowProductCharacteristics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))

		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	productcharacteristicData, err := productcharacteristics.Show(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao localizar os dados da característica. ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao localizar os dados da característica."))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsehelper.Response(true, productcharacteristicData, "Característica do produto localizada com sucesso!"))
}
