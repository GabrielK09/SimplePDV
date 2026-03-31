import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function getAll(perPage: number): Promise<any>
{
    try {
        const res = await api.get(`products/all${perPage > 0 ? `?per_page=${perPage}` : ''}`);
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

export async function createProductCharacteristics(payLoads: ProductCharacteristicsContract[], isUpdate?: boolean): Promise<any>
{
    try {
        let data;
        let res;

        for (let i = 0; i < payLoads.length; i++) {
            const payLoad = payLoads[i];


            if(!isUpdate)
            {
                res = await api.post('products/create/characteristics', payLoad);
            } else {
                res = await api.put('products/update/characteristics', payLoad);
            };

            data = res.data;
        };

        return apiResponse(
            true,
            data.message || 'Grade do produto cadastrada com sucesso!',
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

export async function getProductCharacteristicsById(productGridId: number): Promise<any>
{
    console.log('called getProductCharacteristicsById');

    try {
        const res = await api.get(`products/find/characteristics/${productGridId}`);

        const data = res.data;

        return apiResponse(
            true,
            data.message || 'Grade do produto localizada com sucesso!',
            data.data    || []
        );

    } catch (error) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response
        );
    };
};

export async function getProductCharacteristicsByGridIds({gridId, productGridId}): Promise<any>
{
    try {
        const res = await api.get(`products/${productGridId}/find/characteristics/${gridId}`);

        const data = res.data;

        return apiResponse(
            true,
            data.message || 'Grade do produto localizada com sucesso!',
            data.data    || []
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

export async function manageProductService(id: number, operation: 'active'|'delete'): Promise<any>
{
    try {
        const res = operation === 'delete' ? await api.delete(`products/${operation}/${id}`) : await api.patch(`products/${operation}/${id}`);

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

export async function deleteProductCharacteristics({gridId, productGridId}): Promise<any>
{
    try {
        const res = await api.delete(`products/${productGridId}/delete/characteristics/${gridId}`);
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

        return apiResponse(
            true,
            data.message || 'Dados do produto',
            data
        );

    } catch (error) {
        return null;

    };
};

export async function findByName(name: string): Promise<any>
{
    try {
        const res = await api.get(`products/find-by-name?name=${name}`);
        const data = res.data.data;

        return apiResponse(
            true,
            data.message || 'Dados do produto',
            data
        );

    } catch (error) {
        return null;

    };
};
