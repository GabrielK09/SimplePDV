<template>
    <q-checkbox
        right-label
        v-model="registeredCustomer"
        label="Cliente cadastrado"
    />

    <div v-if="registeredCustomer">
        <q-select
            outlined
            v-model="model"
            :options="options"
        />
    </div>

    <div v-else>
        <q-input
            outlined
            v-model="model"
        />
    </div>
</template>

<script setup lang="ts">
    import { ref } from 'vue';

    const model = defineModel<any>();

    const options = ref<any[]>([]);
    const registeredCustomer = ref<boolean>(false);

    const emits = defineEmits<{
        (e: 'selected:customer', value: any),

    }>();

</script>
