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
                            label="Desassociar"
                            :disable="associateProdutcs.length === 0"
                            @click="disassociateCheckedProdutcs"
                            class="mr-6"
                        />
                    </div>

                    <div>
                        <q-btn
                            no-caps
                            color="primary"
                            label="Associar"
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
                        <template v-slot:body-cell-select="props">
                            <q-td :props="props">
                                <q-checkbox
                                    v-model="selectedProducts"
                                    :val="props.row"
                                    dense
                                />
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

                        <template v-slot:no-data>
                            <div class="ml-auto mr-auto">
                                <q-icon name="warning" size="30px"/>
                                <span class="mt-auto mb-auto ml-2 text-xs">Sem produtos associados</span>

                            </div>
                        </template>
                    </q-table>
                </div>

                <div class="actions">
                    <q-btn
                        color="red"
                        no-caps
                        label="Cancelar"
                        @click="cancelShopping"
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
                        @click="submitShopping"
                        :disable="associateProdutcs.length === 0"
                    />
                </div>
            </section>
        </main>
    </q-page>

    <CreateProductComponent
        v-if="showCreateProductComponent"
        @close="showCreateProductComponent = !$event"
    />

    <InformLoad
        v-if="showInformLoadComponent"
        :shopping-data="shoppingPrePayLoad"
    />
</template>

<script setup lang="ts">
    import { QTableColumn, SessionStorage } from 'quasar';
    import * as ProductsService from 'src/modules/products/services/productsService';
    import { onMounted, ref } from 'vue';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import CreateProductComponent from 'src/components/Products/CreateProductComponent.vue';
    import InformLoad from 'src/components/Shopping/InformLoad.vue';

    const { notify } = useNotify();

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
        },
    ];

    const associateProductsColumns: QTableColumn[] = [
        {
            name: 'select',
            label: '',
            field: 'select',
            align: 'left'
        },
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
    ];

    const productsStockData = ref<ProductContract[]>([]);
    const allProductsStockData = ref<ProductContract[]>([]);

    const associateProdutcs = ref<ShoppingItemContract[]>([]);

    const searchInput = ref<string>('');

    const pagination = ref({
        sortBy: 'id'
    });

    const showCreateProductComponent = ref<boolean>(false);
    const showInformLoadComponent = ref<boolean>(false);

    const shoppingPrePayLoad = ref<ShoppingContract>();

    const selectedProducts = ref<ProductContract[]>([]);

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

    const disassociateCheckedProdutcs = () => {
        selectedProducts.value.forEach(p => {
            associateProdutcs.value = associateProdutcs.value.filter(ap => p.id !== ap.product_id);
        });

        selectedProducts.value = [];
        searchInput.value = '';
    };

    const getAllProductsStokc = async () => {
        const res = await ProductsService.getAll();
        const data = res.data;

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

    };

    const submitShopping = async () => {
        showInformLoadComponent.value = !showInformLoadComponent.value;

        const totalShopping = associateProdutcs.value.reduce((total, a) => total + (a.purchased_value * a.qtde_purchased), 0);

        console.log(`Total da compra: ${totalShopping}`);

        shoppingPrePayLoad.value = {
            id: 0,
            load: 0,
            shopping_itens: associateProdutcs.value,
            totalShopping: totalShopping
        };
    };

    onMounted(() => {
        SessionStorage.remove('shopping_id');
        getAllProductsStokc();
    });
 </script>

 <style>
    .actions {
        display: flex;
        justify-content: end;
        margin: 15px 0 0 0;

    }
</style>
