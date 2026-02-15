<template>
    <q-dialog v-model="internalDialog" persistent>
        <div class="bg-white p-8 rounded-xl shadow-lg flex flex-col items-center gap-4">
            <q-form
                @submit="finallySale"
            >
                <q-card class=" border border-black mb-5 p-6 bg-white shadow-md rounded" >
                    <q-card-section>
                        <div class="text-h6">Formas de Pagamento - Venda: {{ props.saleId }}</div>

                    </q-card-section>

                    <q-separator />

                    <q-card-section>
                        <q-form @submit.prevent="finallySale">
                            <q-list bordered separator class="bg-white text-black rounded-borders">
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
                                            reverse-fill-mask
                                        />
                                    </q-item-section>
                            </q-item>
                        </q-list>
                            <q-card-section class="bg-white text-black rounded-borders">
                                <div class="row q-gutter-sm mb-2">
                                    <q-chip color="red-6" text-color="white">
                                        Valor faltante: R$

                                    </q-chip>

                                    <q-chip color="green-7" text-color="white">
                                        Valor pago: R$ {{ calculatePayMent.amountPaid }}
                                    </q-chip>

                                    <q-chip color="blue-6" text-color="white">
                                        Troco: R$ {{  }}
                                    </q-chip>
                                </div>

                                <q-banner class="bg-gray-300 q-mb-sm rounded-xl">
                                    <div class="text-subtitle2 font-semibold">
                                        Total da venda: R$ {{ props.totalSale }}
                                    </div>
                                </q-banner>

                            </q-card-section>

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
                                    />
                                </div>
                            </q-card-actions>
                        </q-form>
                    </q-card-section>
                </q-card>
            </q-form>
        </div>
    </q-dialog>
</template>

<script setup lang="ts">
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { getAllPayMentFormsService } from 'src/modules/PDV/services/payMentFormsService';
    import { computed, onMounted, ref } from 'vue';

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
        (e: 'paide', value: boolean),

    }>();

    const payMentForms = ref<PayMentFormContract[]>([]);
    const payMentValues = ref<PayMentValue[]>([]);
    const payMentSale = ref<PaySaleContract[]>([]);

    const internalDialog = ref<boolean>(true);

    const calculatePayMent = computed(() => {
        let subTotalSale: number;
        let totalPaid = payMentValues.value.reduce((acc, payment) => {
            const value = Number(payment.amount.replace(',', '.') || 0);
            return acc + value;

        }, 0);

        return {
            subTotal: subTotalSale,
            amountPaid: 0,
            totalPaid: totalPaid

        };
    });

    const syncPayMent = () => {
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

        if(calculatePayMent.value.amountPaid < 0) return;

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
