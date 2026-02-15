<template>
    <q-page padding>
        <div class="m-2 text-3xl" >
            <div class="flex justify-between">
                <h2 class="text-gray-600 m-2">Listagem PDV</h2>

                <div class="mt-auto mb-auto">
                    <q-btn
                        no-caps
                        color="blue"
                        to="/admin/pdv"
                        label="Cadastrar uma venda"

                    />
                </div>
            </div>

            <div>
                <q-table
                    borded
                    :rows="pdvs"
                    :columns="columns"
                    row-key="name"
                    class="rounded-xl"
                    :filter="searchInput"
                >
                    <template v-slot:no-data>
                        <div class="ml-auto mr-auto">
                            <q-icon name="warning" size="30px"/>
                            <span class="mt-auto mb-auto ml-2 text-xs">Sem vendas cadastrados</span>

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
    import { getAll } from '../services/pdvService';
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
            name: 'status',
            label: 'Status',
            field: 'status',
            align: 'center'
        },
        {
            name: 'sale_value',
            label: 'Valor da venda',
            field: 'sale_value',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
        {
            name: 'actions',
            label: 'Ações',
            field: 'actions',
            align: 'right'
        }
    ];

    const allPDVs = ref<PDVContract[]>([]);
    const pdvs = ref<PDVContract[]>([]);

    const searchInput = ref<string>('');

    const getAllshopping = async () => {
        const res = await getAll();
        const data = res.data;

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
        };

        pdvs.value = data;
        allPDVs.value = [...pdvs.value];
    };

    onMounted(() => {
        getAllshopping();
    });
</script>
