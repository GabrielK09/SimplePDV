<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card>
            <main class="flex flex-center text-xl">
                <section class="w-[80vh] rounded-lg shadow px-4 bg-white">
                    <header class="border-gray-100 flex">
                        <h2 class="text-gray-600 flex flex-1 justify-center">Cadastrar um novo produto</h2>

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
                                    :disable="product.use_grid && product.product_with_characteristics.length <= 0"
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

    <CreateGridProduct
        v-if="showCreateGrid"
        :selected-sizes="product.product_with_characteristics.map(c => (c.size))"
        @return:grids="getReturnedGrid($event)"
        @close="showCreateGrid = !$event"
    />  
</template>

<script setup lang="ts">
    import { computed, ref, watch } from 'vue'    
    import { createProduct, createProductCharacteristics } from '../../services/productsService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import CreateGridProduct from 'src/components/Products/UseGrid/Create/CreateGridProduct.vue';
    import QGridTable from 'src/components/Products/UseGrid/QTable/QGridTable.vue';
    import productSchema from '../../schema/productSchema';

    const emits = defineEmits<{
        (e: 'close', value: boolean): void
    }>();

    const priceInput = ref<string>('');

    const product = ref<ProductContract>({
        id: null,
        name: null,
        price: null,
        qtde: null,
        commission: 0,
        use_grid: false,
        product_with_characteristics: []
    });

    const formErrors = ref<Record<string, string>>({});

    const { notify } = useNotify();

    const loadingLogin = ref<boolean>(false);
    const showCreateGrid = ref<boolean>(false);
    const internalDialog = ref<boolean>(true);  

    const nameUpper = computed({
        get: () => product.value.name,
        set: (val: string) => {
            product.value.name = val.toUpperCase();
        }
    });

    watch(
        () => product.value.use_grid,
        (use) => {
            if(use) return product.value.qtde = null;
        }
    );

    const calculateQtde = computed(() => {
        if(!product.value.use_grid) return product.value.qtde;
        
        const list = product.value.product_with_characteristics;

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
            if(product.value.use_grid && product.value.product_with_characteristics.length <= 0)
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

            const validated = await productSchema().validate(formData, {
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
                product_with_characteristics: product.value.product_with_characteristics
            };

            const res = await createProduct(payLoad);

            const productId = res.data;

            if(res.success)
            {
                if(product.value.use_grid && productId > 0)
                {
                    const newProductCharacteristics = product.value.product_with_characteristics.map(c => ({
                        id: 0,
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

                closeCreate();
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
        product.value.product_with_characteristics.push({
            grid_qtde: Number(grid.grid_qtde) || 0,
            id: null,
            product_id: null,
            size: grid.size
        });
    };

    const closeCreate = () => {
        emits('close', true);
        internalDialog.value = false;
    };
</script>
