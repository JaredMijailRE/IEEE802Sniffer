<template>
  <div class="main-container">
    <h1>Main Page</h1>
    <div class="content">
      <p>Welcome to the IEEE 802 Sniffer application.</p>
    </div>
  </div>
</template>

<script setup>
// Component logic here
</script>

<style scoped>
.main-container {
  padding: 2rem;
  background-color: var(--primary-color);
  color: var(--text-color);
}

h1 {
  margin-bottom: 1rem;
  color: var(--accent-color);
}

.content {
  background-color: var(--secondary-color);
  padding: 1rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
}
</style> 