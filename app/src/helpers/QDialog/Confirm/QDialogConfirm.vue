<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card>
            <q-card-section class="row items-center">
                <span class="q-ml-sm">
                    {{ props.text }}
                </span>

            </q-card-section>

            <q-card-actions align="right">
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
    }>();

    const emits = defineEmits<{
        (e: 'close', value: boolean),
        (e: 'confirm', value: boolean)
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
