<template>
  <div class="post">
    <div class="post-info">
      <span class="post-info-title">Merhaba Dünya</span>
      <!-- ilk 20 karakter tablet pc-->
      <!-- ilk 40 karakter mobil -->
      <div>
        <span class="post-info-writer">by Furkanpala</span>
        <span class="post-info-date">29/10/2019 UTC</span>
      </div>
    </div>
    <p
      class="post-content"
    >elit, sed do eiusmodelit, sed do eiusmodelit, sed do eiusmodelit, sed do eiusmodelit, sed do eiusmodelit, sed do eiusmodelit, sed do eiusmodelit, sed</p>
    <!-- ilk 100 karakter tablet pc-->
    <!-- ilk 150 karakter mobil-->
  </div>
</template>

<script>
export default {
  name: "Post"
};
</script>

<style scoped>
.post {
  display: flex;
  padding: 1rem;
  flex-direction: column;
  box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.1);
  cursor: pointer;
}

.post-info {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.post-info-title {
  font-size: 1.5em;
  font-weight: 400;
}

.post-info div {
  display: flex;
  flex-direction: row;
  width: 100%;
  justify-content: space-around;
}

.post-info-writer {
  font-size: 0.75em;
  color: rgba(0, 0, 0, 0.5);
}

.post-info-date {
  font-size: 0.75em;
  color: rgba(0, 0, 0, 0.5);
}

.post-content {
  color: rgba(0, 0, 0, 0.8);
}

@media only screen and (min-width: 700px) {
  .post {
    width: 20rem;
    height: 10rem;
  }
  .post:hover {
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
  }
}
</style>