<template>
    <q-dialog v-model="confirm" persistent>
        <q-card v-show="showInternal" class="text-base">
            <q-card-section>
                <header class="text-gray-600 text-center">
                    <h2>Detalhes da venda N° {{ props.saleId }}</h2>
                </header>

                <article class="flex flex-col gap-4">
                    <div>
                        <span class="font-bold ml-2">
                            Cliente:
                        </span>
                         {{ saleData.customer }}
                    </div>

                    <div>
                        <span class="font-bold ml-2">
                            Valor total R$:
                        </span>
                        {{ saleData.sale_value.toFixed(2).toString().replace('.', ',') }}
                    </div>

                    <div v-if="Number(totalCommission.replace(',', '.')) > 0">
                        <span class="font-bold ml-2">
                            Valor total de comissão da venda R$:
                        </span>
                        {{ totalCommission }}
                    </div>

                    <div>
                        <span class="font-bold ml-2">
                            Produtos:
                        </span>

                        <div class="mt-4 p-2">
                            <q-table
                                :rows="saleData.products"
                                :columns="columnsForSale"
                                :hide-bottom="saleData.products.length < 10"
                                row-key="name"
                            />
                        </div>
                    </div>

                    <div>
                        <span class="font-bold ml-2">
                            Dados da comissão:
                        </span>

                        <div class="mt-4 p-2" v-if="commissionData.length > 0">
                            <q-table
                                :rows="commissionData"
                                :columns="columnsForCommission"
                                :hide-bottom="commissionData.length < 10"
                                row-key="name"
                            />
                        </div>
                    </div>

                    <div>
                        <span class="font-bold ml-2">
                            Status:
                        </span>
                        
                        <span
                            :class="{
                                'text-green-600': saleData.status === 'Concluída',
                                'text-red-600': saleData.status === 'Cancelado'
                            }"
                        >
                            {{ saleData.status }}
                        </span>
                    </div>
                </article>
            </q-card-section>

            <q-card-actions align="right">
                <q-btn 
                    color="black"
                    flat
                    icon="print"
                    @click="printer"
                />
                
                <q-btn
                    color="red"
                    icon="close"
                    @click="emits('close', true)"
                />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import { computed, onMounted, ref } from 'vue';
    import { getSaleDetailsById } from '../../services/pdvService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { QTableColumn } from 'quasar';

    interface DetailSaleContract {
        readonly id: number;
        readonly customer_id: number;
        customer: string;
        products: SaleItemContract[];
        sale_value: number;
        status: 'Pendente'|'Concluída'|'Cancelado'
    };

    const printer = () => {
        window.print()
    };

    interface CommissionSaleContract {
        name: string;
        sale_value: number;
        commission_by_produtc: number;
        commission_generated: number;
    };

    const columnsForSale: QTableColumn[] = [ 
        {
            field: 'name',
            label: 'Produto',
            name: 'name',
            align: 'center'
        },
        {
            field: 'price',
            label: 'Valor da venda',
            name: 'price',
            align: 'center',
            format(val): string {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
    ];

    const columnsForCommission: QTableColumn[] = [
        {
            field: 'name',
            label: 'Nome',
            name: 'name',
            align: 'left',
            format(val): string {
                return val.substring(0, 10) + '...'
            }
        },
        {
            field: 'commission_by_produtc',
            label: 'Comissão do produto',
            name: 'commission_by_produtc',
            align: 'center',
            format(val): string {
                return `${val} %`
            }
        },
        {
            field: 'commission_generated',
            label: 'Comissão gerada',
            name: 'commission_generated',
            align: 'center',
            format(val): string {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
    ];

    const { notify } = useNotify();

    const confirm = ref<boolean>(true);
    const showInternal = ref<boolean>(false);

    const emits = defineEmits<{
        (e: 'close', value: boolean)
    }>();

    const props = defineProps<{
        saleId: number
    }>();

    const commissionData = ref<CommissionSaleContract[]>([]);

    const saleData = ref<DetailSaleContract>({
        id: 0,
        customer_id: 0,
        customer: '',
        products: [],
        sale_value: 0,
        status: 'Pendente'
    });

    const totalCommission = computed(() => {
        let subTotal: number = 0;

        commissionData.value.map(c => {
            subTotal += c.sale_value * (c.commission_by_produtc / 100)
        });

        return subTotal.toFixed(2).toString().replace('.', ',');
    });

    onMounted(async() => {
        const res = await getSaleDetailsById(props.saleId);

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
            emits('close', true);
        };

        const data = res.data;

        console.log(data);

        saleData.value = data.sale || [];
        commissionData.value = data.commission || [];

        showInternal.value = !showInternal.value;
    });
</script>