<script setup>
import { APP_TITLE } from '@/config/constants'
const year = new Date().getFullYear()
</script>

<template>
  <footer
    class="sidebar-root h-14 flex items-center justify-center text-xs border-t shrink-0"
    
  >
    &copy; {{ year }} {{ APP_TITLE }}
  </footer>
</template>

<style scoped>
.sidebar-root {
  background: var(--sidebar-bg);
  color: var(--sidebar-text);
}
</style>