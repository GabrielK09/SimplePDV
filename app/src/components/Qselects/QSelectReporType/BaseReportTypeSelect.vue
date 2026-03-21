<template>
    <q-select
        outlined
        v-model="reportType"
        :options="reportTypes"
        option-label="name"
        emit-value
        label="Tipo do relatório"
        stack-label
        label-slot
        map-options
    />
</template>

<script setup lang="ts">
    import { ref } from 'vue';

    const reportTypesSelect = [
        {name: "Caixa", value: "cash-register"},
        {name: "Formas de pagamento", value: "pay-ment-forms"},
        {name: "Itens vendidos", value: "saled-itens"},
        {name: "Itens comprados", value: "shopping-itens"},
        {name: "Compras", value: "shoppings"}
    ];

    const reportTypes = ref<any>(reportTypesSelect);
    const reportType = defineModel<string | null>()

</script>
