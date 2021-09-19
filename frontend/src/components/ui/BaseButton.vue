<template>
  <button v-if="!link" :class="mode">
    <slot></slot>
  </button>
  <router-link v-else :to="to" :class="mode">
    <slot></slot>
  </router-link>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";

@Options({
  props: {
    mode: {
      type: String,
      required: false,
      default: null,
    },
    link: {
      type: Boolean,
      required: false,
      default: false,
    },
    to: {
      type: String,
      required: false,
      default: "/",
    },
  },
})
export default class BaseButton extends Vue {
  readonly mode!: string;
  readonly link!: boolean;
  readonly to!: string;
}
</script>

<style lang="scss" scoped>
button,
a {
  text-decoration: none;
  padding: 0.75rem 1.5rem;
  font: inherit;
  background-color: $color-primary;
  border: 1px solid $color-primary;
  color: $color-text-secondary;
  cursor: pointer;
  border-radius: 30px;
  margin-right: 0.5rem;
  display: inline-block;
  &:hover,
  &:active {
    background-color: $color-primary-dark;
    border-color: $color-primary-dark;
  }
}

.flat {
  background-color: transparent;
  color: $color-primary;
  border: none;
  &:hover,
  &:active {
    background-color: $color-primary-light;
    color: $color-text-secondary;
  }
}

.outline {
  background-color: transparent;
  border-color: $color-primary-dark;
  color: $color-primary-dark;
  &:hover,
  &:active {
    background-color: $color-primary-light;
    color: $color-text-secondary;
  }
}
</style>
