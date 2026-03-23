<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card>
            <main class="rounded-md flex flex-center mt-4 bg-white text-xl">
                <section class="rounded-lg p-6 flex flex-col">
                    <span class="font-bold mb-2">Última carga: {{ props.lastShoppingId }}</span>

                    <q-input
                        v-model.number="informedLoad"
                        type="text"
                        stack-label
                        label-slot
                        input-class="text-lg"
                        :rules="[
                            val => val > 0 || 'A carga da compra não pode ser menor zero.'
                        ]"
                    >
                        <template v-slot:label>
                            <span class="font-bold text-lg">
                                N° Carga da compra
                            </span>
                        </template>
                    </q-input>

                    <div class="mt-4 flex gap-4">
                        <q-btn
                            color="red"
                            label="Fechar"
                            no-caps
                            @click="emits('close', true)"
                        />

                        <q-btn
                            color="primary"
                            label="Confirmar carga da compra"
                            no-caps
                            @click="emits('return:informed-load', informedLoad)"
                        />

                    </div>
                </section>
            </main>
        </q-card>
    </q-dialog>

</template>

<script setup lang="ts">
    import { ref } from 'vue';

    const emits = defineEmits<{
        (e: 'close', value: boolean),
        (e: 'return:informed-load', value: number),
    }>();

    const props = defineProps<{
        lastShoppingId: number;
    }>();

    const informedLoad = ref<number>(0);
    const internalDialog = ref<boolean>(true);
</script>
