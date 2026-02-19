<template>
    <q-dialog v-model="internalDialog" persistent>
        <div class="flex flex-col bg-white w-full phone:bg-black shadow-lg">
            <div class="bg-white p-4">
                <div class="bg-white p-4">
                    <div class="text-h6">Formas de Pagamento - Venda: {{ props.saleId }}</div>

                    <q-list bordered separator class="bg-white text-black">
                        <q-item v-for="(payment, i) in payMentForms" :key="i">
                            <q-item-section>
                                {{ payment.specie }}
                            </q-item-section>

                            <q-item-section side>
                                <div class="flex">
                                    <span
                                        v-if="Number(payMentValues[i].amount.replace(',', '.') || 0) > 0"
                                        class="mt-auto mb-auto mr-2"
                                        @click="resetValues(payMentValues[i].id)"
                                    >
                                        <q-btn
                                            color="red"
                                            icon="delete"
                                            dense
                                            size="7px"
                                        />
                                    </span>

                                    <q-input
                                        v-model="payMentValues[i].amount"
                                        input-class="text-right"
                                        class="w-24"
                                        dense
                                        outlined
                                        placeholder="0,00"
                                        mask="##,##"
                                        fill-mask="0"
                                        :disable="calculatePayMent.totalPaid >= props.totalSale || payMentValues[i].specie.toLocaleLowerCase() === 'pix' && payMentValues[i].pix_key === ''"
                                        reverse-fill-mask
                                    />
                                </div>
                            </q-item-section>
                        </q-item>
                    </q-list>
                </div>
            </div>

            <div class="px-8">
                <div class="flex flex-col gap-2 mb-6">
                    <q-chip color="red-6" text-color="white">
                        Valor faltante: R$ {{
                            (
                                (props.totalSale - calculatePayMent.totalPaid > 0)
                                    ? props.totalSale - calculatePayMent.totalPaid
                                    : 0
                            ).toFixed(2)
                            .toString()
                            .replace('.', ',')
                        }}

                    </q-chip>

                    <q-chip color="green-7" text-color="white">
                        Valor pago: R$ {{ calculatePayMent.totalPaid.toFixed(2).toString().replace('.', ',') }}
                    </q-chip>

                    <q-chip color="blue-6" text-color="white">
                        Troco: R$ {{ calculateChange.toFixed(2).toString().replace('.', ',') }}
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
                            @click="confirmValues"
                            :disable="calculatePayMent.totalPaid < props.totalSale"
                        />
                    </div>
                </q-card-actions>
            </div>
        </div>
    </q-dialog>

    <QRCodePix
        v-if="showQrCodePix"
        :pix-key="getPixKey"
        :total-sale="valueForPix"
        @close="showQrCodePix = !$event"
        @confirm="handlePayPix($event)"
    />
</template>

<script setup lang="ts">
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { getAllPayMentFormsService } from 'src/modules/PDV/services/payMentFormsService';
    import { computed, onMounted, ref, watch } from 'vue';
    import QRCodePix from '../PIX/QRCodePix.vue';
    import { paySaleService } from 'src/modules/PDV/services/pdvService';

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
    const payMentPayLoad = ref<any>();

    const internalDialog = ref<boolean>(true);

    // Necessidades do PIX
    const showQrCodePix = ref<boolean>(false);
    const getPixKey = ref<string>('');
    const havePix = ref<boolean>(false);
    const valueForPix = ref<number>(0);

    const calculateChange = computed((): number =>
    {
        return Math.abs(calculatePayMent.value.totalPaid - props.totalSale) > 0 ? Math.abs(calculatePayMent.value.totalPaid - props.totalSale) : 0;
    });

    watch(payMentValues, (values) => {
        const pixPayment = values.find(p =>
            p.id === 2 && parseFloat(p.amount) > 0
        );

        console.log(pixPayment);

        if(pixPayment)
        {
            havePix.value = true;
            getPixKey.value = pixPayment.pix_key;
            valueForPix.value = Number(pixPayment.amount.replace(',', '.'));

            console.log('Payment 163', pixPayment);
        } else {
            showQrCodePix.value = false;
            havePix.value = false;
            valueForPix.value = 0;

            console.log('Payment 169', pixPayment);
        };

    }, { deep: true });

    const syncPayMent = () => {
        payMentValues.value = payMentForms.value.map((f) => ({
            id: f.id,
            specie: f.specie,
            amount: '0,00',
            pix_key: f.pix_key
        }));
    };

    const calculatePayMent = computed(() => {
        let totalPaid = payMentValues.value.reduce((acc, payment) => {
            const value = Number(payment.amount.replace(',', '.') || 0);

            return acc + value;

        }, 0);

        return {
            totalPaid: totalPaid

        };
    });

    const resetValues = (id: number): void => {
        console.log('Teste: ', {
            id: id
        });

        const payMentForm = payMentValues.value.find(p => p.id === id);

        payMentForm.amount = '0,00';
    };

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

    const handlePayPix = (val: boolean): void => {
        if(val)
        {
            notify(
                'positive',
                'Pagamento confirmado com sucesso!'
            );

            finallySale();
        };
    };

    const confirmValues = () => {
        if(calculatePayMent.value.totalPaid < 0) return;

        if(havePix.value)
        {
            showQrCodePix.value = true;
            return;
        };

        if(calculatePayMent.value.totalPaid > 0 && !havePix.value)
        {
            finallySale();
            return;
        };
    };

    const finallySale = async () => {
        internalDialog.value = false;
        showQrCodePix.value = false;

        payMentPayLoad.value = payMentValues.value;

        const res = await paySaleService(payMentValues.value, props.saleId);
        console.log(res);

        if(res.success) emits('paide', true);
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
