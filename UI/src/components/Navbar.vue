<template>
  <nav class="navbar">
    <div class="navbar-brand">
      <button class="menu-toggle" @click="$emit('toggle-sidebar')">
        <i class="fas fa-bars"></i>
      </button>
      <h1>IEEE 802 Sniffer</h1>
    </div>
    <div class="navbar-menu">
      <div class="status-indicators">
        <div class="status-item" :class="{ 'active': backendStatus.devices }">
          <i class="fas fa-network-wired"></i>
          <span>{{ backendStatus.devices ? 'Online' : 'Offline' }}</span>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted } from 'vue'
const backendStatus = ref({
  devices: false,
  monitor: false
})

const fetchStatus = async () => {
  try {
    const response = await fetch('http://localhost:3000/status')
    const data = await response.json()
    backendStatus.value = data
  } catch (error) {
    console.error('Error fetching backend status:', error)
  }
}

onMounted(() => {
  fetchStatus()
  // Actualizar el estado cada 5 segundos
  setInterval(fetchStatus, 5000)
})
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

.status-indicators {
  display: flex;
  gap: 1rem;
  margin-right: 1rem;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  background-color: var(--error-color, #ff4444);
  color: white;
  transition: all 0.3s ease;
  font-weight: 500;
}

.status-item.active {
  background-color: var(--success-color, #00C851);
}

.status-item i {
  font-size: 1rem;
}

.status-item span {
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
</style> 