<template>
    <q-page padding class="bg-grey-2">
        <div class="q-gutter-lg">
            <q-card class="q-pa-md mr-6">
                <div class="row q-col-gutter-md items-end">
                    <div class="col-12 col-md-3">
                        <q-input 
                            v-model="filterDashBoardValues.startDate" 
                            type="date" 
                            label="Data inicial" 
                            outlined
                            min="0001-01-01"
                            max="9999-12-31"
                            hide-bottom-space
                            bottom-slots
                            :error="!!formErrors.filterDashBoardValues"
                            :error-message="formErrors.filterDashBoardValues"
                        />
                    </div>

                    <div class="col-12 col-md-3">
                        <q-input 
                            v-model="filterDashBoardValues.endDate" 
                            type="date" 
                            label="Data final" 
                            outlined 
                            min="0001-01-01"
                            max="9999-12-31"
                            :rules="[
                                val => validateYearDate(filterDashBoardValues.startDate, val) || 'A data final deve ser maior ou igual à inicial.'
                            ]"
                            hide-bottom-space
                            bottom-slots
                            class="my-auto"
                            :error="!!formErrors.filterDashBoardValues"
                            :error-message="formErrors.filterDashBoardValues"
                        />
                    </div>
                    
                    <div class="my-auto">
                        <q-btn
                            @click="filterDashBoard"
                            no-caps
                            label="Filtrar"
                            class="full-width bg-[#2563EB] text-white"
                        />
                    </div>

                    <div class="my-auto">
                        <q-btn
                            no-caps
                            class="full-width bg-[#1D4ED8] text-white"
                            label="Relatórios"
                            @click="showReports = !showReports"
                        />
                    </div>
                </div>
            </q-card>

            <div class="row q-col-gutter-md px-4 ml-2 mr-6">
                <div
                    v-for="card in dashBoardCards"
                    :key="card.label"
                    class="col-12 col-sm-4 col-md-3"
                >
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-gray-400">{{ card.label }}</div>

                        <div :class="card.color">
                            {{ card.value }}
                        </div>
                    </q-card>
                </div>

                <div class="col-12 col-sm-6 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-gray-400">Saldo total do caixa</div>

                        <div class="flex justify-between items-center">
                            <span
                                :class="{
                                'text-[#16A34A]': totalBalance > 0,
                                'text-[#DC2626]': totalBalance < 0
                                }"
                            >
                                {{ formatValueToMoney(totalBalance) }}
                            </span>
                            
                            <q-icon 
                                name="info"
                                class="text-blue-500"
                            
                            >
                                <q-tooltip >
                                    Valor independente de filtros!
                                </q-tooltip>
                            </q-icon>
                        </div>
                    </q-card>
                </div>
            </div>

            <div class="mr-6">
                <q-card class="q-pa-md">
                    <q-table
                        title="Itens mais vendidos"
                        :rows="popularItensTableData"
                        :columns="popularItensTableColumn"
                        row-key="product_id"
                    />
                </q-card>
            </div>
        </div>
    </q-page>

    <ManagementReport
        v-if="showReports"
        @close="showReports = !$event"
    />
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { computed, onMounted, ref } from 'vue';
    import { filterPopularItensData, getDashBoardData } from '../services/dashBoardService';
    import { getAll } from 'src/modules/cashRegister/services/cashRegisterService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import ManagementReport from 'src/components/Reports/ManagementReport.vue';
    import { formatValueToMoney, formatValueToNumber } from 'src/helpers/FormatValue/FormatMoney';
    import reportFilterSchema from '../schema/reportFilterSchema';
    
    type ReportFilterDashBoard = {
        startDate: string;
        endDate: string;
    };

    const { notify } = useNotify();

    const filterDashBoardValues = ref<ReportFilterDashBoard>({
        startDate: '',
        endDate: ''
    });

    const formErrors = ref<Record<string, string>>({});

    const dashBoardData = ref<DashBoardContract>();
    const totalBalance = ref<number>(0);
    const showReports = ref<boolean>(false);

    const dashBoardCards = computed(() => [
        {
            label: 'Total vendido',
            value: formatValueToMoney(dashBoardData.value?.total_saled),
            color: 'text-[#2563EB]'
        },
        {
            label: 'Comissão',
            value: formatValueToMoney(dashBoardData.value?.commission),
            color: 'text-[#16A34A]'
        },
        {
            label: 'Quantidade de vendas',
            value: formatValueToNumber(dashBoardData.value?.amount_saled),
            color: 'text-[#7C3AED]'
        },
        {
            label: 'Melhor cliente',
            value: dashBoardData.value?.best_customer || '-',
            color: 'text-[#111827]'
        },
        {
            label: 'Total comprado',
            value: formatValueToMoney(dashBoardData.value?.total_shopping),
            color: 'text-[#3B82F6]'
        },
        {
            label: 'Total de itens comprados',
            value: formatValueToNumber(dashBoardData.value?.amount_shopping_itens),
            color: 'text-[#3B82F6]'
        },
        {
            label: 'Quantidade de compras',
            value: formatValueToNumber(dashBoardData.value?.amount_shopping),
            color: 'text-blue-800'
        }
    ]);

    const popularItensTableColumn: QTableColumn[] = [
        {
            name: 'produto_id',
            label: 'Cód produto',
            field: 'produto_id',
            align: 'left',
        },
        {
            name: 'produto',
            label: 'Produto',
            field: 'produto',
            align: 'left',
        },
        {
            name: 'item_sale_value',
            label: 'Valor do produto',
            field: 'item_sale_value',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }

        },
        {
            name: 'qtde',
            label: 'Qtde',
            field: 'qtde',
            align: 'center',
        },
    ];

    const popularItensTableData = ref<any[]>([]);

    const filterPopularItens = async () => {
        const res = await filterPopularItensData(20);

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
        };

        popularItensTableData.value = res.data;
    };

    const filterDashBoard = async () => {
        try {
            await reportFilterSchema().validate(filterDashBoardValues.value, { abortEarly: false });
            
            const res = await getDashBoardData(filterDashBoardValues.value.startDate, filterDashBoardValues.value.startDate);

            if(!res.success)
            {
            notify(
                    'negative',
                    res.message
                );
            };

            dashBoardData.value = res.data;
            
        } catch (error: any) {
            console.error('Erro:', error.inner);

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
                    error.response?.data?.message || 'Erro na criação do cliente!'
                );
            };

        };
    };

    const validateYearDate = (startDate: string, endDate: string): boolean => {
        const startYear = startDate.split('-')[0];
        const endYear = endDate.split('-')[0];

        if (endYear > startYear || startYear < endYear) return true;

        return false;
    };

    onMounted(async() => {
        await filterPopularItens();

        const res = await getAll();

        if(!res.success) return;
        const cashRegisterData: CashRegisterContract[] = res.data;

        totalBalance.value = cashRegisterData.reduce((total, a) => total + (a.input_value - a.output_value), 0);
    });
</script>
