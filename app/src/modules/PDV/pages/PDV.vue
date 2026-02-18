<template>
    <q-page padding>
        <main class="px-4 max-w" id="sale-page">
            <section class="flex flex-col laptop:flex-row items-start gap-4">
                <div class="w-full laptop:max-w-2xl h-[75vh] flex flex-col bg-white rounded-lg p-4">
                    <div class="flex items-center gap-2">
                        <q-btn
                            icon="menu"
                            flat
                            dense
                            class="shadow-sm !bg-none"
                            @click="showBaseSearchAllProducs = !showBaseSearchAllProducs"
                        />

                        <BaseInputSearchProducts
                            @emit:selected-product="pushProducts([$event])"
                        />
                    </div>

                    <div class="mt-4 overflow-y-auto h-full flex-1 scrollbar-thin border">
                        <q-table
                            :rows="data"
                            :columns="columns"
                            v-model:pagination="pagination"
                            hide-bottom
                        >
                            <template v-slot:body-cell-qtde="props">
                                <q-td :props="props">
                                    <q-input
                                        v-model.number="props.row.qtde"
                                        type="number"
                                        class="w-12 flex ml-auto mr-auto"
                                        input-class="text-center"
                                        dense
                                        @update:model-value="val => validateQtde(Number(val), props.row)"
                                    />
                                </q-td>
                            </template>

                            <template v-slot:body-cell-total="props">
                                <q-td :props="props">
                                    R$ {{
                                        (props.row.price * props.row.qtde)
                                        .toFixed(2)
                                        .replace('.', ',')
                                    }}

                                </q-td>

                            </template>

                            <template v-slot:body-cell-actions="props">
                                <q-td :props="props">
                                    <q-btn
                                        color="red"
                                        icon="delete"
                                        dense
                                        size="7px"
                                        @click="deleteProduct(props.row)"

                                    />
                                </q-td>
                            </template>
                        </q-table>
                    </div>
                </div>

                <!-- Barra da direita-->
                <div class="bg-white rounded-lg h-auto laptop:h-[75vh] p-4 w-full laptop:w-[25rem] laptop:mr-6 flex flex-col">
                    <div class="flex-1 overflow-y-auto">
                        <BaseCustomerSelect
                            v-model="pdvData.customer"
                        />
                    </div>

                    <div class="mt-4 border p-2 rounded">
                        Total R$ {{ calculateTotal }}
                    </div>

                    <div class="mt-4 flex justify-center gap-4">
                        <q-btn
                            icon="payments"
                            dense
                            color="primary"
                            @click="showConfigPayMentForm = !showConfigPayMentForm"
                        />

                        <q-btn
                            dense
                            color="red"
                            icon="delete"
                            :disable="data.length <= 0"
                            @click="deleteSale"
                        />

                        <q-btn
                            dense
                            color="grey"
                            icon="save"
                            :disable="data.length <= 0"
                            @click="showConfirmSaveSaleDialog"
                        />

                        <q-btn
                            color="green"
                            dense
                            no-caps
                            class="mt-auto text-lg font-semibold"
                            label="Finalizar venda"
                            :disable="data.length <= 0"
                            @click="saveSaleForPay(false)"
                        />
                    </div>
                </div>
            </section>
        </main>
    </q-page>

    <BaseSearchAllProducts
        v-if="showBaseSearchAllProducs"
        :type-search="'multiple'"
        @close="showBaseSearchAllProducs = !$event"
        @emit:selected-products="pushProducts($event)"
    />

    <QDialogConfirm
        v-if="showConfirmDialog"
        v-model:show="showConfirmDialog"
        :text="textOperation"
        :operation="operation"
        @close="showConfirmDialog = !$event"
        @confirm="handleConfirmDialog(operation, $event)"

    />

    <PayMentSale
        v-if="showPayMentForms"
        :sale-id="returningSaleId"
        :total-sale="totalSale"
        @close="showPayMentForms = $event"
        @paide="resetSale(!$event)"

    />

    <PayMentForms
        v-if="showConfigPayMentForm"

    />
</template>

