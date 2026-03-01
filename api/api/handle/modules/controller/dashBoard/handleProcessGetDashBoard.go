package controller

import (
	"encoding/json"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"myApi/interface/dashBoard"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool

func HandleProcessGetDashBoard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	loc, err := time.LoadLocation("America/Sao_Paulo")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao definir o timezone:", err)
		resp := responsehelper.Response(false, err, "Erro ao definir o timezone.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	layout := "2006-01-02"

	var dashBoardData dashBoard.DashBoardContract

	if err := json.NewDecoder(r.Body).Decode(&dashBoardData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao processar os dados:", err)
		resp := responsehelper.Response(false, err, "Erro ao processar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	startDate, err := time.ParseInLocation(layout, dashBoardData.StartDateStr, loc)

	if err != nil {
		u.ErrorLogger.Println("Erro ao processar o parse da data inicial:", err)
		w.WriteHeader(http.StatusInternalServerError)

		resp := responsehelper.Response(false, err, "Erro ao processar o parse da data inicial.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	endDate, err := time.ParseInLocation(layout, dashBoardData.EndDateStr, loc)

	if err != nil {
		u.ErrorLogger.Println("Erro ao processar o parse da data final:", err)
		w.WriteHeader(http.StatusInternalServerError)

		resp := responsehelper.Response(false, err, "Erro ao processar o parse da data final.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	dashBoardData.StartDate = startDate
	dashBoardData.EndDate = endDate

	u.GeneralLogger.Printf("Dados para geração do dashboard, inicio: %s - fim %s: ", dashBoardData.StartDate, dashBoardData.EndDate)

	dashBoardData, err = dashBoard.ShowDashBoard(dashBoardData.StartDate, dashBoardData.EndDate)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := responsehelper.Response(false, nil, "Erro ao processar os dados do DashBoard.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, dashBoardData, "Dados do período sucesso!")

	json.NewEncoder(w).Encode(resp)
}
