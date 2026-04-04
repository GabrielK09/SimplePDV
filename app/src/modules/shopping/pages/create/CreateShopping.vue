<template>
    <q-page padding>
        <header class="border-gray-100 flex">
            <span class="text-black cursor-pointer my-auto">
                <router-link to="/admin/shopping">
                    <q-avatar size="30px" icon="arrow_back" />
                </router-link>
            </span>
        </header>

        <main class="rounded-md flex flex-center text-xl mt-4">
            <section class="rounded-lg px-4">
                <div>
                    <q-table
                        v-model:pagination="pagination"
                        title="Estoque"
                        :rows="productsStockData"
                        :columns="productsColumns"
                    >
                        <template v-slot:top-right>
                            <q-input
                                outlined
                                v-model="searchInput"
                                type="text"
                                label=""
                                @update:model-value="filterProductsStock"
                            >
                                <template v-slot:append>
                                    <q-icon name="search" />
                                </template>
                                <template v-slot:label>
                                    <span class="text-xs">Buscar por um produto ...</span>
                                </template>
                            </q-input>
                        </template>

                        <template v-slot:body-cell-select="props">
                            <q-td :props="props">
                                <q-checkbox
                                    v-model="selectedProductsIds"
                                    :val="props.row.id"
                                    :disable="disableCheckBox(props.row.id)"
                                    dense
                                />
                            </q-td>

                        </template>

                        <template v-slot:no-data>
                            <div class="ml-auto mr-auto">
                                <q-icon name="warning" size="30px"/>
                                <span class="mt-auto mb-auto ml-2 text-xs">Sem produtos cadastrados</span>

                            </div>
                        </template>
                    </q-table>
                </div>

                <div class="flex justify-end py-4 ">
                    <div>
                        <q-btn
                            no-caps
                            color="primary"
                            label="Desassociar todos"
                            :disable="associateProdutcs.length === 0"
                            @click="associateProdutcs = []"
                            class="mr-6"
                        />
                    </div>

                    <div>
                        <q-btn
                            no-caps
                            color="primary"
                            label="Adicionar"
                            :disable="selectedProductsIds.length === 0"
                            @click="associateCheckedProdutcs"
                        />
                    </div>
                </div>

                <div class="associate_product">
                    <q-table
                        v-model:pagination="pagination"
                        title="Produtos associados"
                        :rows="associateProdutcs"
                        :columns="associateProductsColumns"
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

                        <template v-slot:body-cell-purchased_value="props">
                            <q-td :props="props">
                                <q-input
                                    v-model.number="props.row.purchased_value"
                                    type="number"
                                    class="w-12 flex ml-auto mr-auto"
                                    input-class="text-center"
                                    dense
                                    @update:model-value="val => validatePrice(Number(val), props.row)"
                                />
                            </q-td>
                        </template>

                        <template v-slot:body-cell-qtde_purchased="props">
                            <q-td :props="props">
                                <q-input
                                    v-model.number="props.row.qtde_purchased"
                                    type="number"
                                    class="w-12 flex ml-auto mr-auto"
                                    input-class="text-center"
                                    dense
                                    @update:model-value="val => validateQtde(Number(val), props.row)"
                                />
                            </q-td>
                        </template>

                        <template v-slot:body-cell-actions="props">
                            <q-td :props="props">
                                <q-btn
                                    size="10px"
                                    color="black"
                                    icon="edit"
                                    flat
                                    @click="updateAssociateProdutcs(props.row.product_id)"
                                />

                                <q-btn
                                    size="10px"
                                    color="red"
                                    icon="delete"
                                    flat
                                    @click="disassociateCheckedProdutcs(props.row.product_id)"
                                />
                            </q-td>
                        </template>

                        <!--template v-slot:body="props">
                            <q-tr :props="props">
                                <q-td v-for="col in props.cols">
                                    <template v-if="col.name === 'name'">                                    
                                        <span>
                                            {{ `${props.row.name.substring(0, 20)}...` }}

                                            <q-tooltip>
                                                {{ props.row.name }}
                                            </q-tooltip>
                                        </span>
                                    </template>

                                    <template v-else-if="col.name === 'purchased_value'">
                                        <div
                                            class="text-center flex flex-center"
                                        >
                                            <q-input
                                                v-model.number="props.row.purchased_value"
                                                type="number"
                                                class="w-12 flex ml-auto mr-auto"
                                                input-class="text-center"
                                                dense
                                                @update:model-value="val => validatePrice(Number(val), props.row)"
                                            />
                                        </div>
                                    </template>

                                    <template v-else-if="col.name === 'qtde_purchased'">
                                        <div
                                            class="text-center flex flex-center"
                                        >
                                            <q-input
                                                v-model.number="props.row.qtde_purchased"
                                                type="number"
                                                class="w-12 flex ml-auto mr-auto"
                                                input-class="text-center"
                                                dense
                                                @update:model-value="val => validateQtde(Number(val), props.row)"
                                            />
                                        </div>
                                    </template>

                                    <template v-else-if="col.name === 'actions'">
                                        <div v-if="props.row.status === 'Concluída'">
                                            <q-btn
                                                size="10px"
                                                color="red"
                                                icon="delete"
                                                flat
                                                @click="disassociateCheckedProdutcs(props.row.product_id)"
                                            />
                                        </div>
                                    </template>

                                    <template v-else>
                                        <div class="text-center">
                                            {{ col.value }}

                                        </div>
                                    </template>
                                </q-td>
                            </q-tr>
                        </template-->

                        <template v-slot:no-data>
                            <div class="ml-auto mr-auto">
                                <q-icon name="warning" size="30px"/>
                                <span class="mt-auto mb-auto ml-2 text-xs">Sem produtos associados</span>

                            </div>
                        </template>
                    </q-table>
                </div>

                <div class="actions_">
                    <q-btn
                        color="red"
                        no-caps
                        label="Cancelar"
                        @click="cancelShopping"
                        :disable="associateProdutcs.length === 0"
                    />

                    <q-btn
                        outline
                        no-caps
                        class="ml-4"
                        label="Continuar depois"
                        @click="submitShopping(true)"
                        :disable="associateProdutcs.length === 0"
                    />

                    <q-btn
                        color="blue"
                        no-caps
                        class="ml-4"
                        label="Cadastrar produto"
                        @click="showCreateProductComponent = !showCreateProductComponent"
                    />

                    <q-btn
                        color="green"
                        no-caps
                        class="ml-4"
                        label="Confirmar compra"
                        @click="submitShopping(false)"
                        :disable="associateProdutcs.length === 0"
                    />
                </div>
            </section>
        </main>
    </q-page>

    <CreateProductComponent
        v-if="showCreateProductComponent"
        @close="reloadProcuts(!$event)"
    />

    <InformLoad
        v-if="showInformLoadComponent"
        @return:informed-load="saveShopping($event)"
        :shopping-data="shoppingPrePayLoad"
        :last-shopping-id="lastShoppingId"
    />

    <PayMentSale
        v-if="showPayMentForms"
        :shopping-id="shoppingPrePayLoad.id"
        :total-sale="shoppingPrePayLoad.total_shopping"
        @close="showPayMentForms = !$event"
        @paide="finallyShopping(!$event)"
    />

    <QSelectGridTable
        v-if="showSizeGrid"
        :is-just-list="true"
        :characteristics="productFullData.productWithCharacteristics"
        @return:selected-grid="handelSelectedGrid($event)"
    />

    <UpdateAssociateProduct
        v-if="showUpdateAssociateProdutc"
        :product-data="producDataForUpdate"
        @close="showUpdateAssociateProdutc = !$event"
    />
