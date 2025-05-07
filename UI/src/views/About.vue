<template>
  <div class="about-container">
    <h1>About</h1>
    <div class="content">
      <section>
        <h2>Project Information</h2>
        <p>Information about the IEEE 802 Sniffer project will be added here.</p>
      </section>
      <section>
        <h2>Version</h2>
        <p>Version information will be added here.</p>
      </section>
      <section>
        <h2>Contributors</h2>
        <p>Contributor information will be added here.</p>
      </section>
      <section>
        <h2>License</h2>
        <p>License information will be added here.</p>
      </section>
    </div>
  </div>
</template>

<style scoped>
.about-container {
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