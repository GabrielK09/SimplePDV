<template>
    <q-page padding>
        <div class="m-2 text-3xl">
            <div class="flex justify-between">
                <header class="text-gray-600 flex">
                    <h2 class="m-2">
                        Caixa -

                        <span>
                            Saldo total: R$ {{ totalBalance.toFixed(2).toString().replace('.', ',') }}
                        </span>
                    </h2>

                </header>

                <div class="mt-auto mb-auto">
                    <q-btn
                        no-caps
                        color="blue"
                        to="/admin/cash-register/create"
                        label="Cadastrar uma movimentação"
                        class="phone:mb-5"
                    />
                </div>
            </div>

            <div>
                <q-table
                    borded
                    :rows="cashRegister"
                    :columns="columns"
                    row-key="name"
                    class="rounded-xl"
                    :filter="searchInput"
                >
                    <template v-slot:no-data>
                        <div class="ml-auto mr-auto">
                            <q-icon name="warning" size="30px"/>
                            <span class="mt-auto mb-auto ml-2 text-xs">Sem movimentações</span>

                        </div>
                    </template>
                </q-table>
            </div>
        </div>
    </q-page>
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { onMounted, ref } from 'vue';
    import { getAll } from '../services/cashRegisterService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';

    const { notify } = useNotify();

    const columns: QTableColumn[] = [
        {
            name: 'id',
            label: 'ID',
            field: 'id',
            align: 'center'
        },
        {
            name: 'customer',
            label: 'Cliente',
            field: 'customer',
            align: 'center'
        },
        {
            name: 'description',
            label: 'Descrição',
            field: 'description',
            align: 'center'
        },
        {
            name: 'input_value',
            label: 'Valor de entrada',
            field: 'input_value',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
        {
            name: 'output_value',
            label: 'Valor de saída',
            field: 'output_value',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
        {
            name: 'total_balance',
            label: 'Saldo total',
            field: 'total_balance',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },

    ];

    const allCashRegister = ref<CashRegisterContract[]>([]);
    const cashRegister = ref<CashRegisterContract[]>([]);
    const totalBalance = ref<number>(0);
    const searchInput = ref<string>('');

    const getAllCashRegister = async () => {
        const res = await getAll();
        const data = res.data;

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
        };

        cashRegister.value = data;
        allCashRegister.value = [...cashRegister.value];
        totalBalance.value = allCashRegister.value[allCashRegister.value.length - 1].total_balance;
    };

    onMounted(() => {
        getAllCashRegister();
    });
</script>
