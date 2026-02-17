<template>
    <q-dialog v-model="internalDialog" persistent>
        <div class="flex flex-col bg-white w-full phone:bg-black shadow-lg">
            <div class="bg-white p-4">
                <div class="bg-white p-4">
                    <div class="text-h6">Formas de Pagamento - Venda: {{ props.saleId }}</div>
                    <q-form @submit.prevent="finallySale" class="mt-6">
                        <q-list bordered separator class="bg-white text-black">
                            <q-item v-for="(payment, i) in payMentForms" :key="i">
                                <q-item-section>
                                    {{ payment.specie }}
                                </q-item-section>

                                <q-item-section side>
                                    <q-input
                                        v-model="payMentValues[i].amount"
                                        input-class="text-right"
                                        class="w-24"
                                        dense
                                        outlined
                                        placeholder="0,00"
                                        mask="##,##"
                                        fill-mask="0"
                                        :disable="calculatePayMent.totalPaid >= props.totalSale"
                                        reverse-fill-mask
                                    />

                                </q-item-section>
                            </q-item>
                        </q-list>
                    </q-form>
                </div>
            </div>

            <div class="px-8">
                <div class="flex flex-col gap-2 mb-6">
                    <q-chip color="red-6" text-color="white">
                        Valor faltante: R$ {{ totalPaid_.toFixed(2).toString().replace('.', ',') }}

                    </q-chip>

                    <q-chip color="green-7" text-color="white">
                        Valor pago: R$ {{ calculatePayMent.totalPaid.toFixed(2).toString().replace('.', ',') }}
                    </q-chip>

                    <q-chip color="blue-6" text-color="white">
                        Troco: R$ {{ totalChange.toFixed(2).toString().replace('.', ',') }}
                    </q-chip>
                </div>

                <q-banner class="bg-gray-300 q-mb-sm rounded-xl">
                    <div class="text-subtitle2 font-semibold">
                        Total da venda: R$ {{ props.totalSale.toFixed(2).toString().replace('.', ',') }}
                    </div>
                </q-banner>

                <q-card-actions align="right">
                    <div class="flex gap-6">
                        <q-btn
                            label="Cancelar"
                            color="negative"
                            @click="emits('close', false)"

                        />

                        <q-btn
                            color="primary"
                            label="Finalizar venda"
                            type="submit"
                            :disable="calculatePayMent.totalPaid <= 0"
                        />
                    </div>
                </q-card-actions>
            </div>
        </div>
    </q-dialog>

    <QRCodePix
        v-if="showQrCodePix"
        :pix-key="getPixKey"
        :total-sale="props.totalSale"
    />
</template>

<script setup lang="ts">
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { getAllPayMentFormsService } from 'src/modules/PDV/services/payMentFormsService';
    import { computed, onMounted, ref } from 'vue';
    import QRCodePix from '../PIX/QRCodePix.vue';

    type PayMentValue = {
        id: number;
        specie: string;
        amount: string;
        pix_key: string;
    };

    const { notify } = useNotify();

    const props = defineProps<{
        saleId: number,
        totalSale: number;

    }>();

    const emits = defineEmits<{
        (e: 'close', value: boolean),
        (e: 'paide', value: boolean)

    }>();

    const payMentForms = ref<PayMentFormContract[]>([]);
    const payMentValues = ref<PayMentValue[]>([]);
    const payMentSale = ref<PaySaleContract[]>([]);

    const internalDialog = ref<boolean>(true);
    const totalChange = ref<number>(0);

    const showQrCodePix = ref<boolean>(false);
    const getPixKey = ref<string>('');

    function calculateTotalPaid(): number
    {
        let totalPaid = payMentValues.value.reduce((acc, payment) => {
            console.log('Forma usada: ', payment);

            if(payment.pix_key !== '' && Number(payment.amount) > 0 && payment.id === 2)
            {
                showQrCodePix.value = true;
                getPixKey.value = payment.pix_key;
                // Pensar em fazer lógica separada para PIX
                return;

            };

            console.log(payment.pix_key !== '' && Number(payment.amount) > 0 && payment.id === 2);

            const value = Number(payment.amount.replace(',', '.') || 0);

            return acc + value;

        }, 0);

        return totalPaid;
    };

    function calculateChange(totalPaid: number): void
    {
        console.log(`totalPaid: ${totalPaid}: - props.totalSale ${props.totalSale}: `, Math.abs(props.totalSale - totalPaid));

        totalChange.value = Math.abs(totalPaid - props.totalSale) > 0 ? Math.abs(totalPaid - props.totalSale) : 0;
    };

    const syncPayMent = () => {
        payMentValues.value = payMentForms.value.map((f) => ({
            id: f.id,
            specie: f.specie,
            amount: '0,00',
            pix_key: f.pix_key
        }));
    };

    const calculatePayMent = computed(() => {
        return {
            totalPaid: calculateTotalPaid()

        };
    });

    const totalPaid_ = computed((): number => {
        const result = props.totalSale - calculateTotalPaid();
        if(props.totalSale <= result) return result;

        calculateChange(calculateTotalPaid());
        return 0;
    });

    const getPayMentForms = async () => {
        const res = await getAllPayMentFormsService();

        if(res.success)
        {
            payMentForms.value = res.data;
            syncPayMent();

            return;
        };

        notify(
            'negative',
            res.message
        );
    };

    const finallySale = async () => {
        console.log(payMentValues.value);

        if(calculatePayMent.value.totalPaid < 0) return;

        emits('paide', false);
    };

    onMounted(() => {
        getPayMentForms();
    });
</script>

<style lang="scss">
    @media (max-width: 1536px) {
        .pay-ment-form {
            overflow-y: auto !important;
        }
    }

</style>
