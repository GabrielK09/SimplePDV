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

export async function createshopping(payLoad: any): Promise<any>
{
    const formatedItens = payLoad.shopping_itens.map((i: ShoppingItemContract) => ({
        product_id: i.product_id,
        name: i.name,
        purchased_value: i.purchased_value,
        qtde_purchased: i.qtde_purchased,
        product_with_characteristics: i.product_with_characteristics === null ? null : i.product_with_characteristics[0]
    }));

    const newPayLoad = {
        id: payLoad.id,
        load: payLoad.load,
        shopping_itens: formatedItens,
        total_shopping: payLoad.total_shopping,
        status: payLoad.status
    };

    try {
        const res = await api.post('shopping/create', newPayLoad);
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

export async function getLastShoppingLoad(): Promise<any> 
{
    try {
        const res = await api.get('/shopping/return-last-load');
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

