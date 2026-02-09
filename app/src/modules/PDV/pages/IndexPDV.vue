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
                >
                    <template v-slot:top-right>
                        <q-input
                            outlined
                            v-model="searchInput"
                            type="text"
                            label=""
                            @update:model-value="filterPDV"
                        >
                            <template v-slot:append>
                                <q-icon name="search" />
                            </template>
                            <template v-slot:label>
                                <span class="text-xs">Buscar por uma venda ...</span>
                            </template>
                        </q-input>
                    </template>

                    <template v-slot:body="props">
                        <q-tr
                            :props="props"
                        >
                            <q-td
                                v-for="(col, i) in props.cols"
                            >
                                <template v-if="col.name === 'actions'">
                                    <div
                                        class="text-center"
                                    >

                                    </div>
                                </template>

                                <template v-else>
                                    {{ col.value }}
                                </template>
                            </q-td>
                        </q-tr>
                    </template>

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
    import { QTableColumn, useQuasar } from 'quasar';
    import { onMounted, ref } from 'vue';
    import camelcaseKeys from 'camelcase-keys';
    import { getAll } from '../services/pdvService';

    const $q = useQuasar();

    const columns: QTableColumn[] = [
        {
            name: 'id',
            label: 'ID',
            field: 'id',
            align: 'center'
        },
        {
            name: 'price',
            label: 'Preço',
            field: 'price',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
        {
            name: 'qtde',
            label: 'Qtde',
            field: 'qtde',
            align: 'center'
        },
        {
            name: 'actions',
            label: '',
            field: 'actions',
            align: 'right'
        }
    ];

    let allPDVs = ref<ProductContract[]>([]);
    let pdvs = ref<ProductContract[]>([]);

    let searchInput = ref<string>('');

    const getAllshopping = async () => {
        const res = await getAll();
        const data = camelcaseKeys(res.data, { deep: true });

        pdvs.value = data;
        allPDVs.value = [...pdvs.value];

    };

    const filterPDV = () => {
        pdvs.value = allPDVs.value.filter(shoop => shoop.name.toLowerCase().includes(searchInput.value));

    };

    onMounted(() => {
        getAllshopping();
    });
</script>
