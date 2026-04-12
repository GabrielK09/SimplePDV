import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function getAll(): Promise<any>
{
    try {
        const res = await api.get('/shopping/all');
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
            error.response
        );
    };
};

export async function createshopping(payLoad: ShoppingContract): Promise<any>
{
    try {
        const res = await api.post('shopping/create', payLoad);
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
            error.response
        );
    };
};

export async function cancelShopping(shoppingId: number): Promise<any>
{
    try {
        const res = await api.put(`shopping/delete/${shoppingId}`);
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
            error.response
        );
    };
};

export async function getShoppingById(id: number): Promise<any> 
{
    try {
        const res = await api.get(`/shopping/details/${id}`);
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data || 0
        );
    } catch (error: any) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response
        );
    };
};

export async function updateShoppingDetails(payLoad: ShoppingContract): Promise<any> 
{
    try {
        const res = await api.put(`/shopping/update`, payLoad);
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data || 0
        );
    } catch (error: any) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response
        );
    };
};

