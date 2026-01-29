import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function getAll(): Promise<any>
{
    try {
        const res = await api.get('shopping/all');
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
            error.response
        );
    };
};

export async function createshopping(payLoad: any): Promise<any>
{
    try {
        const res = await api.post('shopping/create', payLoad);
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
            error.response
        );
    };
};

export async function deleteshopping(shoppingId: number): Promise<any>
{
    try {
        const res = await api.delete(`shopping/delete/${shoppingId}`);
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
            error.response
        );
    };
};
