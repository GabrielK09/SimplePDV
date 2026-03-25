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
                            <div class="p-4">
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

                                    <template v-slot:body="props">
                                        <q-tr :props="props">
                                            <q-td v-for="col in props.cols">
                                                <div class="flex flex-center">
                                                    <template v-if="col.name === 'actions'">
                                                        <q-btn 
                                                            color="primary" 
                                                            icon="edit"
                                                            @click="buildUpdateGrid(props.row)"
                                                        />
                                                    </template>

                                                    <template v-else>
                                                        {{ col.value }}
                                                    </template>
                                                </div>
                                            </q-td>
                                        </q-tr>
                                    </template>
                                </q-table>
                            </div>
                        </div>

                        <div class="flex flex-center">
                            <q-btn
                                color="primary"
                                type="submit"
                                label="Salvar dados do produto"
                                no-caps
                                :loading="loadingLogin"
                            />
                        </div>
                    </div>
                </q-form>
            </section>
        </main>

        <UpdateGridProduct
            v-if="showUpdateGrid"
            :grid-id="selectedGridId"
            :product-id="Number(route.params.id)"
            :grid-full-object="gridFullObject"
            :selected-sizes="product.productWithCharacteristics"
            @return:grids="handleUpdateGrid($event)"
            @close="showUpdateGrid = !$event"
        />  

        <CreateGridProduct
            v-if="showCreateGrid"
            :selected-sizes="product.productWithCharacteristics.map(c => (c.size))"
            @return:grids="getReturnedGrid($event)"
            @close="showCreateGrid = !$event"
        />  
    </q-page>
</template>

<script setup lang="ts">
    import { computed, onMounted, ref, watch } from 'vue';
    import { useRoute, useRouter } from 'vue-router';
    import * as Yup from 'yup';
    import { createProductCharacteristics, findById, updateProduct } from '../../services/productsService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { QTableColumn } from 'quasar';
    import UpdateGridProduct from 'src/components/Products/UseGrid/Update/UpdateGridProduct.vue';
    import CreateGridProduct from 'src/components/Products/UseGrid/Create/CreateGridProduct.vue';

    function formatGridDataForPush(gridData: ProductCharacteristicsContract): ProductCharacteristicsContract {
        return {
            grid_qtde: Number(gridData.grid_qtde) || 0,
            id: gridData.id ?? null,
            product_id: gridData.product_id ?? null,
            size: gridData.size
        };
    };

    const gridTableColumn: QTableColumn[] = [
        {
            name: 'id',
            label: 'ID',
            field: 'id',
            align: 'center'
        },
        {
            name: 'grid_qtde',
            label: 'Qtde',
            field: 'grid_qtde',
            align: 'center'
        },
        {
            name: 'size',
            label: 'Tamanho',
            field: 'size',
            align: 'center'
        },
        {
            name: 'actions',
            label: 'Ações',
            field: 'actions',
            align: 'center'
        }
    ];

    const productSchema = computed(() =>
        Yup.object({
            name: Yup
                .string()
                .required('O nome do produto é obrigatório!'),

            price: Yup
                .number()
                .required('O valor do produto é obrigatório!'),

            qtde: Yup
                .number()
                .required('A quantia do produto é obrigatório!'),

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
        commission: 0,
        use_grid: false,
        productWithCharacteristics: []
    });

    const formErrors = ref<Record<string, string>>({});

    const showUpdateGrid = ref<boolean>(false);
    const showCreateGrid = ref<boolean>(false);

    const selectedGridId = ref<number | null>(); // Usado quando for feito a edição de uma grade já cadsatrada.
    const gridFullObject = ref<any>(); // Usado quando for cadastrado uma nova grade.

    const router = useRouter();
    const route = useRoute();

    const productId = Number(route.params.id);
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
            await productSchema.value.validate(product.value, { abortEarly: false });

            const res = await updateProduct(product.value);

            if(res.success)
            {
                notify(
                    'positive',
                    res.message

                );

                if(product.value.use_grid && productId > 0)
                {
                    const newProductCharacteristics = product.value.productWithCharacteristics.map(c => ({
                        id: c.id,
                        product_id: productId,
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
                    error.response?.data?.message || 'Erro na edição do produto!'
                );
            };
        };
    };

    const handleUpdateGrid = (newGrid: ProductCharacteristicsContract) => {
        const oldGrid = product.value.productWithCharacteristics.find(c => c.id === newGrid.id);

        if (!oldGrid)
        {
            notify(
                'negative',
                'Ocorreu um erro ao alterar a grade.'
            );
            return;
        };

        const index = product.value.productWithCharacteristics.indexOf(oldGrid);

        if (index > -1)
        {
            product.value.productWithCharacteristics.splice(index, 1);
            product.value.productWithCharacteristics.push(formatGridDataForPush(newGrid));
            
        } else {
            return;  
        };
    };

    const buildUpdateGrid = (row: any) => {        
        showUpdateGrid.value = true;
        if(!row?.id) {
            console.warn('Grid não encontrada:', row?.id);

            const newGrid = {
                id: null,
                product_id: productId,
                grid_qtde: row.grid_qtde,
                size: row.size,
                have_register: false
            };

            gridFullObject.value = newGrid;

        } else {
            selectedGridId.value = row?.id;
        };
    };

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

    const getReturnedGrid = (grid: ProductCharacteristicsContract) => {            
        product.value.productWithCharacteristics.push(formatGridDataForPush(grid));

    };

    onMounted(async() => {
        if(!productId) return;

        const res = await findById(productId);

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
            return;
        };

        const productData: ProductContract = res.data.product;
        const productCharacteristicsData: ProductCharacteristicsContract[] = res.data.characteristics;

        product.value = {
            id: productData.id,
            name: productData.name,
            price: productData.price,
            qtde: productData.qtde,
            commission: productData.commission,
            use_grid: productData.use_grid,
            productWithCharacteristics: productCharacteristicsData
        };
    });
</script>
