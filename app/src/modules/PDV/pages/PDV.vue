<template>
    <q-page padding>
        <main class="px-4" id="sale-page">
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

                    <div class="mt-4 overflow-y-auto h-full flex-1 scrollbar-thin border w-[100%]">
                        <q-table
                            :rows="productsSale"
                            :columns="columns"
                            v-model:pagination="pagination"
                            hide-bottom
                        >
                            <template v-slot:body-cell-name="props">
                                <q-td :props="props">
                                    <span>
                                        {{ `${props.row.name.substring(0, 20)}...` }}

                                        <q-tooltip>
                                            {{ props.row.name }}
                                        </q-tooltip>
                                    </span>
                                </q-td>
                            </template>

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
                        <q-checkbox
                            right-label
                            v-model="registeredCustomer"
                            label="Cliente cadastrado"
                        />

                        <BaseCustomerSelect
                            v-model="pdvData.customer_id"
                            @return:customer-name="pdvData.customer = $event"
                            :is-registered-customer="registeredCustomer"
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
                            :disable="disableButtons.editPayMentsForms"
                            @click="showConfigPayMentForm = !showConfigPayMentForm"
                        />

                        <q-btn
                            dense
                            color="red"
                            icon="delete"
                            :disable="disableButtons.deleteSale"
                            @click="deleteSale"
                        />

                        <q-btn
                            dense
                            color="grey"
                            icon="save"
                            :disable="disableButtons.saveSale"
                            @click="showConfirmSaveSaleDialog"
                        />

                        <q-btn
                            color="green"
                            dense
                            no-caps
                            class="mt-auto text-lg font-semibold"
                            label="Finalizar venda"
                            :disable="disableButtons.finallySale"
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
    <!--@emit:selected-products="pushProducts($event)"-->

    <QDialogConfirm
        v-if="showConfirmDialog"
        v-model:show="showConfirmDialog"
        :text="textOperation"
        :operation="operation"
        :on-leave="onLeaveConfirmDialog"
        @close="showConfirmDialog = !$event"
        @confirm="handleConfirmDialog(operation, $event)"

        @leave:cancel-and-leave="cancelAndLeave($event)"
        @leave:save-and-leave="saveSaleForPay($event)"

    />

    <PayMentSale
        v-if="showPayMentForms"
        :sale-id="returningSaleId"
        :total-sale="totalSale"
        @close="resetForCancelPay(!$event)"
        @paide="resetSale(!$event)"
    />

    <PayMentForms
        v-if="showConfigPayMentForm"
        @close="showConfigPayMentForm = !$event"
    />

</template>

