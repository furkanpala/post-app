<template>
  <router-view></router-view>
</template>

<script>
export default {
  name: "Main",
};
</script>

<style>
.posts-main,
.form-main {
  flex-grow: 1;
  flex-shrink: 0;

  display: flex;
  flex-direction: column;
}

.posts {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}
.form-main {
  align-items: center;
}

.form-group {
  margin-top: 2.5rem;
  display: flex;
  flex-direction: column;
  width: 20rem;
}

.form-group input {
  font-size: 1.2em;
  border-top: none;
  border-right: none;
  border-left: none;
  border-bottom: 1px solid rgba(0, 0, 0, 0.5);
  margin-bottom: 0.25rem;
}
.submit-button {
  margin-top: 0.25rem;
  font-size: 1em;
  color: rgba(0, 0, 0, 0.7);
  align-self: center;
  height: 2.5rem;
  width: 5rem;
  background-color: rgba(0, 0, 0, 0.1);
}

.disabled-submit-button {
  margin-top: 0.25rem;
  font-size: 1em;
  color: rgba(0, 0, 0, 0.2);
  align-self: center;
  height: 2.5rem;
  width: 5rem;
  background-color: rgba(0, 0, 0, 0.1);
}

.content-area {
  min-height: 10rem;
  max-height: 20rem;
  width: 20rem;
  font-size: 1.2em;
  resize: vertical;
  color: rgba(0, 0, 0, 0.8);
  padding: 0.5rem;
  line-height: 1.75rem;
  margin-bottom: 0.25rem;
}

.submit-button:hover {
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.submit-button:active {
  box-shadow: 0 -4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.form-name {
  font-size: 1.8rem;
}

.form-group form {
  display: flex;
  flex-direction: column;
}

.pagination-row {
  display: flex;
  flex-direction: row;
  margin: 1rem 0;
  align-self: center;
}

.pagination-row span {
  margin: 0 1.5rem;
  box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.1);
  padding: 0.2rem 0.4rem;
  cursor: pointer;
}

.pagination-row span:hover {
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.pagination-row span:active {
  box-shadow: 0 -4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

@media only screen and (min-width: 700px) {
  .posts {
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>
