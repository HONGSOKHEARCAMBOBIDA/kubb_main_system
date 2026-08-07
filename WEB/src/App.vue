<script setup>
// App.vue is intentionally minimal: it only hosts the router view and the
// single global overlay (loading). All layout logic lives in
// src/layouts/*.vue so pages can opt into AuthLayout or DashboardLayout.
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app'
import LoadingOverlay from '@/components/loading/LoadingOverlay.vue'

const appStore = useAppStore()

onMounted(() => {
  // Apply persisted theme (light/dark) as soon as the app boots to avoid
  // a flash of the wrong theme.
  appStore.applyTheme()
})
</script>

<template>
  <router-view v-slot="{ Component, route }">
    <transition name="app-fade" mode="out-in">
      <component :is="Component" :key="route.fullPath" />
    </transition>
    
  </router-view>

  <!-- Global loading overlay, toggled by the loading store / axios interceptors -->
  <LoadingOverlay />
</template>
