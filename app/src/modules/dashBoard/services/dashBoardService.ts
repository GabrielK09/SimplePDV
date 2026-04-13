import { api } from "src/boot/axios";
import apiResponse from "src/helpers/response/apiResponse";

export async function getDashBoardData(startDate: string, endDate: string): Promise<any>
{
    try {
        const res = await api.post('/dash-board/totales', {
            start_date: startDate,
            end_date: endDate
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
            error.response
        );
    };
};

export async function filterPopularItensData(per_page: number): Promise<any>
{
    try {
        const res = await api.post('/dash-board/popular-itens', {
            per_page: per_page
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
            error.response
        );
    };
};
