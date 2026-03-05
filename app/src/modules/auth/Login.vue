<template>
    <main class="min-h-screen flex items-center justify-center px-4 py-10 bg-gray-50">
        <section class="w-full max-w-2xl">
            <article class="rounded-2xl border border-gray-200 bg-white shadow-sm overflow-hidden">
                <header class="px-6 py-6 border-b border-gray-100">
                    <h2 class="text-center text-2xl font-semibold text-gray-800">
                        Bem-vindo(a) de volta!
                    </h2>
                </header>

                <div class="px-6 py-6">
                    <q-form
                        @submit="loginSubmit"
                        class="q-gutter-md mt-4"
                    >
                        <div class="flex justify-center">
                            <q-input
                                v-model="auth.login"
                                type="text"
                                stack-label
                                label-slot
                                outlined
                                color="grey"
                                class="w-[120%] mb-4 rounded-xl"
                                hide-bottom-space

                            >
                                <template v-slot:label>
                                    <div class="border-b">Login <span class="text-red-500 text-xs relative bottom-1">*</span></div>
                                </template>

                                <template v-slot:prepend>
                                    <div class="mt-2 ml-2">
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-5">
                                            <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 12a4.5 4.5 0 1 1-9 0 4.5 4.5 0 0 1 9 0Zm0 0c0 1.657 1.007 3 2.25 3S21 13.657 21 12a9 9 0 1 0-2.636 6.364M16.5 12V8.25" />
                                        </svg>

                                    </div>
                                </template>
                            </q-input>

                            <q-input
                                v-model="auth.password"
                                :type="show ? 'text' : 'password'"
                                stack-label
                                label-slot
                                outlined
                                color="grey"
                                class="w-[120%] mb-4 rounded-xl"
                                required
                            >
                                <template v-slot:label>
                                    <div class="border-b">Senha <span class="text-red-500 text-xs relative bottom-1">*</span></div>
                                </template>

                                <template v-slot:prepend>
                                    <div class="flex mt-2 ml-4">
                                        <ShowPassword
                                            @show="show = $event"
                                        />

                                    </div>
                                </template>

                            </q-input>

                        </div>

                        <div class="w-40 ml-auto mr-auto">
                            <q-btn 
                                label="Entrar" 
                                type="submit" 
                                class="w-full" 
                                color="primary"
                            />

                        </div>
                    </q-form>
                </div>
            </article>
        </section>
    </main>
</template>

<script setup lang="ts">
    import { LocalStorage, useQuasar } from 'quasar';
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { onMounted, ref } from 'vue';
    import { loginService } from './services/loginService';
    import { useRouter } from 'vue-router';

    const show = ref<boolean>(false);

    const auth = ref<AuthContract>({
        login: '',
        password: ''
    });

    const router = useRouter();
    
    const { notify } = useNotify();

    const loginSubmit = async () => {
        notify(
            'positive',
            'Validando dados de login ... ',
        );

        const res = await loginService(auth.value);

        if(res.success)
        {
            notify(
                'positive',
                'Login bem sucedido!',
            );

            LocalStorage.set("authToken", res.data);

            router.replace({ path: `/admin` });

        } else {
            notify(
                'negative',
                res.message,
            );
            return;
        };
    };

    onMounted(() => {
        LocalStorage.remove("authToken");
        LocalStorage.remove("isAttendant");
        LocalStorage.remove("siteName");
        LocalStorage.remove("lastCheck");
        LocalStorage.remove("lastURL");
    });
</script>