import * as Yup from 'yup';

export default function reportFilterSchema() 
{
    return Yup.object({
        startDate: Yup
            .string()
            .required('A data de inicio é obrigatória.'),
            
        endDate: Yup
            .string()
            .required('A data de fim é obrigatória.'),
    });
};