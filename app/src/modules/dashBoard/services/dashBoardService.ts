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

    } catch (error) {
        return apiResponse(
            false,
            error.response?.data?.message,
            error.response
        );
    };
};
