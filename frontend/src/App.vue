<template>
  <the-header></the-header>
  <router-view v-slot="slotProps">
    <transition name="route" mode="out-in">
      <component :is="slotProps.Component"></component>
    </transition>
  </router-view>
  <div id="upfooter"></div>
  <the-footer></the-footer>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";

import TheHeader from "@/components/layout/TheHeader.vue";
import TheFooter from "@/components/layout/TheFooter.vue";

@Options({
  components: { TheHeader, TheFooter },
  watch: {
    didAutoLogout(curValue, oldValue) {
      if (curValue && curValue != oldValue) {
        this.$router.replace("/");
      }
    },
  },
})
export default class App extends Vue {
  get didAutoLogout(): boolean {
    return this.$store.getters.didAutoLogout;
  }

  created(): void {
    this.$store.dispatch("tryLogin");
  }
}
</script>

<style lang="scss">
@import url("https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap");
* {
  box-sizing: border-box;
}

html {
  font-family: "Roboto", sans-serif;
}

body {
  margin: 0;
  position: relative;
  min-height: 100vh;
}

#upfooter {
  padding-bottom: 3rem;
}

.route-enter-from {
  opacity: 0;
  transform: translateY(-30px);
}

.route-leave-to {
  opacity: 0;
  transform: translateY(30px);
}

.route-enter-active {
  transition: all 0.3s ease-out;
}

.route-leave-active {
  transition: all 0.3s ease-out;
}

.route-enter-to,
.route-leave-from {
  opacity: 1;
  transform: translateY(0);
}
</style>
