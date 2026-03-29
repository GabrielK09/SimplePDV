<template>
    <div class="border p-4">
        <q-table
            title="Grade"
            :rows="props.productData.productWithCharacteristics"
            hide-bottom
            :columns="gridTableColumn"
            row-key="name"
        >
            <template v-slot:top-right>
                <q-btn 
                    color="primary" 
                    no-caps
                    label="Cadastar uma grade" 
                    @click="emists('showCreateGrid', true)" 
                />
            </template>

             <template v-slot:body="props">
                <q-tr :props="props">
                    <q-td v-for="col in props.cols">
                        <div class="flex flex-center">
                            <template v-if="col.name === 'actions'">
                                <div class="flex gap-2">
                                    <q-btn 
                                        color="primary" 
                                        icon="edit"
                                        dense
                                        @click="buildUpdateGrid(props.row)"
                                    />

                                    <q-btn 
                                        color="red" 
                                        icon="delete"
                                        dense
                                        @click="deleteGrid(props.row)"
                                    />

                                </div>
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

    <UpdateGridProduct
        v-if="showUpdateGrid"
        :grid-id="selectedGridId"
        :product-id="props.productId"
        :grid-full-object="gridFullObject"
        :selected-sizes="props.productData.productWithCharacteristics"
        @return:grids="handleUpdateGrid($event)"
        @close="showUpdateGrid = !$event"
    /> 
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import UpdateGridProduct from 'src/components/Products/UseGrid/Update/UpdateGridProduct.vue';
    import { ref } from 'vue';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import formatGridDataForPush from 'src/helpers/FormatValue/Grid/formatGridDataForPush';
import { deleteProductCharacteristics } from 'src/modules/products/services/productsService';

    const props = defineProps<{
        productData: ProductContract,
        productId?: number
    }>();

    const emists = defineEmits<{
        (e: 'showCreateGrid', value: boolean)
    }>();

    const showUpdateGrid = ref<boolean>(false);
    const selectedGridId = ref<number | null>(); // Usado quando for feito a edição de uma grade já cadsatrada.
    const gridFullObject = ref<any>(); // Usado quando for cadastrado uma nova grade.
    const { notify } = useNotify();

    const gridTableColumn: QTableColumn[] = [
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
            label: '',
            field: 'actions',
            align: 'center'
        }
    ];

    const buildUpdateGrid = (row: any) => {        
        showUpdateGrid.value = true;
        if(!row?.id) {
            console.warn('Grid não encontrada:', row?.id);

            const newGrid = {
                id: null,
                product_id: props.productId,
                grid_qtde: row.grid_qtde,
                size: row.size,
                have_register: false
            };

            gridFullObject.value = newGrid;

        } else {
            selectedGridId.value = row?.id;
        };
    };

    const deleteGrid = async (row: any) => {
        const index = props.productData.productWithCharacteristics.indexOf(row);
        console.log(index);
        
        if(index > -1)
        {
            const item = props.productData.productWithCharacteristics[index];

            if(item?.product_id && item?.id)
            {
                const res = await deleteProductCharacteristics({
                    gridId: item.id,
                    productGridId: item.product_id
                });

                if(!res.success)
                {
                    notify(
                        'negative',
                        res.message
                    );
                    return;
                };
            };

            props.productData.productWithCharacteristics.splice(index, 1);
            
        };
    };

    const handleUpdateGrid = (newGrid: ProductCharacteristicsContract) => {
        const oldGrid = props.productData.productWithCharacteristics.find(c => c.id === newGrid.id);

        if (!oldGrid)
        {
            notify(
                'negative',
                'Ocorreu um erro ao alterar a grade.'
            );
            return;
        };

        const index = props.productData.productWithCharacteristics.indexOf(oldGrid);

        if (index > -1)
        {
            props.productData.productWithCharacteristics.splice(index, 1);
            props.productData.productWithCharacteristics.push(formatGridDataForPush(newGrid));
            
        } else {
            return;  
        };
    };    
</script>