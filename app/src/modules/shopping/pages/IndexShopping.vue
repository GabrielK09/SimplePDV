<template>
    <q-page padding>
        <div class="m-2 text-3xl" >
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

            <div>
                <q-table
                    v-model:pagination="pagination"
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
                        <q-tr :props="props">
                            <q-td v-for="col in props.cols">
                                <template v-if="col.name === 'actions'">
                                    <div class="text-center flex flex-center">
                                        <div v-if="props.row.status === 'Concluída' && props.row.status !== 'Cancelado'">
                                            <q-btn
                                                size="10px"
                                                no-caps
                                                color="red"
                                                icon="cancel"
                                                flat
                                                @click="showCancelShopping(props.row.id)"
                                            />
                                        </div>

                                        <div v-else-if="props.row.status !== 'Concluída' && props.row.status !== 'Cancelado'">
                                            <q-btn
                                                size="10px"
                                                no-caps
                                                color="black"
                                                icon="upload"
                                                flat
                                                @click="importShopping(props.row.id)"
                                            />
                                        </div>

                                        <q-btn
                                            size="10px"
                                            color="black"
                                            icon="visibility"
                                            flat                                        
                                        />
                                    </div>
                                </template>

                                <template v-if="col.name === 'status'">
                                    <div
                                        class="text-center flex flex-center"
                                        :class="{
                                            'text-green-600': props.row.status === 'Concluída',
                                            'text-red-600': props.row.status === 'Cancelado'
                                        }"
                                    >
                                        {{ col.value }}
                                    </div>
                                </template>

                                <template v-else>
                                    <div class="text-center">
                                        {{ col.value }}

                                    </div>
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

    <QDialogConfirm
        v-if="showConfirmDialog"
        :text="'Deseja realmente cancelar essa compra?'"
        :show="showConfirmDialog"
        @confirm="handleConfirmDialog($event)"
        @close="showConfirmDialog = !$event"
    />
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { onMounted, ref } from 'vue';
    import { getAll } from '../services/shoppingService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { cancelShoppingOrSale } from 'src/modules/PDV/services/payMentFormsService';
    import QDialogConfirm from 'src/helpers/QDialog/Confirm/QDialogConfirm.vue';
    import { useRouter } from 'vue-router';

    const { notify } = useNotify();
    const showConfirmDialog = ref<boolean>(false);
    const router = useRouter();

    const pagination = ref({
        sortBy: 'id' 
    });

    const columns: QTableColumn[] = [
        {
            name: 'id',
            label: 'ID',
            field: 'id',
            align: 'center'
        },
        {
            name: 'load',
            label: 'Carga da compra',
            field: 'load',
            align: 'center'
        },
        {
            name: 'status',
            label: 'Status',
            field: 'status',
            align: 'center'
        },
        {
            name: 'total_shopping',
            label: 'Total da compra',
            field: 'total_shopping',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
        {
            name: 'actions',
            label: '',
            field: 'actions',
            align: 'right'
        }
    ];

    const allshopping = ref<ProductContract[]>([]);
    const shopping = ref<ProductContract[]>([]);

    const searchInput = ref<string>('');
    const selectedShoppingId = ref<number | null>(0);

    const getAllshopping = async () => {
        const res = await getAll();
        const data = res.data;

        shopping.value = data;
        allshopping.value = [...shopping.value];
    };

    const filtershopping = () => {
        //shopping.value = allshopping.value.filter(shoop => shoop.name.toLowerCase().includes(searchInput.value));
    };

    const handleConfirmDialog = async (event: boolean): Promise<void> => {
        if(!event) return;

        const res = await cancelShoppingOrSale({
            shopping_id: selectedShoppingId.value,
            route: 'shopping'
        });

        if(!res.success)
        {
            notify(
                'negative',
                res.message || 'Erro ao realizar a operação.'
            );

            showConfirmDialog.value = false;
            return;
        };

        notify(
            'positive',
            res.message
        );

        showConfirmDialog.value = false;
        
        await getAllshopping();

        return;
    };

    const showCancelShopping = (shoppingId: number): void => {
        showConfirmDialog.value = true;
        selectedShoppingId.value = shoppingId;
    };

    const importShopping = (shoppingId: number) => {
        router.replace({
            name: 'shopping.create',
            query: {
                id: shoppingId
            }
        });
    };

    onMounted(() => {
        getAllshopping();
    });
</script>
