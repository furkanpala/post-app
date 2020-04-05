import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";

const refreshTokenAxios = axios.create();
refreshTokenAxios.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    const originalRequest = error.config;

    if (error.response.status === 401 && originalRequest.url === "/token") {
      return Promise.reject(error);
    }
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      return axios
        .post("/token", null, { withCredentials: true })
        .then((res) => {
          if (res.status === 200) {
            const token = res.data.access_token;
            store.commit("setToken", token);
            originalRequest.headers["Authorization"] =
              "bearer " + store.state.token;
            return refreshTokenAxios(originalRequest);
          }
        });
    }
    return Promise.reject(error);
  }
);

Vue.use(Vuex);

export const store = new Vuex.Store({
  state: {
    token: localStorage.getItem("access_token") || null,
  },

  mutations: {
    setToken(state, token) {
      state.token = token;
    },
    clearToken(state) {
      state.token = null;
    },
  },

  getters: {
    isLoggedIn(state) {
      return state.token !== null;
    },
  },

  actions: {
    checkAuth(context) {
      axios
        .post("/token", null, {
          withCredentials: true,
        })
        .then((response) => {
          const token = response.data.access_token;
          context.commit("setToken", token);
        })
        .catch(() => {
          context.commit("clearToken");
        });
    },

    login(context, credentials) {
      return new Promise((resolve, reject) => {
        axios
          .post(
            "/login",
            {
              username: credentials.username,
              password: credentials.password,
            },
            {
              withCredentials: true,
            }
          )
          .then((response) => {
            const token = response.data.access_token;
            localStorage.setItem("access_token", token);

            context.commit("setToken", token);
            resolve(response);
          })
          .catch((err) => reject(err));
      });
    },
    register(_, credentials) {
      return new Promise((resolve, reject) => {
        axios
          .post(
            "/register",
            {
              username: credentials.username.trim(),
              password: credentials.password.trim(),
            },
            {
              withCredentials: true,
            }
          )
          .then((response) => {
            resolve(response);
          })
          .catch((err) => reject(err));
      });
    },

    logout(context) {
      if (context.getters.isLoggedIn) {
        return new Promise((resolve, reject) => {
          axios
            .post("/token/logout", null, { withCredentials: true })
            .then((response) => {
              localStorage.removeItem("access_token");
              context.commit("clearToken");
              resolve(response);
            })
            .catch((err) => {
              reject(err);
            });
        });
      }
    },
    addPost(_, post) {
      return new Promise((resolve, reject) => {
        refreshTokenAxios
          .post(
            "/posts",
            { title: post.title, content: post.content },
            {
              headers: {
                Authorization: "bearer " + this.state.token,
              },
            }
          )
          .then((res) => resolve(res))
          .catch((err) => {
            reject(err);
          });
      });
    },
  },
});
