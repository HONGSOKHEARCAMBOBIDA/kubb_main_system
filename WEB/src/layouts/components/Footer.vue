<script setup>
import { APP_TITLE } from '@/config/constants'
const year = new Date().getFullYear()
</script>

<template>
  <footer
    class="h-14 flex items-center justify-center text-xs border-t shrink-0"
    style="color: var(--color-text-secondary); border-color: var(--color-border)"
  >
    &copy; {{ year }} {{ APP_TITLE }}
  </footer>
</template>
