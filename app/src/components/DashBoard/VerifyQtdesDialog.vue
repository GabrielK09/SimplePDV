<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card class="text-base dialog">
            <q-card-section>
                <header class="text-gray-600 flex">  
                    <div class="my-auto">
                        <q-btn 
                            color="red" 
                            icon="close" 
                            @click="emits('close', true)" 
                        />
                    </div>

                    <h3 class="text-xl mx-auto mt-4">Resumo de qtdes</h3>
                </header>

                <div class="bg-white p-4 flex">
                    <div class="m-4">
                        Qtde futura de todos os produtos: 
                        <span class="font-bold">
                            {{ qtdesData.total_future }}
                        </span>
                    </div>

                    <div class="m-4 flex">
                        Status:                    
                        <span class="font-bold">
                            <span class="my-auto flex">
                                <div 
                                    class="w-4 h-4 border my-auto mx-2"
                                    :class="{
                                        'bg-red-500': qtdesData.total_future <= 0,
                                        'bg-green-500': qtdesData.total_future > 10
                                    }"
                                ></div>
                                {{ qtdesData.total_future <= 0 ? 'Atenção' : 'Ok' }}
                            </span>
                        </span>
                    </div>

                    <div class="m-4">
                        Qtde reservada de todos os produtos: 
                        <span class="font-bold">
                            {{ qtdesData.total_reservate }}
                        </span>
                    </div>

                    <div class="m-4 flex">
                        Status:                    
                        <span class="font-bold">
                            <span class="my-auto flex">
                                <div 
                                    class="w-4 h-4 border my-auto mx-2"
                                    :class="{
                                        'bg-red-500': qtdesData.total_future >= 0,
                                        'bg-green-500': qtdesData.total_future < 10
                                    }"
                                ></div>
                                {{ qtdesData.total_future >= 0 ? 'Ok' : 'Atenção' }}
                            </span>
                        </span>
                    </div>
                </div>
            </q-card-section>
        </q-card>
    </q-dialog>
</template>


<script setup lang="ts"> 
    import { returnVerifyQtdes } from 'src/modules/products/services/productsService';
    import { onMounted, ref } from 'vue';
    import { useNotify } from 'src/helpers/QNotify/useNotify';

    interface VerifyQtde {
        total_future: number;
        total_reservate: number;
    };

    const emits = defineEmits<{
        (e: 'close', value: boolean): void
    }>();

    const { notify } = useNotify();

    const internalDialog = ref<boolean>(true);

    const qtdesData = ref<VerifyQtde>({
        total_future: 0,
        total_reservate: 0 
    });

    onMounted(async() => {
        const res = await returnVerifyQtdes();

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );

            internalDialog.value = !internalDialog.value;

            return;
        };

        const data: VerifyQtde = res.data;

        qtdesData.value = {
            total_future: data.total_future,
            total_reservate: data.total_reservate
        };
    });
</script>

<style>
    .dialog {
        width: 100%;
        max-width: 350px;
        min-width: 220px;
        border-radius: 18px;
    }
</style>