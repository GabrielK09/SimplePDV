package controller

import (
	"encoding/json"
	reports "myApi/api/handle/modules/controller/report/internal"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"net/http"
	"time"
)

type ReportBody struct {
	ReportType   string `json:"report_type"`
	StartDateStr string `json:"start_date"`
	EndDateStr   string `json:"end_date"`

	StartDate time.Time `json:"-"`
	EndDate   time.Time `json:"-"`
}

const (
	cashRegister = "cash-register"
	payMentForms = "pay-ment-f orms"
	layout       = "2006-01-02"
)

func HandlePostReports(w http.ResponseWriter, r *http.Request) {
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
	} // Erro de método da rota

	var report ReportBody

	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao processar os dados:", err)

		resp := responsehelper.Response(false, err, "Erro ao processar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	startDate, err := time.ParseInLocation(layout, report.StartDateStr, loc)

	if err != nil {
		u.ErrorLogger.Println("Erro ao processar o parse da data inicial:", err)
		w.WriteHeader(http.StatusInternalServerError)

		resp := responsehelper.Response(false, err, "Erro ao processar o parse da data inicial.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	endDate, err := time.ParseInLocation(layout, report.EndDateStr, loc)

	if err != nil {
		u.ErrorLogger.Println("Erro ao processar o parse da data final:", err)
		w.WriteHeader(http.StatusInternalServerError)

		resp := responsehelper.Response(false, err, "Erro ao processar o parse da data final.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	report.StartDate = startDate
	report.EndDate = endDate

	if report.ReportType != cashRegister && report.ReportType != payMentForms {
		u.ErrorLogger.Println("Tipo de relatório incorreto.")
		w.WriteHeader(http.StatusNotFound)

		resp := responsehelper.Response(false, nil, "Tipo de relatório incorreto.")

		json.NewEncoder(w).Encode(resp)
		return

	} else {
		_ = reports.BuildReport(report.ReportType)

		w.WriteHeader(http.StatusOK)

		resp := responsehelper.Response(true, report, "Dados do relatório.")

		json.NewEncoder(w).Encode(resp)
		return
	}
}
