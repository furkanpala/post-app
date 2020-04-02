import Vue from "vue";
import VueRouter from "vue-router";

import App from "./App.vue";
import Dashboard from "./components/Dashboard.vue";
import Login from "./components/auth/Login.vue";
import Register from "./components/auth/Register.vue";

Vue.config.productionTip = false;

Vue.use(VueRouter);

const router = new VueRouter({
  routes: [
    { path: "/", component: Dashboard },
    { path: "/login", component: Login },
    { path: "/register", component: Register }
  ],
  mode: "history"
});

new Vue({
  router: router,
  render: h => h(App)
}).$mount("#app");
