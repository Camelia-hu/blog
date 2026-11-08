<template>
  <main class="max-w-screen-xl mx-auto px-4 sm:px-6 py-10">
    <!-- Hero 区域 -->
    <section class="text-center py-12 mb-12">
      <div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-ctp-teal/30 to-ctp-blue/30
        border border-ctp-teal/30 flex items-center justify-center mx-auto mb-4
        text-2xl font-bold text-ctp-teal">
        M
      </div>
      <h1 class="text-3xl sm:text-4xl font-bold text-ctp-text mb-3">camelia</h1>
      <p class="text-ctp-subtext text-base max-w-md mx-auto leading-relaxed">
        打怪升级的记录吧 — 记录技术与生活的点点滴滴
      </p>
      <div class="flex items-center justify-center gap-2 mt-5">
        <span class="inline-flex items-center gap-1.5 text-xs text-ctp-overlay0 bg-ctp-surface0/50
          px-3 py-1.5 rounded-full border border-ctp-surface1/50">
          <span class="w-1.5 h-1.5 bg-ctp-green rounded-full animate-pulse"></span>
          {{ totalPosts }} 篇文章
        </span>
        <span class="inline-flex items-center gap-1.5 text-xs text-ctp-overlay0 bg-ctp-surface0/50
          px-3 py-1.5 rounded-full border border-ctp-surface1/50">
          {{ categories.length }} 个分类
        </span>
      </div>
    </section>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center py-16">
      <div class="flex items-center gap-3 text-ctp-subtext">
        <svg class="w-5 h-5 animate-spin text-ctp-teal" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
        加载中...
      </div>
    </div>

    <template v-else>
      <!-- 最新文章 -->
      <section class="mb-12">
        <div class="flex items-center gap-3 mb-6">
          <div class="w-1 h-5 bg-ctp-teal rounded-full"></div>
          <h2 class="text-lg font-semibold text-ctp-text">最新文章</h2>
        </div>
        <div v-if="recentPosts.length" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <PostCard v-for="post in recentPosts" :key="`${post.category}/${post.slug}`" :post="post" />
        </div>
        <p v-else class="text-ctp-subtext text-sm">暂无文章</p>
      </section>

      <!-- 分类导航 -->
      <section>
        <div class="flex items-center gap-3 mb-6">
          <div class="w-1 h-5 bg-ctp-blue rounded-full"></div>
          <h2 class="text-lg font-semibold text-ctp-text">分类</h2>
        </div>
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
          <router-link
            v-for="cat in categories"
            :key="cat.slug"
            :to="`/category/${cat.slug}`"
            class="group bg-ctp-surface0/30 hover:bg-ctp-surface0/60 border border-ctp-surface1/40
              hover:border-ctp-teal/30 rounded-xl p-4 transition-all duration-200
              hover:shadow-md hover:shadow-ctp-crust/40 hover:-translate-y-0.5"
          >
            <div class="flex items-start justify-between mb-2">
              <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-ctp-teal/20 to-ctp-blue/20
                border border-ctp-teal/20 flex items-center justify-center text-sm">
                {{ categoryIcon(cat.slug) }}
              </div>
              <span class="text-xs text-ctp-overlay0 bg-ctp-surface1/50 px-2 py-0.5 rounded-full">
                {{ cat.count }}
              </span>
            </div>
            <h3 class="text-sm font-medium text-ctp-subtext group-hover:text-ctp-teal transition-colors">
              {{ cat.name }}
            </h3>
          </router-link>
        </div>
      </section>
    </template>
  </main>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getCategories, getPosts } from '../api.js'
import PostCard from '../components/PostCard.vue'

const loading = ref(true)
const categories = ref([])
const posts = ref([])

const recentPosts = computed(() =>
  [...posts.value]
    .sort((a, b) => new Date(b.date) - new Date(a.date))
    .slice(0, 6)
)

const totalPosts = computed(() => posts.value.length)

const iconMap = {
  algorithms: '🧮',
  algorithm: '🧮',
  database: '🗄️',
  network: '🌐',
  os: '⚙️',
  architecture: '🏗️',
  tools: '🔧',
  projects: '🚀',
  blog: '📝',
  'system design': '📐',
  'system-design': '📐',
  go: '🐹',
  mind: '🌸',
  ai: '🤖',
}

function categoryIcon(slug) {
  return iconMap[slug.toLowerCase()] || '📄'
}

onMounted(async () => {
  try {
    const [cats, ps] = await Promise.all([getCategories(), getPosts()])
    categories.value = cats
    posts.value = ps
  } catch (e) {
    console.error('Failed to load data:', e)
  } finally {
    loading.value = false
  }
})
</script>
