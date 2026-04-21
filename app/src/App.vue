<template>
  <router-view />
</template>

<script setup lang="ts">
    import { onMounted } from 'vue';
    import { api } from './boot/axios';
    import { useRouter } from 'vue-router';
    import { useNotify } from './helpers/QNotify/useNotify';
    import { LocalStorage } from 'quasar';

    const router = useRouter();
    const { notify } = useNotify();

    onMounted(async() => {
        const res = await api.get('/ping');

        if (res.data.data !== 'pong' && !res.data.success)
        {
            notify(
                'negative',
                'Erro interno'

            );

            LocalStorage.remove("authToken");
            
            router.replace({
                path: '/auth/login'
            });
        };
    });
</script>