</template>

<script setup lang="ts">
    import { QTableColumn, SessionStorage } from 'quasar';
    import * as ProductsService from 'src/modules/products/services/productsService';
    import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import CreateProductComponent from 'src/components/Products/CreateComponent/CreateProductComponent.vue';
    import InformLoad from 'src/components/Shopping/InformLoad.vue';
    import { useRoute, useRouter } from 'vue-router';
    import PayMentSale from 'src/components/PayMent/Pay/PayMentSale.vue';
    import QSelectGridTable from 'src/components/Products/UseGrid/QTable/QSelectGridTable.vue';
    import { createshopping, getLastShoppingLoad, getShoppingById } from '../../services/shoppingService';
    import UpdateAssociateProduct from 'src/components/Shopping/UpdateItem/UpdateAssociateProduct.vue';

    type ProductResponse = {
        product: ProductContract,
        characteristics: ProductCharacteristicsContract
    };

    const { notify } = useNotify();

    const router = useRouter();
    const route = useRoute();

    const productsColumns: QTableColumn[] = [
        {
            name: 'select',
            label: '',
            field: 'select',
            align: 'left'
        },
        {
            name: 'id',
            label: 'ID',
            field: 'id',
            align: 'center'
        },
        {
            name: 'name',
            label: 'Produto atual',
            field: 'name',
            align: 'center'
        },
        {
            name: 'price',
            label: 'Preço atual',
            field: 'price',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
        {
            name: 'qtde',
            label: 'Qtde atual',
            field: 'qtde',
            align: 'center'
        }
    ];

    const associateProductsColumns: QTableColumn[] = [
        {
            name: 'product_id',
            label: 'ID',
            field: 'product_id',
            align: 'center'
        },
        {
            name: 'name',
            label: 'Produto de entrada',
            field: 'name',
            align: 'center'
        },
        {
            name: 'purchased_value',
            label: 'Preço de entrada',
            field: 'purchased_value',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
        {
            name: 'qtde_purchased',
            label: 'Qtde de entrada',
            field: 'qtde_purchased',
            align: 'center'
        },
        {
            name: 'actions',
            label: '',
            field: 'actions',
            align: 'right'
        }
    ];

    const productsStockData = ref<ProductContract[]>([]);
    const allProductsStockData = ref<ProductContract[]>([]);
    
    const associateProdutcs = ref<ShoppingItemContract[]>([]);
    const intermediaryProductItemData = ref<ShoppingItemContract>();

    const showPayMentForms = ref<boolean>(false);
    const isSavingRef = ref<boolean>(false);
    const showSizeGrid = ref<boolean>(false);
    const lastShoppingId = ref<number | null>(null);

    const producDataForUpdate= ref<ShoppingItemContract>({
        name: null,
        product_id: null,
        purchased_value: null,
        qtde_purchased: null,
        productWithCharacteristics: null,
    });

    const shoppingPrePayLoad = ref<ShoppingContract>({
        id: null,
        load: null,
        shopping_itens: [],
        total_shopping: 0
    });

    const productFullData = ref<ProductContract>({
        id: null,
        name: null,
        price: null,
        qtde: null,
        use_grid: null,
        commission: null,
        productWithCharacteristics: [],
    });

    const indexOfProductHaveCharacteristics = ref<number|null>(null);
    const selectedProductsIds = ref<number[]>([]);

    const searchInput = ref<string>('');

    const pagination = ref({
        sortBy: 'id',
        rowsPerPage: 20
    });

    const showCreateProductComponent = ref<boolean>(false);
    const showInformLoadComponent = ref<boolean>(false);
    const showUpdateAssociateProdutc = ref<boolean>(false);

    const normalizeProduct = (p: ProductContract): ShoppingItemContract => ({
        product_id: p.id,
        name: p.name,
        purchased_value: p.price,
        qtde_purchased: 1,
        //@ts-ignore
        product_with_characteristics: p.product_with_characteristics,
    });

    const removeSessionData = (key: string): void => {
        SessionStorage.remove(key);
    };

    const replaceToShoppingIndex = (): void => {
        router.replace({
            name: 'shopping.index',
        });
    };

    watch(
        () => pagination.value.rowsPerPage,
        async (newRowsPerPage) => {
            await ProductsService.getAll(newRowsPerPage);
        }
    );

    const validateQtde = (val: number, row: SaleItemContract) => {
        if(!val || val <= 0) {
            row.qtde = 1;
            return;
        };

        row.qtde = val;
    };

    const validatePrice = (val: number, row: SaleItemContract) => {
        const afterPrice = productsStockData.value.find(i => i.id === row.id)?.price || 1;

        if(!val || val <= 0) {
            row.price = afterPrice;
            return;
        };

        row.price = val;
    };

    const filterProductsStock = (): void => {
        productsStockData.value = allProductsStockData.value.filter(product => product.name ? product.name.toLowerCase().includes(searchInput.value) : []);
    };

    const associateCheckedProdutcs = () => {
        selectedProductsIds.value.forEach(id => {
            const productData = productsStockData.value.find(p => p.id === id) as ProductContract;

            if(productData === undefined)
            {
                notify(
                    'negative',
                    'Ocorreu um erro ao associar o produto, produto não localizado.'
                );
            };

            const newProductToAssociate: ShoppingItemContract = {
                product_id: productData.id,
                name: productData.name,
                purchased_value: productData.price,
                qtde_purchased: 1,
                productWithCharacteristics: productData.productWithCharacteristics
            };

            associateProdutcs.value.push(newProductToAssociate);
        });

        searchInput.value = '';
        productsStockData.value = allProductsStockData.value;
    };

    watch(associateProdutcs.value, (newProducts) => {
        newProducts.forEach(p => {            
            if(p.productWithCharacteristics !== null) 
            {   
                const index = newProducts.findIndex(item => item.product_id === p.product_id);

                console.log(index);
                
                if (index > -1)
                {
                    indexOfProductHaveCharacteristics.value = index;
                    
                    productFullData.value = {
                        commission: 0,
                        id: p.product_id,
                        name: p.name,
                        price: p.purchased_value,
                        qtde: 1,
                        //@ts-ignore
                        productWithCharacteristics: p.productWithCharacteristics
                    };

                    intermediaryProductItemData.value = {
                        name: productFullData.value.name,
                        product_id: productFullData.value.id,
                        purchased_value: productFullData.value.price,
                        qtde_purchased: 1,
                        productWithCharacteristics: productFullData.value.productWithCharacteristics
                    };

                    showSizeGrid.value = true;
                };
            };
        });
    }, {    
        immediate: true 
    });

    const disassociateCheckedProdutcs = (id: number) => {
        if(associateProdutcs.value.length <= 1)
        {
            associateProdutcs.value = []
            return
        };

        associateProdutcs.value = associateProdutcs.value.filter(ap => ap.product_id !== id);

        searchInput.value = '';
    };

    const getAllProductsStock = async (rowsPerPage?: number) => {
        const res = await ProductsService.getAll(rowsPerPage);
        
        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
        };

        console.log(res.data);

        productsStockData.value = res.data.map((c: ProductResponse) => ({
            id: c.product.id,
            name: c.product.name,
            price: c.product.price,
            qtde: c.product.qtde,
            commission: c.product.commission,
            productWithCharacteristics: c.characteristics
        }));
                
        allProductsStockData.value = [...productsStockData.value];
    };

    const cancelShopping = (): void => {
        replaceToShoppingIndex();
    };

    const handelSelectedGrid = (grid: any) => {
        if(!intermediaryProductItemData.value) return;

        const parsedProduct: ShoppingItemContract = {
            product_id: intermediaryProductItemData.value.product_id,
            name: intermediaryProductItemData.value.name,
            purchased_value: intermediaryProductItemData.value.purchased_value,
            qtde_purchased: intermediaryProductItemData.value.qtde_purchased,
            productWithCharacteristics: grid
        };

        const originalProduct = associateProdutcs.value.findIndex(ap => ap.product_id === intermediaryProductItemData.value?.product_id);

        if(originalProduct > -1)
        {
            associateProdutcs.value.splice(originalProduct, 1);
            associateProdutcs.value.push(parsedProduct);

            showSizeGrid.value = false;
            return;
        };

        notify(
            'negative',
            'Erro ao validar os dados, produto não localizado.'
        );
    };

    const submitShopping = async (isSaving: boolean) => {
        isSavingRef.value = isSaving;

        const existingLoad = SessionStorage.getItem('inform_load');

        if(existingLoad || shoppingPrePayLoad.value.id)
        {
            saveShopping(Number(existingLoad));
            return;
        };

        showInformLoadComponent.value = true;

        //@ts-ignore
        const total_shopping = associateProdutcs.value.reduce((total, a) => total + (a.purchased_value * a.qtde_purchased), 0);

        shoppingPrePayLoad.value = {
            id: 0,
            load: 0,
            shopping_itens: associateProdutcs.value,
            total_shopping: total_shopping
        };
    };

    const saveShopping = async (informLoad: number): Promise<void> => {
        let existingShoppingId: number;

        if (routeShoppingId.value)
        {
            existingShoppingId = routeShoppingId.value;
        } else {
            existingShoppingId = SessionStorage.getItem('shopping_id') as number;
        };

        showInformLoadComponent.value = false;
        SessionStorage.set('inform_load', informLoad);

        const existingLoad = SessionStorage.getItem('inform_load') as number;

        const payload: ShoppingContract = {
            id: null,
            load: informLoad,
            shopping_itens: shoppingPrePayLoad.value.shopping_itens,
            total_shopping: shoppingPrePayLoad.value.total_shopping
        };

        if (existingShoppingId && existingLoad)
        {
            console.log(`existingShoppingId = ${existingShoppingId} && existingLoad = ${existingLoad} foi true`);
            shoppingPrePayLoad.value.id = existingShoppingId;
            shoppingPrePayLoad.value.load = existingLoad;

            console.log('PayLoad: ', shoppingPrePayLoad.value);

            if(!isSavingRef.value)
            {
                showPayMentForms.value = true;
                return;
            };

            notify(
                'positive',
                'Compra salva com sucesso!'
            );

            replaceToShoppingIndex();

            return;
        };

        const res = await createshopping(payload);
        const data: number = res.data;

        if(res.success)
        {
            console.log('res foi true');
            SessionStorage.set('shopping_id', data);

            notify(
                'positive',
                res.message
            );

            if(!isSavingRef.value)
            {
                showPayMentForms.value = true;
                shoppingPrePayLoad.value.id = data;
                return;
            };

            notify(
                'positive',
                'Compra salva com sucesso!'
            );

            replaceToShoppingIndex();

            return;
        };

        notify(
            'negative',
            res.message
        );

        return;
    };

    const finallyShopping = (event: boolean) => {
        showPayMentForms.value = event;

        replaceToShoppingIndex();

        removeSessionData('shopping_id');
        removeSessionData('inform_load');
    };

    const routeShoppingId = computed(() => {
        const id = route.query.id;

        if (Array.isArray(id)) return Number(id[0]) || null;
        if (id === null || id === undefined || id === '') return null;

        const parsed = Number(id);
        return Number.isNaN(parsed) ? null : parsed;
    });

    const reloadProcuts = async (event: boolean): Promise<void> => {
        showCreateProductComponent.value = event;
        await getAllProductsStock();
    };

    const updateAssociateProdutcs = (id: number) => {
        const producDataById = associateProdutcs.value.find(p => p.product_id === id);

        if(!producDataById)
        {
            notify(
                'negative',
                'Um erro ocorreu, produto não localizado'
            );
            return;
        };

        producDataForUpdate.value = producDataById;
        showUpdateAssociateProdutc.value = true;
    };

    const disableCheckBox = (productId: number): boolean => {
        if(!associateProdutcs.value.map(a => a.product_id).includes(productId))
            return false;

        return true;
    };

    onMounted(async () => {
        await getAllProductsStock();

        const res = await getLastShoppingLoad();

        if(!res.success)
        {
            notify(
                'negative',
                res.message || 'Erro interno'
            );

            return;
        };

        if (routeShoppingId.value)
        {
            notify(
                'info',
                'Carregando dados da compra...'
            );

            const res = await getShoppingById(routeShoppingId.value);

            if(!res.success)
            {
                notify(
                    'warning',
                    res.message || 'Erro ao carregar os dados da compra.'
                );

                return;
            };

            const shoppingData: ShoppingContract = res.data.shopping;
            const shoppingItensData: ShoppingItemContract[] = res.data.shoppingWithProducts;

            associateProdutcs.value = shoppingItensData;
            shoppingPrePayLoad.value = {
                id: shoppingData.id,
                load: shoppingData.load,
                shopping_itens: associateProdutcs.value,
                total_shopping: shoppingData.total_shopping

            };

            lastShoppingId.value = shoppingData.id;

            return;
        };

        lastShoppingId.value = res.data;
    });

    onUnmounted(() => {
        removeSessionData('shopping_id');
        removeSessionData('inform_load');
        productsStockData.value = [];
        isSavingRef.value = false;
    });
 </script>

 <style>
    .actions_ {
        display: flex;
        justify-content: end;
        margin: 15px 0 0 0;

    }
</style>
