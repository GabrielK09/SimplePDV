<template>
    <div class="border p-4">
        <q-table
            title="Grade"
            :rows="props.productWithCharacteristics"
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
        </q-table>
    </div>
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';

    const props = defineProps<{
        productWithCharacteristics: ProductCharacteristicsContract[]|null
    }>();

    const emists = defineEmits<{
        (e: 'showCreateGrid', value: boolean)
    }>();

    const gridTableColumn: QTableColumn[] = [
        {
            name: 'grid_qtde',
            label: 'Qtde',
            field: 'grid_qtde',
        },
        {
            name: 'size',
            label: 'Tamanho',
            field: 'size',
        },
        {
            name: 'action',
            label: '',
            field: 'action'
        }
    ];
</script>