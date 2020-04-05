<template>
  <div class="overlay-post">
    <div class="post-info">
      <span class="overlay-post-info-title">{{ title }}</span>
      <div>
        <span class="overlay-post-info-writer">by {{ writer }}</span>
        <span class="overlay-post-info-date">{{ date }}</span>
      </div>
    </div>
    <p class="overlay-post-content">
      {{ content }}
    </p>
  </div>
</template>

<script>
export default {
  name: "PostOverlay",
  props: {
    post: Object,
  },

  data() {
    return {
      title: this.post.title,
      content: this.post.content,
      writer: this.post.user,
      windowWidth: window.innerWidth,
    };
  },

  computed: {
    date() {
      return new Date(this.post.date * 1000).toDateString();
    },
  },
};
</script>

<style scoped>
.overlay-post {
  height: 40rem;
  width: 40rem;
  display: flex;
  flex-direction: column;
  background-color: rgb(218, 211, 211);
  padding: 1rem;
  box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.1);
}

.overlay-post-info-title {
  font-size: 1.75em;
  font-weight: 400;
}

.overlay-post-info-writer,
.overlay-post-info-date {
  font-size: 1em;
  color: rgba(0, 0, 0, 0.5);
}

.overlay-post-content {
  font-size: 1.25em;
  color: rgba(0, 0, 0, 0.8);
  overflow: auto;
}
</style>
