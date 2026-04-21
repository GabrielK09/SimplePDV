<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card class="pay-ment-forms-dialog">
            <q-card-section class="dialog">
                <header class="border-gray-100">
                    <h2 class="text-gray-600 text-center">Formas de pagamento</h2>

                </header>

                <div v-for="pay in payMentForms" class="flex flex-col flex-center">
                    <div class="flex flex-col flex-center">
                        <span
                            :class="{
                                'border-b border-blue-500 font-bold cursor-pointer': pay.specie === 'Pix'
                            }"
                            @click.prevent="pay.specie === 'Pix' && pixKey === '' ? showChangePixKey = !showChangePixKey : null"
                        >
                            {{ pay.specie }}
                        </span>

                        <div class="mt-4" v-if="showChangePixKey && pay.specie === 'Pix'">
                            <q-input
                                v-model="pixKey"
                                type="text"
                                class="rounded-md"
                                input-class="text-center"
                                dense
                                outlined
                                maxlength="255"
                            />
                        </div>
                    </div>
                </div>

                <div class="flex justify-end mt-4">
                    <q-btn 
                        dense
                        color="red" 
                        icon="close" 
                        @click="emits('close', true)" 
                    />
                    
                    <q-btn
                        dense
                        icon="save"
                        class="ml-4"
                        @click="savePayMentForm"
                    />
                </div>
            </q-card-section>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { getAllPayMentFormsService, updatePayMentFormService } from 'src/modules/PDV/services/payMentFormsService';
    import { onMounted, onUnmounted, ref } from 'vue';

    const showChangePixKey = ref<boolean>(false);
    const internalDialog = ref<boolean>(true);
    const payMentForms = ref<PayMentFormContract[]>();
    const { notify } = useNotify();

    const emits = defineEmits<{
        (e: 'close', value: boolean): void
    }>();

    const pixKey = ref<string>('');

    const getPayMentForms = async () => {
        const res = await getAllPayMentFormsService();
        console.log(res);

        payMentForms.value = res.data;

        payMentForms.value.map((p) => {
            if(p.pix_key !== "")
            {
                pixKey.value = p.pix_key;
            };
        });
    };

    const savePayMentForm = async () => {
        const res = await updatePayMentFormService(pixKey.value);

        if(res.success)
        {
            notify(
                'positive',
                res.message
            );

            emits('close', true)
        } else {
            notify(
                'negative',
                res.message
            );
        };
    };

    const onKeyDownEnter = (e: KeyboardEvent) =>
    {
        if(e.key.toLocaleLowerCase() !== 'enter') return;
        if(!internalDialog.value) return;

        if(pixKey.value === '') return;

        savePayMentForm();
    };    
    
    onMounted(() => {
        getPayMentForms();
        document.addEventListener('keydown', onKeyDownEnter);
    });

    onUnmounted(() => {
        document.removeEventListener('keydown', onKeyDownEnter);
    });
</script>

<style scoped>
    .pay-ment-forms-dialog {
        width: 90%;
        max-width: 250px;
        min-width: 120px;
        border-radius: 18px;
    }

    .dialog {
        padding: 10px 14px;
        background: linear-gradient(to right, #f8fafc, #ffffff);
    }
</style>