<template>
  <nav class="navbar">
    <div class="navbar-brand">
      <button class="menu-toggle" @click="$emit('toggle-sidebar')">
        <i class="fas fa-bars"></i>
      </button>
      <h1>IEEE 802 Sniffer</h1>
    </div>
    <div class="navbar-menu">
      <div class="navbar-end">
        <a href="#" class="navbar-item">
          <i class="fas fa-cog"></i>
        </a>
        <a href="#" class="navbar-item">
          <i class="fas fa-question-circle"></i>
        </a>
      </div>
    </div>
  </nav>
</template>

<script setup>
// Component logic here
</script>

<style scoped>
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1rem;
  background-color: var(--secondary-color);
  color: var(--text-color);
  height: 60px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.menu-toggle {
  background: none;
  border: none;
  color: var(--text-color);
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0.5rem;
}

.menu-toggle:hover {
  background-color: var(--hover-color);
  border-radius: 4px;
}

.navbar-menu {
  display: flex;
  align-items: center;
}

.navbar-end {
  display: flex;
  gap: 1rem;
}

.navbar-item {
  color: var(--text-color);
  text-decoration: none;
  padding: 0.5rem;
  border-radius: 4px;
}

.navbar-item:hover {
  background-color: var(--hover-color);
}

h1 {
  font-size: 1.5rem;
  margin: 0;
  color: var(--text-color);
}
</style> 