<template>
    <q-form
        class="flex-1"
        @submit="searchProduct"
    >
        <q-input
            outlined
            v-model="id"
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

    const id = ref<number | string | null>(null);

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
        },
    ];

    const itensData = ref<SaleItemContract[]>([]);

    const { notify } = useNotify();

    const emits = defineEmits<{
        (e: 'emit:selected-product', value: SaleItemContract): void
    }>();

    const searchProduct = async (): Promise<void> => {
        let productQtde: number = 0;

        if(!id.value) return;
       
        const product = await findById(Number(id.value));

        if(!product)
        {
            notify('warning', 'Produto não localizado');
            return;
        };

        const input = id.value?.toString().split('') ?? '';

        if(input[1] === '*') productQtde = Number(input[0]) ?? 1;

        const productData: SaleItemContract = {
            id: product.data?.id,
            name: product.data?.name,
            price: product.data?.price,
            product_id: product.data?.product_id,
            qtde: productQtde
        };

        emits('emit:selected-product', productData);
        id.value = '';
    };

    watch(
        () => id.value,
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
        if(!habilitStringSearchInput.value) return;

        const search = id.value.toString().slice(1);

        if(!search) return;
        
        const res = await findByName(search);

        if(!res.success) return;
        
        itensData.value = res.data;
    };

    const selectProduct = (evt: Event, row: SaleItemContract) => {
        emits('emit:selected-product', row);

        id.value = null;
        habilitStringSearchInput.value = false;
        itensData.value = [];
    };
</script>
