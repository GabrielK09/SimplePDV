<template>
    <div v-if="props.isRegisteredCustomer">
        <q-select
            outlined
            v-model="customerId"
            :options="customers"
            option-label="name"
            option-value="id"
            emit-value
            map-options
        />
    </div>

    <div v-else>
        <q-input
            outlined
            v-model="customerName"
        />
    </div>
</template>

<script setup lang="ts">
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { getAll } from 'src/modules/customer/services/customerService';
    import { onMounted, ref, watch } from 'vue';

    const props = defineProps<{
        isRegisteredCustomer: boolean
    }>();

    const emits = defineEmits<{
        (e: 'selected:customer', value: any),
        (e: 'return:customer-name', value: any)
    }>();

    const { notify } = useNotify();
    const customers = ref<CustomerContract[]>([]);
    const customerId = defineModel<number | null>()
    const customerName = ref<string>('Consumidor padrão');

    watch(
        () => props.isRegisteredCustomer,
        (newVal) => {
            console.log(newVal);

            if(!newVal) 
            {
                customerId.value = 1;

            };
        }
    );
    
    watch(
        customerName,
        () => {        
            if(customerName.value !== 'Consumidor padrão' && customerName.value !== '')
            {
                emits('return:customer-name', customerName);

            } else {
                customerName.value = 'Consumidor padrão';
                emits('return:customer-name', customerName.value);
            };
        }
    );

    watch(
        customerId,
        () => {
            console.log('watch customerId');

            if(customerId.value < 1)
            {
                emits('return:customer-name', 'Consumidor padrão');

            } else {
                customerName.value = customers.value.find(c => c.id === customerId.value)?.name;

                emits('return:customer-name', customerName.value);
            };
        }
    );

    onMounted(async() => {
        const res = await getAll();
        const data = res.data;

        if(!res.success)
        {
            notify(
                'negative',
                res.message
            );
            return;
        };
        customers.value = data;

    });
</script>