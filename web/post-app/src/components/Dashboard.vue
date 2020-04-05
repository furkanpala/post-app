<template>
  <main class="posts-main">
    <div v-if="openOverlayPost" class="overlay">
      <span class="overlay-close-button" @click="closeOverlay">&times;</span>
      <PostOverlay :post="openOverlayPost" />
    </div>
    <div class="posts">
      <Post
        v-for="post in posts"
        :key="post.id"
        :post="post"
        @openOverlay="openOverlay"
      />
    </div>
    <div class="pagination-row">
      <span v-if="currentPage !== 1" @click="decrementPage" class="left-arrow"
        >&lt;</span
      >
      <span v-if="previousPage" class="previous-page" @click="decrementPage">{{
        previousPage
      }}</span>
      <span class="current-page">{{ currentPage }}</span>
      <span v-if="afterPage" class="after-page" @click="incrementPage">{{
        afterPage
      }}</span>
      <span
        v-if="currentPage !== totalPages"
        @click="incrementPage"
        class="right-arrow"
        >&gt;</span
      >
    </div>
  </main>
</template>

<script>
import axios from "axios";

import Post from "./Post";
import PostOverlay from "./PostOverlay";

axios.defaults.baseURL = "http://localhost:3000";

export default {
  name: "Dashboard",
  components: {
    Post,
    PostOverlay,
  },
  data() {
    return {
      currentPage: 1,
      totalPages: null,
      postsPerPage: 6,
      posts: [],
      openOverlayPost: null,
    };
  },
  computed: {
    previousPage() {
      return this.currentPage - 1 <= 0 ? null : this.currentPage - 1;
    },
    afterPage() {
      return this.currentPage + 1 > this.totalPages
        ? null
        : this.currentPage + 1;
    },
  },
  watch: {
    async currentPage() {
      const {
        data: { posts },
      } = await axios.get("/posts/" + this.currentPage);
      this.posts = posts;
    },
  },
  async mounted() {
    const {
      data: { count },
    } = await axios.get("/posts/amount");
    this.totalPages = Math.ceil(count / this.postsPerPage);

    const {
      data: { posts },
    } = await axios.get("/posts/" + this.currentPage);
    this.posts = posts;
  },
  methods: {
    incrementPage() {
      this.currentPage++;
    },
    decrementPage() {
      this.currentPage--;
    },
    openOverlay(id) {
      this.openOverlayPost = this.posts.find((post) => post.id === id);
    },
    closeOverlay() {
      this.openOverlayPost = null;
    },
  },
};
</script>

<style scoped>
.overlay {
  width: 100%;
  height: 100%;
  z-index: 1000;
  position: fixed;
  left: 0;
  top: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
.overlay-close-button {
  position: absolute;
  top: 1rem;
  right: 2rem;
  font-size: 5rem;
  padding: 8px;
  text-decoration: none;
  color: rgba(0, 0, 0, 0.8);
  display: block;
  transition: 0.3s;
  cursor: pointer;
}

.overlay-close-button:hover,
.overlay-close-button a:focus {
  color: #f1f1f1;
}
</style>
