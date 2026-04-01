package productcontroller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/product"
	_ "myApi/interface/product/productCharacteristics"
	productcharacteristics "myApi/interface/product/productCharacteristics"

	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5"
)

type ProductWithCharacteristics struct {
	Product         product.ProductContract                                 `json:"product"`
	Characteristics []productcharacteristics.ProductCharacteristicsContract `json:"characteristics"`
}

func HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var perPage interface{}

	//perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	perPageStr := r.URL.Query().Get("per_page")

	u.InfoLogger.Println("perPageStr: ", perPageStr)

	if perPageStr != "" {
		perPage = "all"
	}

	if v, err := strconv.Atoi(perPageStr); err == nil {
		perPage = v
	}

	u.InfoLogger.Println("per_page: ", perPage)

	// params := mux.Vars(r)
	// productId, err := strconv.Atoi(params["id"])

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	var productsWithCharacteristics []ProductWithCharacteristics

	products, err := product.GetAll(perPage)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao retornar todos os produtos."))
		return
	}

	productCharacteristics, err := productcharacteristics.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao retornar todos os produtos."))
		return
	}

	characteristicsMap := make(map[int][]productcharacteristics.ProductCharacteristicsContract)

	for _, c := range productCharacteristics {
		characteristicsMap[c.ProductId] = append(characteristicsMap[c.ProductId], c)
	}

	for _, p := range products {
		productsWithCharacteristics = append(productsWithCharacteristics, ProductWithCharacteristics{
			Product:         p,
			Characteristics: characteristicsMap[p.Id],
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responsehelper.Response(true, productsWithCharacteristics, "Todos os produtos cadastrados."))
}

func HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
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

	if err := product.Delete(id, time.Now()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao deletar o produto."),
		)

		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsehelper.Response(true, nil, "Produto deletado com sucesso!"))
}

func HandleActiveProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPatch {
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

	if err := product.Active(id, time.Now()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao ativar o produto."),
		)

		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsehelper.Response(true, nil, "Produto ativado com sucesso!"))
}

func HandleGetByNameProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	u.InfoLogger.Println("Name - recebido: ", r.URL.Query().Get("name"))

	var productsWithCharacteristics []ProductWithCharacteristics
	productData, err := product.ShowByName(r.URL.Query().Get("name"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao localizar o produto pelo nome: ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao retornar o produtos."))
		return
	}

	productCharacteristics, err := productcharacteristics.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao retornar todos os produtos."))
		return
	}

	characteristicsMap := make(map[int][]productcharacteristics.ProductCharacteristicsContract)

	for _, c := range productCharacteristics {
		characteristicsMap[c.ProductId] = append(characteristicsMap[c.ProductId], c)
	}

	for _, p := range productData {
		productsWithCharacteristics = append(productsWithCharacteristics, ProductWithCharacteristics{
			Product:         p,
			Characteristics: characteristicsMap[p.Id],
		})
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsehelper.Response(true, productsWithCharacteristics, "Produto"))
}

func HandleGetByIdProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var productsWithCharacteristics ProductWithCharacteristics

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))
		return
	}

	params := mux.Vars(r)
	productId, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	if productId <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Id inválido: ", err)
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, nil, "Id inválido."),
		)

		return
	}

	productData, err := product.Show(productId)

	productsWithCharacteristics = ProductWithCharacteristics{
		Product:         *productData,
		Characteristics: nil,
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao retornar o produtos."))
		return
	}

	if productData.UseGrid {
		productCharacteristicsData, err := productcharacteristics.Show(productId)

		u.InfoLogger.Println("productCharacteristics: ", productCharacteristicsData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao retornar o produtos."))
			return
		}

		productsWithCharacteristics = ProductWithCharacteristics{
			Product:         *productData,
			Characteristics: *productCharacteristicsData,
		}
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsehelper.Response(true, productsWithCharacteristics, "Produto"))
}

func HandlePostProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(responsehelper.Response(false, nil, "Método não permetido."))
		return
	}

	var payload product.ProductContract

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.ErrorLogger.Println("Erro ao processar os dados: ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao processar os dados."))
		return
	}

	if err := payload.Validate(); len(err) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		u.ErrorLogger.Println("Campos obrigatórios ausentes.", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Campos obrigatórios ausentes."))
		return
	}

	id, err := payload.Create()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao gravar o produto: ", err)

		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Erro ao gravar o produto."))
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(responsehelper.Response(true, id, "Produto cadastrado com sucesso!"))
}

func HandlePutProduct(w http.ResponseWriter, r *http.Request) {
	var productData *product.ProductContract

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
		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err, "Id inválido."),
		)

		return
	}

	productData, err = product.Show(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		u.ErrorLogger.Println("Produto não localizado:", err)
		json.NewEncoder(w).Encode(responsehelper.Response(false, err, "Produto não localizado."))

		return
	}

	u.InfoLogger.Println("Produto:", productData)

	if err := json.NewDecoder(r.Body).Decode(&productData); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		u.ErrorLogger.Println("Erro ao processar dados: ", err)

		json.NewEncoder(w).Encode(
			responsehelper.Response(false, err.Error(), "Erro ao processar dados"),
		)

		return
	}

	productData.Id = id

	if err := productData.Update(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		u.ErrorLogger.Println("Erro ao alterar o produto: ", err)
		json.NewEncoder(w).Encode(responsehelper.Response(false, err.Error(), "Erro ao atualizar"))

		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsehelper.Response(true, nil, "Produto alterado com sucesso!"))
}
