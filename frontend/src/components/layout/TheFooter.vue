<template>
  <footer>
    <nav>
      <h1>
        <img src="@/assets/images/logo.png" alt="remrratality" />
        <router-link to="/">remrratality</router-link>
      </h1>
      <ul>
        <li>© 2021 hackfeed</li>
        <li><a href="mailto:hackfeed@yandex.ru">Contact me</a></li>
        <li>
          <a href="https://github.com/hackfeed/remrratality">Project Page</a>
        </li>
      </ul>
    </nav>
  </footer>
</template>

<script lang="ts">
import { Vue } from "vue-class-component";

export default class TheFooter extends Vue {}
</script>

<style lang="scss" scoped>
footer {
  width: 100%;
  height: 3rem;
  border-top: 1px solid rgba(0, 0, 0, 0.26);
  position: absolute;
  bottom: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  nav {
    width: 90%;
    margin: auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  ul {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    justify-content: center;
    align-items: center;
  }
}

li {
  margin: 0 0.5rem;
  a {
    color: black;
    text-decoration: none;
    &:hover {
      text-decoration: underline;
    }
  }
}

h1 {
  margin: 0;
  display: flex;
  align-items: center;
  color: black;
  img {
    height: 2rem;
  }
  a {
    color: black;
    text-decoration: none;
    margin: 0;
    padding: 0.75rem 1.5rem;
    display: inline-block;
    &:hover,
    &:active,
    &.router-link-active {
      border-color: transparent;
    }
  }
}
</style>
