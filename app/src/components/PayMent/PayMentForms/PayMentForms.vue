<template>
    <q-dialog v-model="confirm" persistent>
        <q-card>
            <div class=" fixed inset-0 z-50 flex items-center justify-center bg-opacity-40 backdrop-blur-sm">
                <div class="bg-white p-4 rounded-lg">
                    <header class="flex justify-between">
                        <q-btn 
                            color="red" 
                            flat
                            icon="close" 
                            @click="confirm = !confirm"
                            
                        />
                        <h2 class="text-gray-600 text-center ml-4">Formas de pagamento</h2>
                    </header>
                
                    <div v-for="pay in data">
                        <ul class="m-4">
                            <li class="text-center flex justify-center">
                                <span
                                    class=""
                                    :class="{
                                        'border-b border-blue-500 font-bold cursor-pointer': pay.specie === 'Pix'
                                    }"
                                    @click.prevent="pay.specie === 'Pix' ? showChangePixKey = !showChangePixKey : null"
                                >
                                    {{ pay.specie }}

                                </span>
                            </li>

                            <div class="w-full flex flex-center mt-4" v-if="showChangePixKey && pay.specie === 'Pix'">
                                <q-input 
                                    v-model="specie.pix_key" 
                                    type="text" 
                                    class="border rounded-md"
                                    stack-label
                                    label-slot
                                    dense
                                    borderless
                                    maxlength="255"
                                >   
                                    <template v-slot:label>
                                        <span class="ml-2">
                                            Chave PIX
                                        </span>
                                    </template>
                                </q-input>  
                                <q-btn  
                                    dense
                                    icon="save"
                                    class="ml-4"
                                    @click="savePayMentForm"
                                />
                            </div>              
                        </ul>          
                    </div>
                </div>
            </div>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { getAllPayMentFormsService, updatePayMentFormService } from 'src/modules/PDV/services/payMentFormsService';
    import { onMounted, ref } from 'vue';
    
    const showChangePixKey = ref<boolean>(false);
    const confirm = ref<boolean>(true);
    const data = ref<PayMentFormContract[]>();
    const { notify } = useNotify();

    const specie = ref<PayMentFormContract>({
        id: 0,
        specie: '',
        pix_key: ''
    });

    const getPayMentForms = async () => {
        const res = await getAllPayMentFormsService();

        data.value = res.data;
    };

    const savePayMentForm = async () => {
        const res = await updatePayMentFormService(specie.value);

        console.log(res);
        
        if(res.success)
        {
            notify(
                'positive',
                res.message
            );

            confirm.value = false;
        } else {
            notify(
                'negative',
                res.message
            );
        };
    };

    onMounted(() => {
        getPayMentForms();
    });
</script>