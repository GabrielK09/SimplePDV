import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function generateReport(payLoad: ReportContract, fileName: string)
{
    try {
        const res = await api.post('/report', payLoad, {
            responseType: 'blob'
        });

        const blob = res.data;

        const url = window.URL.createObjectURL(blob); // Cria a URL com o arquivo

        const link = document.createElement('a'); // Cria um elemento <a>

        link.href = url; // Atribui a url retornada
        link.download = `${fileName}.pdf`; // Atribui a propriedade de download
        document.body.appendChild(link); // Adiciona ao elemento

        link.click();  // Simula o click

        document.body.removeChild(link); // Remove o elemento
        window.URL.revokeObjectURL(blob); // Remove a URL do arquivo

        return {
            success: true
        };
    } catch (error) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response
        );
    };
};

async function downloadReport() {
  const response = await fetch('/reports', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      report_type: 'cash-register',
      start_date: '2026-03-01',
      end_date: '2026-03-08'
    })
  })

  const blob = await response.blob()

  const url = window.URL.createObjectURL(blob)

  const a = document.createElement('a')
  a.href = url
  a.download = 'report.pdf'
  a.click()

  window.URL.revokeObjectURL(url)
}
