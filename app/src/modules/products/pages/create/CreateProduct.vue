<template>
    <q-page padding>
        <main class="min-h-[60vh] flex flex-center text-xl">
            <section class="w-[80vh] rounded-lg shadow px-4 bg-white">
                <header class="border-gray-100 flex">
                    <span class="text-black cursor-pointer my-auto">
                        <router-link to="/admin/products">
                            <q-avatar size="30px" icon="arrow_back" />

                        </router-link>
                    </span>

                    <h2 class="text-gray-600 text-center">Cadastrar um novo produto</h2>

                </header>

                <q-form
                    @submit.prevent="submitProduct"
                    class="q-gutter-md mt-4"
                >
                    <div class="p-4 inputs">
                        <q-input
                            v-model="nameUpper"
                            maxlength="255"
                            type="text"
                            label-slot
                            stack-label
                            outlined
                            dense
                            class="mb-4"
                            :error="!!formErrors.name"
                            :error-message="formErrors.name"
                        >
                            <template v-slot:label>
                                <div class="text-sm">
                                    Nome <span class="text-red-500">*</span>
                                </div>
                            </template>
                        </q-input>

                        <q-input
                            v-model="product.price"
                            type="text"
                            mask="###,##"
                            label-slot
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
                            type="text"
                            label-slot
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
                        
                        <div class="flex flex-col mb-4">
                            <q-input
                                v-model="product.commission"
                                type="text"
                                mask="##,##"
                                label-slot
                                stack-label
                                outlined
                                dense
                                :error="!!formErrors.commission"
                                :error-message="formErrors.commission"
                            >
                                <template v-slot:label>
                                    <div class="text-sm">
                                        % Comissão <span class="text-red-500">*</span>
                                    </div>
                                </template>
                            </q-input>

                            <div class="mx-auto flex gap-4">
                                <q-btn 
                                    color="primary" 
                                    no-caps 
                                    @click="product.commission = 15"
                                    label="15%"
                                />

                                <q-btn 
                                    color="primary" 
                                    no-caps 
                                    @click="product.commission = 25"
                                    label="25%"
                                />

                                <q-btn 
                                    color="primary" 
                                    no-caps 
                                    @click="product.commission = 35"
                                    label="35%"
                                />
                            </div>
                        </div>

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
            </section>
        </main>
    </q-page>
</template>

<script setup lang="ts">
    import { computed, ref, watch } from 'vue';
    import { useRouter } from 'vue-router';
    import * as Yup from 'yup';
    import { createProduct } from '../../services/productsService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';

    const productSchema = computed(() =>
        Yup.object({
            name: Yup.string().required('O nome do produto é obrigatório!'),

            price: Yup
                    .number()
                    .min(1, 'O valor do produto não pode ser menor que zero.')
                    .required('O valor do produto é obrigatório!'),
            
            qtde: Yup
                    .number()
                    .min(1, 'A qtde do produto não pode ser menor que zero.')
                    .required('A quantia do produto é obrigatório!'),

            commission: Yup
                    .number()
                    .min(0, 'O valor de comissão não pode ser menor que zero.')
                    .max(100, 'O valor de comissão não pode ser maior que 100%.')
                    .required('A quantia do produto é obrigatório!')
        })
    );

    const product = ref<ProductContract>({
        id: 0,
        name: '',
        price: null,
        qtde: null,
        commission: null
    });

    const formErrors = ref<Record<string, string>>({});

    const router = useRouter();
    const { notify } = useNotify();

    let loadingLogin = ref<boolean>(false);

    const nameUpper = computed({
        get: () => product.value.name,
        set: (val: string) => {
            product.value.name = val.toUpperCase();
        }
    });

    watch(
        () => product.value.price,
        (val) => {
            console.log(val);
            
            product.value.price = isNaN(Number(val.toString().replace(',', '.'))) ? 0 : val;
        }
    );

    watch(
        () => product.value.qtde,
        (val) => {
            console.log(val);
            
            product.value.qtde = isNaN(Number(val.toString().replace(',', '.'))) ? 0 : val;
        }
    );

    watch(
        () => product.value.commission,
        (val) => {                    
            const commissionVal = isNaN(val) ? 0 : val;

            product.value.commission = commissionVal > 100 ? 100 : commissionVal;
        }
    );

    const submitProduct = async () => {
        try {
            await productSchema.value.validate(product.value, { abortEarly: false });

            const res = await createProduct(product.value);

            if(res.success)
            {
                notify(
                    'positive',
                    res.data.message

                );

                router.replace({
                    name: 'products.index'

                });

            } else {
                notify(
                    'negative',
                    res.message

                );
            };

        } catch (error: any) {
            console.error('Erro:', error.inner);

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
                    error.response?.data?.message || 'Erro na criação do produto!'
                );
            };
        };
    };
</script>
