import * as Yup from 'yup';

const parsePtBrNumber = (value: string | number | null | undefined): number => {
    if (value === null || value === undefined || value === "") return 0;

    if (typeof value === 'number')
    {
        return Number.isFinite(value) ? value : 0;
    };

    const normalized = value
        .trim()
        .replace(/\./g, '')
        .replace(',', '.');

    const parsed = Number(normalized);
    return Number.isFinite(parsed) ? parsed : 0;
};

export default function productSchema() {
    return Yup.object({
            name: Yup.string()
                .trim()
                .required('O nome do produto é obrigatório!'),

            price: Yup.number()
                .transform((_, originalValue) => parsePtBrNumber(originalValue))
                .min(1, 'O valor do produto não pode ser menor que zero.')
                .required('O valor do produto é obrigatório!'),

            qtde: Yup
                .number()
                .min(1, 'A qtde do produto não pode ser menor que zero.')
                .required('A quantia do produto é obrigatório!'),

            commission: Yup
                .number()
                .min(0, 'O valor de comissão não pode ser menor que zero.')
                .max(100, 'O valor de comissão não pode ser maior que 100%.')
                .required('A quantia do produto é obrigatório!')
        })
};