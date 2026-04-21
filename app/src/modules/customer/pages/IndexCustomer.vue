<template>
    <q-page padding>
        <div class="m-2 text-3xl" >
            <div class="flex justify-between">
                <h2 class="text-gray-600 m-2">Clientes</h2>

                <div class="mt-auto mb-auto">
                    <q-btn
                        no-caps
                        color="blue"
                        @click="manageCustomerModal.create.show = !manageCustomerModal.create.show"
                        class="max-phone:mb-5"
                        label="Cadastrar novo cliente"
                    />
                </div>
            </div>

            <div>
                <q-table
                    v-model:pagination="pagination"
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
                                                    @click="buildShowCustomerUpdate(props.row.id)"
                                                    :disable="props.row.id === 1"
                                                >
                                                    <q-item-section avatar>
                                                        <q-icon name="edit" color="black" size="20px" />
                                                    </q-item-section>
                                                    <q-item-section>
                                                        <q-item-label>Editar</q-item-label>
                                                    </q-item-section>
                                                </q-item>
                                                
                                                <q-item 
                                                    clickable 
                                                    v-close-popup                                          
                                                    @click="showDialogDeleteProduct(props.row.id)"
                                                    :disable="props.row.id === 1"
                                                >
                                                    <q-item-section avatar>
                                                        <q-icon name="delete" color="red" size="20px" />
                                                    </q-item-section>
                                                    <q-item-section>
                                                        <q-item-label>Deletar</q-item-label>
                                                    </q-item-section>
                                                </q-item>
                                            </q-list>
                                        </q-menu>
                                    </q-btn>
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

    <UpdateCustomer
        v-if="manageCustomerModal.update.show"
        :customer-id="manageCustomerModal.update.customerId"
        @close="manageCustomerModal.update.show = !$event"
    />

    <CreateCustomer
        v-if="manageCustomerModal.create.show"
        @close="manageCustomerModal.create.show = !$event"
    />
</template>

<script setup lang="ts">
    import { QTableColumn, useQuasar } from 'quasar';
    import { onMounted, ref } from 'vue';
    import { deleteCustomer, getAll } from '../services/customerService';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import UpdateCustomer from './update/UpdateCustomer.vue';    
    import CreateCustomer from './create/CreateCustomer.vue';
    
    type UpdateCustomer = {
        show: boolean;
        customerId: number|null;
    };

    type CreateCustomer = {
        show: boolean;
    };

    type ManageCustomer = {
        update: UpdateCustomer;
        create: CreateCustomer
    };

    const pagination = ref({
        sortBy: 'id',
        rowsPerPage: 20
    });

    const manageCustomerModal = ref<ManageCustomer>({
        create: {
            show: false,
        },
        update: {
            show: false,
            customerId: null
        }
    });

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
            label: '',
            field: 'actions',
            align: 'right'
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

    const buildShowCustomerUpdate = (customerId: number) => { 
        console.log('Cliente da edição: ', customerId);
        
        manageCustomerModal.value.update = {
            show: true,
            customerId: customerId
        };        
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
        customers.value = allCustomers.value.filter(product => product.name.toLowerCase().includes(searchInput.value));

    };

    onMounted(() => {
        getAllCustomers();
    });
</script>
