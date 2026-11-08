<template>
  <div class="max-w-screen-xl mx-auto px-4 sm:px-6 py-10">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center py-32">
      <div class="flex items-center gap-3 text-ctp-subtext">
        <svg class="w-5 h-5 animate-spin text-ctp-teal" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
        加载中...
      </div>
    </div>

    <!-- 404 -->
    <div v-else-if="error" class="text-center py-32">
      <div class="text-6xl mb-4">😿</div>
      <h2 class="text-xl font-semibold text-ctp-text mb-2">文章未找到</h2>
      <p class="text-ctp-subtext mb-6">{{ error }}</p>
      <router-link to="/" class="text-ctp-teal hover:underline text-sm">← 返回首页</router-link>
    </div>

    <!-- 文章内容 -->
    <template v-else-if="post">
      <!-- 面包屑 -->
      <nav class="flex items-center gap-2 text-sm text-ctp-overlay0 mb-8">
        <router-link to="/" class="hover:text-ctp-teal transition-colors">Home</router-link>
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
        <router-link :to="`/category/${post.category}`" class="hover:text-ctp-teal transition-colors">
          {{ categoryDisplay }}
        </router-link>
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
        <span class="text-ctp-subtext truncate max-w-xs">{{ post.title }}</span>
      </nav>

      <div class="flex gap-8 items-start">
        <!-- 主内容区 -->
        <article class="flex-1 min-w-0">
          <!-- 文章头部 -->
          <header class="mb-8 pb-6 border-b border-ctp-surface1/40">
            <h1 class="text-2xl sm:text-3xl font-bold text-ctp-text leading-tight mb-4">
              {{ post.title }}
            </h1>
            <div class="flex flex-wrap items-center gap-3 text-sm text-ctp-overlay0">
              <!-- 日期 -->
              <span class="flex items-center gap-1.5">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                {{ post.date }}
              </span>
              <!-- 阅读时间 -->
              <span class="flex items-center gap-1.5">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ post.readTime }} 分钟阅读
              </span>
              <!-- 作者 -->
              <span v-if="post.author" class="flex items-center gap-1.5">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
                {{ post.author }}
              </span>
            </div>
            <!-- 标签 -->
            <div v-if="post.tags?.length" class="flex flex-wrap gap-2 mt-4">
              <span
                v-for="tag in post.tags"
                :key="tag"
                class="text-xs px-2.5 py-1 rounded-full bg-ctp-surface0/80
                  border border-ctp-surface1/60 text-ctp-subtext hover:border-ctp-teal/40
                  hover:text-ctp-teal transition-colors cursor-default"
              >
                # {{ tag }}
              </span>
            </div>
          </header>

          <!-- Markdown 内容 -->
          <div
            ref="contentRef"
            class="prose prose-sm sm:prose max-w-none"
            v-html="post.content"
          ></div>

          <!-- 底部导航 -->
          <div class="mt-12 pt-6 border-t border-ctp-surface1/40">
            <router-link
              :to="`/category/${post.category}`"
              class="inline-flex items-center gap-2 text-sm text-ctp-teal hover:text-ctp-sky
                transition-colors group"
            >
              <svg class="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
              查看更多 {{ categoryDisplay }} 文章
            </router-link>
          </div>
        </article>

        <!-- TOC 侧边栏（桌面端） -->
        <aside v-if="toc.length" class="hidden xl:block w-56 flex-shrink-0">
          <div class="sticky top-32">
            <h4 class="text-xs font-semibold text-ctp-overlay0 uppercase tracking-wider mb-3 px-3">
              目录
            </h4>
            <nav class="space-y-0.5">
              <a
                v-for="item in toc"
                :key="item.id"
                :href="`#${item.id}`"
                class="toc-link block text-sm py-1 px-3 border-l-2 border-ctp-surface1/50
                  text-ctp-overlay0 hover:text-ctp-text hover:border-ctp-surface2
                  transition-colors truncate"
                :class="{
                  'pl-5 text-xs': item.level === 3,
                  'active': activeId === item.id
                }"
                @click.prevent="scrollTo(item.id)"
              >
                {{ item.text }}
              </a>
            </nav>
          </div>
        </aside>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getPost } from '../api.js'
import hljs from 'highlight.js'

const route = useRoute()
const loading = ref(true)
const error = ref(null)
const post = ref(null)
const contentRef = ref(null)
const toc = ref([])
const activeId = ref('')

const categoryDisplay = computed(() => {
  const s = post.value?.category || ''
  return s.charAt(0).toUpperCase() + s.slice(1)
})

// 从渲染后的 HTML 中提取 TOC
function buildToc() {
  if (!contentRef.value) return
  const headings = contentRef.value.querySelectorAll('h2, h3')
  toc.value = Array.from(headings).map(h => ({
    id: h.id || h.textContent.toLowerCase().replace(/\s+/g, '-').replace(/[^\w-]/g, ''),
    text: h.textContent,
    level: parseInt(h.tagName[1]),
  }))
  // 确保每个 heading 有 id
  headings.forEach((h, i) => {
    if (!h.id) h.id = toc.value[i].id
  })
}

// 代码高亮
function highlightCode() {
  if (!contentRef.value) return
  contentRef.value.querySelectorAll('pre code').forEach(block => {
    hljs.highlightElement(block)
  })
}

// 滚动到指定锚点
function scrollTo(id) {
  const el = document.getElementById(id)
  if (el) {
    const y = el.getBoundingClientRect().top + window.scrollY - 100
    window.scrollTo({ top: y, behavior: 'smooth' })
  }
}

// 监听滚动，更新激活的 TOC 项
let observer = null
function setupObserver() {
  if (!contentRef.value || !toc.value.length) return
  observer?.disconnect()
  observer = new IntersectionObserver(entries => {
    const visible = entries.filter(e => e.isIntersecting)
    if (visible.length) {
      activeId.value = visible[0].target.id
    }
  }, { rootMargin: '-80px 0px -60% 0px', threshold: 0 })

  toc.value.forEach(item => {
    const el = document.getElementById(item.id)
    if (el) observer.observe(el)
  })
}

async function loadPost() {
  loading.value = true
  error.value = null
  post.value = null
  try {
    post.value = await getPost(route.params.category, route.params.slug)
    await nextTick()
    buildToc()
    highlightCode()
    setupObserver()
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

onMounted(loadPost)
watch(() => route.params, loadPost, { deep: true })
onUnmounted(() => observer?.disconnect())
</script>