<script setup lang="ts">
    import { SessionStorage, QTableColumn } from 'quasar';
    import { computed, onMounted, reactive, ref, watch } from 'vue';
    import BaseInputSearchProducts from 'src/components/Qinputs/BaseInputSearchProducts.vue';
    import BaseCustomerSelect from 'src/components/Qselects/BaseCustomerSelect.vue';
    import BaseSearchAllProducts from 'src/components/Qtables/BaseSearchAllProducts.vue';
    import PayMentForms from 'src/components/PayMent/PayMentForms/PayMentForms.vue';
    import QDialogConfirm from 'src/helpers/QDialog/Confirm/QDialogConfirm.vue';
    import PayMentSale from 'src/components/PayMent/Pay/PayMentSale.vue';
    import { getSaleDetailsById, insertNewItens, saveSaleService } from '../services/pdvService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router';

    type DisableButtons = {
        editPayMentsForms: boolean;
        deleteSale: boolean;
        saveSale: boolean;
        finallySale: boolean;
    };

    const pagination = ref<TPagination>({
        rowsPerPage: 0,
        sortBy: 'product_id'
    });

    let nextRef: any = null;

    const { notify } = useNotify();

    const registeredCustomer = ref<boolean>(false);
    const hoveredProductId = ref<boolean>(false);

    const disableButtons = reactive<DisableButtons>({
        editPayMentsForms: false,
        deleteSale: true,
        saveSale: true,
        finallySale: true
    });

    /**data is products for sale */
    const productsSale = ref<SaleItemContract[]>([]);

    const originalProductsSale = ref<SaleItemContract[]>([]);

    const route = useRoute();
    const router = useRouter();

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
    const showBaseSearchAllProducs = ref<boolean>(false);
    const showConfirmDialog = ref<boolean>(false);
    const onLeaveConfirmDialog = ref<boolean>(false);
    const showPayMentForms = ref<boolean>(false);
    const textOperation = ref<string>('');
    const operation = ref<'save'|'delete'|''>('');
    const totalSale = ref<number>(0);

    const pdvData = ref<SaleContract>({
        id: 0,
        customer_id: 1,
        customer: '',
        specie: '',
        products: []
    });

    const removeSessionData = (key: string): void => {
        SessionStorage.remove(key);
    };

    const routeSaleId = computed(() => {
        const id = route.query.id;

        if (Array.isArray(id)) return Number(id[0]) || null;
        if (id === null || id === undefined || id === '') return null;

        const parsed = Number(id);
        return Number.isNaN(parsed) ? null : parsed;
    });

    watch(
        () => productsSale.value?.length ?? 0,
        (length) => {

            pagination.value.rowsPerPage = length;
        },
        { immediate: true }
    );

    const deleteProduct = (row: SaleItemContract) => {
        const index = productsSale.value.indexOf(row);

        if(index > -1)
        {
            productsSale.value.splice(index, 1);
        };

        if (productsSale.value.length <= 0)
        {
            disableButtons.editPayMentsForms = false;
            disableButtons.deleteSale = true;
            disableButtons.saveSale = true;
            disableButtons.finallySale = true;
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
        if(!Array.isArray(productsSale.value)) {
            productsSale.value = [];
            return;
        };

        selectedProducts.forEach(p => {
            const exisit = productsSale.value.find(i => i.product_id === p.id);

            if(exisit)
            {
                exisit.qtde += p.qtde;

            } else {
                productsSale.value.push({
                    id: null,
                    product_id: p.id,
                    name: p.name,
                    price: p.price,
                    qtde: p.qtde ?? 1,
                    product_with_characteristics: p.product_with_characteristics

                });
            };
        });

        disableButtons.editPayMentsForms = true;
        disableButtons.deleteSale = false;
        disableButtons.saveSale = false;
        disableButtons.finallySale = false;
    };

    const cloneProducts = (items: SaleItemContract[]) =>
        items.map(item => ({ ...item }));

    const calculateTotal = computed(() => {
        const subTotal = productsSale.value.reduce((total, p) => {
            return total + (p.price * p.qtde);
        }, 0);

        totalSale.value = subTotal;
        return subTotal.toFixed(2).replace('.', ',');
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

    const hasProductChanged = (): boolean => {
        notify(
            'info',
            'Conferindo itens da venda.'
        );

        if (productsSale.value.length !== originalProductsSale.value.length) {
            return true;
        };

        const currentMap = new Map(
            productsSale.value.map(item => [
                item.product_id,
                { qtde: item.qtde, price: item.price }
            ])
        );

        for (const oldItem of originalProductsSale.value) {
            const current = currentMap.get(oldItem.product_id);

            if (!current) return true;

            console.log('Itens: ', {
                current: current,
                oldItem: oldItem
            });

            if (
                current.qtde !== oldItem.qtde || current.price !== oldItem.price
            ) return true;
        };

        return false;
    };

    const handleConfirmDialog = (operation: 'save'|'delete'|'', confirmed: boolean) => {
        if(confirmed && operation === 'save')
        {
            if(SessionStorage.getItem('sale') && routeSaleId.value !== 0)
            {
                console.log('Tem uma venda ainda salva: ', SessionStorage.getItem('sale'), ' e ou é uma venda importada: ', routeSaleId.value, ' valor: ', SessionStorage.getItem('sale') || routeSaleId.value !== 0);

                notify(
                    'positive',
                    'Dados salvos com sucesso!'
                );

                if (SessionStorage.getItem('sale'))
                {
                    removeSessionData('sale_id');
                    removeSessionData('sale');
                };

                if(routeSaleId.value !== 0) router.replace({query: {}});

                productsSale.value = [];
                showConfirmDialog.value = false;

                disableButtons.editPayMentsForms = false;
                disableButtons.deleteSale = true;
                disableButtons.saveSale = true;
                disableButtons.finallySale = true;

                return;
            };

            console.log('Não tinha uma venda salva, vai salvar e resetar os dados do Session');

            saveSaleForPay(true);
            resetSale(false);
        };

        if(confirmed && operation === 'delete')
        {
            notify('positive', 'Venda cancelada com sucesso!');
            removeSessionData('sale_id');
            removeSessionData('sale');

            router.replace({query: {}});

            productsSale.value = [];
            showConfirmDialog.value = false;

            disableButtons.editPayMentsForms = false;
            disableButtons.deleteSale = true;
            disableButtons.saveSale = true;
            disableButtons.finallySale = true;

            return;
        };

        showConfirmDialog.value = false;

        resetSale(false);
    };

    const saveSaleForPay = async (isSave?: boolean) => {
        disableButtons.finallySale = true;

        try {
            const payload: SaleContract = {
                id: null,
                customer_id: pdvData.value.customer_id,
                customer: pdvData.value.customer,
                specie: pdvData.value.specie,
                products: productsSale.value,
            };

            const existingSale = SessionStorage.getItem('sale_id');
            const saleId = routeSaleId.value;

            if(isSave)
            {
                console.log('Foi apenas para salvar');
                const res = await saveSaleService(payload);

                if(res.success)
                {
                    if(hasProductChanged()) {
                        const res = await insertNewItens({
                            id: Number(existingSale || saleId),
                            customer: pdvData.value.customer,
                            customer_id: pdvData.value.customer_id,
                            products: productsSale.value,
                            specie: pdvData.value.specie
                        });

                        if (!res.success) {
                            notify('negative', res.message);
                            console.error('Erro no insertNewItens');

                            resetBtns();

                            return;
                        };

                        console.log('Possui alterações nos produtos.');
                    };

                    returningSaleId.value = res.data.id ?? routeSaleId.value;

                    SessionStorage.set('sale_id', returningSaleId.value);

                    if(!res.data.id || res.data.id === 0)
                    {
                        notify(
                            'negative',
                            'Erro ao finalizar a venda. Identificador inválido!'
                        );

                        return;
                    };

                    notify('positive', 'Dados salvos com sucesso!');

                    removeSessionData('sale_id');
                    removeSessionData('sale');

                    productsSale.value = [];

                    originalProductsSale.value = [];

                    router.replace({query: {}});

                    return;
                };
            };

            // Confirma se a venda não foi reaberta / importada
            if((existingSale || saleId) && hasProductChanged()) {
                const res = await insertNewItens({
                    id: Number(existingSale || saleId),
                    customer: pdvData.value.customer,
                    customer_id: pdvData.value.customer_id,
                    products: productsSale.value,
                    specie: pdvData.value.specie
                });

                if (!res.success) {
                    console.error('Erro no insertNewItens');

                    notify('negative', res.message);
                    resetBtns();

                    return;
                };

                originalProductsSale.value = cloneProducts(productsSale.value)
            };

            if (saleId) {
                showPayMentForms.value = true;
                returningSaleId.value = Number(saleId);
                return;
            };

            SessionStorage.set('sale', payload);

            if(!isSave && existingSale)
            {
                returningSaleId.value = Number(existingSale);
                showPayMentForms.value = true;
                return;
            };

            notify('positive', 'Processando dados da venda.');

            const res = await saveSaleService(payload);

            if(res.success)
            {
                returningSaleId.value = res.data.id ?? routeSaleId.value;

                SessionStorage.set('sale_id', returningSaleId.value);

                if(!res.data.id || res.data.id === 0)
                {
                    notify(
                        'negative',
                        'Erro ao finalizar a venda. Identificador inválido!'
                    );
                    return;
                };

                showPayMentForms.value = true;

            } else {
                isSave
                    ? null
                    : notify('negative', res.message);

                resetBtns();
            };
        } catch (error) {
            notify(
                'negative',
                error.message
            );

        } finally {
            disableButtons.finallySale = false;

        };
    };

    /**
     * Event please report false if using in emits or final saveSaleForPay
     */
    const resetSale = (event: boolean) => {
        showPayMentForms.value = event;

        removeSessionData('sale_id');
        removeSessionData('sale');

        productsSale.value = [];
        originalProductsSale.value = [];
        pdvData.value.products = productsSale.value;

        registeredCustomer.value = false;

        disableButtons.editPayMentsForms = false;
        disableButtons.deleteSale = true;
        disableButtons.saveSale = true;
        disableButtons.finallySale = true;

        router.replace({query: {}});

        pdvData.value.customer = 'Consumidor padrão';
    };

    const resetBtns = (): void => {
        disableButtons.editPayMentsForms = false;
        disableButtons.deleteSale = true;
        disableButtons.saveSale = true;
        disableButtons.finallySale = true;
    };

    onMounted(async () => {
        productsSale.value = [];

        if(routeSaleId.value)
        {
            notify(
                'positive',
                'Carregando dados da venda ...'
            );

            const res = await getSaleDetailsById(Number(route.query.id));
            const resData: SaleContract = res.data.sale;

            if(!res.success)
            {
                notify(
                    'negative',
                    res.message
                );
            };

            productsSale.value = resData.products || [];
            originalProductsSale.value = cloneProducts(productsSale.value);

            pdvData.value = {
                customer: resData.customer,
                id: resData.id,
                customer_id: resData.customer_id,
                products: productsSale.value,
                specie: ''
            };

            disableButtons.editPayMentsForms = true;
            disableButtons.deleteSale = false;
            disableButtons.saveSale = false;
            disableButtons.finallySale = false;

            return;
        };

        const existingSaleId: number = SessionStorage.getItem('sale_id');
        const existingSale: SaleContract = SessionStorage.getItem('sale');

        if(!existingSaleId && !existingSale) return;

        productsSale.value = existingSale.products || [];

        originalProductsSale.value = cloneProducts(existingSale.products);

        pdvData.value = {
            id: existingSaleId,
            customer: existingSale.customer,
            customer_id: existingSale.customer_id,
            products: productsSale.value,
            specie: existingSale.specie
        };

        disableButtons.editPayMentsForms = true;
        disableButtons.deleteSale = false;
        disableButtons.saveSale = false;
        disableButtons.finallySale = false;

    });

    const resetForCancelPay = (event: boolean) => {
        showPayMentForms.value = event;

        disableButtons.editPayMentsForms = true;
        disableButtons.deleteSale = false;
        disableButtons.saveSale = false;
        disableButtons.finallySale = false;
    };

    const cancelAndLeave = (confirm: boolean) => {
        if (confirm && nextRef)
        {
            nextRef();
        } else if (nextRef){
            nextRef(false);
        };

        nextRef = null;
    };

    onBeforeRouteLeave((to, from, next) => {
        const hasOpenSale = SessionStorage.getItem('sale_id') || SessionStorage.getItem('sale') || productsSale.value?.length > 0

        if (hasOpenSale) {
            textOperation.value = 'Existe uma venda aberta!'
            showConfirmDialog.value = true;
            onLeaveConfirmDialog.value = true;

            nextRef = next;
            return;
        };

        next();
    });
</script>

<style lang="scss">
    @media (min-width: 1536px) {
        body {
            overflow-y: hidden !important;
        }
    }
</style>
