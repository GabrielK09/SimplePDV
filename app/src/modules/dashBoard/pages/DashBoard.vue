<template>
    <q-page padding class="bg-grey-2">
        <div class="q-gutter-lg">
            <q-card class="q-pa-md mr-6">
                <div class="row q-col-gutter-md items-end">
                    <div class="col-12 col-md-3">
                        <q-input v-model="startDate" type="date" label="Data inicial" outlined />
                    </div>

                    <div class="col-12 col-md-3">
                        <q-input v-model="endDate" type="date" label="Data final" outlined />
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

                    <div class="my-auto">
                        <q-btn 
                            no-caps
                            color="primary" 
                            label="Resumo de qtdes"
                            class="full-width bg-[#1b2747] text-white"
                            @click="showResumeQtde = !showResumeQtde"
                        />
                    </div>
                </div>
            </q-card>

            <div class="row q-col-gutter-md ml-2 mr-6">
                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Total vendido</div>
                        <div class="text-h5 text-[#2563EB]">
                            R$ {{ dashBoardData?.total_saled.toFixed(2).toString().replace('.', ',') || '0,00' }}
                        </div>
                    </q-card>
                </div>

                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Comissão</div>
                        <div class="text-h5 text-[#16A34A]">
                            R$ {{ dashBoardData?.commission.toFixed(2).toString().replace('.', ',') || '0,00' }}
                        </div>
                    </q-card>
                </div>

                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Quantidade de vendas</div>
                        <div class="text-h5 text-[#7C3AED]">
                            {{ dashBoardData?.amount_saled.toString().replace('.', ',') || '0' }}
                        </div>
                    </q-card>
                </div>

                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Melhor cliente</div>
                        <div class="text-h6 text-[#111827]">
                            <span class="text-[#9CA3AF]">
                                {{ dashBoardData?.best_customer || '-' }}
                            </span>
                        </div>
                    </q-card>
                </div>

                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Total comprado</div>
                        <div class="text-h5 text-[#3B82F6]">
                            R$ {{ dashBoardData?.total_shopping.toFixed(2).toString().replace('.', ',') || '0,00' }}
                        </div>
                    </q-card>
                </div>

                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Total de itens comprados</div>
                        <div class="text-h5 text-[#3B82F6]">
                           {{ dashBoardData?.amount_shopping_itens.toString().replace('.', ',') || '0' }}
                        </div>
                    </q-card>
                </div>

                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Quantidade de compras</div>
                        <div class="text-h5 text-blue-800">
                            {{ dashBoardData?.amount_shopping.toString().replace('.', ',') || '0' }}
                        </div>
                    </q-card>
                </div>

                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Saldo total do caixa</div>

                        <div class="flex justify-between">
                            <div class="text-h5 text-primary">
                                <span
                                    :class="{
                                        'text-[#16A34A]': totalBalance > 0,
                                        'text-[#DC2626]': totalBalance < 0
                                    }"
                                >
                                    R$ {{ totalBalance.toFixed(2).toString().replace('.', ',') || '0,00' }}

                                </span>
                            </div>

                            <div
                                v-if="showCashRegisterInformation"
                                class="absolute-top-right text-white bg-blue-600 font-bold text-dark q-pa-sm rounded-borders shadow-3 text-caption"
                            >
                                Valor independente de filtros!
                            </div>

                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke-width="1.5"
                                stroke="currentColor"
                                class="size-6 text-primary cursor-pointer"
                                @mouseenter="showCashRegisterInformation = true"
                                @mouseleave="showCashRegisterInformation = false"
                            >
                                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z" />
                            </svg>
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

    <VerifyQtdesDialog
        v-if="showResumeQtde"
        @close="showResumeQtde = !$event"
    />
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { onMounted, ref } from 'vue';
    import { filterPopularItensData, getDashBoardData } from '../services/dashBoardService';
    import { getAll } from 'src/modules/cashRegister/services/cashRegisterService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import ManagementReport from 'src/components/Reports/ManagementReport.vue';
    import VerifyQtdesDialog from 'src/components/DashBoard/VerifyQtdesDialog.vue';

    const { notify } = useNotify();
    const startDate = ref<string>('');
    const endDate = ref<string>('');
    const dashBoardData = ref<DashBoardContract>();
    const totalBalance = ref<number>(0);
    const showCashRegisterInformation = ref<boolean>(false);
    const showReports = ref<boolean>(false);
    const showResumeQtde = ref<boolean>(false);

    const popularItensTableColumn: QTableColumn[] = [
        {
            name: 'produto_id',
            label: 'Cód. Produto',
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
        const res = await getDashBoardData(startDate.value, endDate.value);

        if(!res.success)
        {
           notify(
                'negative',
                res.message
            );
        };

        dashBoardData.value = res.data;
    };

    onMounted(async() => {
        await filterPopularItens();

        const res = await getAll();

        if(!res.success) return;
        const cashRegisterData: CashRegisterContract[] = res.data;

        totalBalance.value = cashRegisterData.reduce((total, a) => total + (a.input_value - a.output_value), 0);
    });
</script>
