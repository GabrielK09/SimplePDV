<template>
    <q-select
        outlined
        v-model="size"
        :disable="props.isUpdate"
        :options="filterSizeOptions"
        option-label="name"
        emit-value
        label="Tamanho"
        stack-label
        label-slot
        map-options
    />
</template>

<script setup lang="ts">
    import { computed } from 'vue';

    type SizeOption = {
        name: string;
        value: string
    };

    const props = defineProps<{
        selectedSizes: any[],
        isUpdate?: boolean
    }>();

    const size = defineModel<string | null>();

    const sizeOptions: SizeOption[] = [
        {name: "PP", value: "PP"},
        {name: "P", value: "P"},
        {name: "M", value: "M"},
        {name: "G", value: "G"},
        {name: "XG", value: "XG"},
        {name: "XGG", value: "XGG"},
        {name: "EG", value: "EG"},
        {name: "EGG", value: "EGG"},
        {name: "O", value: "O"},
    ];

    const filterSizeOptions = computed(() => {
        const formated = sizeOptions.filter(o => !props.selectedSizes.includes(o.name));

        return formated;
    });    
</script>