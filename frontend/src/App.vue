<script setup lang="ts">
import { RouterView } from "vue-router";
import Toast from "primevue/toast";
</script>

<template>
  <Toast />
  <RouterView />
</template>

<style>
/* Global reset or minimal styles if needed */
body {
  margin: 0;
  font-family: var(--font-family, sans-serif);
}
</style>
