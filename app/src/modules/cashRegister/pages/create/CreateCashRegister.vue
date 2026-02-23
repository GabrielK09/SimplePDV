<template>
    <q-page padding>
        <main class="min-h-[60vh] flex flex-center text-xl">
            <section class="w-[80vh] rounded-lg shadow px-4 bg-white">
                <header class="border-gray-100 flex">
                    <span class="text-black cursor-pointer">
                        <router-link to="/admin/cash-register">
                            <q-avatar size="30px" icon="arrow_back" />

                        </router-link>
                    </span>
                    <h2 class="text-gray-600 text-center">Cadastrar um novo movimento finaceiro</h2>

                </header>

                <q-form
                    @submit.prevent="submitCashRegister"
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
                            :error="!!formErrors.description"
                            :error-message="formErrors.description"
                        >
                            <template v-slot:label>
                                <div class="text-sm">
                                    Nome <span class="text-red-500">*</span>
                                </div>
                            </template>
                        </q-input>

                        <q-input
                            v-model.number="cashRegisterData.input_value"
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
                            :error="!!formErrors.input_value"
                            :error-message="formErrors.input_value"
                            :disable="cashRegisterData.output_value > 0"
                        >
                            <template v-slot:label>
                                <div class="text-sm">
                                    Valor de entrada <span class="text-red-500">*</span>
                                </div>
                            </template>
                        </q-input>

                        <q-input
                            v-model.number="cashRegisterData.output_value"
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
                            :error="!!formErrors.output_value"
                            :error-message="formErrors.output_value"
                            :disable="cashRegisterData.input_value > 0"
                        >
                            <template v-slot:label>
                                <div class="text-sm">
                                    Valor de saída <span class="text-red-500">*</span>
                                </div>
                            </template>
                        </q-input>

                        <div class="flex flex-center">
                            <q-btn
                                color="primary"
                                type="submit"
                                label="Cadastrar movimento"
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
    import { computed, ref } from 'vue';
    import { useRouter } from 'vue-router';
    import * as Yup from 'yup';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { createManualCashRegister } from '../../services/cashRegisterService';

    const cashRegisterSchema = computed(() =>
        Yup.object({
            name: Yup.string().required('O nome do produto é obrigatório!'),
            price: Yup.number().required('O valor do produto é obrigatório!'),
            qtde: Yup.number().required('A quantia do produto é obrigatório!'),
        })
    );

    const cashRegisterData = ref<CashRegisterContract>({
        id: 0,
        sale_id: 0,
        shopping_id: 0,
        customer: '',
        description: '',
        input_value: 0,
        output_value: 0,
        specie: '',
        specie_id: 0,
        total_balance: 0
    });

    const formErrors = ref<Record<string, string>>({});

    const router = useRouter();
    const { notify } = useNotify();

    let loadingLogin = ref<boolean>(false);

    const nameUpper = computed({
        get: () => cashRegisterData.value.description,
        set: (val: string) => {
            cashRegisterData.value.description = val.toUpperCase();
        }
    });

    const submitCashRegister = async () => {
        try {
            await cashRegisterSchema.value.validate(cashRegisterData.value, { abortEarly: false });

            const res = await createManualCashRegister(cashRegisterData.value);

            if(res.success)
            {
                notify(
                    'positive',
                    res.data.message

                );

                router.replace({
                    name: 'cash-register.index'

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
                    error.response?.data?.message || 'Erro na criação do registro do caixa.'
                );
            };
        };
    };
</script>
