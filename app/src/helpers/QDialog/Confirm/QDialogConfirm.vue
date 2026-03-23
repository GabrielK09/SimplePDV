<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card>
            <q-card-section class="row items-center">
                <span>
                    {{ props.text }}
                </span>

            </q-card-section>

            <q-card-actions align="right">
                <div v-if="!onLeave">
                    <q-btn
                        flat
                        dense
                        label="Cancelar"
                        color="primary"
                        no-caps
                        @click="emits('close', true)"
                    />

                    <q-btn
                        label="Confirmar"
                        dense
                        no-caps
                        :color="props.operation === 'delete' ? 'red' : 'primary'"
                        type="submit"
                        :loading="startLoading"
                        @click="emitEvent()"
                    />
                </div>

                <div v-else class="flex gap-4">
                    <q-btn
                        label="Cancelar e sair"
                        dense
                        no-caps
                        color="red"
                        type="submit"
                        :loading="startLoading"
                        @click="emits('leave:cancel-and-leave', true)"
                    />

                    <q-btn
                        label="Cancelar e ficar"
                        dense
                        no-caps
                        color="red"
                        flat
                        type="submit"
                        :loading="startLoading"
                        @click="emits('confirm', true)"
                    />

                    <q-btn
                        label="Salvar e sair"
                        dense
                        no-caps
                        color="primary"
                        type="submit"
                        :loading="startLoading"
                        @click="emits('leave:save-and-leave', true)"
                    />
                </div>
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import { ref, watch } from 'vue';

    const props = defineProps<{
        show: boolean;
        text: string;
        operation?: 'save'|'delete'|'';
        onLeave?: boolean;
    }>();

    const emits = defineEmits<{
        (e: 'close', value: boolean),
        (e: 'confirm', value: boolean),

        (e: 'leave:cancel-and-leave', value: boolean)
        (e: 'leave:save-and-leave', value: boolean)
    }>();

    const internalDialog = ref<boolean>(props.show);
    const startLoading = ref<boolean>(false);

    watch(
        () => props.show,
        (val) => {
            internalDialog.value = val;
        }
    );

    const emitEvent = () => {
        startLoading.value = true;
        emits('confirm', true);
    };
</script>
