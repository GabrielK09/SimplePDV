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
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { findById, findByName } from 'src/modules/products/services/productsService';
    import { ref, watch } from 'vue';

    const searchInput = ref<number | string | null>(null);

    const habilitStringSearchInput = ref<boolean>(false);

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

        let productQtde = [];
    
        const splitedSearchInput = searchInput.value?.toString().split('');

        if(splitedSearchInput.includes('*'))
        {
            //const removeSpaces = splitedSearchInput.join('');
            splitedSearchInput.map(s => {
                console.log(s);
        
                if(s !== "*")
                {
                    productQtde.push(s);
                    return;
                };
            });
            
            console.log('productQtde: ', productQtde);
            return;
        };

        console.log('Number(searchInput.value): ', searchInput.value);
        
        const product = await findById(Number(searchInput.value));

        console.log('product localizado: ', product);
    
        if(!product)
        {
            notify('warning', 'Produto não localizado');
            return;
        };

        const productData: SaleItemContract = {
            id: product.data.product?.id,
            name: product.data.product?.name,
            price: product.data.product?.price,
            product_id: product.data.product?.id,
            qtde: productQtde ?? 1
        };

        console.log('Vai enviar: ', productData);

        emits('emit:selected-product', productData);
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

        emits('emit:selected-product', row);

        searchInput.value = null;
        habilitStringSearchInput.value = false;
        itensData.value = [];
    };
</script>
