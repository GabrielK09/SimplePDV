<template>
    <q-page padding>
        <div class="m-2 text-3xl" >
            <div class="flex justify-between">
                <h2 class="text-gray-600 m-2">Produtos</h2>

                <div class="mt-auto mb-auto">
                    <q-btn
                        no-caps
                        color="blue"
                        @click="manageProductModal.create.show = !manageProductModal.create.show"
                        class="max-phone:mb-5"
                        label="Cadastrar novo produto"

                    />
                </div>
            </div>

            <div>
                <q-table
                    v-model:pagination="pagination"
                    borded
                    :rows="products"
                    :columns="columns"
                    row-key="name"
                    class="rounded-xl"
                >
                    <template v-slot:top-right>
                        <q-input
                            outlined
                            v-model="searchInput"
                            type="text"
                            dense
                            label=""
                            @update:model-value="filterProducts"
                        >
                            <template v-slot:append>
                                <q-icon name="search" />
                            </template>
                            <template v-slot:label>
                                <span class="text-xs">Buscar por um produto ...</span>
                            </template>
                        </q-input>
                    </template>

                    <template v-slot:body="props">
                        <q-tr :props="props">
                            <q-td v-for="col in props.cols">
                                <template v-if="col.name === 'actions'">
                                    <q-btn 
                                        dense
                                        flat
                                        icon="more_vert"
                                    >
                                        <q-menu
                                            anchor="bottom right"
                                            self="top right"
                                            class="rounded shadow-xl bg-white"
                                            transition-show="jump-down"
                                        >
                                            <q-list style="min-width: 90px">
                                                <q-item 
                                                    clickable 
                                                    v-close-popup  
                                                    v-if="props.row.deleted_at === null"
                                                    @click="buildUpdateProduct(props.row.id)"
                                                >
                                                    <q-item-section avatar>
                                                        <q-icon name="edit" color="primary" size="20px" />
                                                    </q-item-section>
                                                    <q-item-section>
                                                        <q-item-label>Editar</q-item-label>
                                                    </q-item-section>
                                                </q-item>

                                                <q-item
                                                    clickable 
                                                    v-close-popup  
                                                    v-if="props.row.deleted_at === null"
                                                    @click="showDialogActionProduct(props.row.id, 'delete')"
                                                >
                                                    <q-item-section avatar>
                                                        <q-icon name="delete" color="red" size="20px" />
                                                    </q-item-section>
                                                    <q-item-section>
                                                        <q-item-label>Deletar produto</q-item-label>
                                                    </q-item-section>
                                                </q-item>

                                                <q-item
                                                    clickable 
                                                    v-close-popup  
                                                    v-if="props.row.deleted_at !== null"
                                                    @click="showDialogActionProduct(props.row.id, 'active')"
                                                >
                                                    <q-item-section avatar>
                                                        <q-icon name="rotate_left" color="green" size="20px" />
                                                    </q-item-section>
                                                    <q-item-section>
                                                        <q-item-label>Ativar produto</q-item-label>
                                                    </q-item-section>
                                                </q-item>
                                            </q-list>
                                        </q-menu>
                                    </q-btn>
                                </template>

                                <template v-else>
                                    <div class="text-center">
                                        <span v-if="props.row.deleted_at !== null" class="text-gray-400">
                                            {{ col.value }}

                                        </span>
                                        <div v-else>
                                            {{ col.value }}

                                        </div>
                                    </div>
                                </template>
                            </q-td>
                        </q-tr>
                    </template>

                    <template v-slot:no-data>
                        <div class="ml-auto mr-auto">
                            <q-icon name="warning" size="30px"/>
                            <span class="mt-auto mb-auto ml-2 text-xs">Sem produtos cadastrados</span>

                        </div>
                    </template>
                </q-table>
            </div>
        </div>
    </q-page>

    <UpdateProduct
        v-if="manageProductModal.update.show"
        :product-id="manageProductModal.update.productId"
        @close="closeManageProductModal(!$event)"
    />

    <CreateProduct
        v-if="manageProductModal.create.show"
        @close="closeManageProductModal(!$event)"
    />
</template>

<script setup lang="ts">
    import { QTableColumn, useQuasar } from 'quasar';
    import { onMounted, ref, watch } from 'vue';
    import { getAll, manageProductService } from '../services/productsService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import UpdateProduct from './update/UpdateProduct.vue';
    import CreateProduct from './create/CreateProduct.vue';

    type UpdateProduct = {
        show: boolean;
        productId: number|null;
    };

    type CreateProduct = {
        show: boolean;
    };

    type ManageProduct = {
        update: UpdateProduct;
        create: CreateProduct
    };

    const $q = useQuasar();
    const { notify } = useNotify();

    const pagination = ref({
        sortBy: 'id',
        rowsPerPage: 20
    });

    const manageProductModal = ref<ManageProduct>({
        create: {
            show: false,
        },
        update: {
            show: false,
            productId: null
        }
    });

    const columns: QTableColumn[] = [
        {
            name: 'id',
            label: 'Cód produto',
            field: 'id',
            align: 'center'
        },
        {
            name: 'name',
            label: 'Produto',
            field: 'name',
            align: 'center'
        },
        {
            name: 'price',
            label: 'Preço de venda',
            field: 'price',
            align: 'center',
            format(val: number) {
                return `R$ ${val.toFixed(2).toString().replace('.', ',')}`
            }
        },
        {
            name: 'qtde',
            label: 'Qtde',
            field: 'qtde',
            align: 'center'
        },
        {
            name: 'actions',
            label: '',
            field: 'actions',
            align: 'right'
        }
    ];

    const allProducts = ref<ProductContract[]>([]);
    const products = ref<ProductContract[]>([]);

    const searchInput = ref<string>('');

    const showDialogActionProduct = (productId: number, operation: 'active'|'delete') => {
        $q.dialog({
            title: `${operation === 'delete' ? 'Excluir' : 'Ativar'} produto`,
            message: `Deseja realmente ${operation === 'delete' ? 'deletear' : 'ativar'} esse produto (${productId})?`,
            cancel: {
                push: true,
                label: 'Não',
                color: operation === 'delete' ? 'red' : 'green'
            },

            ok: {
                push: true,
                label: 'Sim',
                color: 'green',
            },

        }).onOk(() => {
            manageProduct(productId, operation);

        }).onCancel(() => {
            return;
        });
    };

    const filterProducts = (): void => {
        products.value = allProducts.value.filter(product => product.name ? product.name.toLowerCase().includes(searchInput.value) : null);
    };

    watch(
        () => pagination.value.rowsPerPage,
        async (newRowsPerPage) => {
            console.log('newRowsPerPage: ', newRowsPerPage);

            await getAllProducts(newRowsPerPage);
        }
    );

    const getAllProducts = async (rowsPerPage?: number) => {
        const res = await getAll(rowsPerPage);

        const data = res.data;

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );

            return;
        };
    
        products.value = data;
        allProducts.value = [...products.value];
    };

    const manageProduct = async (productId: number, operation: 'active'|'delete') => {
        const res = await manageProductService(productId, operation);

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
            return;
        };

        notify(
            'positive',
            res.message
        );

        await getAllProducts();
    };

    const buildUpdateProduct = (productId: number): void => {
        console.log(productId);
        
        manageProductModal.value.update = {
            show: true,
            productId: productId
        };
    };

    const closeManageProductModal = (_: boolean): void => {        
        manageProductModal.value = {
            create: {
                show: false
            },

            update: {
                productId: null,
                show: false
            }
        };

        getAllProducts();
    };

    onMounted(() => {
        getAllProducts();
    });
</script>
