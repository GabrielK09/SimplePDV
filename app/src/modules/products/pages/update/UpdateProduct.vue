<template>
    <q-page padding>
        <main class="min-h-[60vh] flex flex-center text-xl">
            <section class="w-[80vh] rounded-lg shadow px-4 bg-white">
                <header class="border-gray-100 flex">
                    <span class="text-black cursor-pointer">
                        <router-link to="/admin/products">
                            <q-avatar size="30px" icon="arrow_back" />

                        </router-link>
                    </span>
                    <h2 class="text-gray-600 text-center">Edição do produto</h2>

                </header>

                <q-form
                    @submit.prevent="submitProduct"
                    class="q-gutter-md mt-4"
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
                            v-model.number="product.price"
                            type="number"
                            label-slot
                            stack-label
                            outlined
                            dense
                            placeholder="0,00"
                            mask="##,##"
                            fill-mask="0"
                            reverse-fill-mask
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
                            v-model.number="product.qtde"
                            type="number"
                            label-slot
                            stack-label
                            outlined
                            dense
                            placeholder="0"
                            mask="#.###"
                            fill-mask="0"
                            reverse-fill-mask
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
                                v-model.number="product.commission"
                                type="number"
                                label-slot
                                stack-label
                                outlined
                                dense
                                placeholder="0,00"
                                mask="##,##"
                                fill-mask="0"
                                reverse-fill-mask
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
                                label="Editar produto"
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
    import { computed, onMounted, ref } from 'vue';
    import { useRoute, useRouter } from 'vue-router';
    import * as Yup from 'yup';
    import { findById, updateProduct } from '../../services/productsService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';

    const productSchema = computed(() =>
        Yup.object({
            name: Yup.string().required('O nome do produto é obrigatório!'),
            price: Yup.number().required('O valor do produto é obrigatório!'),
            qtde: Yup.number().required('A quantia do produto é obrigatório!'),
            commission: Yup
                            .number()
                            .min(0, 'O valor de comissão não pode ser menor que zero.')
                            .max(100, 'O valor de comissão não pode ser maior que 100%.')
                            .required('A quantia do produto é obrigatório!'),
        })
    );
    
    const product = ref<ProductContract>({
        id: 0,
        name: '',
        price: null,
        qtde: 0,
        commission: 0
    });

    const formErrors = ref<Record<string, string>>({});

    const router = useRouter();
    const route = useRoute();
    const { notify } = useNotify();

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

            const res = await updateProduct(product.value);

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
                    error.response?.data?.message || 'Erro na edição do produto!'
                );
            };
        };
    };

    onMounted(async() => {
        const productId = Number(route.params.id);
        if(!productId) return;
        
        const res = await findById(productId);

        if(!res.success) return;
        
        product.value = {
            id: res.data.id,
            name: res.data.name,
            price: res.data.price,
            qtde: res.data.qtde,
            commission: res.data.commission
        };
    });
</script>
