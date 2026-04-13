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

                        <div class="flex justify-end p-5">
                            <q-btn
                                color="primary"
                                no-caps
                                label="Confirmar"
                                @click="emitProducts"
                            />
                        </div>
                    </div>
                </div>
            </q-card-section>
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
    const indexOfProductHaveCharacteristics = ref<number | null>();
    
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
            ?  selectedProductsId.value = null
            : selectedProductsIds.value = [];

        showSizeGrid.value = event;
    };

    const normalizeProduct = (p: ProductContract): SaleItemContract => ({
        id: p.id,
        product_id: p.id,
        name: p.name,
        price: p.price,
        product_with_characteristics: p.product_with_characteristics,
        qtde: 1
    });

    const getAllProducts = async (perPager: number) => {
        const res: any[] = (await getAll(perPager)).data;

        const formatedProducts: ProductContract[] = res.map(r => ({
            id: r.product.id,
            name: r.product.name,
            commission: r.product.commission,
            price: r.product.price,
            qtde: 1,
            deleted_at: r.product.deleted_at,
            use_grid: r.product.use_grid,
            product_with_characteristics: r.characteristics,

        })).filter((r: ProductContract) => r.deleted_at === null);

        products.value = formatedProducts;

        allProducts.value = [...products.value];
    };

    const filterProducts = () => {
        products.value = allProducts.value.filter(product => product.name ? product.name.toLowerCase().includes(searchInput.value) : null);
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

    watch(selectedProducts, (newProducts) => {
        newProducts.forEach(p => {
            if(p.product_with_characteristics !== null) 
            {   
                const index = newProducts.findIndex(item => item.id === p.id);

                if (index > -1)
                {
                    indexOfProductHaveCharacteristics.value = index;

                    productCharacteristics.value = p.product_with_characteristics;

                    showSizeGrid.value = true;
                };
            };
        });
    }, {    
        immediate: true 
    });

    const handelSelectedGrid = (grid: ProductCharacteristicsContract) => {    
        if(indexOfProductHaveCharacteristics.value < 0) return;

        const oldProductData = selectedProducts.value[indexOfProductHaveCharacteristics.value];

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

        selectedProducts.value[indexOfProductHaveCharacteristics.value] = normalizedProduct;
        
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
</style>