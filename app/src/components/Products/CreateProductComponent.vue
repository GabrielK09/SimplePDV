<template>
    <q-dialog>
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
                            v-model="priceInput"
                            type="text"
                            label-slot
                            stack-label
                            outlined
                            dense
                            class="mb-4"
                            :error="!!formErrors.price"
                            :error-message="formErrors.price"
                        >
                            <template #label>
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
                                v-model.number="product.commission"
                                type="number"
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
    </q-dialog>
</template>

<script setup lang="ts">
    import { computed, ref } from 'vue';
    import { useRouter } from 'vue-router';
    import * as Yup from 'yup';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { createProduct } from 'src/modules/products/services/productsService';

    const emits = defineEmits<{
        (e: 'close', value: boolean)
    }>();

    const priceInput = ref<string>('');

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

    const loadingLogin = ref<boolean>(false);

    const nameUpper = computed({
        get: () => product.value.name,
        set: (val: string) => {
            product.value.name = val.toUpperCase();
        }
    });

    const parsePtBrNumber = (value: string | number | null | undefined): number => {
        if (value === null || value === undefined || value === "") return 0;

        if (typeof value === 'number') 
        {
            return Number.isFinite(value) ? value : 0;
        };

        const normalized = value
            .trim()
            .replace(/\./g, '')
            .replace(',', '.');

        const parsed = Number(normalized);
        return Number.isFinite(parsed) ? parsed : 0;
    };

    const productSchema = computed(() =>
        Yup.object({
            name: Yup.string()
                .trim()
                .required('O nome do produto é obrigatório!'),

            price: Yup.number()
                .transform((_, originalValue) => parsePtBrNumber(originalValue))
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

    const submitProduct = async () => {
        try {
            const formData = {
                id: product.value.id,
                name: product.value.name,
                price: priceInput.value,
                qtde: product.value.qtde,
                commission: product.value.commission
            };

            const validated = await productSchema.value.validate(formData, {
                abortEarly: false, 
                stripUnknown: true 
            });

            const payLoad: ProductContract = {
                id: 0,
                name: validated.name,
                price: Number(validated.price.toFixed(2)),
                commission: Number(validated.commission),
                qtde: Number(validated.qtde)
            };

            const res = await createProduct(payLoad);

            if(res.success)
            {
                notify(
                    'positive',
                    res.data.message

                );

                router.replace({
                    name: 'products.index'

                });

                return;
            };
            
            notify(
                'negative',
                res.message

            );
        } catch (error: any) {
            console.error('Erro:', error);
            console.error('Erro:', error?.inner);

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
