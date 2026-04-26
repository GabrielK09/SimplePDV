<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card class="flex flex-center text-xl overflow-hidden">
            <q-card-section>
                <section class="w-[60vh] rounded-lg shadow px-4 bg-white">
                    <header class="flex justify-center">
                        <h2 class="text-gray-600 ">Cadastrar um novo produto</h2>

                         <q-btn
                            color="red"
                            icon="close"
                            class="w-2 h-2 my-auto ml-4"
                            @click="emits('close', true)"
                        />
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
            </q-card-section>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import { computed, ref } from 'vue';
    import * as Yup from 'yup';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { createProduct } from 'src/modules/products/services/productsService';
import productSchema from 'src/modules/products/schema/productSchema';

    const emits = defineEmits<{
        (e: 'close', value: boolean): void
    }>();

    const internalDialog = ref<boolean>(true);
    const priceInput = ref<string>('');

    const product = ref<ProductContract>({
        id: null,
        name: '',
        price: null,
        qtde: null,
        commission: 0
    });

    const formErrors = ref<Record<string, string>>({});

    const { notify } = useNotify();

    const loadingLogin = ref<boolean>(false);

    const nameUpper = computed({
        get: () => product.value.name,
        set: (val: string) => {
            product.value.name = val.toUpperCase();
        }
    });

    const submitProduct = async () => {
        try {
            const formData = {
                id: product.value.id,
                name: product.value.name,
                price: priceInput.value,
                qtde: product.value.qtde,
                commission: product.value.commission
            };

            const validated = await productSchema().validate(formData, {
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
                    res.message

                );

                internalDialog.value = false;

                return;
            };

            notify(
                'negative',
                res.message

            );

            emits('close', true);

        } catch (error: any) {
            if(error.inner)
            {
                console.error('Erro:', error?.inner);

                formErrors.value = {};

                error.inner.forEach((err: any) => {
                    formErrors.value[err.path] = err.message;

                    notify(
                        'negative',
                        err.message

                    );
                });
            } else {
                console.error('Erro:', error);
                notify(
                    'negative',
                    error.response?.data?.message || 'Erro na criação do produto!'
                );
            };
        };
    };
</script>
