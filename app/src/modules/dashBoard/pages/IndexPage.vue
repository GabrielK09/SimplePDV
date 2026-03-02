<template>
    <q-page padding class="bg-grey-2">
        <div class="q-gutter-lg">
            <q-card class="q-pa-md">
                <q-form
                    @submit.prevent="filterDashBoard"
                >
                    <div class="row q-col-gutter-md items-end">
                        <div class="col-12 col-md-3">
                            <q-input v-model="startDate" type="date" label="Data inicial" outlined />
                        </div>

                        <div class="col-12 col-md-3">
                            <q-input v-model="endDate" type="date" label="Data final" outlined />
                        </div>

                        <div class="col-12 col-md-2">
                            <q-btn
                                color="primary"
                                type="submit"
                                no-caps
                                label="Filtrar"
                                class="full-width"
                            />
                        </div>
                    </div>
                </q-form>
            </q-card>

            <!-- CARDS DE MÉTRICAS -->
            <div class="row q-col-gutter-md mr-6">
                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Saldo total do caixa</div>

                        <div class="flex justify-between">
                            <div class="text-h5 text-primary">
                                R$ {{ totalBalance.toFixed(2).toString().replace('.', ',') || '0,00' }}
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

                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Total Vendido</div>
                        <div class="text-h5 text-blue-800">
                            R$ {{ dashBoardData?.total_saled.toFixed(2).toString().replace('.', ',') || '0,00' }}
                        </div>
                    </q-card>
                </div>

                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Comissão</div>
                        <div class="text-h5 text-green">
                            R$ {{ dashBoardData?.commission.toFixed(2).toString().replace('.', ',') || '0,00' }}
                        </div>
                    </q-card>
                </div>

                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Quantidade de Vendas</div>
                        <div class="text-h5 text-indigo">
                            {{ dashBoardData?.amount_saled.toFixed(2).toString().replace('.', ',') || '0,00' }}
                        </div>
                    </q-card>
                </div>

                <div class="col-12 col-md-3">
                    <q-card class="q-pa-md shadow-2">
                        <div class="text-grey">Melhor Cliente</div>
                        <div class="text-h6 text-deep-orange">
                            {{ dashBoardData?.best_customer || '-' }}
                        </div>
                    </q-card>
                </div>
            </div>
        </div>
    </q-page>
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { computed, onMounted, ref } from 'vue';
    import ApexChart from 'vue3-apexcharts';
    import { getDashBoardData } from '../services/dashBoardService';
    import dayjs from 'dayjs';
    import { getAll } from 'src/modules/cashRegister/services/cashRegisterService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';

    dayjs('2006-01-02')

    const { notify } = useNotify();
    const startDate = ref<string>('');
    const endDate = ref<string>('');
    const dashBoardData = ref<DashBoardContract>();
    const totalBalance = ref<number>(0);
    const showCashRegisterInformation = ref<boolean>(false);

    const filterDashBoard = async () => {
        const res = await getDashBoardData(startDate.value, endDate.value);

        if(res.success)
        {
            dashBoardData.value = res.data;

        } else {
            notify(
                'negative',
                res.message
            );
        };
    };

    onMounted(async() => {
        const res = await getAll();

        if(!res.success) return;
        const cashRegisterData = res.data;

        totalBalance.value = cashRegisterData[cashRegisterData.length - 1].total_balance;
    });
</script>
