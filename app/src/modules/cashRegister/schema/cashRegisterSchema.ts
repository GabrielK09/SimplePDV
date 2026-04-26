import * as Yup from 'yup';

export default function cashRegisterSchema()
{
    return Yup.object({
        description: Yup.string().required('A descrição do movimento financeiro é obrigatória!'),
        input_value: Yup.number().min(0, 'O valor de entrada não pode ser menor que zero.'),
        output_value: Yup.number().min(0, 'O valor de saída não pode ser menor que zero.')
    })
    .test(
        'input-or-output',
        'Informe apenas um dos valores, entrada ou saída',
        (obj) => {
            const inputValue = Number(obj.input_value ?? 0);
            const outputValue = Number(obj.output_value ?? 0);

            const haveInput = inputValue > 0;
            const haveOutput = outputValue > 0;

            return (!haveInput && haveOutput) || (haveInput && !haveOutput);
        }
    );
};