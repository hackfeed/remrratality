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
  background-color: #389948;
  border: 1px solid #389948;
  color: white;
  cursor: pointer;
  border-radius: 30px;
  margin-right: 0.5rem;
  display: inline-block;
  &:hover,
  &:active {
    background-color: #24642f;
    border-color: #24642f;
  }
}

.flat {
  background-color: transparent;
  color: #389948;
  border: none;
  &:hover,
  &:active {
    background-color: #70b87c;
    color: white;
  }
}

.outline {
  background-color: transparent;
  border-color: #24642f;
  color: #24642f;
  &:hover,
  &:active {
    background-color: #70b87c;
    color: white;
  }
}
</style>
