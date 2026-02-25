import { apiCnpj } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function getCnpjDataService(cnpj: string): Promise<any>
{
    try {    
        const res = await apiCnpj.get(`/${cnpj}`);
        const data = res.data;

        return apiResponse(
            true,
            'Dados do CNPJ.',
            data
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