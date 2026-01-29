export default function apiResponse(status: boolean, message: string, data: unknown[]|unknown)
{
    return {
        success: status,
        message: message,
        data: data
    };
};
