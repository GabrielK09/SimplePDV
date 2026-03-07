package reportController

import (
	"encoding/json"

	reportservices "myApi/api/services"
	reportsdata "myApi/api/services/reports"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"net/http"
	"time"
)

const (
	cashRegister = "cash-register"
	payMentForms = "pay-ment-forms"
	saledItens   = "saled-itens"
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

	var report reportsdata.ReportBody

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

	if report.ReportType != cashRegister && report.ReportType != payMentForms && report.ReportType != saledItens {
		u.ErrorLogger.Println("Tipo de relatório incorreto.")
		u.ErrorLogger.Println("Tipo recebido: ", report.ReportType)
		u.ErrorLogger.Printf("Tipos esperados: %s, %s, %s", cashRegister, payMentForms, saledItens)

		w.WriteHeader(http.StatusNotFound)

		resp := responsehelper.Response(false, nil, "Tipo de relatório incorreto.")

		json.NewEncoder(w).Encode(resp)
		return

	} else {
		data, err := report.BuildDataReport()

		if err != nil {
			u.ErrorLogger.Println("Erro ao processar o relatório: ", err)
			w.WriteHeader(http.StatusInternalServerError)

			resp := responsehelper.Response(false, err, "Erro ao processar o relatório.")

			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)

		_ = reportservices.CreatePDFMaroto(data)

		resp := responsehelper.Response(true, data, "Dados do relatório.")

		json.NewEncoder(w).Encode(resp)
		return
	}
}

/*
func GeneratePDFReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=Relatório.pdf")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	filePath, err := reportservices.CreateReport()

	if err != nil {
		u.ErrorLogger.Println("Erro ao gerar o relatório: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		resp := responsehelper.Response(false, err, "Erro ao gerar o relatório.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := responsehelper.Response(true, filePath, "Dados do relatório.")

	json.NewEncoder(w).Encode(resp)
}
*/
