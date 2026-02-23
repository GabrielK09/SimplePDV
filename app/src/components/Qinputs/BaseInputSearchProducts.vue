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

        <q-btn
            type="submit"
            v-show="false"
        />
    </q-form>
</template>

<script setup lang="ts">
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { findById } from 'src/modules/products/services/productsService';
    import { ref } from 'vue';

    const id = ref<number|any>(null);
    const { notify } = useNotify();

    const emits = defineEmits<{
        (e: 'emit:selected-product', value: SaleItemContract): void
    }>();

    const searchProduct = async (): Promise<void> => {
        if(!id.value) return;

        const product = await findById(id.value);

        id.value = null;

        if(!product)
        {
            notify('warning', 'Produto não localizado');
            return;
        };

        emits('emit:selected-product', product);
    };
</script>
