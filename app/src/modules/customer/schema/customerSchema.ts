import * as Yup from 'yup';

export default function customerSchema() {
    return Yup.object({
        name: Yup.string().required('O nome do produto é obrigatório!'),
    });
};