export function formatValueToMoney(val?: number) {
    if (!val) return 'R$ 0,00';

    return new Intl.NumberFormat('pt-BR', {
        style: 'currency',
        currency: 'BRL'
    }).format(val);
};

export function formatValueToNumber(val?: number) {
    if (!val) return '0';

    return new Intl.NumberFormat('pt-BR').format(val);
};