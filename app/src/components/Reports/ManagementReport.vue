<template>
    <q-dialog v-model="internalDialog" persistent>
        <div class="flex flex-col bg-white w-full phone:bg-black shadow-lg">
            <header class="text-gray-600 text-center">
                <h3 class="text-xl mt-4">Relatórios</h3>
            </header>

            <div class="bg-white p-4">
                <BaseReportTypeSelect
                    v-model="reportData.report_type"
                    :error="!!formErrors.report_type"
                    :error-message="formErrors.report_type"
                />

                <div class="my-4">
                    <q-input
                        v-model="reportData.start_date"
                        type="date"
                        label="Data inicial"
                        outlined
                        :error="!!formErrors.start_date"
                        :error-message="formErrors.start_date"
                    />
                </div>

                <div>
                    <q-input
                        v-model="reportData.end_date"
                        type="date"
                        label="Data final"
                        outlined
                        :error="!!formErrors.end_date"
                        :error-message="formErrors.end_date"
                    />
                </div>
            </div>

            <div class="flex flex-center mb-4">
                <q-btn
                    color="red"
                    icon="close"
                    no-caps
                    @click="emits('close', true)"
                    class="mr-4"
                />

                <q-btn
                    color="primary"
                    label="Gerar relatório"
                    outline
                    no-caps
                    @click="submitReport"
                />
            </div>
        </div>
    </q-dialog>
</template>

<script setup lang="ts">
    import { computed, ref } from 'vue';
    import BaseReportTypeSelect from 'src/components/Qselects/QSelectReporType/BaseReportTypeSelect.vue';
    import * as Yup from 'yup';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { generateReport } from 'src/services/Report/reportService';

    enum ReportTypes {
        CashRegister = "cash-register",
        PayMentForms = "pay-ment-forms",
        SaledItens = "saled-itens"
    };

    const emits = defineEmits<{
        (e: 'close', value: boolean)
    }>();

    const internalDialog = ref<boolean>(true);
    const formErrors = ref<Record<string, string>>({});
    const { notify } = useNotify();

    const reportSchema = computed(() =>
        Yup.object({
            report_type: Yup
                .mixed<ReportTypes>()
                .oneOf(Object.values(ReportTypes))
                .required('O tipo do relatório é obrigatório.'),

            start_date: Yup
                .string()
                .required('A data inicial é obrigatória.'),

            end_date: Yup
                .string()
                .required('A data final é obrigatória.')
        })
    );

    const reportData = ref<ReportContract>({
        report_type: null,
        start_date: null,
        end_date: null,
    });

    const submitReport = async (): Promise<any> => {
        try {
            let fileName: string;
            await reportSchema.value.validate(reportData.value, { abortEarly: false });

            switch (reportData.value.report_type) {
                case "cash-register":
                    fileName = "Relatório_caixa.pdf";
                    break;

                case "pay-ment-forms":
                    fileName = "Relatório_formas_de_pagamento.pdf";

                    break;

                case "saled-itens":
                    fileName = "Relatório_itens_vendidos.pdf"
                    break;

                default:
                    fileName = "Relatório.pdf";
                    break;
            }

            const res = await generateReport(reportData.value, fileName);

            if (!res.success)
            {
                notify(
                    'negative',
                    'Erro ao baixar o arquivo'

                );
            };

            console.log(res);

        } catch (error) {
            if(error.inner)
            {
                formErrors.value = {};

                error.inner.forEach((err: any) => {
                    formErrors.value[err.path] = err.message;

                    notify(
                        'negative',
                        err.message

                    );
                });
            } else {
                notify(
                    'negative',
                    error.response?.data?.message || 'Erro na geração do relatório!'
                );
            };
        };
    };
</script>
