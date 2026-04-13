<template>
    <q-dialog v-model="confirm" persistent>
        <q-card class="dialog">
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
                                v-model:pagination="pagination"
                                :rows="saleData.products"
                                :columns="columnsForSale"
                                :hide-bottom="saleData.products.length < 10"
                                row-key="product_id"
                            >   
                                <template v-slot:body="props">
                                    <q-tr :props="props">
                                        <q-td key="name" :props="props">
                                            <span>
                                                {{ `${props.row.name.substring(0, 20)}...` }}

                                                <q-tooltip>
                                                    {{ props.row.name }}
                                                </q-tooltip>
                                            </span>
                                        </q-td>

                                        <q-td key="qtde" :props="props">
                                            <span>
                                                {{ props.row.qtde }}
                                            </span>
                                        </q-td>

                                        <q-td key="qtde" :props="props">
                                            <span>
                                                {{ props.row.sale_value.toFixed(2).replace('.', ',') }}
                                            </span>
                                        </q-td>

                                        <q-td key="actions" :props="props">
                                            <div class="flex flex-row items-center gap-1">
                                                <q-btn
                                                    v-if="hasCharacteristics(props.row)"
                                                    size="10px"
                                                    color="black"
                                                    :icon="isExpanded(props.row.product_id) ? 'expand_less' : 'grid_on'"
                                                    flat
                                                    @click="toggleExpanded(props.row.product_id)"
                                                >
                                                    <q-tooltip>
                                                        {{
                                                            isExpanded(props.row.product_id)
                                                                ? 'Ocultar'
                                                                : 'Ver grade'
                                                        }}
                                                    </q-tooltip>
                                                </q-btn>
                                            </div>
                                        </q-td>
                                    </q-tr>

                                    <q-tr
                                        v-if="isExpanded(props.row.product_id) && hasCharacteristics(props.row)"
                                        :props="props"
                                    >
                                        <q-td colspan="100%" class="bg-gray-200">
                                            <div class="q-pa-md">
                                                <div class="text-subtitle2 text-weight-bold q-mb-sm">
                                                    Grade do produto
                                                </div>

                                                <div class="row q-col-gutter-sm">
                                                    <div
                                                        v-for="(characteristic, i) in props.row.product_with_characteristics"
                                                        :key="`${props.row.product_id}-${characteristic.size}`"
                                                        class="col-12 col-sm-6 col-md-3"
                                                    >
                                                        <q-card flat bordered>
                                                            <q-card-section class="q-pa-pmd">
                                                                <div class="text-caption text-gray-700">
                                                                    Tamanho

                                                                </div>
                                                                <div class="text-body2 text-weight-bold">
                                                                    {{ characteristic.size }}
                                                                </div>

                                                                <div class="text-caption text-grey-7 q-mt-sm">Quantidade</div>
                                                                <div>
                                                                    <q-input
                                                                        v-model.number="characteristic.grid_qtde"
                                                                        type="number"
                                                                        class="w-12 flex ml-auto mr-auto"
                                                                        input-class="text-center"
                                                                        dense
                                                                        disable
                                                                    />
                                                                </div>
                                                            </q-card-section>
                                                        </q-card>
                                                    </div>
                                                </div>
                                            </div>
                                        </q-td>
                                    </q-tr>
                                </template>
                            </q-table>
                        </div>
                    </div>

                    <div v-if="commissionData.length > 0">
                        <span class="font-bold ml-2">
                            Dados da comissão:
                        </span>

                        <div class="mt-4 p-2">
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

    interface CommissionSaleContract {
        name: string;
        sale_value: number;
        commission_by_produtc: number;
        commission_generated: number;
    };
    
    const printer = () => {
        window.print()
    };
    
    const pagination = ref({
        sortBy: 'id',
        rowsPerPage: 20
    });

    const columnsForSale: QTableColumn[] = [
        {
            field: 'name',
            label: 'Produto',
            name: 'name',
            align: 'center'
        },
        {
            field: 'qtde',
            label: 'Qtde vendida',
            name: 'qtde',
            align: 'center'
        },
        {
            field: 'sale_value',
            label: 'Valor da venda',
            name: 'sale_value',
            align: 'center'
        },
        {
            name: 'actions',
            label: '',
            field: 'actions',
            align: 'right'
        }
    ];

    const columnsForCommission: QTableColumn[] = [
        {
            field: 'name',
            label: 'Nome',
            name: 'name',
            align: 'left'
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
        (e: 'close', value: boolean): void
    }>();

    const props = defineProps<{
        saleId: number
    }>();

    const commissionData = ref<CommissionSaleContract[]>([]);
    const expdandeRows = ref<number[]>([]);

    const saleData = ref<DetailSaleContract>({
        id: 0,
        customer_id: 0,
        customer: '',
        products: [],
        sale_value: 0,
        status: 'Pendente'
    });

    const hasCharacteristics = (row: any): boolean => {
        return Array.isArray(row.product_with_characteristics) && row.product_with_characteristics.length > 0;
    };

    const isExpanded = (productId: number) => {
        return expdandeRows.value.includes(productId);
    };

    const toggleExpanded = (productId: number): void => {
        if (isExpanded(productId)) {
            expdandeRows.value = expdandeRows.value.filter(id => id !== productId);
            return;
        };

        expdandeRows.value.push(productId);
    };

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

        const saleDetails: DetailSaleContract = data.sale;
        const saleWithProducts: SaleItemContract[] = data.sale_with_products;
        commissionData.value = data.commission || [];

        saleData.value = {
            id: saleDetails.id,
            customer: saleDetails.customer,
            customer_id: saleDetails.customer_id,
            sale_value: saleDetails.sale_value,
            products: saleWithProducts,
            status: saleDetails.status
        };

        showInternal.value = !showInternal.value;
    });
</script>

<style>
    .dialog {
        width: 130%;
        max-width: 1200px;
        min-width: 620px;
        border-radius: 18px;
    }
</style>
