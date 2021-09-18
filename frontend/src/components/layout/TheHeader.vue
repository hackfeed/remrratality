<template>
  <header>
    <nav>
      <h1>
        <img src="@/assets/logo.png" alt="mrratality" />
        <router-link to="/">mrratality</router-link>
      </h1>
      <ul>
        <li v-if="isLoggedIn">
          <router-link to="/analytics">Analytics</router-link>
        </li>
        <li v-else>
          <router-link to="/auth">Login</router-link>
        </li>
        <li v-if="isLoggedIn">
          <base-button @click="logout">Logout</base-button>
        </li>
      </ul>
    </nav>
  </header>
</template>

<script lang="ts">
import { Vue } from "vue-class-component";

export default class TheHeader extends Vue {
  get isLoggedIn(): boolean {
    return this.$store.getters.isAuthenticated;
  }

  logout(): void {
    this.$store.dispatch("logout");
    this.$router.replace("/");
  }
}
</script>

<style lang="scss" scoped>
header {
  width: 100%;
  height: 4rem;
  background-color: #389948;
  display: flex;
  justify-content: center;
  align-items: center;
  a {
    text-decoration: none;
    color: #ffffff;
    display: inline-block;
    padding: 0.75rem 1.5rem;
    border: 1px solid transparent;
  }
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

a {
  :active,
  :hover,
  .router-link-active {
    border: 1px solid #ffffff;
  }
}

h1 {
  margin: 0;
  display: flex;
  align-items: center;
  img {
    height: 3rem;
  }
  a {
    color: white;
    margin: 0;
    :hover,
    :active,
    .router-link-active {
      border-color: transparent;
    }
  }
}

li {
  margin: 0 0.5rem;
}
</style>
