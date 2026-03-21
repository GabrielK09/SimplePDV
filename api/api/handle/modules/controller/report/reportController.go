package reportController

import (
	"encoding/json"

	reportservices "myApi/api/services"
	reportsdata "myApi/api/services/reports"
	u "myApi/helpers/logger"
	responsehelper "myApi/helpers/response"
	"net/http"
)

const (
	cashRegister = "cash-register"
	payMentForms = "pay-ment-forms"
	saledItens   = "saled-itens"
	layout       = "2006-01-02"
)

func HandlePostReports(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := responsehelper.Response(false, nil, "Método não permetido.")

		json.NewEncoder(w).Encode(resp)
		return
	} // Erro de método da rota

	var report reportsdata.ReportBody

	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		u.ErrorLogger.Println("Erro ao processar os dados:", err)

		resp := responsehelper.Response(false, err, "Erro ao processar os dados.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	u.InfoLogger.Println("Valor de report:", report)

	if report.ReportType != cashRegister && report.ReportType != payMentForms && report.ReportType != saledItens {
		u.ErrorLogger.Println("Tipo de relatório incorreto.")
		u.ErrorLogger.Println("Tipo recebido: ", report.ReportType)
		u.ErrorLogger.Printf("Tipos esperados: %s, %s, %s", cashRegister, payMentForms, saledItens)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		resp := responsehelper.Response(false, nil, "Tipo de relatório incorreto.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	data, err := report.BuildDataReport()

	if err != nil {
		u.ErrorLogger.Println("Erro ao processar os dados do relatório: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		resp := responsehelper.Response(false, err, "Erro ao processar o relatório.")

		json.NewEncoder(w).Encode(resp)
		return
	}

	u.InfoLogger.Println("Dados a serem carregados: ", data)

	pdfBytes, err := reportservices.CreateReport(data)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		resp := responsehelper.Response(false, err, "Erro ao gerar o relatório.")
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=report.pdf")

	w.Write(pdfBytes)
}
