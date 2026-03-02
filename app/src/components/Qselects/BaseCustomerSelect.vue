<template>
    <q-checkbox
        right-label
        v-model="internalRegisteredCustomer"
        label="Cliente cadastrado"
    />

    <div v-if="internalRegisteredCustomer">
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
        (e: 'update:customer', value: any)
    }>();

    const { notify } = useNotify();
    const customers = ref<CustomerContract[]>([]);

    const customerId = defineModel<number | null>()
    const customerName = ref<string>('Consumidor padrão');
    const registeredCustomerByProps = ref<boolean>(props.isRegisteredCustomer);
    const internalRegisteredCustomer = registeredCustomerByProps;

    const selectedCustomer = ref<CustomerContract | null>(null);

    watch(
        selectedCustomer,
        (val) => {
            emits('update:customer', val.id);
            emits('selected:customer', val);
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
        console.log(customers.value);

    });

</script>
