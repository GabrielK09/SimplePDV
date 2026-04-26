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
                        <div class="flex">
                            <div class="mr-4">
                                <q-select 
                                    v-model="byStatus" 
                                    :options="statusOptions" 
                                    option-label="Status"
                                    emit-value
                                    map-options
                                    outlined
                                    dense
                                    :display-value="selectedLabel"
                                    :clearable="true"
                                    @update:model-value="applyFilters"
                                />
                            </div>

                            <div>
                                <q-input
                                    outlined
                                    v-model="searchInput"
                                    type="text"
                                    dense
                                    label=""
                                    @update:model-value="applyFilters"
                                >
                                    <template v-slot:append>
                                        <q-icon name="search" />
                                    </template>
                                    <template v-slot:label>
                                        <span class="text-xs">Buscar por uma compra ...</span>
                                    </template>
                                </q-input>
                            </div>
                        </div>
                    </template>

                    <template v-slot:body="props">
                        <q-tr :props="props">
                            <q-td v-for="col in props.cols">
                                <template v-if="col.name === 'actions'">
                                    <q-btn 
                                        dense
                                        flat
                                        icon="more_vert"
                                    >
                                        <q-menu
                                            anchor="bottom right"
                                            self="top right"
                                            class="rounded shadow-xl bg-white"
                                            transition-show="jump-down"
                                        >
                                            <q-list style="min-width: 90px">
                                                <q-item 
                                                    clickable 
                                                    v-close-popup  
                                                    @click="buildShowShoppingDetails(props.row.id)"
                                                >
                                                    <q-item-section avatar>
                                                        <q-icon name="visibility" color="primary" size="20px" />
                                                    </q-item-section>
                                                    <q-item-section>
                                                        <q-item-label>Ver detalhes</q-item-label>
                                                    </q-item-section>
                                                </q-item>

                                                <q-item  
                                                    v-if="props.row.status !== 'Concluída' && props.row.status !== 'Cancelada'" 
                                                    clickable 
                                                    v-close-popup 
                                                    @click="importShopping(props.row.id)"
                                                >                                                        
                                                    <q-item-section avatar>
                                                        <q-icon name="upload" size="20px"/>
                                                    </q-item-section>
                                                    <q-item-section>
                                                        <q-item-label>Importar venda</q-item-label>
                                                    </q-item-section>
                                                </q-item>

                                                <q-item 
                                                    v-if="props.row.status === 'Concluída'"
                                                    clickable 
                                                    v-close-popup
                                                    @click="showCancelShopping(props.row.id)"
                                                >
                                                    <q-item-section avatar>
                                                        <q-icon name="cancel" color="red" size="20px"/>
                                                    </q-item-section>
                                                    <q-item-section>
                                                        <q-item-label>Cancelar venda</q-item-label>
                                                    </q-item-section>
                                                </q-item>
                                            </q-list>
                                        </q-menu>
                                    </q-btn>
                                </template>                                    

                                <template v-if="col.name === 'status'">
                                    <div
                                        class="text-center flex flex-center"
                                        :class="{
                                            'text-green-600': props.row.status === 'Concluída',
                                            'text-red-600': props.row.status === 'Cancelada'
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

    <ShoppingDetails
        v-if="showShoppingDetails"
        :shopping-id="selectedShoppingId"
        @close="showShoppingDetails = !$event"
    />
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { computed, onMounted, ref } from 'vue';
    import { getAll } from '../services/shoppingService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { cancelShoppingOrSale } from 'src/modules/PDV/services/payMentFormsService';
    import QDialogConfirm from 'src/helpers/QDialog/Confirm/QDialogConfirm.vue';
    import { useRouter } from 'vue-router';
    import ShoppingDetails from './Show/ShoppingDetails.vue';

    const statusOptions: Exclude<FilterByStatus, null>[] = [
        'Pendente', 
        'Concluída', 
        'Cancelada'
    ];

    const byStatus = ref<FilterByStatus>(null);

    const { notify } = useNotify();
    const showConfirmDialog = ref<boolean>(false);
    const router = useRouter();

    const pagination = ref({
        sortBy: 'id'
    });

    const columns: QTableColumn[] = [
        {
            name: 'id',
            label: 'Cód compra',
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

    const allshopping = ref<ShoppingContract[]>([]);
    const shopping = ref<ShoppingContract[]>([]);

    const searchInput = ref<string>('');
    const selectedShoppingId = ref<number | null>(0);
    const showShoppingDetails = ref<boolean>(false);

    const getAllshopping = async () => {
        const res = await getAll();
        const data = res.data as ShoppingContract[];

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
        };
    
        allshopping.value = data;
        applyFilters();
    };

    const applyFilters = () => {
        let filtred = [...allshopping.value];

        if (byStatus.value) {
            filtred = filtred.filter(s => s.status === byStatus.value);
        };

        if(searchInput.value.trim()) {
            const search = searchInput.value.trim().toLocaleLowerCase();

            filtred = filtred.filter(s => 
                String(s.id).includes(search) ||
                String(s.load ?? '').toLowerCase().includes(search) ||
                String(s.status ?? '').toLowerCase().includes(search)
            );
        };

        shopping.value = filtred;
    };

    const selectedLabel = computed(() => {
        return byStatus.value ?? 'Todos';
    }); 

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

    const buildShowShoppingDetails = (shoppingId: number): void => {
        showShoppingDetails.value = !showShoppingDetails.value;
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
