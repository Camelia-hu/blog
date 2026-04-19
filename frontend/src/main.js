import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import './style.css'

// 路由配置
import HomeView from './views/HomeView.vue'
import CategoryView from './views/CategoryView.vue'
import PostView from './views/PostView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomeView },
    { path: '/category/:slug', component: CategoryView },
    { path: '/post/:category/:slug', component: PostView },
  ],
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    }
    return savedPosition || { top: 0 }
  },
})

createApp(App).use(router).mount('#app')
