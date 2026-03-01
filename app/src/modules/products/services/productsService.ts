import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function getAll(): Promise<any>
{
    try {
        const res = await api.get('products/all');
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
            error.response
        );
    };
};

export async function createProduct(payLoad: ProductContract): Promise<any>
{
    try {
        const res = await api.post('products/create', payLoad);
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
            error.response
        );
    };
};

export async function updateProduct(payLoad: ProductContract): Promise<any>
{
    try {
        const res = await api.put(`products/update/${payLoad.id}`, payLoad);
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
            error.response
        );
    };
};

export async function deleteProduct(id: number): Promise<any>
{
    try {
        const res = await api.delete(`products/delete/${id}`)
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data || []
        );

    } catch (error) {
        console.error(error.response.data);

        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data
        );
    };
};

export async function findById(id: number): Promise<any>
{
    try {
        const res = await api.get(`products/find/${id}`);
        const data = res.data.data;

        console.log(data);

        return apiResponse(
            true,
            data.message || 'Dados do produto',
            data
        );

    } catch (error) {
        return null;

    };
};
