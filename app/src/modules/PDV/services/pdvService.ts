import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";
import { QRCodePix } from ''

export async function getAll(): Promise<any>
{

};

export async function saveSaleService(payLoad: SaleContract): Promise<any>
{
    try {
        const res = await api.post('/sale/create', payLoad);
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data
        );
    } catch (error) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data
        );
    };
};

export async function generateQRCode(saleId: number): Promise<any>
{
    try {
        //const res = await 
        
    } catch (error) {
        return apiResponse(
            false,
            'Erro ao gerar QR Code do pagamento.',
            []
        );      
    };
};