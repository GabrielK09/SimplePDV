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

                        <template v-slot:no-data>
                            <div class="ml-auto mr-auto">
                                <q-icon name="warning" size="30px"/>
                                <span class="mt-auto mb-auto ml-2 text-xs">Sem produtos cadastrados</span>

                            </div>
                        </template>
                    </q-table>
                </div>

                <div class="flex items-center gap-4 py-4">
                    <div class="flex-1">
                        <q-input 
                            v-model="nameUpper" 
                            type="text" 
                            stack-label
                            label-slot
                            :label="filterOption === 'name' ? 'Digite o nome do produto aqui' : 'Digite o ID do produto aqui'"
                            outlined
                            dense
                        />
                    </div>

                    <div class="flex-1">
                        <q-select   
                            v-model="filterOption" 
                            :options="filterOptions" 
                            option-label="name"
                            option-value="field"
                            emit-value
                            map-options
                            outlined
                            clearable
                            label="Filtro de busca" 
                            stack-label
                            label-slot
                            dense
                            :display-value="filterSelected"
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
                        <template v-slot:body-cell-price="props">
                            <q-td :props="props">
                                <q-input
                                    v-model.number="props.row.price"
                                    type="number"
                                    class="w-12 flex ml-auto mr-auto"
                                    input-class="text-center"
                                    dense
                                    @update:model-value="val => validatePrice(Number(val), props.row)"
                                />
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
                        @click="confirmShopping" 
                    />
                </div>
            </section>
        </main>
    </q-page>
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import * as ProductsService from 'src/modules/products/services/productsService';
    import { computed, onMounted, ref, watch } from 'vue';
    import { useNotify } from 'src/helpers/QNotify/useNotify';

    const { notify } = useNotify();

    const filterOptions = [
        { field: 'id', name: 'ID' },
        { field: 'name', name: 'Nome' },
    ];

    const filterOption = ref<'id'|'name'>('name');

    const productsColumns: QTableColumn[] = [
        {
            name: 'id',
            label: 'ID',
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
            label: 'Preço',
            field: 'price',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
        {
            name: 'qtde',
            label: 'Qtde',
            field: 'qtde',
            align: 'center'
        },
    ];

    const associateProductsColumns: QTableColumn[] = [
        {
            name: 'id',
            label: 'ID',
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
            label: 'Preço',
            field: 'price',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
        {
            name: 'qtde',
            label: 'Qtde',
            field: 'qtde',
            align: 'center'
        },
    ];

    const productsStockData = ref<ProductContract[]>([]);
    const allProductsStockData = ref<ProductContract[]>([]);

    const associateProdutcs = ref<ProductContract[]>([]);
    const allAssociateProdutcs = ref<ProductContract[]>([]);

    const searchInput = ref<string>('');

    const pagination = ref({
        sortBy: 'id' 
    });

    const searchProduct = ref<string>('');

    const showCreateProductComponent = ref<boolean>(false);

    const filterSelected = computed(() => {
        if(filterOption.value === undefined) return 'Nome';

        return filterOptions.find(o => o.field === filterOption.value)?.name;
    });

    const nameUpper = computed({
        get: () => searchProduct.value,
        set: (val: string) => {
            searchProduct.value = val.toUpperCase();
        }
    });

    watch(
        () => searchProduct.value,
        (newVal) => {            
            let productStock: ProductContract;

            switch (filterOption.value) {
                case 'id':
                    productStock = productsStockData.value.find(i => i.id === Number(newVal));
                    break;
            
                case 'name':
                    productStock = productsStockData.value.find(i => i.name === newVal.toUpperCase());
                    break;
            
                default:
                    // Por padrão busca por nome
                    productStock = productsStockData.value.find(i => i.name === newVal.toUpperCase());
                    break;
            };
            
            if(!productStock) return;
        
            if(associateProdutcs.value.find(p => p.id === productStock.id || p.name === productStock.name))
            {
                nameUpper.value = '';
                
                notify(
                    'info',
                    'Produto já associado!'
                );

                return;
            };
    
            const newProductStock: ProductContract = {
                id: productStock.id,
                name: productStock.name,
                price: productStock.price,
                qtde: 1,
                commission: 0,
            };

            associateProdutcs.value.push(newProductStock);

            nameUpper.value = '';
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
        productsStockData.value = allProductsStockData.value.filter(product => product.name.toLowerCase().includes(searchInput.value));

    };

    const cancelShopping = (): void => {

    };

    const confirmShopping = (): void => {

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

    onMounted(() => {
        getAllProductsStokc();
    });
 </script>

 <style>
    body {
        background-color: #fff !important;
    }

    .actions {
        display: flex;
        justify-content: end;
        margin: 15px 0 0 0;

    }
</style>