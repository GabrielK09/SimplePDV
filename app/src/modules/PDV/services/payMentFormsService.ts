import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function getAllPayMentFormsService(): Promise<any>
{
    try {
        const res = await api.get('/sale/pay-ment-forms');
        const data = res.data;

        return apiResponse(
            false,
            data.message,
            data.data
        );
        
    } catch (error) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data
        );  
    };  
};

export async function updatePayMentFormService(payLoad: PayMentFormContract): Promise<any> 
{
    try {
        const res = await api.put(`/sale/update/pay-ment-forms`, payLoad);
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data
        );
        
    } catch (error) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response?.data
        );  
    };  
};