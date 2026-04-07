import type { RouteRecordRaw } from 'vue-router';
import Session from '@/mythicaldash/Session';
import Permissions from '@/mythicaldash/Permissions';

const timePackagesRoutes: RouteRecordRaw[] = [
    {
        path: '/mc-admin/time-packages',
        name: 'Time Packages',
        component: () => import('@/views/admin/TimePackages.vue'),
        meta: {
            requiresAuth: true,
            requiresAdmin: true,
        },
        beforeEnter: (to, from, next) => {
            if (Session.hasOrRedirectToErrorPage(Permissions.ADMIN_SERVERS_LIST)) {
                next();
            } else {
                next('/errors/403');
            }
        },
    },
];

export default timePackagesRoutes;
