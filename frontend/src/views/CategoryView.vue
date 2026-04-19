<template>
  <main class="max-w-screen-xl mx-auto px-4 sm:px-6 py-10">
    <!-- 面包屑 -->
    <nav class="flex items-center gap-2 text-sm text-ctp-overlay0 mb-8">
      <router-link to="/" class="hover:text-ctp-teal transition-colors">Home</router-link>
      <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
      <span class="text-ctp-subtext">{{ categoryName }}</span>
    </nav>

    <!-- 分类标题 -->
    <div class="mb-8">
      <div class="flex items-center gap-3 mb-2">
        <span class="text-3xl">{{ categoryIcon }}</span>
        <h1 class="text-2xl font-bold text-ctp-text">{{ categoryName }}</h1>
      </div>
      <p class="text-ctp-subtext text-sm ml-12">
        共 {{ posts.length }} 篇文章
      </p>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="flex items-center gap-3 text-ctp-subtext">
        <svg class="w-5 h-5 animate-spin text-ctp-teal" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
        加载中...
      </div>
    </div>

    <!-- 文章列表 -->
    <template v-else>
      <div v-if="posts.length" class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <PostCard v-for="post in sortedPosts" :key="`${post.category}/${post.slug}`" :post="post" />
      </div>
      <div v-else class="text-center py-20">
        <div class="text-4xl mb-4">📭</div>
        <p class="text-ctp-subtext">该分类暂无文章</p>
      </div>
    </template>
  </main>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getPosts } from '../api.js'
import PostCard from '../components/PostCard.vue'

const route = useRoute()
const loading = ref(true)
const posts = ref([])

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
  'system-design': '📐',
  go: '🐹',
  mind: '🌸',
}

const categorySlug = computed(() => route.params.slug)

const categoryName = computed(() => {
  const s = categorySlug.value || ''
  return s.charAt(0).toUpperCase() + s.slice(1).replace(/-/g, ' ')
})

const categoryIcon = computed(() =>
  iconMap[categorySlug.value?.toLowerCase()] || '📄'
)

const sortedPosts = computed(() =>
  [...posts.value].sort((a, b) => new Date(b.date) - new Date(a.date))
)

async function loadPosts() {
  loading.value = true
  try {
    posts.value = await getPosts(categorySlug.value)
  } catch (e) {
    console.error('Failed to load posts:', e)
    posts.value = []
  } finally {
    loading.value = false
  }
}

onMounted(loadPosts)
watch(categorySlug, loadPosts)
</script>
