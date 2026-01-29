<template>
    <q-page padding>
        <div
            class="m-2"
        >
            <div class="flex justify-between">
                <h2 class="text-gray-600 m-2">Compras</h2>

                <div class="mt-auto mb-auto">
                    <q-btn
                        no-caps
                        color="blue"
                        to="/admin/shopping/create"
                        label="Cadastrar novo compra"

                    />
                </div>
            </div>

            <div class="">
                <q-table
                    borded
                    :rows="shopping"
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
                            @update:model-value="filtershopping"
                        >
                            <template v-slot:append>
                                <q-icon name="search" />
                            </template>
                            <template v-slot:label>
                                <span class="text-xs">Buscar por uma compra ...</span>
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
                                        <q-btn size="10px" no-caps color="red" icon="delete" flat @click="showDialogDeleteshopping(props.row.id)"/>

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
                            <span class="mt-auto mb-auto ml-2 text-xs">Sem compras cadastrados</span>

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
    import { getAll, deleteshopping } from '../services/shoppingService';

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

    let allshopping = ref<ProductContract[]>([]);
    let shopping = ref<ProductContract[]>([]);

    let searchInput = ref<string>('');

    const getAllshopping = async () => {
        const res = await getAll();
        const data = camelcaseKeys(res.data, { deep: true });

        shopping.value = data;
        allshopping.value = [...shopping.value];

    };

    const showDialogDeleteshopping = (shoppingId: number) => {
        $q.dialog({
            title: 'Excluir compra',
            message: `Deseja realmente remover esse compra (${shoppingId})?`,
            cancel: {
                push: true,
                label: 'Não',
                color: 'red',
            },

            ok: {
                push: true,
                label: 'Sim',
                color: 'green',
            },

        }).onOk(() => {
            deleteShoopByDialog(shoppingId);

        }).onCancel(() => {
            return;
        });
    };

    const deleteShoopByDialog = async (shoppingId: number) => {
        const data = await deleteshopping(shoppingId);

        console.log(data);

        if(data.success)
        {
            $q.notify({
                type: 'positive',
                message: data.message,
                position: 'top',
                timeout: 1200

            });
        } else {
            $q.notify({
                type: 'negative',
                message: data.data.data,
                position: 'top',
                timeout: 1200

            });
        };

        getAllshopping();
    };

    const filtershopping = () => {
        shopping.value = allshopping.value.filter(shoop => shoop.name.toLowerCase().includes(searchInput.value));

    };

    onMounted(() => {
        getAllshopping();
    });
</script>
