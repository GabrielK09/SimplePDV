export default function formatNumber(val?: number) {
    if (!val) return '0';

    return new Intl.NumberFormat('pt-BR').format(val);
};