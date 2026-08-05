<script setup>
// Collapsible sidebar with nested (2-level) menu support, driven by
// src/config/menu.js. Works both as a fixed sidebar (desktop) and inside
// an el-drawer (mobile) — the `mobile` prop disables the collapse behavior.
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { menuConfig } from '@/config/menu'
import { APP_TITLE } from '@/config/constants'

const props = defineProps({
  mobile: { type: Boolean, default: false },
})
const emit = defineEmits(['navigate'])

const appStore = useAppStore()
const authStore = useAuthStore()
const route = useRoute()

const collapsed = computed(() => !props.mobile && appStore.sidebarCollapsed)
const activePath = computed(() => route.path)

const visibleMenu = computed(() =>
  menuConfig.filter((item) => authStore.hasPermission(item.permission))
)

function onSelect() {
  if (props.mobile) emit('navigate')
}
</script>

<template>
  <aside
    class="sidebar-root fixed top-0 left-0 h-screen z-40 flex flex-col transition-all duration-200"
    :class="{ 'sidebar-mobile': mobile }"
    :style="{ width: mobile ? '260px' : collapsed ? 'var(--sidebar-width-collapsed)' : 'var(--sidebar-width)' }"
  >
    <!-- Logo -->
    <div class="flex items-center gap-2 h-16 px-4 shrink-0">
      <span
        v-s how="!collapsed"
        class="text-white font-semibold text-[15px] whitespace-nowrap overflow-hidden"
      >
        {{ APP_TITLE }}
      </span>
    </div>

    <!-- Menu -->
    <el-scrollbar class="flex-1">
      <el-menu
        :default-active="activePath"
        :collapse="collapsed"
        :collapse-transition="false"
        background-color="transparent"
        text-color="var(--sidebar-text)"
        active-text-color="var(--sidebar-text-active)"
        class="sidebar-menu border-none"
        router
        @select="onSelect"
      >
        <template v-for="item in visibleMenu" :key="item.title">
          <el-sub-menu v-if="item.children?.length" :index="item.title">
            <template #title>
              <el-icon><component :is="item.icon" /></el-icon>
              <span>{{ item.title }}</span>
            </template>
            <el-menu-item
              v-for="child in item.children"
              :key="child.path"
              :index="child.path"
            >
              <el-icon><component :is="child.icon" /></el-icon>
              <template #title>{{ child.title }}</template>
            </el-menu-item>
          </el-sub-menu>

          <el-menu-item v-else :index="item.path">
            <el-icon><component :is="item.icon" /></el-icon>
            <template #title>{{ item.title }}</template>
          </el-menu-item>
        </template>
      </el-menu>
    </el-scrollbar>

    <!-- Collapse toggle (desktop only) -->
    <div
      v-if="!mobile"
      class="h-12 flex items-center justify-center border-t cursor-pointer shrink-0"
      style="border-color: rgba(255, 255, 255, 0.08)"
      @click="appStore.toggleSidebar()"
    >
      <el-icon class="text-white/60 hover:text-white transition-colors">
        <component :is="collapsed ? 'Expand' : 'Fold'" />
      </el-icon>
    </div>
  </aside>
</template>

<style scoped>
.sidebar-root {
  background: var(--sidebar-bg);
}
.sidebar-mobile {
  position: static;
  height: 100%;
}

.sidebar-menu {
  --el-menu-item-height: 44px;
}
:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  margin: 2px 10px;
  border-radius: 8px;
  width: calc(100% - 20px);
}
:deep(.el-menu-item.is-active) {
  background: rgba(255, 255, 255, 0.08) !important;
  font-weight: 600;
}
:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background: rgba(255, 255, 255, 0.06) !important;
}
</style>
