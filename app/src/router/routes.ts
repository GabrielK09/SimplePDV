import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/admin'
    },
    {
        path: '/admin',
        component: () => import('src/layouts/MainLayout.vue'),
        props: true,
        children: [
            {
                path: '',
                component: () => import('src/modules/dashBoard/pages/IndexPage.vue')
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
                    },
                    {
                        path: 'create',
                        name: 'products.create',
                        component: () => import('src/modules/products/pages/create/CreateProduct.vue')
                    },
                    {
                        path: 'edit/:id',
                        name: 'products.edit',
                        props: true,
                        component: () => import('src/modules/products/pages/update/UpdateProduct.vue')
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
                    },
                    {
                        path: 'create',
                        name: 'customers.create',
                        component: () => import('src/modules/customer/pages/create/CreateCustomer.vue')
                    },
                    {
                        path: 'edit/:id',
                        name: 'customers.edit',
                        props: true,
                        component: () => import('src/modules/customer/pages/update/UpdateCustomer.vue')
                    }
                ]
            },
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
