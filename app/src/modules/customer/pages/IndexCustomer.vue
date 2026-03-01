<template>
    <q-page padding>
        <div class="m-2 text-3xl" >
            <div class="flex justify-between">
                <h2 class="text-gray-600 m-2">Clientes</h2>

                <div class="mt-auto mb-auto">
                    <q-btn
                        no-caps
                        color="blue"
                        to="/admin/customers/create"
                        label="Cadastrar novo cliente"

                    />
                </div>
            </div>

            <div>
                <q-table
                    borded
                    :rows="customers"
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
                            @update:model-value="filterCustomers"
                        >
                            <template v-slot:append>
                                <q-icon name="search" />
                            </template>
                            <template v-slot:label>
                                <span class="text-xs">Buscar por um cliente ...</span>
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
                                                :disable="props.row.id === 1"
                                                :to="`customers/edit/${props.row.id}`"
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
                                                :disable="props.row.id === 1"
                                            />
                                        </div>
                                    </div>
                                </template>

                                <template v-else>
                                    <div
                                        class="text-center"
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
                            <span class="mt-auto mb-auto ml-2 text-xs">Sem clientes cadastrados</span>

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
    import { deleteCustomer, getAll } from '../services/customerService';
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
            label: 'Cliente',
            field: 'name',
            align: 'center'
        },
        {
            name: 'cpf_cnpj',
            label: 'Documento',
            field: 'cpf_cnpj',
            align: 'center'
        },
        {
            name: 'actions',
            label: 'Ações',
            field: 'actions',
            align: 'center'
        }
    ];

    const allCustomers = ref<CustomerContract[]>([]);
    const customers = ref<CustomerContract[]>([]);

    const searchInput = ref<string>('');

    const getAllCustomers = async () => {
        const res = await getAll();
        const data = res.data;

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
        };

        customers.value = data;
        allCustomers.value = [...customers.value];

    };

    const showDialogDeleteProduct = (customerId: number) => {
        if(customerId === 1)
        {
            notify(
                'negative',
                'O cliente padrão não pode ser desativado.'
            )
            return;  
        };

        $q.dialog({
            title: 'Excluir cliente',
            message: `Deseja realmente remover esse cliente (${customerId})?`,
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
            deleteCustomerByDialog(customerId);

        }).onCancel(() => {
            return;
        });
    };

    const deleteCustomerByDialog = async (customerId: number) => {
        const res = await deleteCustomer(customerId);

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

        getAllCustomers();
    };

    const filterCustomers = () => {
        console.log(searchInput.value);

        customers.value = allCustomers.value.filter(product => product.name.toLowerCase().includes(searchInput.value));
        console.log(allCustomers.value);

    };

    onMounted(() => {
        getAllCustomers();
    });
</script>
