<template>
  <header class="header">
    <span class="header-title">
      <router-link to="/">PostApp</router-link>
    </span>
    <nav class="header-nav">
      <span class="nav-item" v-if="!isLoggedIn">
        <router-link to="/login">Login</router-link>
      </span>
      <span class="nav-item" v-if="!isLoggedIn">
        <router-link to="/register">Register</router-link>
      </span>
      <span class="nav-item" v-if="isLoggedIn">
        <router-link to="/add">Add Post</router-link>
      </span>
      <span class="nav-item nav-logout" v-if="isLoggedIn" @click="logout"
        >Logout</span
      >
    </nav>
  </header>
</template>

<script>
export default {
  name: "Header",
  computed: {
    isLoggedIn() {
      return this.$store.getters.isLoggedIn;
    },
  },
  methods: {
    logout() {
      this.$store.dispatch("logout").then(() => {
        this.$router.push("/");
      });
    },
  },
};
</script>

<style scoped>
.header {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 1rem 0;
  text-align: center;
}

.header-title {
  font-size: 2em;
}

.header-title a {
  color: rgba(0, 0, 0, 0.84);
  font-weight: 500;
}

.header-nav {
  display: flex;
  justify-content: space-evenly;
}

.nav-item a {
  color: rgba(0, 0, 0, 0.84);
}

.nav-logout {
  cursor: pointer;
}

@media only screen and (min-width: 992px) {
  .header {
    flex-direction: row;
    padding: 1.5rem;
    text-align: start;
  }

  .header-title {
    flex-grow: 6;
  }

  .header-nav {
    flex-grow: 1;
    justify-content: space-around;
    align-items: center;
  }
}
</style>
