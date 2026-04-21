<template>
    <q-page padding id="q-page">
        <div class="m-2 text-3xl" >
            <main class="rounded-md flex flex-center text-xl mt-4">
                <section class="rounded-lg px-4">
                    <div>
                        <q-table
                            v-model:pagination="pagination"
                            :rows="productsStockData"
                            :columns="productsColumns"
                        >
                            <template v-slot:top>
                                <div class="w-full">
                                    <div class="mx-auto flex justify-between">
                                        <router-link to="/admin/shopping">
                                            <q-avatar size="45px" icon="arrow_back" />
                                        </router-link>
    
                                        <q-input
                                            outlined
                                            v-model="searchInput"
                                            type="text"
                                            label=""
                                            dense
                                            @update:model-value="filterProductsStock"
                                        >
                                            <template v-slot:append>
                                                <q-icon name="search" />
                                            </template>
                                            <template v-slot:label>
                                                <span class="text-xs">Buscar por um produto ...</span>
                                            </template>
                                        </q-input>
    
                                    </div>
                                </div>
                            </template>

                            <template v-slot:body="props">
                                <q-tr :props="props">
                                    <q-td v-for="col in props.cols">                           
                                        <template v-if="col.name === 'select'">
                                            <q-radio
                                                id="selected-product-radio"
                                                v-model="selectedProductsId"
                                                :val="props.row.id"
                                                dense
                                            />
                                        </template>

                                        <template v-if="col.name === 'name'">
                                            {{ props.row.name }}
                                        </template>

                                        <template v-else>
                                            <span>
                                                {{ col.value }}
                                            </span>
                                        </template>
                                    </q-td>
                                </q-tr>
                            </template>

                            <template v-slot:no-data>
                                <div class="ml-auto mr-auto">
                                    <q-icon name="warning" size="30px"/>
                                    <span class="mt-auto mb-auto ml-2 text-xs">Sem produtos cadastrados</span>

                                </div> 
                            </template>
                        </q-table>
                    </div>

                    <div class="flex justify-end py-4">
                        <div>
                            <q-btn
                                no-caps
                                color="red"
                                label="Desassociar todos"
                                :disable="associateProducts.length === 0"
                                @click="associateProducts = []"
                                class="mr-6"
                            />
                        </div>

                        <div>
                            <q-btn
                                no-caps
                                color="primary"
                                label="Adicionar"
                                :disable="selectedProductsId === 0"
                                @click="associateCheckedProdutcs"
                            />
                        </div>
                    </div>

                    <div class="associate_product">
                        <q-table
                            v-model:pagination="pagination"
                            title="Produtos associados"
                            :rows="associateProducts"
                            :columns="associateProductsColumns"
                            row-key="product_id"
                        >
                            <template v-slot:body="props">
                                <q-tr :props="props">
                                    <q-td key="name" :props="props">
                                        <span>
                                            {{ `${props.row.name.substring(0, 20)}...` }}

                                            <q-tooltip>
                                                {{ props.row.name }} - {{ props.row?.product_with_characteristics[0].size }}
                                            </q-tooltip>
                                        </span>
                                    </q-td>

                                    <q-td key="purchased_value" :props="props">
                                        <q-input
                                            v-model.number="props.row.purchased_value"
                                            type="number"
                                            class="w-12 flex ml-auto mr-auto"
                                            input-class="text-center"
                                            dense
                                            @update:model-value="val => validatePrice(Number(val), props.row)"
                                        />
                                    </q-td>

                                    <q-td key="qtde_purchased" :props="props">
                                        <q-input
                                            v-model.number="props.row.qtde_purchased"
                                            type="number"
                                            class="w-12 flex ml-auto mr-auto"
                                            input-class="text-center"
                                            dense
                                            @update:model-value="val => validateQtde(Number(val), props.row)"
                                            :disable="props.row.product_with_characteristics && props.row.product_with_characteristics.length > 0"
                                        />
                                    </q-td>

                                    <q-td key="actions" :props="props">
                                        <div class="flex flex-row items-center gap-1">
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

                                            <q-btn
                                                size="10px"
                                                color="red"
                                                icon="delete"
                                                flat
                                                @click="disassociateCheckedProdutc(props.row.product_id)"
                                            />
                                        </div>  
                                    </q-td>
                                </q-tr>

                                <q-tr 
                                    v-if="isExpanded(props.row.product_id) && hasCharacteristics(props.row)"
                                    :props="props"
                                >
                                    <q-td colspan="100%" class="bg-gray-200">
                                        <div class="q-pa-md">
                                            <div class="text-subtitle2 text-weight-bold q-mb-sm">
                                                Grade do produto
                                            </div>

                                            <div class="row q-col-gutter-sm">
                                                <div
                                                    v-for="(characteristic, i) in props.row.product_with_characteristics"
                                                    :key="`${props.row.product_id}-${characteristic.size}`"
                                                    class="col-12 col-sm-6 col-md-3"
                                                >   
                                                    <q-card flat bordered>
                                                        <q-card-section class="q-pa-pmd">
                                                            <div class="text-caption text-gray-700">
                                                                Tamanho

                                                                <q-btn
                                                                    size="10px"
                                                                    color="red"
                                                                    icon="close"
                                                                    flat
                                                                    @click="disassociateCheckedGridProdutc(props.rowIndex, Number(i))"
                                                                />
                                                            </div>
                                                            <div class="text-body2 text-weight-bold">
                                                                {{ characteristic.size }}
                                                            </div>

                                                            <div class="text-caption text-grey-7 q-mt-sm">Quantidade</div>
                                                            <div>
                                                                <q-input
                                                                    v-model.number="characteristic.grid_qtde"
                                                                    type="number"
                                                                    class="w-12 flex ml-auto mr-auto"
                                                                    input-class="text-center"
                                                                    dense
                                                                    @update:model-value="val => validateQtde(Number(val), props.row, Number(i))"
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

                            <template v-slot:no-data>
                                <div class="ml-auto mr-auto">
                                    <q-icon name="warning" size="30px"/>
                                    <span class="mt-auto mb-auto ml-2 text-xs">Sem produtos associados</span>

                                </div>
                            </template>
                        </q-table>
                    </div>

                    <div class="bg-white rounded-lg h-auto p-4 w-full laptop:mr-6 flex flex-col mt-4">
                        <div class="flex-1 overflow-y-auto">
                            <div class="mt-4 border p-2 rounded">
                                Total da compra R$ {{ totalShopping }}
                            </div>
                        </div>
                    </div>

                    <div class="actions_">
                        <q-btn
                            color="red"
                            no-caps
                            label="Cancelar"
                            @click="replaceToShoppingIndex"
                            :disable="associateProducts.length === 0"
                        />

                        <q-btn
                            outline
                            no-caps
                            class="ml-4"
                            label="Continuar depois"
                            @click="submitShopping(true)"
                            :disable="associateProducts.length === 0"
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
                            :disable="associateProducts.length === 0"
                        />
                    </div>
                </section>
            </main>
        </div>

        <div id="action_bar">
            <div class="bg-[#03202e] w-[60%] mb-4 p-3 text-white flex justify-end rounded-lg">
                <q-btn 
                    no-caps
                    color="green" 
                    label="Confirmar compra" 
                    class="mr-6"
                    @click="scrollToDown" 
                    v-if="associateProducts.length > 0"
                />

                <q-btn
                    no-caps
                    dense
                    color="red"
                    label="Cancelar seleção"
                    :disable="selectedProductsId === 0"
                    @click="selectedProductsId = 0"
                    class="mr-6"
                />

                <q-btn
                    v-if="associateProducts.length > 0"
                    no-caps
                    dense
                    color="blue"
                    label="Desassociar todos"
                    @click="associateProducts = []"
                    class="mr-6"
                />

                <q-btn
                    no-caps
                    color="primary"
                    dense
                    label="Adicionar"
                    :disable="selectedProductsId === 0"
                    @click="associateCheckedProdutcs"
                />
            </div>
        </div>
    </q-page>

    <CreateProductComponent
        v-if="showCreateProductComponent"
        @close="reloadProcuts(!$event)"
    />

    <InformLoad
        v-if="showInformLoadComponent"
        @return:informed-load="saveShopping($event)"
        @close="showInformLoadComponent = !$event"
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
        :characteristics="productCharacteristics"
        @return:selected-grid="handelSelectedGrid($event)"
        @close="showSizeGrid = !$event"
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
    import { createshopping, getShoppingById, updateShoppingDetails } from '../../services/shoppingService';

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
            label: 'Cód produto',
            field: 'id',
            align: 'center'
        },
        {
            name: 'name',
            label: 'Produto',
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
            label: 'Ações',
            field: 'actions',
            align: 'right'
        }
    ];

    const productsStockData = ref<ProductContract[]>([]);
    const allProductsStockData = ref<ProductContract[]>([]);
    
    const originalProductsShopping = ref<ShoppingItemContract[]>([]);
    const associateProducts = ref<ShoppingItemContract[]>([]);
    const intermediaryProductItemData = ref<ShoppingItemContract>();

    const showPayMentForms = ref<boolean>(false);
    const isSavingRef = ref<boolean>(false);
    const showSizeGrid = ref<boolean>(false);

    const shoppingPrePayLoad = ref<ShoppingContract>({
        id: null,
        load: null,
        shopping_itens: [],
        total_shopping: 0
    });

    const productCharacteristics = ref<ProductCharacteristicsContract[]>([]);
    
    const expdandeRows = ref<number[]>([]);
    const selectedProductsId = ref<number|null>(null);

    const searchInput = ref<string>('');

    const pagination = ref({
        sortBy: 'id',
        rowsPerPage: 20
    });

    const showCreateProductComponent = ref<boolean>(false);
    const showInformLoadComponent = ref<boolean>(false);

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

    watch(
        () => selectedProductsId.value,
        (newId) => {
            const actionBar = document.getElementById('action_bar') as HTMLElement;

            if(newId !== 0 && newId !== null)
            {                
                actionBar.style.display = 'flex';
            } else {
                actionBar.style.display = 'none';  
            };

            return;
        }
    );

    const scrollToDown = () => {
        const actionBar = document.getElementById('action_bar') as HTMLElement;

        window.scrollBy({
            top: document.documentElement.scrollHeight,
            behavior: 'smooth'
        });

        actionBar.style.display = 'none';
        selectedProductsId.value = null;
    };

    const validateToReturnQtde = (val: number): number => {
        if(!val || val <= 0) {
            return 1;
        };

        return val;
    };

    const validateQtde = (val: number, row: ShoppingItemContract, index?: number) => { 
        const qtde = validateToReturnQtde(val);

        if(index < 0)
        {   
            row.qtde_purchased = qtde;
        } else {          
            row.product_with_characteristics[index].grid_qtde = qtde;

            row.qtde_purchased = row.product_with_characteristics.reduce((total, a) => total + (a.grid_qtde), 0);

        };
    };

    const validatePrice = (val: number, row: ShoppingItemContract) => {
        const product = productsStockData.value.find(i => i.id === row.product_id);

        if(!val || val <= 0) {            
            row.purchased_value = product.price || 1;
            return;
        };
        
        row.purchased_value = val;
    };

    const filterProductsStock = (): void => {
        productsStockData.value = allProductsStockData.value.filter(product => product.name ? product.name.toLowerCase().includes(searchInput.value) : []);
    };

    const checkEveryGridsRegistred = (currentGrids: ProductCharacteristicsContract[], associateProductById: ShoppingItemContract): ProductCharacteristicsContract[] => {
        return currentGrids.filter(c => !associateProductById.product_with_characteristics.map(ap => ap.size).includes(c.size));
    };

    const associateCheckedProdutcs = () => {
        const productStockData: ProductContract = productsStockData.value.find(p => p.id === selectedProductsId.value);

        if(productStockData === undefined)
        {
            notify(
                'negative',
                'Ocorreu um erro ao associar o produto, produto não localizado.'
            );

            return;
        };

        notify(
            'info',
            'Produto adicionado com sucesso!'
        );
            
        if (productStockData.product_with_characteristics !== null)
        {
            const characteristics = productStockData.product_with_characteristics;
            const associateProduct = associateProducts.value.find(ap => ap.product_id === productStockData.id);

            if(associateProduct)
            {   
                productCharacteristics.value = checkEveryGridsRegistred(characteristics, associateProduct);

            } else {
                productCharacteristics.value = characteristics;
            };

            intermediaryProductItemData.value = {
                name: productStockData.name,
                product_id: productStockData.id,
                purchased_value: productStockData.price,
                qtde_purchased: 1,
                product_with_characteristics: productStockData.product_with_characteristics
            };

            showSizeGrid.value = true;

            return;
        };

        const newProductToAssociate: ShoppingItemContract = {
            product_id: productStockData.id,
            name: productStockData.name,
            purchased_value: productStockData.price,
            qtde_purchased: 1,
            product_with_characteristics: productStockData.product_with_characteristics ?? null
        };

        associateProducts.value.push(newProductToAssociate);

        selectedProductsId.value = null;

        searchInput.value = '';
    };

    const disassociateCheckedProdutc = (id: number) => {
        associateProducts.value = associateProducts.value.filter(ap => ap.product_id !== id);

        searchInput.value = '';
    };

    const disassociateCheckedGridProdutc = (indexProduct: number, indexGrid: number): void => {
        const product = associateProducts.value[indexProduct];

        if(!associateProducts.value[indexProduct] && !product.product_with_characteristics[indexGrid])
        {
            notify(
                'negative',
                'Erro ao deletar a grade, grade não localizada.'
            );

            return;
        };

        associateProducts.value[indexProduct].product_with_characteristics.splice(indexGrid, 1);

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

        productsStockData.value = res.data.map((c: ProductResponse) => ({
            id: c.product.id,
            name: c.product.name,
            price: c.product.price,
            qtde: c.product.qtde,
            commission: c.product.commission,
            product_with_characteristics: c.characteristics
        }));
                
        allProductsStockData.value = [...productsStockData.value];
    };

    const handelSelectedGrid = (grid: ProductCharacteristicsContract) => {
        showSizeGrid.value = false;
        
        if(!intermediaryProductItemData.value)     
            return;
        
        const parsedProduct: ShoppingItemContract = {
            product_id: intermediaryProductItemData.value.product_id,
            name: intermediaryProductItemData.value.name,
            purchased_value: intermediaryProductItemData.value.purchased_value,
            qtde_purchased: intermediaryProductItemData.value.qtde_purchased,
            product_with_characteristics: [{
                id: grid.id, 
                product_id: intermediaryProductItemData.value.product_id, 
                size: grid.size,
                grid_qtde: 1, 
            }]
        };

        const existingAssociateProduct = associateProducts.value.find(p => p.product_id === intermediaryProductItemData.value.product_id);

        if (!existingAssociateProduct)
        {
            associateProducts.value.push(parsedProduct);
            return;
        };

        existingAssociateProduct.product_with_characteristics.push({
            id: grid.id, 
            product_id: intermediaryProductItemData.value.product_id, 
            size: grid.size,
            grid_qtde: 1, 
        });

        selectedProductsId.value = null;
    };

    const hasProductChanged = (): boolean => {
        notify(
            'info',
            'Conferindo itens da compra.'
        );

        if (associateProducts.value.length !== originalProductsShopping.value.length) {
            return true;
        };

        const currentMap = new Map(
            associateProducts.value.map(item => [
                item.product_id,
                { qtde_purchased: item.qtde_purchased, purchased_value: item.purchased_value }
            ])
        );

        for (const oldItem of originalProductsShopping.value) {
            const current = currentMap.get(oldItem.product_id);

            if (!current) return true;
            
            if (
                current.qtde_purchased !== oldItem.qtde_purchased || current.purchased_value !== oldItem.purchased_value
            ) return true;
        };

        return false;
    };

    const submitShopping = async (isSaving: boolean): Promise<void> => 
    {
        isSavingRef.value = isSaving;
        const exisitInformedLoad = Number(SessionStorage.getItem('inform_load'));
        
        if (exisitInformedLoad > 0)
        {
            saveShopping(exisitInformedLoad);  
            return;
        };

        showInformLoadComponent.value = true;
    };

    const cloneProducts = (items: ShoppingItemContract[]) =>
        items.map(item => ({ ...item }));
        
    /**
     * if isSaving = true, just save and return to shopping.index
     * else isSaving = false, continue to payment
     * 
     */
    const saveShopping = async (informLoad: number): Promise<void> =>
    {        
        // Constroi o prePayload da compra        
        shoppingPrePayLoad.value = {
            id: null,
            load: informLoad,
            shopping_itens: associateProducts.value,
            total_shopping: Number(totalShopping.value.replace(',', '.'))
        };

        showInformLoadComponent.value = false;
        // Precisa conferir antes se a venda vai ser apenas para salvar ou se é para finalizar, seguindo fluxos diferentes para cada.
        // isSaving = true        
        if(isSavingRef.value)
        {            
            const existingShoppingId = Number(SessionStorage.getItem('shopping_id'));

            if((routeShoppingId.value > 0 && shoppingPrePayLoad.value.load > 0) || existingShoppingId > 0)
            {
                shoppingPrePayLoad.value.id = routeShoppingId.value > 0 ? routeShoppingId.value : existingShoppingId;
                if(hasProductChanged())
                {   
                    const res = await updateShoppingDetails(shoppingPrePayLoad.value);

                    if(!res.success)
                    {
                        notify(
                            'negative',
                            res.message
                        );
                        return;
                    };
                };

                const res = await updateShoppingDetails(shoppingPrePayLoad.value);

                if (!res.success)
                {
                    notify(
                        'negative',
                        res.message
                    );

                    removeSessionData('shopping_id');
                    removeSessionData('inform_load');

                    return;
                };

                notify(
                    'positive',
                    res.message
                );
                
                replaceToShoppingIndex();
                return;
            } else {
                const res = await createshopping(shoppingPrePayLoad.value);

                if (!res.success)
                {
                    notify(
                        'negative',
                        res.message
                    );

                    removeSessionData('shopping_id');
                    removeSessionData('inform_load');

                    return;
                };

                notify(
                    'positive',
                    res.message
                );
                
                replaceToShoppingIndex();

                return;
            };    
        } else {
            // Não é para salvar, é para finalizar
            
            const existingShoppingId = Number(SessionStorage.getItem('shopping_id'));

            if (routeShoppingId.value > 0 || existingShoppingId > 0)
            {   
                shoppingPrePayLoad.value.id = routeShoppingId.value > 0 ? routeShoppingId.value : existingShoppingId;
                showPayMentForms.value = true;

                if(hasProductChanged())
                {
                    const res = await updateShoppingDetails(shoppingPrePayLoad.value);

                    if(!res.success)
                    {
                        notify(
                            'negative',
                            res.message
                        );
                        return;
                    };
                };

                return;

            } else {
                SessionStorage.set('inform_load', informLoad);

                const res = await createshopping(shoppingPrePayLoad.value);
                
                if (!res.success)
                {
                    type Errors = {
                        message: string;
                    };

                    const errors: Errors[] = [];
                    const data = res.data.data.data; // Erradíssimo mas funciona
                    
                    Object.entries(data).forEach(([_, message]) => {
                        errors.push({
                            message: String(message)
                        });
                    });

                    notify(
                        'negative',
                        res.message
                    );

                    errors.forEach(e => {
                        notify(
                            'negative',
                            e.message
                        );                        
                    });
                
                    removeSessionData('shopping_id');
                    removeSessionData('inform_load');

                    return;
                };

                const returningSaleId = res.data

                notify(
                    'positive',
                    res.message
                );
                
                SessionStorage.set('shopping_id', returningSaleId);
                
                shoppingPrePayLoad.value.id = returningSaleId;

                showPayMentForms.value = true;
                return;
            
            };
        };   
    };
    
    const finallyShopping = (event: boolean) => {
        showPayMentForms.value = event;

        removeSessionData('shopping_id');
        removeSessionData('inform_load');
        
        replaceToShoppingIndex();
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

    const totalShopping = computed(() => {
        const subTotal = associateProducts.value.reduce((total, p) => {
            return total + (Number(p.purchased_value) * p.qtde_purchased)
        }, 0);

        return subTotal.toFixed(2).replace('.', ',');
    });

    window.addEventListener('scroll', () => {    
        const currentScrool = window.scrollY;
        const screenHeight = window.innerHeight;
        const totalHeight = document.documentElement.scrollHeight;        
        const actionBar = document.getElementById('action_bar') as HTMLElement;

        if (actionBar)
        {
            if(Math.ceil(currentScrool + screenHeight) >= totalHeight)
            {
                actionBar.style.display = 'none';
                return;
            };
            
            if (selectedProductsId.value !== null || selectedProductsId.value > 0)
            {
                actionBar.style.display = 'flex';
                return;
            };
        };
    });

    onMounted(async () => {
        await getAllProductsStock();

        if (routeShoppingId.value)
        {
            notify(
                'info',
                'Carregando dados da compra...'
            );

            const res = await getShoppingById(routeShoppingId.value);

            console.log(res);

            if(!res.success)
            {
                notify(
                    'warning',
                    res.message || 'Erro ao carregar os dados da compra.'
                );

                return;
            };

            associateProducts.value = res.data.shopping_with_products;
                
            SessionStorage.set('inform_load', res.data.shopping_data.load);

            shoppingPrePayLoad.value = {
                id: res.data.shopping_data.id,
                load: res.data.shopping_data.load,
                shopping_itens: associateProducts.value,
                total_shopping: res.data.shopping_data.total_shopping
            };

            originalProductsShopping.value = cloneProducts(associateProducts.value);
            
            return;
        };
    });

    onUnmounted(() => {
        removeSessionData('shopping_id');
        removeSessionData('inform_load');
    });
 </script>

<style>
    .actions_ {
        display: flex;
        justify-content: end;
        margin: 15px 0 0 0;
    }

    @keyframes up_action_bar {
        from {
            transform: translateY(100%);
            opacity: 0;
        }

        to {
            transform: translateY(0);
            opacity: 1;
        }
    }

    @keyframes hidde_action_bar {
        from { transform: translateY(0); }

        to { transform: translateY(100%); }
    }

    #action_bar {
        position: sticky;
        z-index: 3000;
        bottom: 0;
        display: none;
        justify-content: center;
        animation: up_action_bar 0.4s ease-in-out;   
    }

    #action_bar.hidde {
        animation: hidde_action_bar 0.4s ease-in-out;
    }
</style>
