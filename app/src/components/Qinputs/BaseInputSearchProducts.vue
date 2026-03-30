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
        :product-data="productFullData"
        @return:selected-grid="handelSelectedGrid($event)"
    />
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { findById, findByName } from 'src/modules/products/services/productsService';
    import { ref, watch } from 'vue';
    import QSelectGridTable from '../Products/UseGrid/QTable/QSelectGridTable.vue';

    const productFullData = ref<ProductContract>({
        id: null,
        commission: null,
        name: null,
        price: null,
        qtde: null,
        productWithCharacteristics: null,
        use_grid: null,
    });

    const intermediaryProductItemData = ref<SaleItemContract>();

    const searchInput = ref<number | string | null>(null);

    const habilitStringSearchInput = ref<boolean>(false);
    const showSizeGrid = ref<boolean>(false);

    const itensTableColumn: QTableColumn[] = [
        {
            name: 'id',
            label: 'Cód. Produto',
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

    const { notify } = useNotify();

    const emits = defineEmits<{
        (e: 'emit:selected-product', value: SaleItemContract): void
    }>();

    const searchProduct = async (): Promise<void> => {
        if (habilitStringSearchInput.value) return; // Se for busca pelo nome, não valida busca pelo qtde/código

        if(!searchInput.value) return; // Se não for busca pelo nome, precisa possuir algum valor inserido no campo de busca

        let productQtde: number = 1;
        let productId: number = Number(searchInput.value);
    
        if(searchInput.value?.toString().split('').includes('*'))
        {
            const splitedInput = searchInput.value?.toString().split('*');
            productQtde = Number(splitedInput[0]);
            productId = Number(splitedInput[1]);
        };
        
        const resProduct = (await findById(productId)).data;

        const productResData: ProductContract = resProduct.product;
        productResData.productWithCharacteristics = resProduct.characteristics;

        console.log('product localizado: ', resProduct);
    
        if(!productResData)
        {
            notify('warning', 'Produto não localizado');
            return;
        };

        const productData: SaleItemContract = {
            id: productResData?.id,
            name: productResData?.name,
            price: productResData?.price,
            product_id: productResData?.id,
            qtde: productQtde ?? 1,
            product_with_characteristics: null
        };

        intermediaryProductItemData.value = productData;

        if(!productResData.use_grid && productResData.productWithCharacteristics.length <= 0)
        {
            emitProduct(productData);
            return;
        };
    
        productFullData.value = productResData;
        showSizeGrid.value = true;
        return;
    };

    const handelSelectedGrid = (grid: any) => {
        const normalizedProduct: SaleItemContract = {
            id: intermediaryProductItemData.value.id,
            name: intermediaryProductItemData.value.name,
            price: intermediaryProductItemData.value.price,
            product_id: intermediaryProductItemData.value.product_id,
            qtde: intermediaryProductItemData.value.qtde,
            product_with_characteristics: grid
        };

        emitProduct(normalizedProduct);
    };

    const emitProduct = (product: SaleItemContract) => {
        emits('emit:selected-product', product);
        searchInput.value = '';
    };

    watch(
        () => searchInput.value,
        async (idValue) => {
            const input = idValue?.toString().split('') ?? '';

            console.log(input);

            if(input[0] === '/')
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
        console.log('call getProductByName');

        if(!habilitStringSearchInput.value) return;

        const search = searchInput.value.toString().slice(1);

        console.log('Vai buscar por: ', search);

        if(!search) return;

        const res = await findByName(search);

        if(!res.success) return;

        itensData.value = res.data;
    };

    const selectProduct = (evt: Event, row: SaleItemContract) => {
        console.log('evt: ', evt);

        const normalizedProduct: SaleItemContract= {
            id: row.id,
            name: row.name,
            price: row.price,
            product_id: row.product_id,
            product_with_characteristics: row.product_with_characteristics,
            qtde: 1
        };

        emits('emit:selected-product', normalizedProduct);

        searchInput.value = null;
        habilitStringSearchInput.value = false;
        itensData.value = [];
    };
</script>
