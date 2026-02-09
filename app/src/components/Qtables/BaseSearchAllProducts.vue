<template>
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-opacity-40 backdrop-blur-sm">
        <div class="w-[80vh] bg-white">
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
            >
                <template v-slot:body-cell-select="props">
                    <q-td :props="props">
                        <div v-if="propsComponent.typeSearch === 'single'">
                            <q-radio
                                v-model="selectedProducts"
                                :val="props.row"
                                dense
                            />
                        </div>

                        <div v-else>
                            <q-checkbox
                                v-model="selectedProducts"
                                :val="props.row"
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
                    label="Ok"
                    @click="emitProducts"
                />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
    import camelcaseKeys from 'camelcase-keys';
    import { QTableColumn } from 'quasar';
    import { getAll } from 'src/modules/products/services/productsService';
    import { onMounted, ref } from 'vue';

    type TypeSearch = 'multiple' | 'single';

    const propsComponent = defineProps<{
        typeSearch: TypeSearch;
    }>();

    const emits = defineEmits<{
        (e: 'close', value: boolean): void,
        (e: 'emit:selected-products', value: SaleItemContract[]): void
    }>();

    const products = ref<ProductContract[]>([]);
    const allProducts = ref<ProductContract[]>([]);

    const selectedProducts = ref<SaleItemContract[]>([]);

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

    let searchInput = ref<string>('');

    const getAllProducts = async () => {
        const res = await getAll();
        const data = camelcaseKeys(res.data, { deep: true });

        products.value = data;
        allProducts.value = [...products.value];
    };

    const filterProducts = () => {
        console.log(searchInput.value);

        products.value = allProducts.value.filter(product => product.name.toLowerCase().includes(searchInput.value));

    };

    const emitProducts = () => {
        emits('emit:selected-products', selectedProducts.value);
        emits('close', true);
    };

    onMounted(() => {
        getAllProducts();
    });
</script>
