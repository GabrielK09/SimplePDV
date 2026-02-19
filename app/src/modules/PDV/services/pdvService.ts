import { api } from "src/boot/axios";
import formatValueToNumber from "src/helpers/FormatValue/FormatToNumber";
import apiResponse from "src/helpers/response/apiResponse";

export async function getAll(): Promise<any>
{
    try {
        const res = await api.get('/sale/all');
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

export async function paySaleService(payMentValues: PayMentValue[], saleId: number): Promise<any>
{
    try {
        const payLoad: PayMentPayLoadContract = {
            sale_id: saleId,
            species: payMentValues.map((f) => ({
                id: f.id,
                specie: f.specie,
                amount: formatValueToNumber(f.amount),
            }))
        };

        const res = await api.put('sale/pay', payLoad);
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data
        );
    } catch (error) {
        console.error('Erro: ', error);

        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data?.data
        );
    };
};
