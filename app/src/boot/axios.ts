import { boot } from 'quasar/wrappers';
import axios, { type AxiosInstance } from 'axios';

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
const api = axios.create({
    baseURL: process.env.API_URL

});

export default boot(({ app }) => {
    api.interceptors.request.use(
        (config) => {
            return config;
        },
        (error) => {
            console.error('Erro: ', error);

            return Promise.reject(error);
        }
    );

    api.interceptors.response.use(
        (response) => response,
        (error) => {
            if (error.response) {
                console.error('Status:', error.response.status);
                console.error('Mensagem:', error.response.data.message);
                console.error('Erro:', error.response.data.data);

            } else {
                console.error('Erro de rede:', error.message);
            }

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

export { api };
