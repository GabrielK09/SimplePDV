 <template>
    <q-page padding>
        <main class="px-4" id="sale-page">
            <section class="section_pdv">
                <div class="left-bar rounded-lg">
                    <div class="flex items-center gap-2 mb-6">
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

                    <div class="sale_products border">
                        <q-table
                            class="full-width"
                            :rows="productsSale"
                            :columns="columns"
                            v-model:pagination="pagination"
                            hide-bottom
                        >
                            <template v-slot:body="props">
                                <q-tr :props="props">
                                    <q-td key="product_id" :props="props">
                                        <span>{{ props.row.product_id }}</span>

                                    </q-td>

                                    <q-td key="name" :props="props">
                                        <span class="product_name">
                                            {{ props.row.name }}

                                            <span v-if="hasCharacteristics(props.row)">
                                                <q-btn 
                                                    v-if="hasCharacteristics(props.row)"
                                                    size="10px"
                                                    color="black"
                                                    :icon="isExpanded(props.row.product_id) ? 'expand_less' : 'grid_on'"
                                                    flat
                                                    @click="toggleExpanded(props.row.product_id)"
                                                >
                                                    <q-tooltip>
                                                        {{ 
                                                            isExpanded(props.row.product_id)
                                                                ? 'Ocultar'
                                                                : 'Ver grade'
                                                        }}
                                                    </q-tooltip>
                                                </q-btn>
                                            </span>
                                            
                                        </span>
                                    </q-td>

                                    <q-td key="qtde" :props="props">
                                        <q-input
                                            v-model.number="props.row.qtde"
                                            type="number"
                                            class="w-12 flex ml-auto mr-auto"
                                            input-class="text-center"
                                            dense
                                            :disable="hasCharacteristics(props.row)"
                                            @update:model-value="val => validateQtde(Number(val), props.row)"
                                        />
                                    </q-td>
                                    
                                    <q-td key="price" :props="props">
                                        R$ {{
                                            props.row.price
                                                .toFixed(2)
                                                .replace('.', ',')
                                        }}
                                    </q-td>

                                    <q-td key="total" :props="props">
                                        R$ {{
                                            hasCharacteristics(props.row) 
                                                ? (props.row.price * props.row.product_with_characteristics.reduce((total: any, a: any) => total + (a.grid_qtde), 0)).toFixed(2).replace('.', ',')
                                                : (props.row.price * props.row.qtde).toFixed(2).replace('.', ',')
                                        }}
                                    </q-td>

                                    <q-td key="actions" :props="props">
                                        <q-btn
                                            color="red"
                                            icon="delete"
                                            dense
                                            size="7px"
                                            @click="deleteProduct(props.row)"

                                        />
                                    </q-td>
                                </q-tr>

                                <q-tr 
                                    v-if="isExpanded(props.row.product_id) && hasCharacteristics(props.row)"
                                    :props="props"
                                >
                                    <q-td colspan="80%" class="bg-gray-200">
                                        <div class="q-pa-md">
                                            <div class="row q-col-gutter-sm">
                                                <div
                                                    v-for="(characteristic, i) in props.row.product_with_characteristics"
                                                    :key="`${props.row.product_id}-${characteristic.size}`"
                                                    class="col-12 col-sm-6 col-md-3"
                                                >       
                                                    <q-card flat bordered class="w-max">
                                                        <q-card-section>
                                                            <div class="text-caption text-gray-700 flex mt-2">
                                                                Tamanho: 
                                                                
                                                                <div class="text-body2 text-weight-bold ml-3">
                                                                    {{ characteristic.size }}
                                                                </div>
                                                            </div>

                                                            <div>
                                                                <q-input
                                                                    v-model.number="characteristic.grid_qtde"
                                                                    type="number"
                                                                    label="Qtde"
                                                                    stack-label
                                                                    label-slot
                                                                    input-class="text-center"
                                                                    dense
                                                                    @update:model-value="val => validateGridQtde(Number(val), props.row, Number(i))"
                                                                />
                                                            </div>
                                                        </q-card-section>
                                                    </q-card>
                                                </div>
                                            </div>
                                        </div>
                                    </q-td>
                                </q-tr>
                            </template>
                        </q-table>
                    </div>
                </div>

                <div class="right-bar rounded-lg">
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
    import PayMentSale from 'src/components/PayMent/Pay/PayMentModal.vue';
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
            align: 'left',
            style: 'width: 60px'
        },
        {
            name: 'name',
            label: 'Produto',
            field: 'name',
            align: 'left'
        },
        {
            name: 'qtde',
            label: 'Qtde',
            field: 'qtde',
            align: 'center',
            style: 'width: 90px'
        },
        {
            name: 'price',
            label: 'Preço',
            field: 'price',
            align: 'right',
            style: 'width: 100px'
        },
        {
            name: 'total',
            label: 'Total',
            field: 'total',
            align: 'right',
            style: 'width: 110px'
        },
        {
            name: 'actions',
            label: '',
            field: 'actions',
            align: 'center',
            style: 'width: 60px'
        }
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

    const expdandeRows = ref<number[]>([]);

    const hasCharacteristics = (row: any): boolean => {
        return Array.isArray(row.product_with_characteristics) && row.product_with_characteristics.length > 0;
    };

    const isExpanded = (productId: number) => {
        return expdandeRows.value.includes(productId);
    };

    const toggleExpanded = (productId: number): void => {
        if (isExpanded(productId)) {
            expdandeRows.value = expdandeRows.value.filter(id => id !== productId);
            return;
        };

        expdandeRows.value.push(productId);
    }; 

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

    const validateQtde = (val: number, row: SaleItemContract): number => {
        if(!val || val <= 0) {
            row.qtde = 1;
        
            return;
        };

        row.qtde = val;
    };

    const validateGridQtde = (val: number, row: SaleItemContract, index?: number) => { 
        if(!val || val <= 0) {
            row.product_with_characteristics[index].grid_qtde = 1;
        
            return;
        };

        row.product_with_characteristics[index].grid_qtde = val;
        row.qtde = row.product_with_characteristics.reduce((total, a) => total + (a.grid_qtde), 0)
    };

    const pushProducts = (selectedProducts: SaleItemContract[]) => {        
        if(!Array.isArray(productsSale.value)) {
            productsSale.value = [];
            return;
        };

        selectedProducts.forEach(p => {
            const exisit = productsSale.value.find(i => i.product_id === p.id);

            if(exisit && exisit.product_with_characteristics === null)
            {                
                exisit.qtde += p.qtde;

            } else if (exisit && exisit.product_with_characteristics !== null) {
                p.product_with_characteristics.forEach(grid => {
                    const existingGrid = exisit.product_with_characteristics.find(characteristic => characteristic.size === grid.size)

                    if (existingGrid) {
                        existingGrid.grid_qtde += grid.grid_qtde;

                    } else {
                        exisit.product_with_characteristics.push(grid);

                    };
                });

                exisit.qtde = exisit.product_with_characteristics.reduce((total, a) => total + a.grid_qtde, 0);
            } else {
                productsSale.value.push({
                    id: null,
                    product_id: p.id,
                    name: p.name,
                    price: p.price,
                    qtde: p.qtde ?? 1,
                    product_with_characteristics: p.product_with_characteristics ?? []

                });
            };
        });

        disableButtons.editPayMentsForms = true;
        disableButtons.deleteSale = false;
        disableButtons.saveSale = false;
        disableButtons.finallySale = false;
    };

    const calculateTotalItem = (item: SaleItemContract): number => {
        hasCharacteristics(item)
            ? item.product_with_characteristics.reduce((total, item) => {
                return total + (item.grid_qtde ?? 0)
            }, 0)

            : item.qtde

        return item.price * item.qtde;
    };

    const cloneProducts = (items: SaleItemContract[]) =>
        items.map(item => ({ ...item }));

    const calculateTotal = computed(() => {
        const subTotal = productsSale.value.reduce((total, p) => {
            return total + (Number(p.price) * p.qtde);
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

                            resetBtns();

                            return;
                        };
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
        } catch (error: any) {
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

        const existingSaleId: number|null = SessionStorage.getItem('sale_id');
        const existingSale = SessionStorage.getItem('sale') as SaleContract;
        
        if(existingSaleId === null && existingSale === null) 
        {
            return;
        };

        productsSale.value = existingSale?.products || [];

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
</script>

<style lang="scss">
    .section_pdv {    
        margin: 1rem 0 0 0;

        .left-bar {            
            background-color: #fff;
            padding: 15px;
        }

        .right-bar {
            background-color: #fff;
            padding: 15px;
            height: auto;
        }
    }

    @media (min-width: 1536px) {
        body {
            overflow-y: hidden !important;
        }
    }

    @media (min-width: 1536px) {
        body {
            overflow-y: hidden !important;
        }

        .section_pdv {
            display: flex;
            
            .left-bar {            
                width: 80%;   
                display: flex;
                flex-direction: column;
                height: 75vh;

                .sale_products {
                    flex: 1;
                    overflow-y: auto;
                    margin-top: 10px;
                }
            }

            .right-bar {
                width: 30%;
                margin: 0 0 0 2rem;
                display: flex;
                flex-direction: column;
                height: 75vh;        
            }
        }
    }

    @media (max-width: 1536px) {
        .left-bar {     
            width: 100%;        
            height: 75vh;
            display: flex;
            flex-direction: column;

            .sale_products {
                flex: 1;
                overflow-y: auto;
                margin-top: 10px;
            }
        }

        .right-bar {            
            margin: 1.5rem 0 0 0;
            display: flex;
            flex-direction: column;
        }
    }
</style>