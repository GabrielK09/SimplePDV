<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card>
            <main class="rounded-md flex flex-center mt-4 bg-white text-xl">
                <section class="rounded-lg p-6 flex flex-col">
                    <span class="font-bold mb-2">Última carga: {{ lastShoppingId }}</span>

                    <q-input
                        v-model.number="shoppingDataPayLoad.load"
                        type="text"
                        stack-label
                        label-slot
                        input-class="text-lg"
                    >
                        <template v-slot:label>
                            <span class="font-bold text-lg">
                                N° Carga da compra
                            </span>
                        </template>
                    </q-input>

                    <div class="mt-4">
                        <q-btn
                            color="primary"
                            label="Confirmar carga da compra"
                            no-caps
                            @click="saveShoppingForPay"
                        />

                    </div>
                </section>
            </main>
        </q-card>
    </q-dialog>

    <PayMentSale
        v-if="showPayMentForms"
        :shopping-id="shoppingId"
        :total-sale="shoppingDataPayLoad.totalShopping"
        @close="showPayMentForms = !$event"
        @paide="resetShopping(!$event)"
    />

</template>

<script setup lang="ts">
    import { SessionStorage } from 'quasar';
    import { onMounted, ref } from 'vue';
    import { createshopping, getLastShoppingId } from 'src/modules/shopping/services/shoppingService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import PayMentSale from '../PayMent/Pay/PayMentSale.vue';
    import { useRouter } from 'vue-router';

    const { notify } = useNotify();

    const router = useRouter();
    const props = defineProps<{
        shoppingData: ShoppingContract
    }>();

    const shoppingDataPayLoad = ref<ShoppingContract>(props.shoppingData);
    const shoppingId = ref<number | null>(null);
    const lastShoppingId = ref<number>(0);
    const internalDialog = ref<boolean>(true);
    const showPayMentForms = ref<boolean>(false);

    const removeSessionData = (key: string): void => {
        SessionStorage.remove(key);
    };

    const saveShoppingForPay = async (): Promise<void> => {
        const existingShoppinId: number = SessionStorage.getItem('shopping_id');

        if(existingShoppinId >= 1)
        {
            shoppingId.value = existingShoppinId;
            showPayMentForms.value = true;

            return;
        };

        const res = await createshopping(shoppingDataPayLoad.value);

        if(res.success)
        {
            shoppingId.value = res.data;

            SessionStorage.set('shopping_id', shoppingId.value);

            showPayMentForms.value = true;

            notify(
                'positive',
                res.message
            );

            return;
        };

        notify(
            'negative',
            res.message
        );

        return;
    };

    const resetShopping = (event: boolean) => {
        showPayMentForms.value = event;
        internalDialog.value = event;

        router.replace({
            name: 'shopping.index'
        });

        removeSessionData('shopping_id');
    };

    onMounted(async () => {
        const res = await getLastShoppingId();
        const data = res.data;

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );   

            internalDialog.value = false;
            return;  
        };

        lastShoppingId.value = data;
    }); 
</script>