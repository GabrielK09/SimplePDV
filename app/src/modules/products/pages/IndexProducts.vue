<template>
    <q-page padding>
        <div class="m-2 text-3xl" >
            <div class="flex justify-between">
                <h2 class="text-gray-600 m-2">Produtos</h2>

                <div class="mt-auto mb-auto">
                    <q-btn
                        no-caps
                        color="blue"
                        to="/admin/products/create"
                        label="Cadastrar novo produto"

                    />
                </div>
            </div>

            <div>
                <q-table
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
                        <q-tr
                            :props="props"
                        >
                            <q-td
                                v-for="col in props.cols"
                            >
                                <template v-if="col.name === 'actions'">
                                    <div
                                        class="text-center flex flex-center"
                                    >
                                        <div>
                                            <q-btn 
                                                size="10px" 
                                                no-caps 
                                                color="black" 
                                                icon="edit" 
                                                flat 
    
                                            />
                                        </div>

                                        <div>
                                            <q-btn 
                                                size="10px" 
                                                no-caps 
                                                color="red" 
                                                icon="delete" 
                                                flat 
                                                @click="showDialogDeleteProduct(props.row.id)"
                                            />
                                        </div>
                                    </div>
                                </template>

                                <template v-else>
                                    <div
                                        class="text-center"
                                        :title="props.row.active !== 1 ? 'Produto desativado!' : ''"
                                    >
                                        {{ col.value }}

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
</template>

<script setup lang="ts">
    import { QTableColumn, useQuasar } from 'quasar';
    import { onMounted, ref } from 'vue';
    import camelcaseKeys from 'camelcase-keys';
    import { getAll, deleteProduct } from '../services/productsService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';

    const $q = useQuasar();
    const { notify } = useNotify();

    const columns: QTableColumn[] = [
        {
            name: 'id',
            label: 'ID',
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
            label: 'Preço',
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
            label: 'Ações',
            field: 'actions',
            align: 'center'
        }
    ];

    let allProducts = ref<ProductContract[]>([]);
    let products = ref<ProductContract[]>([]);

    let searchInput = ref<string>('');

    const getAllProducts = async () => {
        const res = await getAll();
        const data = res.data;

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
        };

        products.value = data;
        allProducts.value = [...products.value];

    };

    const showDialogDeleteProduct = (productId: number) => {
        $q.dialog({
            title: 'Excluir produto',
            message: `Deseja realmente remover esse produto (${productId})?`,
            cancel: {
                push: true,
                label: 'Não',
                color: 'red',
            },

            ok: {
                push: true,
                label: 'Sim',
                color: 'green',
            },

        }).onOk(() => {
            deleteProductByDialog(productId);

        }).onCancel(() => {
            return;
        });
    };

    const deleteProductByDialog = async (productId: number) => {
        const res = await deleteProduct(productId);

        if(res.success)
        {
            notify(
                'positive',
                res.message
            );

        } else {
            notify(
                'positive',
                res.message
            );
        };

        getAllProducts();
    };

    const filterProducts = () => {
        console.log(searchInput.value);

        products.value = allProducts.value.filter(product => product.name.toLowerCase().includes(searchInput.value));
        console.log(allProducts.value);

    };

    onMounted(() => {
        getAllProducts();
    });
</script>
