import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

import store from "@/store/index";
import Analytics from "@/views/Analytics.vue";
import Index from "@/views/Index.vue";
import NotFound from "@/views/NotFound.vue";
import UserAuth from "@/views/UserAuth.vue";

const routes: Array<RouteRecordRaw> = [
  { path: "/", component: Index },
  { path: "/analytics", component: Analytics, meta: { requiresAuth: true } },
  { path: "/auth", component: UserAuth, meta: { requiresUnauth: true } },
  { path: "/:notFound(.*)", component: NotFound },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach((to, _from, next) => {
  if (to.meta.requiresAuth && !store.getters.isAuthenticated) {
    next("/auth");
  } else if (to.meta.requiresUnauth && store.getters.isAuthenticated) {
    next("/");
  } else {
    next();
  }
});

export default router;
