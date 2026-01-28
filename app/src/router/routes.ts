import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
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
                        component: () => import('src/modules/products/pages/IndexPage.vue') 
                    },
                    {   
                        path: 'create',
                        name: 'products.create',
                        component: () => import('src/modules/products/pages/create/CreateProduct.vue') 
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
