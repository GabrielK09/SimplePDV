export default function apiResponse(status: boolean, message: string, data: unknown[]|unknown) 
{
    return {
        status: status,
        message: message,
        data: data
    };
};