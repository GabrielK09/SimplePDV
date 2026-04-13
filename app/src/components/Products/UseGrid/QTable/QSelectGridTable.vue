<template>
    <q-dialog v-model="confirm" persistent>
        <q-card class="product-dialog">
            <q-card-section class="dialog-header">
                <div class="flex justify-end">
                    <q-btn 
                        class="mb-4"
                        color="red" 
                        icon="close"
                        @click="emists('close',true)" 
                    />
                </div>
                <div class="mx-auto">

                    <q-table
                        title="Grades"
                        :rows="props.characteristics || []"
                        hide-bottom
                        :columns="gridTableColumn"
                        row-key="name"
                        @row-click="selectGrid"

                    />
                </div>
            </q-card-section>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { ref } from 'vue';

    const props = defineProps<{
        characteristics: ProductCharacteristicsContract[]

    }>();

    const emists = defineEmits<{
        (e: 'close', value: boolean): void
        (e: 'return:selected-grid', value: ProductCharacteristicsContract): void
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

<style scoped>
    .product-dialog {
        width: 100%;
        max-width: 100px;
        min-width: 320px;
        border-radius: 18px;
    }

    .dialog-header {
        background: linear-gradient(to right, #f8fafc, #ffffff);
    }
</style>