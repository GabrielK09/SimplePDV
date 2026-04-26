<template>
    <q-form
        class="flex-1"
        @submit="searchProduct"
    >
        <q-input
            outlined
            v-model="searchInput"
            type="text"
            stack-label
            label-slot
            label="Busque um produto"
        />

        <div
            v-if="habilitStringSearchInput"
            class="relative-bottom"
        >
            <div class="absolute text-white bg-blue-600 font-bold text-dark q-pa-sm rounded-borders shadow-3 text-caption z-50 w-full">
                <q-table
                    dense
                    :rows="itensData"
                    :columns="itensTableColumn"
                    hide-pagination
                    row-key="id"
                    @row-click="selectProduct"
                    class="cursor-pointer"
                />
            </div>
        </div>

        <q-btn
            type="submit"
            v-show="false"
        />
    </q-form>

    <QSelectGridTable
        v-if="showSizeGrid"
        :is-just-list="true"
        :characteristics="productCharacteristics"
        @return:selected-grid="handelSelectedGrid($event)"
        @close="showSizeGrid = !$event"
    />
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { findById, findByName } from 'src/modules/products/services/productsService';
    import { ref, watch } from 'vue';
    import QSelectGridTable from '../Products/UseGrid/QTable/QSelectGridTable.vue';

    const productCharacteristics = ref<ProductCharacteristicsContract[]>([]);

    const intermediaryProductItemData = ref<SaleItemContract>();

    const searchInput = ref<number | string | null>(null);

    const habilitStringSearchInput = ref<boolean>(false);
    const showSizeGrid = ref<boolean>(false);
    const productQtdeByOperator = ref<number|null>(1);

    const itensTableColumn: QTableColumn[] = [
        {
            name: 'id',
            label: 'Cód produto',
            field: 'id',
            align: 'left',
        },
        {
            name: 'name',
            label: 'Produto',
            field: 'name',
            align: 'center',
        },
        {
            name: 'qtde',
            label: 'Qtde',
            field: 'qtde',
            align: 'center',
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

    const itensData = ref<SaleItemContract[]>([]);

    const emits = defineEmits<{
        (e: 'emit:selected-product', value: SaleItemContract): void
    }>();

    const normalizeProduct = (p: ProductContract): SaleItemContract => ({
        id: p.id,
        product_id: p.id,
        name: p.name,
        sale_value: p.price,
        product_with_characteristics: p.product_with_characteristics,
        qtde: !p.use_grid ? productQtdeByOperator.value : 1
    });

    const searchProduct = async (): Promise<void> => {
        if (habilitStringSearchInput.value) return; // Se for busca pelo nome, não valida busca pelo qtde/código

        if(!searchInput.value) return; // Se não for busca pelo nome, precisa possuir algum valor inserido no campo de busca

        let productId: number = Number(searchInput.value);

        if(searchInput.value?.toString().split('').includes('*'))
        {            
            const splitedInput = searchInput.value?.toString().split('*');

            productQtdeByOperator.value = Number(splitedInput[0]);
            productId = Number(splitedInput[1]);
            
        };
                
        const resProduct = (await findById(productId)).data as ProductContract;
    
        const productResData: ProductContract = resProduct;
            
        if(!productResData.use_grid)
        {
            emitProduct(normalizeProduct(productResData));
            return;
        };

        productCharacteristics.value = resProduct.product_with_characteristics;
        productResData.product_with_characteristics = resProduct.product_with_characteristics;
        
        const productData: SaleItemContract = {
            id: productResData?.id,
            name: productResData?.name,
            sale_value: productResData?.price,
            product_id: productResData?.id,
            qtde: 1,
            product_with_characteristics: null
        };

        intermediaryProductItemData.value = productData;

        showSizeGrid.value = true;
        return;
    };

    const handelSelectedGrid = (grid: ProductCharacteristicsContract) => {
        if(!intermediaryProductItemData.value) return;

        const parsedProduct: ProductContract = {
            id: intermediaryProductItemData.value.id,
            name: intermediaryProductItemData.value.name,
            price: intermediaryProductItemData.value.sale_value,
            qtde: intermediaryProductItemData.value.qtde,
            use_grid: true,
            commission: 0,
            product_with_characteristics: []
        };

        parsedProduct.product_with_characteristics.push({
            grid_qtde: 1,
            id: grid.id,
            size: grid.size,
            product_id: grid.product_id
        });
        
        emitProduct(normalizeProduct(parsedProduct));
        showSizeGrid.value = false;
    };

    const emitProduct = (product: SaleItemContract) => {
        emits('emit:selected-product', product);
        searchInput.value = '';
    };

    watch(
        () => searchInput.value, 
        async (idValue) => {
            if(idValue?.toString().split('')[0] === '/')
            {
                habilitStringSearchInput.value = true;
                await getProductByName();

                return;

            } else {
                habilitStringSearchInput.value = false;
            };
        },

        { immediate: true }
    );

    const getProductByName = async () => {
        if(!habilitStringSearchInput.value) return;
        if(!searchInput.value) return;

        const search = searchInput.value.toString().slice(1);

        if(!search) return;

        const res = await findByName(search);

        if(!res.success) return;

        const data: ProductContract[] = res.data;

        itensData.value = data.map(p => ({
            id: p.id,
            name: p.name,
            sale_value: p.price,
            product_id: p.id,
            qtde: 1,
            use_grid: p.use_grid,
            product_with_characteristics: p.product_with_characteristics,
        }));
    };

    const selectProduct = (_: Event, row: SaleItemContract) => {
        if(row.product_with_characteristics !== null && row.use_grid)
        {
            showSizeGrid.value = true;

            productCharacteristics.value = row.product_with_characteristics;

            intermediaryProductItemData.value = row;
            searchInput.value = null;
            habilitStringSearchInput.value = false;
            itensData.value = [];
            return;
        };

        emits('emit:selected-product', row);

        searchInput.value = null;
        habilitStringSearchInput.value = false;
        itensData.value = [];
    };
</script>