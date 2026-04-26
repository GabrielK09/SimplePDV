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
                        <div class="flex">
                            <div class="mr-4 select-status">
                                <q-select 
                                    v-model="byStatus" 
                                    :options="statusOptions" 
                                    option-label="Status"
                                    emit-value
                                    map-options
                                    outlined
                                    dense
                                    :display-value="selectedLabel"
                                    :clearable="true"
                                    @update:model-value="applyFilters"
                                />
                            </div>
                            
                            <div>
                                <q-input
                                    outlined
                                    v-model="searchInput"
                                    type="text"
                                    label=""
                                    dense
                                    @update:model-value="applyFilters"
                                >
                                    <template v-slot:append>
                                        <q-icon name="search" />
                                    </template>
                                    <template v-slot:label>
                                        <span class="text-xs">Buscar por um cliente ...</span>
                                    </template>
                                </q-input>
                            </div>
                        </div>
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
                                                    :disable="props.row.id === 1 || props.row.deleted_at !== null"
                                                    @click="buildShowCustomerUpdate(props.row.id)"
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
                                                    v-if="props.row.id === 1 || props.row.deleted_at === null"
                                                    :disable="props.row.id === 1"
                                                    @click="showDialogActionCustomer(props.row.id, 'delete')"
                                                >
                                                    <q-item-section avatar>
                                                        <q-icon name="delete" color="red" size="20px" />
                                                    </q-item-section>
                                                    <q-item-section>
                                                        <q-item-label>Deletar</q-item-label>
                                                    </q-item-section>
                                                </q-item>

                                                <q-item
                                                    clickable 
                                                    v-close-popup  
                                                    v-if="props.row.deleted_at !== null"
                                                    @click="showDialogActionCustomer(props.row.id, 'active')"
                                                >
                                                    <q-item-section avatar>
                                                        <q-icon name="rotate_left" color="green" size="20px" />
                                                    </q-item-section>
                                                    <q-item-section>
                                                        <q-item-label>Ativar</q-item-label>
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
        @close="closeManageProductModal($event)"
    />

    <CreateCustomer
        v-if="manageCustomerModal.create.show"
        @close="closeManageProductModal($event)"
    />
</template>

<script setup lang="ts">
    import { QTableColumn, useQuasar } from 'quasar';
    import { computed, onMounted, ref } from 'vue';
    import { getAll, manageCustomerService } from '../services/customerService';
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

    const statusOptions: Exclude<FilterByActiveOrDisable, null>[] = [
        'Ativos',
        'Inativos',
        'Todos'
    ];

    const byStatus = ref<FilterByActiveOrDisable>(null);

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
        const data = res.data as CustomerContract[];

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
        };

        allCustomers.value = data;
        applyFilters();
    };

    const applyFilters = () => {
        let filtred = [...allCustomers.value];

        if(byStatus.value) 
        {
            switch (byStatus.value) {
                case 'Ativos':
                    filtred = filtred.filter(c => c.deleted_at === null)
                    break;  

                case 'Inativos':
                    filtred = filtred.filter(c => c.deleted_at !== null)
                    break;  

                case 'Todos':
                    getAllCustomers();
                    break;  
            
                default:
                    getAllCustomers();
                    break;
            }
        };

        if(searchInput.value.trim())
        {
            const search = searchInput.value.trim().toLowerCase();

            filtred = filtred.filter(c => 
                String(c.name).includes(search) ||
                String(c.id).includes(search) ||
                String(c.cpf_cnpj === search)

            );
        };

        customers.value = filtred;
    };

    const selectedLabel = computed(() => {
        return byStatus.value ?? 'Todos';
    }); 

    const buildShowCustomerUpdate = (customerId: number) => { 
        manageCustomerModal.value.update = {
            show: true,
            customerId: customerId
        };        
    };

    const showDialogActionCustomer = (customerId: number, operation: 'active'|'delete') => {
        if(customerId === 1)
        {
            notify(
                'negative',
                'O cliente padrão não pode ser desativado.'
            )
            return;  
        };

        $q.dialog({
            title: `${operation === 'delete' ? 'Excluir' : 'Ativar'} cliente`,
            message: `Deseja realmente ${operation === 'delete' ? 'deletear' : 'ativar'} esse cliente (${customerId})?`,
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
            manageProduct(customerId, operation);

        }).onCancel(() => {
            return;
        });
    };

    const manageProduct = async (productId: number, operation: 'active'|'delete') => {
        const res = await manageCustomerService(productId, operation);

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

        await getAllCustomers();
    };

    const closeManageProductModal = (_: boolean): void => {        
        manageCustomerModal.value = {
            create: {
                show: false
            },

            update: {
                customerId: null,
                show: false
            }
        };

        getAllCustomers();
    };

    onMounted(() => {
        getAllCustomers();
    });
</script>

