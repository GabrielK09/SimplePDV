import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function loginService(payLoad: AuthContract): Promise<any>
{
    try {
        const res = await api.post("/auth/login", payLoad);
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