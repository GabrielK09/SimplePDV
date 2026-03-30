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
                            :disable="product.use_grid"
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

                        <div v-if="product.use_grid" class="mx-2 my-4">
                            <QGridTable
                                :product-data="product"
                                @show-create-grid="showCreateGrid = $event"
                            />
                        </div>

                        <div class="flex flex-center">
                            <q-btn
                                color="primary"
                                type="submit"
                                label="Cadastrar produto"
                                :disable="product.use_grid && product.productWithCharacteristics.length <= 0"
                                no-caps
                                :loading="loadingLogin"
                            />
                        </div>
                    </div>
                </q-form>
            </section>
        </main>
    </q-page>

    <CreateGridProduct
        v-if="showCreateGrid"
        :selected-sizes="product.productWithCharacteristics.map(c => (c.size))"
        @return:grids="getReturnedGrid($event)"
        @close="showCreateGrid = !$event"
    />  
</template>

<script setup lang="ts">
    import { computed, ref, watch } from 'vue';
    import { useRouter } from 'vue-router';
    import * as Yup from 'yup';
    import { createProduct, createProductCharacteristics } from '../../services/productsService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import CreateGridProduct from 'src/components/Products/UseGrid/Create/CreateGridProduct.vue';
    import QGridTable from 'src/components/Products/UseGrid/QTable/QGridTable.vue';

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

    watch(
        () => product.value.use_grid,
        (use) => {
            if(use) return product.value.qtde = null;
        }
    );

    const calculateQtde = computed(() => {
        if(!product.value.use_grid) return product.value.qtde;
        
        const list = product.value.productWithCharacteristics;

        if (!list || list.length === 0) return null;

        return list.reduce((total, a) => total + (a.grid_qtde || 0), 0);
    });

    watch(calculateQtde, (value) => {
        if (value !== null) {
            product.value.qtde = value;
        };
    });

    const submitProduct = async () => {
        loadingLogin.value = true;

        try {
            if(product.value.use_grid && product.value.productWithCharacteristics.length <= 0)
            {
                notify(
                    'info',
                    'Caso seja usado grade, ao menos uma grade precisa ser cadastrada!'
                );
                return;
            };

            const formData = {
                id: product.value.id,
                name: product.value.name,
                price: priceInput.value,
                qtde: product.value.qtde,
                commission: product.value.commission,
                use_grid: product.value.use_grid
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
                qtde: Number(validated.qtde),
                use_grid: product.value.use_grid,
                productWithCharacteristics: product.value.productWithCharacteristics
            };

            const res = await createProduct(payLoad);

            const productId = res.data;

            if(res.success)
            {
                if(product.value.use_grid && productId > 0)
                {
                    const newProductCharacteristics = product.value.productWithCharacteristics.map(c => ({
                        id: null,
                        product_id: productId,
                        grid_qtde: c.grid_qtde,
                        size: c.size
                    }));

                    const resCharacteristics = await createProductCharacteristics(newProductCharacteristics);
                    
                    if(!resCharacteristics.success)
                    {
                        notify(
                            'negative',
                            resCharacteristics.message
                        );  
                    };
                };
                
                notify(
                    'positive',
                    res.message
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
        } finally {
            loadingLogin.value = false;
        };
    };

    const getReturnedGrid = (grid: ProductCharacteristicsContract) => {
        product.value.productWithCharacteristics.push({
            grid_qtde: Number(grid.grid_qtde) || 0,
            id: null,
            product_id: null,
            size: grid.size
        });
    };
</script>
