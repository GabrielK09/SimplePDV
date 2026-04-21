import * as Yup from 'yup';

export default function productSchema() {
    return Yup.object({
        name: Yup
            .string()
            .required('O nome do produto é obrigatório!'),

        price: Yup
            .number()
            .typeError('O valor do produto precisa ser um número!')
            .required('O valor do produto é obrigatório!'),

        use_grid: Yup
                .boolean(),

        qtde: Yup
            .number()
            .typeError('O valor do produto precisa ser um número!')
            .when('use_grid', {
                is: false,
                then: (schema) => 
                    schema
                        .min(1, 'A qtde do produto não pode ser menor que zero.')
                        .required('A quantia do produto é obrigatório!'),
                        
                    otherwise: (schema) => 
                        schema.notRequired()  
            }),
    
        commission: Yup
            .number()
            .typeError('O valor do produto precisa ser um número!')
            .min(0, 'O valor de comissão não pode ser menor que zero.')
            .max(100, 'O valor de comissão não pode ser maior que 100%.')
    })
};