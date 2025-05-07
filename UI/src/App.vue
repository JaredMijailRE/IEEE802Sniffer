<script setup>
import { ref } from 'vue'
import { RouterView } from 'vue-router'
import Navbar from './components/Navbar.vue'
import Sidebar from './components/Sidebar.vue'

const isSidebarOpen = ref(false)

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value
}
</script>

<template>
  <div class="app">
    <Navbar @toggle-sidebar="toggleSidebar" />
    <Sidebar :is-open="isSidebarOpen" />
    <main class="main-content" :class="{ 'sidebar-open': isSidebarOpen }">
      <RouterView />
    </main>
  </div>
</template>

<style>
:root {
  --primary-color: #1a1a1a;
  --secondary-color: #2d2d2d;
  --accent-color: #4a9eff;
  --text-color: #ffffff;
  --hover-color: rgba(255, 255, 255, 0.1);
  --navbar-height: 60px;
  --sidebar-width: 250px;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  width: 100%;
  overflow-x: hidden;
  font-family: Arial, sans-serif;
  background-color: var(--primary-color);
  color: var(--text-color);
}

.app {
  min-height: 100vh;
  width: 100vw;
  position: relative;
}

.main-content {
  min-height: calc(100vh - var(--navbar-height));
  margin-top: var(--navbar-height);
  padding: 2rem;
  transition: margin-left 0.3s ease;
  background-color: var(--primary-color);
}

.main-content.sidebar-open {
  margin-left: var(--sidebar-width);
}

@media (max-width: 768px) {
  .main-content.sidebar-open {
    margin-left: 0;
  }
}

/* Import Font Awesome */
@import '@fortawesome/fontawesome-free/css/all.min.css';
</style>
