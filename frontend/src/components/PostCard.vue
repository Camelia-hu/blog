<template>
  <router-link
    :to="`/post/${post.category}/${post.slug}`"
    class="post-card block bg-ctp-surface0/30 hover:bg-ctp-surface0/50 border border-ctp-surface1/50
      hover:border-ctp-teal/30 rounded-xl p-5 hover:shadow-lg hover:shadow-ctp-crust/50
      group transition-all duration-200"
  >
    <!-- 分类标签 -->
    <div class="flex items-center justify-between mb-3">
      <span class="text-xs font-medium px-2 py-0.5 rounded-full bg-ctp-teal/10 text-ctp-teal border border-ctp-teal/20">
        {{ categoryDisplay }}
      </span>
      <span class="text-xs text-ctp-overlay0 flex items-center gap-1">
        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {{ post.readTime }} 分钟阅读
      </span>
    </div>

    <!-- 标题 -->
    <h3 class="text-ctp-text font-semibold text-base leading-snug mb-2
      group-hover:text-ctp-teal transition-colors line-clamp-2">
      {{ post.title }}
    </h3>

    <!-- 摘要 -->
    <p v-if="post.excerpt" class="text-ctp-subtext text-sm leading-relaxed mb-3 line-clamp-2">
      {{ post.excerpt }}
    </p>

    <!-- 底部：日期 + 标签 -->
    <div class="flex items-center justify-between mt-3 pt-3 border-t border-ctp-surface1/30">
      <span class="text-xs text-ctp-overlay0">{{ post.date }}</span>
      <div class="flex items-center gap-1 flex-wrap justify-end">
        <span
          v-for="tag in (post.tags || []).slice(0, 3)"
          :key="tag"
          class="text-xs px-1.5 py-0.5 rounded bg-ctp-surface1/50 text-ctp-overlay1"
        >
          {{ tag }}
        </span>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  post: {
    type: Object,
    required: true
  }
})

const categoryDisplay = computed(() => {
  const s = props.post.category || ''
  return s.charAt(0).toUpperCase() + s.slice(1)
})
</script>
