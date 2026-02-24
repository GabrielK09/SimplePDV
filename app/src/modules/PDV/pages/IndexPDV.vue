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

                    :rows="pdvs"
                    :columns="columns"
                    row-key="name"
                    class="rounded-xl"
                    :filter="searchInput"
                >
                    <template v-slot:body="props">
                        <q-tr
                            :props="props"
                        >
                            <q-td
                                v-for="col in props.cols"
                            >
                                <div
                                    class="flex flex-center"
                                >
                                    <template v-if="col.name === 'actions'">
                                        <div v-if="props.row.status === 'Concluída'">
                                            <q-btn
                                                size="10px"
                                                no-caps
                                                color="red"
                                                icon="cancel"
                                                flat
                                                @click="showCancelSale(props.row.id)"
                                            />
                                        </div>

                                        <div v-else-if="props.row.status !== 'Concluída' && props.row.status !== 'Cancelado'">
                                            <q-btn
                                                size="10px"
                                                no-caps
                                                color="black"
                                                icon="upload"
                                                flat
                                                @click="importSale(props.row.id)"
                                            />
                                        </div>

                                        <q-btn
                                            size="10px"
                                            no-caps
                                            color="black"
                                            icon="visibility"
                                            flat
                                            @click="buildShowSaleDetails(props.row.id)"
                                        />
                                    </template>

                                    <template v-if="col.name === 'status'">
                                        <div
                                            :class="{
                                                'text-green-600': props.row.status === 'Concluída',
                                                'text-red-600': props.row.status === 'Cancelado'
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
    import { onMounted, ref } from 'vue';
    import { cancelSale, getAll } from '../services/pdvService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import QDialogConfirm from 'src/helpers/QDialog/Confirm/QDialogConfirm.vue';
    import { useRouter } from 'vue-router';
    import SaleDetails from './Show/SaleDetails.vue';

    const { notify } = useNotify();
    const router = useRouter();
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

    const showCancelSale = (saleId: number): void => {
        showConfirmDialog.value = true;
        selectedSaleId.value = saleId;
    };

    const handleConfirmDialog = async (event: boolean): Promise<void> => {
        if(!event) return;

        const res = await cancelSale(selectedSaleId.value);

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );

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
