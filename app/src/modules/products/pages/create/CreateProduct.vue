<template>    
    <section class="text-2xl">
        <div
            class="m-2"
        >
            <h2 class="text-gray-600 register-title m-2">Cadastrar um(a) novo(a) produto</h2>

            <div class="ml-2 text-xs">
                <div 
                    @click="router.replace({ path: `/admin/products` })"
                    class="flex mb-auto mt-auto cursor-pointer"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="mr-1 back-row">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 15 3 9m0 0 6-6M3 9h12a6 6 0 0 1 0 12h-3" />
                    </svg>      
                    <span class="back-product-label">
                        Voltar para listagem
                    </span>
                </div>
            </div>

            <q-form
                @submit.prevent="submitProduct"
                class="q-gutter-md mt-4 w-[90%]"
            >
                <div class="p-4 inputs">
                    <q-input 
                        v-model="nameUpper" 
                        type="text" 
                        label="E-mail *" 
                        stack-label
                        outlined
                        dense
                        class="mb-4"
                        :error="!!formErrors.trade_name"
                        :error-message="formErrors.trade_name"
                    >
                        <template v-slot:label>
                            <div class="text-sm">
                                Nome <span class="text-red-500">*</span>
                            </div>
                        </template>
                    </q-input>

                    <q-input 
                        v-model="product.price" 
                        type="number" 
                        label="" 
                        stack-label
                        outlined
                        dense
                        class="mb-4"
                        :error="!!formErrors.price"
                        :error-message="formErrors.price"
                    >
                        <template v-slot:label>
                            <div class="text-sm">
                                Preço <span class="text-red-500">*</span>
                            </div>
                        </template>
                    </q-input>

                    <q-input 
                        v-model="product.qtde" 
                        type="number" 
                        label="" 
                        stack-label
                        outlined
                        dense
                        class="mb-4"
                        :error="!!formErrors.qtde"
                        :error-message="formErrors.qtde"
                    >
                        <template v-slot:label>
                            <div class="text-sm">
                                Qtde <span class="text-red-500">*</span>
                            </div>
                        </template>
                    </q-input>

                    <div class="flex flex-center">
                        <q-btn 
                            color="primary" 
                            type="submit" 
                            label="Cadastrar produto"
                            no-caps
                            :loading="loadingLogin"
                        />
                    </div>
                </div>
            </q-form>
        </div>
    </section>
</template>

<script setup lang="ts">
    import { useQuasar } from 'quasar';
    import { computed, ref } from 'vue';
    import { useRouter } from 'vue-router';
    import * as Yup from 'yup';
    import { createProduct } from '../../services/productsService';

    const productSchema = computed(() =>
        Yup.object({
            name: Yup.string().required('O nome do produto é obrigatório!'),
            price: Yup.number().required('O valor do produto é obrigatório!'),
            qtde: Yup.number().required('A quantia do produto é obrigatório!'),
        })
    );

    const product = ref<ProductContract>({
        id: 0,
        name: '',
        price: 0,
        qtde: 0
    });

    const formErrors = ref<Record<string, string>>({});

    const router = useRouter();
    const $q = useQuasar();

    let loadingLogin = ref<boolean>(false);

    const nameUpper = computed({
        get: () => product.value.name,
        set: (val: string) => {
            product.value.name = val.toUpperCase();  
        }
    });

    const submitProduct = async () => {
        try {
            await productSchema.value.validate(product.value, { abortEarly: false });
            
            const res = await createProduct(product.value);

            if(res.success)
            {
                $q.notify({
                    type: 'positive',
                    position: 'top',
                    message: res.data.message
                    
                });
                
                router.replace({
                    name: 'products.index'

                });
                
            } else {
                $q.notify({
                    type: 'negative',
                    position: 'top',
                    message: res.message

                });
            };
            
        } catch (error: any) {
            console.error('Erro:', error.inner);
            
            if(error.inner)
            {
                formErrors.value = {};

                error.inner.forEach((err: any) => {
                    formErrors.value[err.path] = err.message;

                    $q.notify({
                        type: 'negative',
                        position: 'top',
                        message: err.message

                    });
                });  
            } else {
                $q.notify({
                    type: 'negative',
                    position: 'top',
                    message: error.response?.data?.message || 'Erro na criação do produto!'

                });
            };
        };
    };
</script>