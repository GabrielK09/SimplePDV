import { QrCodePix } from 'qrcode-pix';
import apiResponse from "src/helpers/response/apiResponse";

export async function generateQRCode(
    totalSale: number, 
    pixKey: string
): Promise<any>
{
    try {
        const qrCodePix = QrCodePix({
            version: '01',
            city: 'Concórida',
            key: pixKey,
            name: 'Julcineia Kochem',
            transactionId: 'JU0912',
            message: 'Volte sempre!',
            cep: '89711226',
            value: totalSale
        });

        const payLoad = qrCodePix.payload();
        const qrCode = await qrCodePix.base64();

        return {
            base64: qrCode,
            payLoad: payLoad
        };
        
    } catch (error) {
        return apiResponse(
            false,
            'Erro ao gerar QR Code do pagamento.',
            error
        );      
    };
};