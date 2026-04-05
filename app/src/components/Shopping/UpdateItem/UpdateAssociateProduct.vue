<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card v-if="newProduct">
            <main class="rounded-md flex flex-center mt-4 bg-white text-xs">
                <section class="rounded-lg p-6 flex flex-col">
                    <span class="text-center font-bold mb-2">
                        Produto em edição: {{ newProduct.name }}
                    </span>

                    <div v-if="props.productData?.product_with_characteristics !== null" class="mt-4">
                        <BaseSelectGridTypes
                            v-model="newProduct.product_with_characteristics[0].size"
                            :selected-sizes="props.productData?.product_with_characteristics"
                        />
                    </div>

                    <div class="mt-6 flex justify-center">
                        <span class="font-bold mr-4 my-auto ">Grade atual: {{ newProduct.product_with_characteristics[0].size }}</span>

                    </div>

                    <div class="mt-4 flex justify-end gap-2">
                        <q-btn 
                            flat 
                            label="Cancelar" 
                            color="negative" 
                            @click="emit('close', true)" 
                        />

                        <q-btn 
                            label="Salvar" 
                            color="primary" 
                            @click="saveData" 
                        />
                    </div>
                </section>
            </main>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import BaseSelectGridTypes from 'src/components/Products/UseGrid/QSelectGridTypes/BaseSelectGridTypes.vue';
    import { ref, watch } from 'vue';    

    const props = defineProps<{
        productData: ShoppingItemContract | undefined
    }>();

    const emit = defineEmits<{
        (e: 'close', value: boolean): void
        (e: 'update:product', value: ShoppingItemContract): void
    }>();

    const internalDialog = ref<boolean>(true);

    const newProduct = ref<ShoppingItemContract | null>(null);

    watch(
        () => props.productData,
        (product) => {
            if(!product) return;

            newProduct.value = product ? { ...product } : null;

        },
        { immediate: true }
    );

    const saveData = () => {
        if (!newProduct.value) return;

        console.log(newProduct.value);
        
        emit('update:product', newProduct.value);
        emit('close', true);
    };
</script>