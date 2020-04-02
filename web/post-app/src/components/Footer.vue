<template>
  <footer class="footer">
    <span class="me">Furkan Pala</span>
    <span class="pipe">|</span>
    <span class="github-link">
      <a href="https://github.com/furkanpala">Github</a>
    </span>
  </footer>
</template>

<script>
export default {
  name: "Footer"
};
</script>

<style scoped>
.footer {
  flex-shrink: 0;
  display: flex;
  justify-content: center;
  box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.5);
  min-height: 2rem;
}

.pipe {
  margin: 0 0.5rem;
}
</style>