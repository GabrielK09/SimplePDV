import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function getAllPayMentFormsService(): Promise<any>
{
    try {
        const res = await api.get('/pay-ment-forms/pay');
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

export async function updatePayMentFormService(payLoad: string): Promise<any>
{
    try {
        const res = await api.put('/pay-ment-forms/update/pix-key', {
            pix_key: payLoad
        });

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

export async function cancelShoppingOrSale(payLoad: CancelContract): Promise<any> 
{
    try {
        const res = await api.put('/cancel-operation', payLoad);
        const data = res.data;

        return apiResponse(
            true,
            data.message,
            data.data || []
        );            
    } catch (error: any) {
        return false;      
    };
};