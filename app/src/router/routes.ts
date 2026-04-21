import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/admin'
    },
    {
        path: '/auth/login',
        name: 'auth.login',
        component: () => import('src/modules/auth/Login.vue'),
    },
    {
        path: '/admin',
        component: () => import('src/layouts/MainLayout.vue'),
        props: true,
        children: [
            {
                path: '',
                component: () => import('src/modules/dashBoard/pages/DashBoard.vue'),
            },
            {
                path: 'api',
                name: 'api',
                component: () => import('src/modules/api/ApiTest.vue')
            },
            {
                path: 'cash-register',
                children: [
                    {
                        path: '',
                        name: 'cash-register.index',
                        component: () => import('src/modules/cashRegister/pages/IndexCashRegister.vue')
                    },
                    {
                        path: 'create',
                        name: 'cash-register.create',
                        component: () => import('src/modules/cashRegister/pages/create/CreateCashRegister.vue')
                    }
                ]
            },
            {
                path: 'products',
                children: [
                    {
                        path: '',
                        name: 'products.index',
                        component: () => import('src/modules/products/pages/IndexProducts.vue')
                    }
                ]
            },
            {
                path: 'pdv',
                children: [
                    {
                        path: 'list-pdv',
                        name: 'list-pdv.index',
                        component: () => import('src/modules/PDV/pages/IndexPDV.vue')
                    },
                    {
                        path: '',
                        name: 'pdv.index',
                        component: () => import('src/modules/PDV/pages/PDV.vue')
                    }
                ]
            },
            {
                path: 'customers',
                children: [
                    {
                        path: '',
                        name: 'customers.index',
                        component: () => import('src/modules/customer/pages/IndexCustomer.vue')
                    }
                ]
            },
            {
                path: 'shopping',
                name: 'shopping',
                children: [
                    {
                        path: '',
                        name: 'shopping.index',
                        component: () => import('src/modules/shopping/pages/IndexShopping.vue')
                    },
                    {
                        path: 'create',
                        name: 'shopping.create',
                        component: () => import('src/modules/shopping/pages/create/CreateShopping.vue')
                    },
                ]
            }
        ],
    },

    // Always leave this as last one,
    // but you can also remove it
    {
        path: '/:catchAll(.*)*',
        component: () => import('src/pages/ErrorNotFound.vue'),
        },
    ];

export default routes;


