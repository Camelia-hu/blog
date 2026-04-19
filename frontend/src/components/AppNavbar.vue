<template>
  <!-- 顶部导航栏 -->
  <header class="sticky top-0 z-50">
    <!-- 主导航栏 -->
    <div class="bg-ctp-mantle border-b border-ctp-surface0">
      <div class="max-w-screen-xl mx-auto px-4 sm:px-6">
        <div class="flex items-center justify-between h-14">
          <!-- 品牌名 -->
          <router-link to="/" class="flex items-center gap-2 group">
            <div class="w-8 h-8 rounded-lg bg-ctp-teal/20 flex items-center justify-center text-ctp-teal text-base font-bold group-hover:bg-ctp-teal/30 transition-colors">
              M
            </div>
            <span class="text-ctp-text font-semibold text-base tracking-wide">camelia</span>
          </router-link>

          <!-- 右侧图标 -->
          <div class="flex items-center gap-3">
            <!-- 搜索图标（装饰） -->
            <button class="text-ctp-overlay0 hover:text-ctp-text transition-colors p-1.5 rounded-md hover:bg-ctp-surface0" title="搜索">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </button>
            <!-- GitHub 链接 -->
            <a href="https://github.com" target="_blank" rel="noopener noreferrer"
              class="text-ctp-overlay0 hover:text-ctp-text transition-colors p-1.5 rounded-md hover:bg-ctp-surface0" title="GitHub">
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z" />
              </svg>
            </a>
          </div>
        </div>
      </div>
    </div>

    <!-- 分类 Tab 栏（MkDocs Material 风格）-->
    <div class="bg-ctp-crust/80 backdrop-blur-sm border-b border-ctp-surface0/50">
      <div class="max-w-screen-xl mx-auto px-4 sm:px-6">
        <nav class="flex items-center gap-1 overflow-x-auto scrollbar-none">
          <router-link
            to="/"
            class="nav-tab"
            :class="{ 'nav-tab-active': $route.path === '/' }"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
            </svg>
            Home
          </router-link>

          <router-link
            v-for="cat in categories"
            :key="cat.slug"
            :to="`/category/${cat.slug}`"
            class="nav-tab"
            :class="{ 'nav-tab-active': $route.params.slug === cat.slug }"
          >
            {{ cat.name }}
          </router-link>
        </nav>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCategories } from '../api.js'

const categories = ref([])

onMounted(async () => {
  try {
    categories.value = await getCategories()
  } catch (e) {
    console.error('Failed to load categories:', e)
  }
})
</script>

<style scoped>
.nav-tab {
  @apply flex items-center gap-1.5 px-3 py-3 text-sm text-ctp-subtext whitespace-nowrap
    hover:text-ctp-text transition-colors border-b-2 border-transparent
    hover:border-ctp-surface1 -mb-px;
}

.nav-tab-active {
  @apply text-ctp-teal border-ctp-teal;
}
</style>