<script setup lang="ts">
    import { SessionStorage, QTableColumn, event } from 'quasar';
    import { computed, onMounted, ref, watch } from 'vue';
    import BaseInputSearchProducts from 'src/components/Qinputs/BaseInputSearchProducts.vue';
    import BaseCustomerSelect from 'src/components/Qselects/BaseCustomerSelect.vue';
    import BaseSearchAllProducts from 'src/components/Qtables/BaseSearchAllProducts.vue';
    import PayMentForms from 'src/components/PayMent/PayMentForms/PayMentForms.vue';
    import QDialogConfirm from 'src/helpers/QDialog/Confirm/QDialogConfirm.vue';
    import PayMentSale from 'src/components/PayMent/Pay/PayMentSale.vue';
    import { saveSaleService } from '../services/pdvService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';

    type TPagination = {
        rowsPerPage: number;
    };

    const { notify } = useNotify();
    const data = ref<SaleItemContract[]>([]);

    const columns: QTableColumn[] = [
        {
            name: 'product_id',
            label: 'ID',
            field: 'product_id',
            align: 'left'
        },
        {
            name: 'name',
            label: 'Produto',
            field: 'name',
            align: 'center'
        },
        {
            name: 'qtde',
            label: 'Qtde',
            field: 'qtde',
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
            name: 'total',
            label: 'Total',
            field: 'total',
            align: 'center'
        },
        {
            name: 'actions',
            label: '',
            field: 'actions',
            align: 'right'
        },
    ];

    const returningSaleId = ref<number>();
    const showConfigPayMentForm = ref<boolean>(false);

    const pagination = ref<TPagination>({
        rowsPerPage: 0
    });

    const showBaseSearchAllProducs = ref<boolean>(false);
    const showConfirmDialog = ref<boolean>(false);
    const showPayMentForms = ref<boolean>(false);

    const textOperation = ref<string>('');
    const operation = ref<'save'|'delete'|''>('');
    const totalSale = ref<number>(0);

    const pdvData = ref<SaleContract>({
        id: 0,
        customer: 'Consumidor padrão',
        specie: '',
        products: []
    });

    function removeSessionData(key: string): void {
        SessionStorage.remove(key);
    };

    watch(
        data,
        (newVal) => {
            pagination.value.rowsPerPage = newVal.length;
        },
        { deep: true }
    );

    const deleteProduct = (row: SaleItemContract) => {
        const index = data.value.indexOf(row);

        if(index > -1)
        {
            data.value.splice(index, 1);
        };
    };

    const validateQtde = (val: number, row: SaleItemContract) => {
        if(!val || val <= 0) {
            row.qtde = 1;
            return;
        };

        row.qtde = val;
    };

    const pushProducts = (selectedProducts: SaleItemContract[]) => {
        selectedProducts.forEach(p => {
            const exisit = data.value.find(i => i.id === p.id);

            if(exisit)
            {
                exisit.qtde += 1;

            } else {
                data.value.push({
                    id: 0,
                    product_id: p.id,
                    name: p.name,
                    price: p.price,
                    qtde: 1
                });
            };
        });
    };

    const calculateTotal = computed(() => {
        let subTotal: number = 0;

        data.value.map(p => {
            subTotal += p.price * p.qtde;
        });

        totalSale.value = subTotal;
        return subTotal.toFixed(2).toString().replace('.', ',');
    });

    const deleteSale = () => {
        textOperation.value = 'Deseja realmente cancelar essa venda?';
        operation.value = 'delete';
        showConfirmDialog.value = true;
    };

    const showConfirmSaveSaleDialog = () => {
        textOperation.value = 'Deseja salvar essa venda?'
        operation.value = 'save';
        showConfirmDialog.value = true;
    };

    const confirmDelete = (val: boolean) => {
        if(val)
        {
            data.value = [];
            showConfirmDialog.value = false;
        };
    };

    const handleConfirmDialog = (operation: 'save'|'delete'|'', confirmed: boolean) => {
        if(confirmed && operation === 'delete')
        {
            notify('positive', 'Venda cancelada com sucesso!');
            removeSessionData('sale_id');
            removeSessionData('sale');
            confirmDelete(confirmed);

        };

        if(confirmed && operation === 'save')
        {
            // Isso é gambiarra
            console.log('Salvando: ', SessionStorage.getItem('sale'));

            if(SessionStorage.getItem('sale'))
            {
                notify('positive', 'Dados salvos com sucesso!');
                removeSessionData('sale_id');
                removeSessionData('sale');
                data.value = [];
                showConfirmDialog.value = false;
                return;
            };

            saveSaleForPay(true);
        };

        showConfirmDialog.value = false;
    };

    const saveSaleForPay = async (isSave?: boolean) => {
        const payload: SaleContract = {
            id: 0,
            customer: pdvData.value.customer,
            specie: pdvData.value.specie,
            products: data.value,
        };

        SessionStorage.set('sale', payload);

        const existingSale = SessionStorage.getItem('sale_id');

        if(!isSave && existingSale)
        {
            returningSaleId.value = existingSale as number;
            showPayMentForms.value = true;
            return;
        };

        //console.log('Não existe uma venda aberta no SessionStorage');

        notify('positive', 'Processando dados da venda.');

        const res = await saveSaleService(payload);

        if(res.success)
        {
            returningSaleId.value = res.data.id;
            SessionStorage.set('sale_id', returningSaleId.value);

            if(!res.data.id || res.data.id === 0)
            {
                notify(
                    'negative',
                    'Erro ao finalizar a venda. Identificador inválido!'
                );
                return;
            };

            if(isSave)
            {
                console.log('Chegou nesse isSave');

                notify('positive', 'Dados salvos com sucesso!');
                removeSessionData('sale_id');
                removeSessionData('sale');
                data.value = [];
                return;

            };

            showPayMentForms.value = true;

        } else {
            isSave ? null : notify('negative', `Erro ao finalizar a venda: ${res.message}`);

        };
    };

    /**
     * @param event please report false if using in emits or final saveSaleForPay
     */
    const resetSale = (event: boolean) => {
        console.log('Finalizo a venda, go next');

        removeSessionData('sale_id');
        removeSessionData('sale');

        data.value = [];

        showPayMentForms.value = event;
    };

    onMounted(() => {
        const existingSaleId: number = SessionStorage.getItem('sale_id');
        const existingSale: SaleContract = SessionStorage.getItem('sale');

        if(!existingSaleId && !existingSale) return;

        data.value = existingSale.products;

        pdvData.value = {
            customer: existingSale.customer,
            id: existingSaleId,
            products: data.value,
            specie: ''
        };
    });
</script>

<style lang="scss">
    @media (min-width: 1536px) {
        body {
            overflow-y: hidden !important;
        }
    }
</style>
