<template>
    <q-dialog v-model="confirm" persistent>
        <q-card class="text-base dialog">
            <q-card-section>
                <header class="text-gray-600 text-center">
                    <h2>Detalhes da compra N° {{ props.shoppingId }}</h2>
                </header>

                <article class="flex flex-col gap-4">
                    <div>
                        <span class="font-bold ml-2">
                            Valor total R$:
                        </span>
                        {{ shoppingData.total_shopping.toFixed(2).toString().replace('.', ',') }}
                    </div>

                    <div>
                        <span class="font-bold ml-2">
                            Produtos:
                        </span>

                        <div class="mt-4 p-2">
                            <pre>{{ shoppingData.shopping_itens }}</pre>
                            <q-table
                                v-model:pagination="pagination"
                                :rows="shoppingData.shopping_itens"
                                :columns="columnsForShoppingItens"
                                :hide-bottom="shoppingData.shopping_itens.length < 10"
                                row-key="product_id"
                            >
                                <template v-slot:body="props">
                                    <q-tr :props="props">
                                        <q-td key="name" :props="props">
                                            <span>
                                                {{ `${props.row.product.name.substring(0, 20)}...` }}

                                                <q-tooltip>
                                                    {{ props.row.product.name }}
                                                </q-tooltip>
                                            </span>
                                        </q-td>

                                        <q-td key="qtde_purchased" :props="props">
                                            <span>
                                                {{ props.row.product.qtde_purchased }}
                                            </span>
                                        </q-td>

                                        <q-td key="purchased_value" :props="props">
                                             <span>
                                                R$ {{ props.row.product.purchased_value.toFixed(2).replace('.', ',') }}
                                            </span>
                                        </q-td>

                                        <q-td key="actions" :props="props">
                                            <div class="flex flex-row items-center gap-1">
                                                <q-btn 
                                                    v-if="hasCharacteristics(props.row.product)"
                                                    size="10px"
                                                    color="black"
                                                    :icon="isExpanded(props.row.product.product_id) ? 'expand_less' : 'grid_on'"
                                                    flat
                                                    @click="toggleExpanded(props.row.product.product_id)"
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
                                        v-if="isExpanded(props.row.product.product_id) && hasCharacteristics(props.row.product_with_characteristics)"
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

                    <div>
                        <span class="font-bold ml-2">
                            Status:
                        </span>

                        <span
                            :class="{
                                'text-green-600': shoppingData.status === 'Concluída',
                                'text-red-600': shoppingData.status === 'Cancelada'
                            }"
                        >
                            {{ shoppingData.status }}
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
    import { onMounted, ref } from 'vue';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { QTableColumn } from 'quasar';
    import { getShoppingById } from '../../services/shoppingService';

    const printer = () => {
        window.print()
    };

    const shoppingData = ref<ShoppingContract>({
        id: null,
        load: null,
        shopping_itens: [],
        total_shopping: 0
    });

    const columnsForShoppingItens: QTableColumn[] = [
        {
            field: 'name',
            label: 'Produto',
            name: 'name',
            align: 'center'
        },
        {
            field: 'qtde_purchased',
            label: 'Qtde comprada',
            name: 'qtde_purchased',
            align: 'center'
        },
        {
            field: 'purchased_value',
            label: 'Valor da compra',
            name: 'purchased_value',
            align: 'center',
            format(val): string {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
        {
            name: 'actions',
            label: '',
            field: 'actions',
            align: 'right'
        }
    ];

    const { notify } = useNotify();

    const confirm = ref<boolean>(true);
    const showInternal = ref<boolean>(false);

    const pagination = ref({
        sortBy: 'id',
        rowsPerPage: 20
    });

    const expdandeRows = ref<number[]>([]);

    const emits = defineEmits<{
        (e: 'close', value: boolean): void
    }>();

    const props = defineProps<{
        shoppingId: number
    }>();

    const hasCharacteristics = (row: any): boolean => {
        console.log(row);
        
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

    onMounted(async() => {
        const res = await getShoppingById(props.shoppingId);
        
        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
            emits('close', true);
        };

        console.log(res.data);
        

        const shoppingDetails: ShoppingContract = res.data.shopping;
        const shoppingWithProducts: ShoppingItemContract[] = res.data.shoppingWithProducts;

        shoppingData.value = {
            id: shoppingDetails.id,
            load: shoppingDetails.load,
            shopping_itens: shoppingWithProducts,
            total_shopping: shoppingDetails.total_shopping,
            status: shoppingDetails.status
        };

        showInternal.value = !showInternal.value;
    });
</script>

<style>
    .dialog {
        width: 100%;
        max-width: 1150px;
        min-width: 320px;
        border-radius: 18px;
    }
</style>