<template>
  <main class="form-main">
    <div class="form-group">
      <span class="form-name">Login</span>
      <form action="#" autocomplete="off" spellcheck="false">
        <label for="username">Username</label>
        <input type="text" name="username" />
        <label for="password">Password</label>
        <input type="password" name="password" />
        <button type="submit">Submit</button>
      </form>
    </div>
  </main>
</template>

<script>
export default {
  name: "Login"
};
</script>
