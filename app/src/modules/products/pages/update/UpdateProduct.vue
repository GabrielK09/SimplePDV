<template>
    <q-dialog v-model="internalDialog" persistent>
        <q-card>
            <main class="min-h-[60vh] flex flex-center text-xl">
                <section class="w-[80vh] rounded-lg shadow px-4 bg-white">
                    <header class="border-gray-100 flex">
                        <h2 class="text-gray-600 flex flex-1 justify-center">Edição do produto</h2>
                        
                        <q-btn
                            color="red" 
                            icon="close"
                            class="w-12 h-12 my-auto ml-auto"
                            @click="closeUpdate"
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
                                        :flat="product.commission === 15"
                                    />

                                    <q-btn
                                        color="primary"
                                        no-caps
                                        @click="product.commission = 25"
                                        label="25%"
                                        :flat="product.commission === 25"
                                    />

                                    <q-btn
                                        color="primary"
                                        no-caps
                                        @click="product.commission = 35"
                                        label="35%"
                                        :flat="product.commission === 35"
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
                                    :product-id="productId"
                                />
                            </div>

                            <div class="flex flex-center">
                                <q-btn
                                    color="primary"
                                    type="submit"
                                    label="Salvar dados do produto"
                                    no-caps
                                    :disable="product.use_grid && product.product_with_characteristics.length <= 0"
                                    :loading="loadingLogin"
                                />
                            </div>
                        </div>
                    </q-form>
                </section>
            </main>         

            <CreateGridProduct
                v-if="showCreateGrid"
                :selected-sizes="product.product_with_characteristics.map(c => (c.size))"
                @return:grids="product.product_with_characteristics.push(formatGridDataForPush($event))"
                @close="showCreateGrid = !$event"
            />  
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import { computed, onMounted, ref, watch } from 'vue';
    import { createProductCharacteristics, getProductCharacteristicsById, findById, updateProduct } from '../../services/productsService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import QGridTable from 'src/components/Products/UseGrid/QTable/QGridTable.vue';
    import CreateGridProduct from 'src/components/Products/UseGrid/Create/CreateGridProduct.vue';
    import formatGridDataForPush from 'src/helpers/FormatValue/Grid/formatGridDataForPush';
    import productSchema from '../../schema/productSchema';

    const props = defineProps<{
        productId: number;
    }>();

    const emits = defineEmits<{
        (e: 'close', value: boolean): void
    }>();

    const product = ref<ProductContract>({
        id: props.productId,
        name: null,
        price: null,
        qtde: null,
        commission: 0,
        use_grid: false,
        product_with_characteristics: []
    });

    const formErrors = ref<Record<string, string>>({});
    const showCreateGrid = ref<boolean>(false);

    const { notify } = useNotify();

    const loadingLogin = ref<boolean>(false);
    const internalDialog = ref<boolean>(true);

    const nameUpper = computed({
        get: () => product.value.name,

        set: (val: string) => {
            product.value.name = val.toUpperCase();
        }
    });

    const submitProduct = async () => {
        try {
            await productSchema().validate(product.value, { abortEarly: false });
            
            const res = await updateProduct(product.value);

            if(res.success)
            {
                notify(
                    'positive',
                    res.message
                );

                if(product.value.use_grid && props.productId > 0)
                {
                    const newProductCharacteristics = product.value.product_with_characteristics.map(c => ({
                        id: c.id,
                        product_id: props.productId,
                        grid_qtde: c.grid_qtde,
                        size: c.size
                    }));

                    const resCharacteristics = await createProductCharacteristics(newProductCharacteristics, true);
                    
                    if(!resCharacteristics.success)
                    {
                        notify(
                            'negative',
                            resCharacteristics.data.message
                        );  
                    };
                };

                closeUpdate();
    
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
                    error.response?.data?.message || 'Erro na edição do produto!'
                );
            };
        };
    };

    watch(
        () => product.value.use_grid,
        async (use) => {
            if(use)  
            {
                product.value.qtde = null;

                const res = await getProductCharacteristicsById(props.productId);

                if(!res.success)
                {
                    product.value.use_grid = false;
                    notify(
                        'negative',
                        res.message
                    );

                    return;
                };
                                
                product.value.product_with_characteristics = res.data || [];
            };
        }
    );

    const calculateQtde = computed((): number => {
        if(!product.value.use_grid) return product.value.qtde;

        const list = product.value.product_with_characteristics;
        
        if (!list || list.length === 0) return 0;

        return list.reduce((total, a) => total + (a.grid_qtde || 0), 0);
    });

    watch(calculateQtde, (value) => {        
        if (value !== null) {
            product.value.qtde = value;
        };
    });

    onMounted(async() => {
        if(!props.productId) emits('close', true);

        const res = await findById(props.productId);

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
            return;
        };

        const productData: ProductContract = res.data;

        product.value = {
            id: productData.id,
            name: productData.name,
            price: productData.price,
            qtde: productData.qtde,
            commission: productData.commission,
            use_grid: productData.use_grid,
            product_with_characteristics: productData.product_with_characteristics || []
        };
    });

    const closeUpdate = () => {
        emits('close', true);
        internalDialog.value = false;
    };
</script>
