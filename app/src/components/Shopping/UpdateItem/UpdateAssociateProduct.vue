<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card v-if="newProduct">
            <main class="rounded-md flex flex-center mt-4 bg-white text-xs">
                <section class="rounded-lg p-6 flex flex-col">
                    <span class="text-center font-bold mb-2">
                        Produto em edição: {{ newProduct.name }}
                    </span>

                    <div>
                        <q-input v-model.number="newProduct.qtde_purchased" type="number" stack-label label-slot
                            label="Qtde comprada" />
                    </div>

                    <div>
                        <q-input v-model.number="newProduct.purchased_value" type="number" stack-label label-slot
                            label="Qtde comprada" />
                    </div>

                    <div v-if="props.productData?.productWithCharacteristics !== null">
                        <q-btn 
                            color="primary" 
                            label="Alterar grade"
                            no-caps
                            @click="changeGrid = !changeGrid" 
                        />
        
                        <QSelectGridTable   
                            v-if="changeGrid"
                            :characteristics="props.productData?.productWithCharacteristics"
                        />
                        
                    </div>

                    <div class="mt-4 flex justify-end gap-2">
                        <q-btn flat label="Cancelar" color="negative" @click="emit('close', true)" />
                        <q-btn label="Salvar" color="primary" @click="saveData" />
                    </div>
                </section>
            </main>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import QSelectGridTable from 'src/components/Products/UseGrid/QTable/QSelectGridTable.vue';
    import { ref, watch } from 'vue';    

    const props = defineProps<{
        productData: ShoppingItemContract | undefined
    }>();

    const emit = defineEmits<{
        (e: 'close', value: boolean): void
        (e: 'update:product', value: ShoppingItemContract): void
    }>();

    const internalDialog = ref<boolean>(true);
    const changeGrid = ref<boolean>(true);
    const newProduct = ref<ShoppingItemContract | null>(null);

    watch(
        () => props.productData,
        (product) => {
            newProduct.value = product ? { ...product } : null;
        },
        { immediate: true }
    );

    const saveData = () => {
        if (!newProduct.value) return;

        emit('update:product', newProduct.value);
        emit('close', false);
    };
</script>