<template>
  <main class="form-main">
    <div class="form-group">
      <div class="error" v-for="(value, index) in error" :key="index">
        {{ value }}
      </div>
      <span class="form-name">Register</span>
      <form
        action="#"
        autocomplete="off"
        spellcheck="false"
        @submit.prevent="register"
      >
        <label for="username">Username</label>
        <input type="text" name="username" v-model="username" required />
        <label for="password">Password</label>
        <input type="password" name="password" v-model="password" required />
        <label for="passwordRepeat">Repeat Password</label>
        <input
          type="password"
          name="passwordRepeat"
          v-model="passwordRepeat"
          required
        />
        <button v-if="!loading" type="submit" class="submit-button">
          Register
        </button>
        <button v-if="loading" disabled class="disabled-submit-button">
          Registering...
        </button>
      </form>
    </div>
  </main>
</template>

<script>
export default {
  name: "Register",
  data() {
    return {
      username: "",
      password: "",
      passwordRepeat: "",
      error: [],
      loading: false,
    };
  },
  methods: {
    register() {
      this.loading = true;
      this.error = [];

      if (this.password !== this.passwordRepeat) {
        this.error.push("Passwords did not match");
      } else {
        this.$store
          .dispatch("register", {
            username: this.username,
            password: this.password,
          })
          .then(() => {
            this.$router.push("/login");
          })
          .catch((err) => {
            if (err.response.status === 500) {
              this.error.push(err.response.data.message.title);
            } else if (err.response.status === 409) {
              const {
                response: {
                  data: {
                    message: { title },
                  },
                },
              } = err;
              this.error.push(title);
            } else {
              const {
                response: {
                  data: {
                    message: { detail },
                  },
                },
              } = err;
              this.error = detail.split("|");
            }
            this.username = "";
            this.password = "";
            this.passwordRepeat = "";
          });
      }
      this.loading = false;
    },
  },
};
</script>
