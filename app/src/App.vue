<template>
  <router-view />
</template>

<script setup lang="ts">
    import { onMounted } from 'vue';
    import { api } from './boot/axios';
    import { useRouter } from 'vue-router';
    import { useNotify } from './helpers/QNotify/useNotify';

    const router = useRouter();
    const { notify } = useNotify();

    onMounted(async() => {
        const res = await api.get('/ping');
        const data = res.data;

        if (data.data !== 'pong') 
        {
            notify(
                'negative',
                'Erro interno'

            );
            
            router.replace({
                path: '/auth/login'
            });
        };
    });
</script>