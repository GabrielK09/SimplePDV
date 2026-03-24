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

                                <q-checkbox 
                                    v-model="product.use_grid" label="Usa grade" 
                                    
                                />
                            </div>
                        </div>

                        <div v-if="!product.use_grid" class="mx-2 my-4">
                            <div class="border p-4">
                                <q-table
                                    title="Grade"
                                    :rows="product.productWithCharacteristics"
                                    hide-bottom
                                    :columns="gridTableColumn"
                                    row-key="name"
                                >
                                    <template v-slot:top-right>
                                        <q-btn 
                                            color="primary" 
                                            no-caps
                                            label="Cadastar uma grade" 
                                            @click="showCreateGrid = !showCreateGrid" 
                                        />
                                    </template>
                                </q-table>
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
        <pre>{{ product }}</pre>
    </q-page>

    <CreateGridProduct
        v-if="showCreateGrid"
        :registred-grids="[]"
        @return:grids="getReturnedGrid($event)"
        @close="showCreateGrid = !$event"
    />  
</template>

<script setup lang="ts">
    import { computed, ref } from 'vue';
    import { useRouter } from 'vue-router';
    import * as Yup from 'yup';
    import { createProduct, createProductCharacteristics } from '../../services/productsService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { QTableColumn } from 'quasar';
    import CreateGridProduct from 'src/components/Products/UseGrid/CreateGridProduct.vue';

    const gridTableColumn: QTableColumn[] = [
        {
            name: 'id',
            label: 'ID',
            field: 'id'
        },
        {
            name: 'grid_qtde',
            label: 'Qtde',
            field: 'grid_qtde',
        },
        {
            name: 'size',
            label: 'Tamanho',
            field: 'size',
        },
        {
            name: 'action',
            label: '',
            field: 'action'
        }
    ];

    const priceInput = ref<string>('');

    const product = ref<ProductContract>({
        id: null,
        name: null,
        price: null,
        qtde: null,
        commission: 0,
        use_grid: false,
        productWithCharacteristics: []
    });

    const formErrors = ref<Record<string, string>>({});

    const router = useRouter();
    const { notify } = useNotify();

    const loadingLogin = ref<boolean>(false);
    const showCreateGrid = ref<boolean>(false);

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
                id: null,
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

                if(product.value.use_grid)
                {
                };
                
                /*router.replace({
                    name: 'products.index'

                });*/

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

    const getReturnedGrid = (grid: ProductCharacteristicsContract) => {
        console.log('getReturnedGrid call', grid);
        console.log('product.value.productWithCharacteristics', product.value.productWithCharacteristics);
        
        product.value.productWithCharacteristics.push({
            grid_qtde: grid.grid_qtde,
            id: null,
            product_id: null,
            size: grid.size
        });
    };
</script>
