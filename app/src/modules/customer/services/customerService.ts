import apiResponse from "src/helpers/response/apiResponse";
import { api } from "src/boot/axios";

export async function getAll(): Promise<any>
{
    try {
        const res = await api.get('/customer/all');
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

export async function createCustomer(payLoad: CustomerContract): Promise<any>
{
    try {
        const res = await api.post(`/customer/create`, payLoad);
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

/*export async function deleteCustomer(customerId: number): Promise<any>
{
    try {
        const res = await api.delete(`/customer/delete/${customerId}`);
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
};*/

export async function manageCustomerService(id: number, operation: 'active'|'delete'): Promise<any>
{
    try {
        const label = operation === 'delete' ? 'desativado' : 'ativado';
        const res = operation === 'delete' ? await api.delete(`customer/${operation}/${id}`) : await api.patch(`customer/${operation}/${id}`);

        const data = res.data;

        return apiResponse(
            true,
            `Cliente ${label} com sucesso!`,
            data.data || []
        );

    } catch (error: any) {
        console.error(error.response.data);

        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data
        );
    };
};

export async function findCustomerById(customerId: number): Promise<any> 
{
    try {
        const res = await api.get(`/customer/find/${customerId}`);
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

export async function updateCustomerById(customerId: number, payLoad: CustomerContract): Promise<any> 
{
    try {
        const res = await api.put(`/customer/update/${customerId}`, payLoad);
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