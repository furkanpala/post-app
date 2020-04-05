<template>
  <div id="app">
    <Header />
    <Main />
    <Footer />
  </div>
</template>

<script>
import Header from "./components/Header";
import Main from "./components/Main";
import Footer from "./components/Footer";

export default {
  name: "App",
  components: {
    Header,
    Main,
    Footer,
  },
  created() {
    this.$store.dispatch("checkAuth");
  },
};
</script>

<style>
* {
  padding: 0;
  margin: 0;
  box-sizing: border-box;
}

html,
body {
  height: 100%;
}

a {
  text-decoration: none;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  display: flex;
  flex-direction: column;
  min-height: 100%;
}

/* .content-text {
  font-size: 21px;
  font-weight: 400;
  word-break: break-word;
  line-height: 1.58;
  color: rgba(0, 0, 0, 0.84);
}

.title-text {
}

.caption-text {
} */
</style>
