<template>
    <q-dialog v-model="confirm" persistent>
        <q-card>
            <div class="fixed inset-0 z-50 items-center justify-center bg-opacity-40 backdrop-blur-sm flex">
                <q-table
                    title="Grades"
                    :rows="props.characteristics || []"
                    hide-bottom
                    :columns="gridTableColumn"
                    row-key="name"
                    @row-click="selectGrid"
                >
                    <template v-slot:top-right>
                        <q-btn color="red" icon="close" @click="emists('close',true)" />

                    </template>

                </q-table>
            </div>
        </q-card>
    </q-dialog>
</template>

<script setup lang="ts">
    import { QTableColumn } from 'quasar';
    import { ref } from 'vue';

    const props = defineProps<{
        characteristics: ProductCharacteristicsContract[]|undefined

    }>();

    const emists = defineEmits<{
        (e: 'close', value: boolean): void
        (e: 'return:selected-grid', value: any): void
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