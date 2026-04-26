<template>
    <q-dialog v-model="confirm" persistent>
        <q-card class="product-dialog">
            <q-card-section class="dialog-header">
                <div class="mx-auto">
                    <div class="bg-white rounded-sm">
                        <header class="border-gray-100 flex justify-between">
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke-width="1.5"
                                stroke="currentColor"
                                class="size-6 ml-4 mt-auto mb-auto cursor-pointer"
                                @click="emits('close', true)"
                            >
                                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
                            </svg>

                            <h2 class="text-gray-600 text-center ml-4">Seleção de produtos</h2>

                            <q-input
                                outlined
                                v-model="searchInput"
                                type="text"
                                stack-label
                                dense
                                class="mt-auto mb-auto mr-2"
                                label-slot
                                @update:model-value="filterProducts"
                            >
                                <template v-slot:append>
                                    <q-icon name="search" />
                                </template>

                                <template v-slot:label>
                                    <span class="text-xs">Buscar por um produto ...</span>
                                </template>
                            </q-input>
                        </header>

                        <q-table
                            :rows="products"
                            :columns="columns"
                            v-model:pagination="pagination"
                            row-key="id"
                        >
                            <template v-slot:body-cell-select="props">
                                <q-td :props="props">
                                    <div v-if="propsComponent.typeSearch === 'single'">
                                        <q-radio
                                            v-model="selectedProductsId"
                                            :val="props.row.id"
                                            dense
                                        />
                                    </div>

                                    <div v-else>
                                        <q-checkbox
                                            v-model="selectedProductsIds"
                                            :val="props.row.id"
                                            dense
                                        />
                                    </div>
                                </q-td>
                            </template>
                        </q-table>
                    </div>
                </div>
            </q-card-section>

            <div id="action_bar">
                <div class="bg-[#03202e] w-[60%] mb-4 p-3 text-white flex justify-end rounded-lg">
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
                        color="primary"
                        no-caps
                        label="Confirmar"
                        @click="emitProducts"
                    />
                </div>
            </div>

        </q-card>
    </q-dialog>

    <QSelectGridTable
        v-if="showSizeGrid"
        :is-just-list="true"
        :characteristics="productCharacteristics"
        @return:selected-grid="handelSelectedGrid($event)"
        @close="closeQSelectGridTable(!$event)"
    />
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { getAll } from 'src/modules/products/services/productsService';
    import { computed, onMounted, ref, watch } from 'vue';
    import QSelectGridTable from '../Products/UseGrid/QTable/QSelectGridTable.vue';

    type TypeSearch = 'multiple' | 'single';

    const propsComponent = defineProps<{
        typeSearch: TypeSearch;
    }>();

    const pagination = ref({
        sortBy: 'id',
        rowsPerPage: 20
    });

    const emits = defineEmits<{
        (e: 'close', value: boolean): void;
        (e: 'emit:selected-products', value: SaleItemContract[]): void;
    }>();

    const products = ref<ProductContract[]>([]);
    const allProducts = ref<ProductContract[]>([]);

    const selectedProductsIds = ref<number[]>([]);
    const selectedProductsId = ref<number | null>(null);

    const idOfProductHaveCharacteristics = ref<number | null>();
    
    const confirm = ref<boolean>(true);
    const showSizeGrid = ref<boolean>(false);
    
    const productCharacteristics = ref<ProductCharacteristicsContract[]>([]);

    const columns: QTableColumn[] = [
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
        }
    ];

    const searchInput = ref<string>('');

    const closeQSelectGridTable = (event: boolean) => {
        propsComponent.typeSearch === 'single' 
            ? selectedProductsId.value = null
            : selectedProductsIds.value = [];

        showSizeGrid.value = event;
    };

    const normalizeProduct = (p: ProductContract): SaleItemContract => ({
        id: p.id,
        product_id: p.id,
        name: p.name,
        price: p.price,
        product_with_characteristics: p.product_with_characteristics,
        qtde: 1,
        use_grid: p.use_grid
    });

    const getAllProducts = async (perPager: number) => {
        const res = (await getAll(perPager));

        const data = res.data;
        console.log(data);

        const formatedProducts: ProductContract[] = data.map((r: ProductContract) => ({
            id: r.id,
            name: r.name,
            commission: r.commission,
            price: r.price,
            qtde: 1,
            deleted_at: r.deleted_at,
            use_grid: r.use_grid,
            product_with_characteristics: r.product_with_characteristics

        })).filter((r: ProductContract) => r.deleted_at === null);

        products.value = formatedProducts;

        allProducts.value = [...products.value];
    };

    const filterProducts = () => {
        products.value = allProducts.value.filter(product => product.name ? product.name.toLowerCase().includes(searchInput.value) : []);
    };

    const selectedProducts = computed<SaleItemContract[]>(() => {        
        if (propsComponent.typeSearch === 'single') 
        {            
            const product = products.value.find(p => p.id === selectedProductsId.value);

            return product ? [normalizeProduct(product)] : [];
        };
        
        return products.value
            .filter(p => selectedProductsIds.value.includes(Number(p.id)))
            .map(normalizeProduct);
    });

    watch(
        () => selectedProductsIds.value,
        (ids) => {
            const actionBar = document.getElementById('action_bar') as HTMLElement;

            if(actionBar)
            {                
                actionBar.style.display = ids.length > 0 ? 'flex' : 'none';
            };
            
            ids.forEach(i => {
                const productData = products.value.find(p => p.id === i);

                if(productData?.use_grid && productData?.product_with_characteristics !== null)
                {
                    idOfProductHaveCharacteristics.value = productData.id;

                    productCharacteristics.value = productData.product_with_characteristics;

                    showSizeGrid.value = true;   
                };
            });
        }
    );

    const handelSelectedGrid = (grid: ProductCharacteristicsContract) => {    
        if(idOfProductHaveCharacteristics.value < 0) return;

        const oldProductData = selectedProducts.value.find(p => p.id === idOfProductHaveCharacteristics.value);

        const normalizedProduct: SaleItemContract = {
            id: oldProductData.id,
            name: oldProductData.name,
            price: oldProductData.price,
            product_id: oldProductData.product_id,
            qtde: oldProductData.qtde,
            product_with_characteristics: [{
                grid_qtde: 1,
                id: grid.id,
                product_id: grid.product_id,
                size: grid.size
            }]
        };

        let newProduct = selectedProducts.value.find(p => p.id === idOfProductHaveCharacteristics.value);

        newProduct = {
            id: normalizedProduct.id,
            name: normalizedProduct.name,
            price: normalizedProduct.price,
            product_id: normalizedProduct.product_id,
            product_with_characteristics: normalizedProduct.product_with_characteristics,
            qtde: normalizedProduct.qtde,
            use_grid: normalizedProduct.use_grid
        };
        
        showSizeGrid.value = false;
    };

    const emitProducts = () => {
        emits('emit:selected-products', selectedProducts.value);
        emits('close', true);
    };

    onMounted(() => {
        getAllProducts(20);
    });
</script>

<style scoped>
    .product-dialog {
        width: 100%;
        max-width: 750px;
        min-width: 320px;
        border-radius: 18px;
    }

    .dialog-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 16px;
        padding: 20px 24px;
        background: linear-gradient(to right, #f8fafc, #ffffff);
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