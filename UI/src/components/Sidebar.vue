<template>
  <aside class="sidebar" :class="{ 'is-open': isOpen }">
    <nav class="sidebar-menu">
      <router-link to="/" class="sidebar-item">
        <i class="fas fa-home"></i>
        <span>Main</span>
      </router-link>
      <router-link to="/settings" class="sidebar-item">
        <i class="fas fa-cog"></i>
        <span>Settings</span>
      </router-link>
      <router-link to="/manual" class="sidebar-item">
        <i class="fas fa-book"></i>
        <span>Manual</span>
      </router-link>
      <router-link to="/about" class="sidebar-item">
        <i class="fas fa-info-circle"></i>
        <span>About</span>
      </router-link>
    </nav>
  </aside>
</template>

<script setup>
defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.sidebar {
  position: fixed;
  top: 60px;
  left: 0;
  bottom: 0;
  width: 250px;
  background-color: var(--secondary-color);
  color: var(--text-color);
  transform: translateX(-100%);
  transition: transform 0.3s ease;
  z-index: 900;
  box-shadow: 2px 0 4px rgba(0,0,0,0.2);
}

.sidebar.is-open {
  transform: translateX(0);
}

.sidebar-menu {
  padding: 1rem 0;
}

.sidebar-item {
  display: flex;
  align-items: center;
  padding: 0.75rem 1.5rem;
  color: var(--text-color);
  text-decoration: none;
  transition: background-color 0.2s ease;
}

.sidebar-item:hover {
  background-color: var(--hover-color);
}

.sidebar-item i {
  margin-right: 0.75rem;
  width: 20px;
  text-align: center;
}

.router-link-active {
  background-color: var(--hover-color);
  border-left: 4px solid var(--accent-color);
}
</style> 