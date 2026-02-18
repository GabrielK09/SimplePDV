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

export async function paySaleService(payMent: PayMentValue[], saleId: number): Promise<any>
{
    try {
        let res: any;

        for (let i = 0; i < payMent.length; i++) {
            const p = payMent[i];

            const payLoad: PayMentPayLoadContract = {
                amount_paid: formatValueToNumber(p.amount),
                sale_id: saleId,
                specie: p.specie
            };

            console.log(payLoad);

            res = await api.put('/sale/pay', payLoad);

            console.log('Dados: ', {
                res: res,

            });
        };

        /*return apiResponse(
            true,
            data.message,
            data.data
        )*/;
    } catch (error) {
        console.error('Erro: ', error);

        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data?.data
        );
    };
};
