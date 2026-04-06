<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card class="flex flex-center text-xl overflow-hidden">
            <div class="flex flex-col bg-white w-full phone:bg-black shadow-lg">
                <header class="text-gray-600 text-center">
                    <h3 class="text-xl mt-4">Tamanhos</h3>
                </header>
            </div>
            
            <div class="flex flex-col p-4 ">
                <q-form
                    @submit="submitGrid"
                    class="q-gutter-md"
                >
                    <BaseSelectGridTypes
                        v-model="productCharacteristics.size"                    
                        :selected-sizes="props.selectedSizes"
                        :error="!!formErrors.size"
                        :error-message="formErrors.size"
                        :is-update="true"
                    />

                    <q-input 
                        v-model="productCharacteristics.grid_qtde" 
                        type="text" 
                        class="mt-7 mb-4"
                        label="Qtde da grade" 
                        stack-label
                        label-slot
                        outlined
                        dense
                    />
                    
                    <div class="flex flex-center gap-4">
                        <q-btn 
                            color="red"
                            label="Cancelar"
                            no-caps
                            @click="emits('close', true)"
                        />

                        <q-btn 
                            label="Gravar" 
                            color="primary" 
                            no-caps
                            type="submit"
                        />
                    </div>
                </q-form>
            </div>
        </q-card>
    </q-dialog>
</template> 

<script setup lang="ts">
    import { computed, onMounted, ref } from 'vue';
    import BaseSelectGridTypes from 'src/components/Products/UseGrid/QSelectGridTypes/BaseSelectGridTypes.vue';
    import * as Yup from 'yup';
    import { useNotify } from 'src/helpers/QNotify/useNotify';  
    import { getProductCharacteristicsByGridIds } from 'src/modules/products/services/productsService';

    const internalDialog = ref<boolean>(true);
    const { notify } = useNotify();

    const emits = defineEmits<{
        (e: 'return:grids', value: any)
        (e: 'close', value: boolean)
    }>();

    const props = defineProps<{
        productId: number;
        gridId: number;
        gridFullObject?: any;
        selectedSizes: any[];

    }>();

    const gridSchema = computed(() =>
        Yup.object({
            grid_qtde: Yup
                .number()
                .min(1, 'A qtde da grade não pode ser menor que zero.')
                .required('A qtde da grade do produto é obrigatória!'),
        })
    );

    const formErrors = ref<Record<string, string>>({});

    const productCharacteristics = ref<ProductCharacteristicsContract>({
        grid_qtde: null,
        id: null,
        product_id: props.productId || null,
        size: null
    });

    const submitGrid = async (): Promise<any> => {
        try {
            await gridSchema.value.validate(productCharacteristics.value, { abortEarly: false });
           
            emits('return:grids', productCharacteristics.value);
            emits('close', true);
            
        } catch (error: any) {
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

    onMounted(async () => {
        if(props.gridId && props.productId)
        {
            const res = await getProductCharacteristicsByGridIds({gridId: props.gridId, productGridId: props.productId});
        
            if (!res.success)
            {
                notify(
                    'negative',
                    res.message
                );
                return;      
            };
        
            const data: ProductCharacteristicsContract = res.data;
        
            productCharacteristics.value = {
                id: data.id,
                grid_qtde: data.grid_qtde,
                product_id: data.product_id,
                size: data.size
            };  

        } else {
            productCharacteristics.value = {
                id: null,
                grid_qtde: props.gridFullObject?.grid_qtde,
                product_id: props.gridFullObject?.product_id,
                size: props.gridFullObject?.size
            };
        };
    }); 
</script>