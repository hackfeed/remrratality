<template>
  <div class="card">
    <slot></slot>
  </div>
</template>

<style lang="scss" scoped>
.card {
  border-radius: 12px;
  box-shadow: 0 2px 8px $color-border-primary;
  padding: 1rem;
  margin: 2rem auto;
  max-width: 40rem;
}
</style>
