<template>
    <q-page padding>
        <div class="flex flex-center">
            <div>
                <q-input
                    v-model="startDate"
                    type="date"
                    label="Data inicial"
                />

                <q-input
                    v-model="endDate"
                    type="date"
                    label="Data final"
                />

                <q-btn
                    color="primary"
                    no-caps
                    label="Filtrar"
                    @click="filterDashBoard"

                />
            </div>

            <div id="chart" class="w-[50%]">
                <ApexChart
                    type="area"
                    :options="chartOptions"
                    :series="chartSeries"
                />
            </div>

            <div class="customer-table">
                <pre>{{ dashBoardData }}</pre>

            </div>
        </div>
    </q-page>
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { ref } from 'vue';
    import ApexChart from 'vue3-apexcharts';
    import { getDashBoardData } from '../services/dashBoardService';
    import dayjs from 'dayjs';
    import { useNotify } from 'src/helpers/QNotify/useNotify';

    dayjs('2006-01-02')

    const { notify } = useNotify();
    const startDate = ref<string>('');
    const endDate = ref<string>('');
    const dashBoardData = ref<DashBoardContract>();

    const chartOptions: ApexCharts.ApexOptions = {
        chart: {
            width: 10,
            height: 160,
            sparkline: {
                enabled: true

            },
        },
        subtitle: {
            text: 'Vendas',
            offsetX: 30,
            style: {
                fontSize: '14px',
            }
        },
        stroke: {
            curve: 'smooth'
        }
    };

    const chartSeries = [
        {
            name: 'Vendas',
            data: [1, 2, 3]
        }
    ];

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
</script>
