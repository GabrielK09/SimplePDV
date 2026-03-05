import { boot } from 'quasar/wrappers';
import axios, { type AxiosInstance } from 'axios';
import { LocalStorage } from 'quasar';
import { useNotify } from 'src/helpers/QNotify/useNotify';

axios.defaults.withCredentials = false;
declare module 'vue' {
    interface ComponentCustomProperties {
        $axios: AxiosInstance;
        $api: AxiosInstance;
        $apiService: AxiosInstance;
    }
}

// Be careful when using SSR for cross-request state pollution
// due to creating a Singleton instance here;
// If any client changes this (global) instance, it might be a
// good idea to move this instance creation inside of the
// "export default () => {}" function below (which runs individually
// for each client)

const { notify } = useNotify();

const api = axios.create({
    baseURL: process.env.API_URL,
    headers: {
        Accept: 'application/json'
    }
});

const apiCep = axios.create({
    baseURL: process.env.API_CEP,
    headers: {
        Accept: 'application/json'
    }
});

const apiCnpj = axios.create({
    baseURL: process.env.API_CNPJ,
    headers: {
        Accept: 'application/json'
    }
});

export default boot(({ app, router }) => {
    api.interceptors.request.use(
        (config) => {
            const publicRoutes: string[] = [
                '/auth/login',
            ];

            const isPublic = publicRoutes.some(route => config.url?.includes(route));
            const token = LocalStorage.getItem("authToken");

            if (!isPublic && (token === "undefined" || !token))
            {
                router.replace({
                    path: '/auth/login'
                });

                LocalStorage.remove("authToken");
            };

            if(token) config.headers.Authorization = `Bearer ${token}`;

            return config;
        },
        (error) => {      
            return Promise.reject(error);
        }
    );

    api.interceptors.response.use(
        (response) => response,
        (error) => {            
            if (error.response.status === 401) 
            {
                router.replace({
                    path: '/auth/login'
                });

                return Promise.reject(error);
            };

            if (error.response) {
                console.error('Eror:', error);
                console.error('Status:', error.response.status);
                console.error('Mensagem:', error.response.data.message);
                console.error('Erro:', error.response.data.data);

            } else {
                console.error('Erro de rede:', error.message);
            };

            return Promise.reject(error);
        }
    );

    app.config.globalProperties.$axios = axios;
    // ^ ^ ^ this will allow you to use this.$axios (for Vue Options API form)
    //       so you won't necessarily have to import axios in each vue file

    app.config.globalProperties.$api = api;
    // ^ ^ ^ this will allow you to use this.$api (for Vue Options API form)
    //       so you can easily perform requests against your app's API
});

export { api, apiCnpj, apiCep };
