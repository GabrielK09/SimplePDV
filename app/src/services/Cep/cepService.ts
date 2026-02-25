import { apiCep } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function getCepData(cep: string): Promise<any>
{
try {
        const res = await apiCep.get(`/${cep.replace(/\D/, '')}/json`);
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