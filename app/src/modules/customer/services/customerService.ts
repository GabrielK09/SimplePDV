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

    } catch (error) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data
        );
    };
};
