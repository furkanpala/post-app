<template>
  <main class="form-main">
    <div class="form-group">
      <div class="error" v-if="error">{{ error }}</div>
      <span class="form-name">Login</span>
      <form
        action="#"
        autocomplete="off"
        spellcheck="false"
        @submit.prevent="login"
      >
        <label for="username">Username</label>
        <input type="text" name="username" v-model="username" required />
        <label for="password">Password</label>
        <input type="password" name="password" v-model="password" required />
        <button v-if="!loading" type="submit" class="submit-button">
          Login
        </button>
        <button v-if="loading" disabled class="disabled-submit-button">
          Logging In...
        </button>
      </form>
    </div>
  </main>
</template>

<script>
export default {
  name: "Login",
  data() {
    return {
      username: "",
      password: "",
      error: "",
      loading: false,
    };
  },
  methods: {
    login() {
      this.loading = true;
      this.$store
        .dispatch("login", {
          username: this.username,
          password: this.password,
        })
        .then(() => {
          this.loading = false;
          this.$router.push("/");
        })
        .catch((err) => {
          this.loading = false;
          this.error =
            err.response.data.message.title +
            " - " +
            err.response.data.message.detail;
          this.username = "";
          this.password = "";
        });
    },
  },
};
</script>
