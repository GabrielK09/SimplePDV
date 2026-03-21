<template>
    <q-dialog v-model="confirm" persistent>
        <q-card v-show="showInternal" class="text-base">
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
                            <q-table
                                :rows="shoppingData.shopping_itens"
                                :columns="columnsForShoppingItens"
                                :hide-bottom="shoppingData.shopping_itens.length < 10"
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
                                'text-green-600': shoppingData.status === 'Concluída',
                                'text-red-600': shoppingData.status === 'Cancelado'
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
            field: 'purchased_value',
            label: 'Valor da compra',
            name: 'purchased_value',
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
        shoppingId: number
    }>();

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

        const shoppingDetails: ShoppingContract = res.data.shopping;
        const shoppingWithProducts: ShoppingItemContract[] = res.data.shoppingWithProducts;

        console.log(res.data);

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
