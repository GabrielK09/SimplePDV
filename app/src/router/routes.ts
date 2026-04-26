import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/admin',
        meta: {
            title: ''
        }
    },
    {
        path: '/auth/login',
        name: 'auth.login',
        component: () => import('src/modules/auth/Login.vue'),
        meta: {
            title: 'Login'
        }
    },
    {
        path: '/admin',
        component: () => import('src/layouts/MainLayout.vue'),
        props: true,
        children: [
            {
                path: '',
                component: () => import('src/modules/dashBoard/pages/DashBoard.vue'),
                meta: {
                    title: 'DashBoard'
                }
            },
            {
                path: 'cash-register',
                meta: {
                    title: 'Caixa'
                },
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
                meta: {
                    title: 'Produtos'
                },
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
                meta: {
                    title: 'PDV'
                },
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
                meta: {
                    title: 'Clientes'
                },
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
                meta: {
                    title: 'Compras'
                },
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


