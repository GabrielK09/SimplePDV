import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/admin'
    },
    {
        path: '/admin',
        component: () => import('src/layouts/MainLayout.vue'),
        children: [
            {
                path: '',
                component: () => import('src/pages/IndexPage.vue')
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
                    }
                ]
            },
            {
                path: 'shopping',
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
