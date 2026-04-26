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
                        class="max-phone:mb-5"
                        label="Cadastrar uma venda"
                    />
                </div>
            </div>

            <div>
                <q-table
                    v-model:pagination="pagination"
                    :rows="pdvs"
                    :columns="columns"
                    row-key="name"
                    class="rounded-xl"
                    :filter="searchInput"
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

                            <div class="">
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
                                        <span class="text-xs">Buscar por uma venda ...</span>
                                    </template>
                                </q-input>
                            </div>
                        </div>
                    </template>

                    <template v-slot:body="props">
                        <q-tr :props="props">
                            <q-td v-for="col in props.cols">
                                <div class="flex flex-center">
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
                                                        @click="buildShowSaleDetails(props.row.id)"
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
                                                        @click="importSale(props.row.id)"
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
                                                        @click="showCancelSale(props.row.id)"
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
                                            :class="{
                                                'text-green-600': props.row.status === 'Concluída',
                                                'text-red-600': props.row.status === 'Cancelada'
                                            }"
                                        >
                                            {{ col.value }}
                                        </div>
                                    </template>

                                    <template v-else>
                                        {{ col.value }}

                                    </template>
                                </div>
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

    <QDialogConfirm
        v-if="showConfirmDialog"
        :text="'Deseja realmente cancelar essa venda?'"
        :show="showConfirmDialog"
        @confirm="handleConfirmDialog($event)"
        @close="showConfirmDialog = !$event"
    />

    <SaleDetails
        v-if="showSaleDetails"
        :saleId="selectedSaleId"
        @close="showSaleDetails = !$event"
    />
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { computed, onMounted, ref } from 'vue';
    import { getAll } from '../services/pdvService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import QDialogConfirm from 'src/helpers/QDialog/Confirm/QDialogConfirm.vue';
    import { useRouter } from 'vue-router';
    import SaleDetails from './Show/SaleDetails.vue';
    import { cancelShoppingOrSale } from '../services/payMentFormsService';

    const statusOptions: Exclude<FilterByStatus, null>[] = [
        'Pendente', 
        'Concluída', 
        'Cancelada'
    ];

    const byStatus = ref<FilterByStatus>(null);

    const { notify } = useNotify();
    const router = useRouter();

    const pagination = ref<TPagination>({
        sortBy: 'id',
        rowsPerPage: 10
    });

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
            align: 'center'
        }
    ];

    const allPDVs = ref<PDVContract[]>([]);
    const pdvs = ref<PDVContract[]>([]);

    const showConfirmDialog = ref<boolean>(false);
    const showSaleDetails = ref<boolean>(false);
    const selectedSaleId = ref<number>(0);

    const searchInput = ref<string>('');

    const getAllPdv = async () => {
        const res = await getAll();
        const data = res.data as PDVContract[];

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
        };

        allPDVs.value = data;
        applyFilters();
    };

    const applyFilters = () => {
        let filtred = [...allPDVs.value];

        if(byStatus.value) {
            filtred = filtred.filter(s => s.status === byStatus.value);
        };

        if(searchInput.value.trim()) 
        {
            const search = searchInput.value.trim().toLowerCase();

            filtred = filtred.filter(s =>
                String(s.id).includes(search) ||
                String(s.customer).includes(search)
            );
        };

        pdvs.value = filtred;
    };

    const selectedLabel = computed(() => {
        return byStatus.value ?? 'Todos';
    }); 

    const showCancelSale = (saleId: number): void => {
        showConfirmDialog.value = true;
        selectedSaleId.value = saleId;
    };

    const handleConfirmDialog = async (event: boolean): Promise<void> => {
        if(!event) return;

        const res = await cancelShoppingOrSale({
            sale_id: selectedSaleId.value,
            route: 'sale'
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

        await getAllPdv();

        return;
    };

    const buildShowSaleDetails = (saleId: number): void => {
        showSaleDetails.value = !showSaleDetails.value;
        selectedSaleId.value = saleId;
    };

    /**
     * For re-open sale
     */
    const importSale = (saleId: number) => {
        router.replace({
            path: '/admin/pdv',
            query: {
                id: saleId
            }
        });
    };

    onMounted(() => {
        getAllPdv();
    });
</script>
