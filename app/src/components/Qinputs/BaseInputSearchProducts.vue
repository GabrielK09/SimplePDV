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

    const normalizeProduct = (p: ProductContract): SaleItemContract => ({
        id: p.id,
        product_id: p.id,
        name: p.name,
        price: p.price,
        //@ts-ignore
        product_with_characteristics: p.productWithCharacteristics,
        qtde: 1
    });

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
        const parsedProduct: ProductContract = {
            id: intermediaryProductItemData.value.id,
            name: intermediaryProductItemData.value.name,
            price: intermediaryProductItemData.value.price,
            qtde: intermediaryProductItemData.value.qtde,
            commission: 0,
            productWithCharacteristics: grid
        };

        emitProduct(normalizeProduct(parsedProduct));
    };

    const emitProduct = (product: SaleItemContract) => {
        emits('emit:selected-product', product);
        searchInput.value = '';
    };

    watch(
        () => searchInput.value,
        async (idValue) => {
            const input = idValue?.toString().split('') ?? '';

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
        if(!habilitStringSearchInput.value) return;

        const search = searchInput.value.toString().slice(1);

        if(!search) return;

        const res = await findByName(search);

        console.log(res);

        if(!res.success) return;

        const data: any[] = res.data;

        console.log(data);

        itensData.value = data.map(p => ({
            id: p.product.id,
            name: p.product.name,
            price: p.product.price,
            product_id: p.product.product_id,
            qtde: 1,
            use_grid: p.use_grid,
            //@ts-ignore
            product_with_characteristics: p.characteristics,
        }));
    };

    const selectProduct = (_: Event, row: SaleItemContract) => {
        if(row.product_with_characteristics !== null)
        {
            showSizeGrid.value = true;
            productFullData.value = {
                commission: 0,
                id: row.id,
                name: row.name,
                price: row.price,
                qtde: 1,
                use_grid: true,
                 //@ts-ignore
                productWithCharacteristics: row.product_with_characteristics
            };

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
