<template>
    <q-dialog v-model="confirm" persistent>
        <q-card class="text-base">
            <q-card-section>
                <header class="text-gray-600">
                    <h2>Detalhes da venda N° {{ props.saleId }}</h2>
                </header>

                <article class="w-max flex flex-col gap-4">
                    <div>
                        <span class="font-bold">
                            Cliente:
                        </span>
                         {{ saleData.customer }}
                    </div>

                    <div>
                        <span class="font-bold">
                            Valor total R$:
                        </span>
                        {{ saleData.sale_value.toFixed(2).toString().replace('.', ',') }}
                    </div>

                    <div>
                        <span class="font-bold">
                            Produtos:
                        </span>

                        <div class="flex flex-col ml-4">
                            <ul v-for="product in saleData.products">
                                <li class="list-disc">
                                    {{ product.name }}

                                </li>
                            </ul>
                        </div>
                    </div>
                </article>
            </q-card-section>

            <q-card-actions align="right">
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
    import { getSaleDetailsById } from '../../services/pdvService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';

    interface DetailSaleContract {
        readonly id: number;
        readonly customer_id: number;
        customer: string;
        products: SaleItemContract[];
        sale_value: number;
        status: 'Pendente'|'Concluída'|'Cancelada'

    };

    const { notify } = useNotify();
    const confirm = ref<boolean>(true);

    const emits = defineEmits<{
        (e: 'close', value: boolean)
    }>();

    const props = defineProps<{
        saleId: number
    }>();

    const saleData = ref<DetailSaleContract>({
        id: 0,
        customer_id: 0,
        customer: '',
        products: [],
        sale_value: 0,
        status: 'Pendente'
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
        saleData.value = data;
    });

</script>
