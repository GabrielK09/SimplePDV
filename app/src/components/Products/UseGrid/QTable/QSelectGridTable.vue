<template>
    <q-dialog v-model="confirm" persistent>
        <q-card>
            <div class="fixed inset-0 z-50 flex items-center justify-center bg-opacity-40 backdrop-blur-sm">
                <q-table
                    title="Grade"
                    :rows="props.productData.productWithCharacteristics"
                    hide-bottom
                    :columns="gridTableColumn"
                    row-key="name"
                    @row-click="selectGrid"
                />

            </div>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { ref } from 'vue';

    const props = defineProps<{
        productData?: ProductContract,
        productId?: number
    }>();

    const emists = defineEmits<{
        (e: 'return:selected-grid', value: any)
    }>();

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
        }
    ];

    const confirm = ref<boolean>(true);

    const selectGrid = (_: Event, row: any) => {
        emists('return:selected-grid', row);
        confirm.value = false;
    };
</script>