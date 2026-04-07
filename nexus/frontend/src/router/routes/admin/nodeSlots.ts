import type { RouteRecordRaw } from 'vue-router';
import Session from '@/mythicaldash/Session';
import Permissions from '@/mythicaldash/Permissions';

const nodeSlotsRoutes: RouteRecordRaw[] = [
    {
        path: '/mc-admin/node-slots',
        name: 'Node Slots',
        component: () => import('@/views/admin/NodeSlots.vue'),
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

export default nodeSlotsRoutes;
