<template>
  <div class="min-h-screen bg-ctp-base">
    <AppNavbar />
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<script setup>
import AppNavbar from './components/AppNavbar.vue'
</script>
