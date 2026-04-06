import formatValueToNumber from "src/helpers/FormatValue/FormatToNumber";
import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function getAllPayMentFormsService(): Promise<any>
{
    try {
        const res = await api.get('/pay-ment-forms/all');
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data || []
        );

    } catch (error: any) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data
        );
    };
};

export async function updatePayMentFormService(payLoad: string): Promise<any>
{
    try {
        const res = await api.put('/pay-ment-forms/update/pix-key', {
            pix_key: payLoad
        });

        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data || []
        );

    } catch (error: any) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data
        );
    };
};


export async function payMentService(payMentValues: PayMentValue[], saleId: number, shoppingId: number): Promise<any>
{
    try {
        if (saleId <= 0 && shoppingId <= 0) return;

        const payLoad: PayMentPayLoadContract = {
            sale_id: saleId,
            shopping_id: shoppingId,
            species: payMentValues.map((f) => ({
                id: f.id,
                specie: f.specie,
                amount: formatValueToNumber(f.amount),
            }))
        };

        const res = await api.put('/pay-ment-forms/pay', payLoad);
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data || []
        );
    } catch (error: any) {
        console.error('Erro: ', error);

        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data?.data
        );
    };
};
