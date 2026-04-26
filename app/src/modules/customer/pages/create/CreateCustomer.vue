<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card>
            <main class="flex flex-center text-xl">
                <section class="w-[80vh] rounded-lg shadow px-4 bg-white">
                    <header class="border-gray-100 flex">
                        <h2 class="text-gray-600 text-center">Cadastrar um(a) novo(a) cliente</h2>

                        <q-btn
                            color="red" 
                            icon="close"
                            class="w-12 h-12 my-auto ml-auto"
                            @click="closeCreate"
                        />
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
                                v-model="customer.cpf_cnpj"
                                type="text"
                                label-slot
                                stack-label
                                outlined
                                dense
                                @update:model-value="cpfCnpjMask(customer.cpf_cnpj)"
                                :rules="[
                                    val => {
                                        return !val || validateCPF(val) || 'CPF inválido' 
                                    }
                                ]"
                                maxlength="18"
                                class="mb-4"
                            >
                                <template v-slot:label>
                                    <div class="text-sm">
                                        CPF/CNPJ<span class="text-red-500">*</span>
                                    </div>
                                </template>
                            </q-input>

                            <div class="flex flex-center">
                                <q-btn
                                    color="primary"
                                    type="submit"
                                    label="Cadastrar cliente"
                                    no-caps
                                    :loading="loadingLogin"
                                />
                            </div>
                        </div>
                    </q-form>
                </section>
            </main>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import { computed, ref } from 'vue';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { createCustomer } from '../../services/customerService';
    import { getCnpjDataService } from 'src/services/CNPJ/getCnpjData';
    import customerSchema from '../../schema/customerSchema';
    import validateCPF from 'src/helpers/CPF_CNPJ/validateCPF';
    
    const emits = defineEmits<{
        (e: 'close', value: boolean): void
    }>();

    const customer = ref<CustomerContract>({
        id: 0,
        name: '',
        cpf_cnpj: ''
    });

    const formErrors = ref<Record<string, string>>({});

    const { notify } = useNotify();

    const internalDialog = ref<boolean>(true);  
    const loadingLogin = ref<boolean>(false);

    const nameUpper = computed({
        get: () => customer.value.name,
        set: (val: string) => {
            customer.value.name = val.toUpperCase();
        }
    });

    const submitProduct = async () => {
        try {
            await customerSchema().validate(customer.value, { abortEarly: false });

            const res = await createCustomer(customer.value);

            if(res.success)
            {
                notify(
                    'positive',
                    'Cliente cadastrado com sucesso!'

                );

                closeCreate();

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
                    error.response?.data?.message || 'Erro na criação do cliente!'
                );
            };
        };
    };

    const cpfCnpjMask = (val: string) => {
        const cnpjCpf = val.replace(/\D/g, ''); 

        if (cnpjCpf.length === 11) {
            customer.value.cpf_cnpj = cnpjCpf.replace(
                /(\d{3})(\d{3})(\d{3})(\d{2})/g,
                '$1.$2.$3-$4'
            );

            return customer.value.cpf_cnpj;
        }

        if (cnpjCpf.length === 14) {
            customer.value.cpf_cnpj = cnpjCpf.replace(
                /(\d{2})(\d{3})(\d{3})(\d{4})(\d{2})/g,
                '$1.$2.$3/$4-$5'
            );

            getCnpjData(cnpjCpf);

            return customer.value.cpf_cnpj;
        }
    };

    const getCnpjData = async (newCnpj: string) => {
        const res = await getCnpjDataService(newCnpj.replace(/\D/, ''));

        if(!res.success) return;

        customer.value.name = res.data.alias;
    
    };

    const closeCreate = () => {
        emits('close', true);
        internalDialog.value = false;
    };
</script>