<template>
    <q-layout view="hHr LpR lFf">
        <q-btn
            @click="drawerLeft = !drawerLeft"
            flat
            class="rounded"
            icon="menu"

        />

        <q-drawer
            v-model="drawerLeft"
            show-if-above
            :width="210"
            class="bg-[#03202e] text-white rounded-r-md dashboard"
        >
            <q-toolbar>
                <q-list padding class="p-2">
                    <q-item
                        v-for="row in ticketRows"
                        v-ripple
                        clickable
                        :key="row.name"
                        :to="`/admin/${row.path}`"
                        :active="route.path === row.path"
                        class="rounded mt-3"
                        active-class="my-link"
                    >
                        <q-item-section avatar>
                            <q-icon :name="row.icon" />
                        </q-item-section>

                        <q-item-section>
                            {{ row.label }}
                        </q-item-section>
                    </q-item>
                </q-list>
            </q-toolbar>
        </q-drawer>

        <q-page-container>
            <div class="ml-4">
                <router-view />

            </div>
        </q-page-container>
    </q-layout>
</template>

<script setup lang="ts">
    import { ref } from 'vue';
    import { useRoute } from 'vue-router';

    type TicketRows = {
        label: string;
        icon: string;
        name: string;
        path: string;
    };

    const route = useRoute();
    const drawerLeft = ref<boolean>(true);

    const ticketRows = ref<TicketRows[]>([
        { label: 'DashBoard', icon: 'dashboard', name: 'dashboard', path: '' },
        { label: 'Produtos', icon: 'inventory_2', name: 'inventory_2', path: 'products' },
        { label: 'Compras', icon: 'shopping_bag', name: 'shopping_bag', path: 'shopping' },
        { label: 'Caixa', icon: 'request_quote', name: 'request_quote', path: 'cash-register' },
        { label: 'PDV', icon: 'point_of_sale', name: 'point_of_sale', path: 'pdv' },
        { label: 'Listagem de vendas', icon: 'analytics', name: 'analytics', path: 'pdv/list-pdv' }

    ]);
</script>

<style lang="scss">
    .dashboard {
        height: 100vh;

    }

    .my-link {
        color: #fff;
        background: #07425f;
    }

    .fixed-logout-button {
        display: flex;
    }

    @media (max-width: 1100px) {
        .fixed-logout-button {
            position: fixed;
            bottom: 0;
            margin: 0 0 2rem 0;
        }
    }
</style>
