<template>
    <q-dialog v-model="isValidSale" persistent>
        <q-card>
            <q-card-section class="row items-center">
                <q-avatar icon="signal_wifi_off" color="primary" text-color="white" />
                <span class="q-ml-sm">You are currently not connected to any network.</span>
            </q-card-section>
            <q-card-actions align="right">
                <q-btn flat label="Cancel" color="primary" v-close-popup />
                <q-btn flat label="Turn on Wifi" color="primary" v-close-popup />
            </q-card-actions>
        </q-card>
    </q-dialog>

    <div class="fixed inset-0 z-50 flex items-center justify-center bg-opacity-40 backdrop-blur-sm">
        <div class="w-[80vh] bg-white">
            <header class="border-gray-100 flex justify-between">
                <h2 class="text-gray-600 text-center ml-4">Formas de pagamento - ID da venda: {{ props.saleId }}</h2>

                <div v-for="(payMentForm, i) in payMentForms">
                    <q-input 
                        v-model="payMentValues[i]" 
                        type="text" 
                        label="Label" 
                    />
                    
                </div>
            </header>
        </div>
    </div>
</template>

<script setup lang="ts">
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { getAllPayMentFormsService } from 'src/modules/PDV/services/payMentFormsService';
    import { onMounted, ref } from 'vue';

    const { notify } = useNotify();

    const props = defineProps<{
        saleId: number

    }>();

    const isValidSale = ref<boolean>(true);

    const payMentForms = ref<PayMentFormContract[]>([]);

    const payMentValues = ref<any[]>();

    const getPayMentForms = async () => {
        const res = await getAllPayMentFormsService();

        if(res.sucess)
        {
            payMentForms.value = res.data;

        } else {
            notify(
                'negative',
                res.message
            );

            isValidSale.value = !isValidSale.value;
        };
    };  

    onMounted(() => {
        if(props.saleId === 0 || !props.saleId) 
        {
            isValidSale.value = !isValidSale.value;
            notify(
                'negative',
                'Identficador da venda inválido!'
            );
        };

        getPayMentForms();
    });
</script>