<template>
  <div class="post" @click="openOverlay">
    <div class="post-info">
      <span class="post-info-title">{{ title }}</span>
      <!-- ilk 20 karakter tablet pc-->
      <!-- ilk 40 karakter mobil -->
      <div>
        <span class="post-info-writer">by {{ writer }}</span>
        <span class="post-info-date">{{ date }}</span>
      </div>
    </div>
    <p class="post-content">
      {{ content }}
    </p>
    <!-- ilk 100 karakter tablet pc-->
    <!-- ilk 150 karakter mobil-->
  </div>
</template>

<script>
export default {
  name: "Post",
  props: {
    post: Object,
  },
  data() {
    return {
      writer: this.post.user,
      windowWidth: window.innerWidth,
    };
  },

  methods: {
    openOverlay() {
      this.$emit("openOverlay", this.post.id);
    },
  },

  computed: {
    title() {
      if (this.windowWidth < 700) {
        return this.post.title.substr(0, 40);
      }
      return this.post.title.substr(0, 20);
    },
    content() {
      if (this.windowWidth < 700) {
        return this.post.content.substr(0, 150);
      }
      return this.post.content.substr(0, 100);
    },
    date() {
      return new Date(this.post.date * 1000).toDateString();
    },
  },
};
</script>

<style>
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

.post-info-writer,
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
