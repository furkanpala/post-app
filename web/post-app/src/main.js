import Vue from "vue";
import VueRouter from "vue-router";

import App from "./App.vue";
import Dashboard from "./components/Dashboard.vue";
import Login from "./components/auth/Login.vue";
import Register from "./components/auth/Register.vue";
import AddPost from "./components/AddPost.vue";
import { store } from "./store";

Vue.config.productionTip = false;

Vue.use(VueRouter);
const originalPush = VueRouter.prototype.push;
VueRouter.prototype.push = function push(location) {
  return originalPush.call(this, location).catch((err) => err);
};

const router = new VueRouter({
  routes: [
    { path: "/", component: Dashboard },
    { path: "/login", component: Login, meta: { visitorRequired: true } },
    { path: "/register", component: Register, meta: { visitorRequired: true } },
    { path: "/add", component: AddPost, meta: { authRequired: true } },
  ],
  mode: "history",
});

router.beforeEach((to, from, next) => {
  if (to.matched.some((record) => record.meta.authRequired)) {
    if (!store.getters.isLoggedIn) {
      next("/login");
    } else {
      next();
    }
  } else if (to.matched.some((record) => record.meta.visitorRequired)) {
    if (store.getters.isLoggedIn) {
      next("/");
    } else {
      next();
    }
  } else {
    next();
  }
});

new Vue({
  router: router,
  store: store,
  render: (h) => h(App),
}).$mount("#app");
