export default function formatValueToNumber(val?: string) {
    const formatVal = val.replace(',', '.');

    return Number(formatVal);
};