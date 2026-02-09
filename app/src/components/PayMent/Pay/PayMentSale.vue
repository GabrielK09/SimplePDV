<template>
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-opacity-40 backdrop-blur-sm">
        <div class="w-[80vh] bg-white">
            <header class="border-gray-100 flex justify-between">
                <h2 class="text-gray-600 text-center ml-4">Formas de pagamento - ID da venda: {{ props.saleId }}</h2>
            </header>
        </div>
    </div>
</template>

<script setup lang="ts">
    const props = defineProps<{
        saleId: number

    }>();

    
</script>