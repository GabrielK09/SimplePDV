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
            <section class="w-[80vh] rounded-lg px-4">
                <div>
                    <q-table
                        v-model:pagination="pagination"
                        title="Estoque"
                        :rows="productsStockData"
                        :columns="productsColumns"
                    >
                        <template v-slot:top-selection>
                            <span>a</span>
                        </template>

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
                                    v-model="selectedProducts"
                                    :val="props.row"
                                    :disable="associateProdutcs.find(ap => ap.product_id === props.row.id) !== undefined ? true : false"
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
                            :disable="selectedProducts.length === 0"
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
                        <template v-slot:body="props">
                            <q-tr :props="props">
                                <q-td v-for="col in props.cols">
                                    <template v-if="col.name === 'purchased_value'">
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
                        </template>

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

</template>

<script setup lang="ts">
    import { QTableColumn, SessionStorage } from 'quasar';
    import * as ProductsService from 'src/modules/products/services/productsService';
    import { computed, onMounted, onUnmounted, ref } from 'vue';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import CreateProductComponent from 'src/components/Products/CreateComponent/CreateProductComponent.vue';
    import InformLoad from 'src/components/Shopping/InformLoad.vue';
    import { useRoute, useRouter } from 'vue-router';
    import PayMentSale from 'src/components/PayMent/Pay/PayMentSale.vue';
    import { createshopping, getLastShoppingLoad, getShoppingById } from '../../services/shoppingService';

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
    const showPayMentForms = ref<boolean>(false);
    const isSavingRef = ref<boolean>(false);
    const lastShoppingId = ref<number | null>(null);

    const shoppingPrePayLoad = ref<ShoppingContract>({
        id: null,
        load: null,
        shopping_itens: [],
        total_shopping: 0
    });

    const selectedProducts = ref<ProductContract[]>([]);

    const searchInput = ref<string>('');

    const pagination = ref({
        sortBy: 'id'
    });

    const showCreateProductComponent = ref<boolean>(false);
    const showInformLoadComponent = ref<boolean>(false);

    const removeSessionData = (key: string): void => {
        SessionStorage.remove(key);
    };

    const replaceToShoppingIndex = (): void => {
        router.replace({
            name: 'shopping.index',
        });
    };

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
        productsStockData.value = allProductsStockData.value.filter(product => product.name.toLowerCase().includes(searchInput.value));
    };

    const associateCheckedProdutcs = () => {
        selectedProducts.value.forEach(p => {
            if(associateProdutcs.value.find(ap => ap.product_id === p.id))
            {
                notify(
                    'info',
                    'Produto já associado!'
                );

                return;
            };

            const newProductStock: ShoppingItemContract = {
                product_id: p.id,
                name: p.name,
                purchased_value: p.price,
                qtde_purchased: 1,
            };

            associateProdutcs.value.push(newProductStock);
        });

        selectedProducts.value = [];
        searchInput.value = '';
        productsStockData.value = allProductsStockData.value;

    };

    const disassociateCheckedProdutcs = (id: number) => {
        console.log('call disassociateCheckedProdutcs');

        if(associateProdutcs.value.length <= 1)
        {
            associateProdutcs.value = []
            return
        };

        console.log(`Deve remover o id: ${id}`);

        associateProdutcs.value = associateProdutcs.value.filter(ap => ap.product_id !== id);

        selectedProducts.value = [];
        searchInput.value = '';
    };

    const getAllProductsStock = async () => {
        const res = (await ProductsService.getAll()).data;
        const data = res.product;

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
        };

        productsStockData.value = data;
        allProductsStockData.value = [...productsStockData.value];
    };

    const cancelShopping = (): void => {
        replaceToShoppingIndex();
    };

    const submitShopping = async (isSaving: boolean) => {
        isSavingRef.value = isSaving;

        const existingLoad: number = SessionStorage.getItem('inform_load');

        console.log(shoppingPrePayLoad.value);

        if(existingLoad || shoppingPrePayLoad.value.id)
        {
            saveShopping(existingLoad);
            return;
        };

        showInformLoadComponent.value = true;

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
            existingShoppingId = SessionStorage.getItem('shopping_id');
        };

        showInformLoadComponent.value = false;
        SessionStorage.set('inform_load', informLoad);

        const existingLoad: number = SessionStorage.getItem('inform_load');

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
                'Carregando dados da compra.'
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

            console.log(shoppingPrePayLoad.value);

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
