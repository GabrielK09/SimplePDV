export default function formatValueToNumber(str: string): number
{
    return Number(str.replace(',', '.'));
};
