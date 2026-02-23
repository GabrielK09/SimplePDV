<template>
    <q-dialog v-model="confirm" persistent>
        <q-card>
            <div class="fixed inset-0 z-50 flex items-center justify-center bg-opacity-40 backdrop-blur-sm">
                <div class="bg-white p-4 rounded-lg">
                    <header class="flex justify-between w-max">
                        <svg
                            xmlns="http://www.w3.org/2000/svg" fill="none"
                            viewBox="0 0 24 24"
                            stroke-width="1.5"
                            stroke="currentColor"
                            class="size-6 my-auto text-red-600 cursor-pointer"
                            @click="emits('close', true)"
                        >
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
                        </svg>

                        <h2 class="text-gray-600 text-center ml-4">Formas de pagamento</h2>
                    </header>

                    <div v-for="pay in payMentForms">
                        <ul class="m-4">
                            <li class="text-center flex justify-center">
                                <span
                                    :class="{
                                        'border-b border-blue-500 font-bold cursor-pointer': pay.specie === 'Pix'
                                    }"
                                    @click.prevent="pay.specie === 'Pix' ? showChangePixKey = !showChangePixKey : null"
                                >
                                    {{ pay.specie }}
                                </span>
                            </li>

                            <div class="flex flex-center mt-4" v-if="showChangePixKey && pay.specie === 'Pix'">
                                <q-input
                                    v-model="pixKey"
                                    type="text"
                                    class="rounded-md"
                                    stack-label
                                    label-slot
                                    dense
                                    label="Chave PIX"
                                    maxlength="255"
                                />

                                <q-btn
                                    dense
                                    icon="save"
                                    class="ml-4"
                                    @click="savePayMentForm"
                                />
                            </div>
                        </ul>
                    </div>
                </div>
            </div>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { getAllPayMentFormsService, updatePayMentFormService } from 'src/modules/PDV/services/payMentFormsService';
    import { onMounted, ref } from 'vue';

    const showChangePixKey = ref<boolean>(false);
    const confirm = ref<boolean>(true);
    const payMentForms = ref<PayMentFormContract[]>();
    const { notify } = useNotify();

    const emits = defineEmits<{
        (e: 'close', value: boolean)
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

    onMounted(() => {
        getPayMentForms();
    });
</script>
