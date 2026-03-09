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
        if(!id.value) return;

        const product = await findById(Number(id.value));

        id.value = null;

        if(!product)
        {
            notify('warning', 'Produto não localizado');
            return;
        };

        emits('emit:selected-product', product.data);
    };

    watch(
        () => id.value,
        async (idValue) => {
            let splitedInput: any;

            if (idValue) {
                splitedInput = idValue.toString().split('');

            } else {
                splitedInput = '';
            };

            if(typeof idValue === 'string' && splitedInput[0]=== '/')
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

        if(search) return;

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
