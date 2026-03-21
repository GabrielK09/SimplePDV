import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function getAll(): Promise<any>
{
    try {
        const res = await api.get('/sale/all');
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data || []
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
            data.data || []
        );
    } catch (error) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data
        );
    };
};

export async function getSaleDetailsById(saleId: number): Promise<any>
{
    try {
        const res = await api.get(`/sale/details/${saleId}`);
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data || []
        );

    } catch (error) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data
        );
    };
};

export async function insertNewItens(payLoad: SaleContract) : Promise<any>
{
    try {
        const res = await api.put(`/sale/new-itens`, payLoad);
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data || []
        );

    } catch (error) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data
        );
    };
};
