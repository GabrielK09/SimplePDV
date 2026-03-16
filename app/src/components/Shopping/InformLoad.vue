<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card>
            <main class="rounded-md flex flex-center text-xl mt-4 bg-white">
                <section class="w-[80vh] rounded-lg px-4">
                    <q-input
                        v-model.number="shoppingDataPayLoad.load"
                        type="text"
                        label="N° Carga da compra"
                        stack-label
                        label-slot

                    />

                    <q-btn
                        color="primary"
                        label="OK"
                        no-caps
                        @click="saveShoppingForPay"
                    />
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
    import { ref } from 'vue';
    import { createshopping } from 'src/modules/shopping/services/shoppingService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import PayMentSale from '../PayMent/Pay/PayMentSale.vue';

    const { notify } = useNotify();

    const props = defineProps<{
        shoppingData: ShoppingContract
    }>();

    const shoppingDataPayLoad = ref<ShoppingContract>(props.shoppingData);
    const shoppingId = ref<number | null>(null);
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

        removeSessionData('shopping_id');
    };

</script>
