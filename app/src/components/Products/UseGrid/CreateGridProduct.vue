<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card class="flex flex-center text-xl overflow-hidden">
            <div class="flex flex-col bg-white w-full phone:bg-black shadow-lg">
                <header class="text-gray-600 text-center">
                    <h3 class="text-xl mt-4">Tamanhos</h3>
                </header>
            </div>
            
            <div class="flex flex-col p-4 ">
                <BaseSelectGridTypes
                    v-model="productCharacteristics.size"
                    :selected-sizes="props.registredGrids"
                    :error="!!formErrors.size"
                    :error-message="formErrors.size"
                />

                <q-input 
                    v-model="productCharacteristics.grid_qtde" 
                    type="text" 
                    class="my-6"
                    label="Qtde da grade" 
                    stack-label
                    label-slot
                />
                
                <div>
                    <q-btn label="Submit" color="primary" @click="submitGrid"/>

                    <q-btn 
                        color="red"
                        no-caps 
                        icon="close"
                        @click="emits('close', true)"
                    />
                </div>
            </div>
        </q-card>
    </q-dialog>
</template> 

<script setup lang="ts">
    import { computed, ref } from 'vue';
    import BaseSelectGridTypes from './QSelectGridTypes/BaseSelectGridTypes.vue';
    import * as Yup from 'yup';
    import { useNotify } from 'src/helpers/QNotify/useNotify';  

    enum Sizes {
        PP = 'PP', 
        P = 'P', 
        M = 'M', 
        G = 'G', 
        GG = 'GG', 
        XG = 'XG', 
        XGG = 'XGG', 
        EG = 'EG', 
        EGG = 'EGG', 
        O = 'O'
    };

    const internalDialog = ref<boolean>(true);
    const { notify } = useNotify();

    const emits = defineEmits<{
        (e: 'return:grids', value: any)
        (e: 'close', value: boolean)
    }>();

    const props = defineProps<{
        productId?: number;
        registredGrids: ProductCharacteristicsContract[]
    }>();

    const gridSchema = computed(() =>
        Yup.object({
            size: Yup
                .mixed<Sizes>()
                .oneOf(Object.values(Sizes))
                .required('O tamanho da grade é obrigatório.'),

            grid_qtde: Yup
                .number()
                .min(1, 'A qtde da grade não pode ser menor que zero.')
                .required('A qtde da grade do produto é obrigatória!'),
        })
    );

    const formErrors = ref<Record<string, string>>({});

    const productCharacteristics = ref<ProductCharacteristicsContract>({
        grid_qtde: 0,
        id: null,
        product_id: props.productId || null,
        size: null
    });

    const submitGrid = async (): Promise<any> => {
        console.log('called submitGrid');
        
        try {
            await gridSchema.value.validate(productCharacteristics.value, { abortEarly: false });

            console.log(productCharacteristics.value);
            
            emits('return:grids', productCharacteristics.value);
            emits('close', true);
            
        } catch (error) {
            console.log(error);
            
            if(error.inner)
            {
                formErrors.value = {};

                error.inner.forEach((err: any) => {
                    formErrors.value[err.path] = err.message;

                    notify(
                        'negative',
                        err.message

                    );
                });
            } else {
                notify(
                    'negative',
                    error.response?.data?.message || 'Erro na geração da grade!'
                );
            };
        };
    };
</script>