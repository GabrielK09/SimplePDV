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
        
        if(data.id && payLoad.qtde)
        {   
            const formatedCharacteristics = payLoad.productWithCharacteristics.map(p => ({
                id: null,
                product_id: data.id,
                grid_qtde: p.grid_qtde,
                size: p.size
            }));

            createProductCharacteristics(formatedCharacteristics)
        };
    
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

export async function createProductCharacteristics(payLoad: ProductCharacteristicsContract[]): Promise<any>
{
    try {
        const res = await api.post('products/create', payLoad);
        const data = res.data;

        return apiResponse(
            true,
            data.message || 'Produto cadastrado com sucesso!',
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
