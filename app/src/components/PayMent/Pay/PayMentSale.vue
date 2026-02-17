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
                                    <div class="flex">
                                        <span
                                            v-if="calculatePayMent.totalPaid > 0"
                                            class="mt-auto mb-auto mr-2"
                                            @click="resetValues"
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
                                            :disable="calculatePayMent.totalPaid >= props.totalSale"
                                            reverse-fill-mask
                                        />
                                    </div>
                                </q-item-section>
                            </q-item>
                        </q-list>
                    </q-form>
                </div>
            </div>

            <div class="px-8">
                <div class="flex flex-col gap-2 mb-6">
                    <q-chip color="red-6" text-color="white">
                        Valor faltante: R$ {{
                            ((props.totalSale - calculatePayMent.totalPaid > 0)
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
        @close="showQrCodePix = !$event"
        @confirm="confirmByPix($event)"
    />
</template>

<script setup lang="ts">
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { getAllPayMentFormsService } from 'src/modules/PDV/services/payMentFormsService';
    import { computed, onMounted, ref, watch } from 'vue';
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

    const calculateChange = computed((): number =>
    {
        return Math.abs(calculatePayMent.value.totalPaid - props.totalSale) > 0 ? Math.abs(calculatePayMent.value.totalPaid - props.totalSale) : 0;
    });

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

    const resetValues = () => {
        payMentValues.value = payMentForms.value.map((f) => ({
            id: f.id,
            specie: f.specie,
            amount: '0,00',
            pix_key: f.pix_key
        }));
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

    const finallySale = async () => {
        console.log(payMentValues.value);

        if(calculatePayMent.value.totalPaid < 0) return;

        emits('paide', false);
    };

    const confirmByPix = (val: boolean): void => {
        console.log('confirmByPix');

        if(val)
        {
            notify(
                'positive',
                'Pagamento confirmado com sucesso!'
            );

            internalDialog.value = false;
            showQrCodePix.value = false;
            emits('paide', false);
            return;
        };
    };

    watch(payMentValues, (values) => {
        const pixPayment = values.find(p =>
            p.id === 2 && parseFloat(p.amount) > 0
        );

        if(pixPayment)
        {
            showQrCodePix.value = true;
            getPixKey.value = pixPayment.pix_key;

        } else {
            showQrCodePix.value = false;

        };
    }, { deep: true });

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
