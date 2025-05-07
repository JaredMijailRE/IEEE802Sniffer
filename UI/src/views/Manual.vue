<template>
  <div class="manual-container">
    <h1>Manual</h1>
    <div class="content">
      <section>
        <h2>Getting Started</h2>
        <p>Getting started guide will be added here.</p>
      </section>
      <section>
        <h2>Features</h2>
        <p>Features documentation will be added here.</p>
      </section>
      <section>
        <h2>Usage</h2>
        <p>Usage instructions will be added here.</p>
      </section>
      <section>
        <h2>Troubleshooting</h2>
        <p>Troubleshooting guide will be added here.</p>
      </section>
    </div>
  </div>
</template>

<style scoped>
.manual-container {
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

section {
  margin-bottom: 2rem;
}

h2 {
  color: var(--accent-color);
  margin-bottom: 0.5rem;
}
</style> 