<template>
  <main class="form-main">
    <div class="form-group">
      <div class="error" v-if="error">{{ error }}</div>
      <div class="success" v-if="success">{{ success }}</div>
      <span class="form-name">Add Post</span>
      <form
        action="#"
        autocomplete="off"
        spellcheck="true"
        @submit.prevent="addPost"
      >
        <label for="title">Title</label>
        <input type="text" name="title" v-model="title" required />
        <label for="content">Content</label>
        <textarea
          type="text"
          name="content"
          class="content-area"
          v-model="content"
          required
        />
        <button v-if="!adding" type="submit" class="submit-button">
          Add Post
        </button>
        <button v-if="adding" disabled class="disabled-submit-button">
          Adding...
        </button>
      </form>
    </div>
  </main>
</template>

<script>
export default {
  name: "AddPost",
  data() {
    return {
      title: "",
      content: "",
      adding: false,
      success: "",
      error: "",
    };
  },
  methods: {
    addPost() {
      this.adding = true;
      this.$store
        .dispatch("addPost", {
          title: this.title,
          content: this.content,
        })
        .then(() => {
          this.adding = false;
          this.success = "Post created successfully";
        })
        .catch((err) => {
          this.adding = false;
          this.error =
            err.response.data.message.title +
            " - " +
            err.response.data.message.detail;
        });
      this.title = "";
      this.content = "";
    },
  },
};
</script>
