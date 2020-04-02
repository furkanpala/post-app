<template>
  <router-view></router-view>
</template>

<script>
export default {
  name: "Main"
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
  margin-bottom: 0.5rem;
}
.form-group button {
  font-size: 1em;
  color: rgba(0, 0, 0, 0.7);
  align-self: center;
  height: 2.5rem;
  width: 5rem;
  background-color: rgba(0, 0, 0, 0.1);
}

.form-group button:hover {
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.form-group button:active {
  box-shadow: 0 -4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.form-name {
  font-size: 1.8rem;
}

.form-group form {
  display: flex;
  flex-direction: column;
}

@media only screen and (min-width: 700px) {
  .posts-main {
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>