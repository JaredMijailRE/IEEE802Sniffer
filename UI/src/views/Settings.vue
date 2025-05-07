<template>
  <div class="settings-container">
    <h1>Settings</h1>
    <div class="content">
      <div class="devices-section">
        <h2>Network Devices</h2>
        <div class="devices-list" v-if="devices.length > 0">
          <div v-for="(deviceName, index) in devices" :key="index" 
               class="device-item" 
               :class="{ 'selected': selectedDevice === index }"
               @click="selectDevice(index)">
            <div class="device-info">
              <i class="fas fa-network-wired"></i>
              <span class="device-name">{{ deviceName }}</span>
            </div>
            <button class="monitor-btn" 
                    @click.stop="setMonitor(index)"
                    :disabled="selectedDevice === index && isMonitoring">
              {{ selectedDevice === index && isMonitoring ? 'Monitoring' : 'Monitor' }}
            </button>
          </div>
        </div>
        <div v-else class="no-devices">
          <i class="fas fa-exclamation-circle"></i>
          <p>No network devices found</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const devices = ref([])
const selectedDevice = ref(null)
const isMonitoring = ref(false)

const fetchDevices = async () => {
  try {
    const response = await fetch('http://localhost:3000/devices')
    devices.value = await response.json()
  } catch (error) {
    console.error('Error fetching devices:', error)
  }
}

const selectDevice = (index) => {
  selectedDevice.value = index
}

const setMonitor = async (index) => {
  try {
    const response = await fetch(`http://localhost:3000/monitor/${index}`, {
      method: 'POST'
    })
    if (response.ok) {
      isMonitoring.value = true
      selectedDevice.value = index
    }
  } catch (error) {
    console.error('Error setting monitor:', error)
  }
}

onMounted(() => {
  fetchDevices()
})
</script>

<style scoped>
.settings-container {
  padding: 2rem;
  background-color: var(--primary-color);
  color: var(--text-color);
}

h1 {
  margin-bottom: 1rem;
  color: var(--accent-color);
}

h2 {
  color: var(--accent-color);
  margin-bottom: 1rem;
  font-size: 1.25rem;
}

.content {
  background-color: var(--secondary-color);
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
}

.devices-section {
  margin-top: 1rem;
}

.devices-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.device-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background-color: var(--primary-color);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.device-item:hover {
  transform: translateX(5px);
  background-color: var(--hover-color);
}

.device-item.selected {
  border-left: 4px solid var(--accent-color);
}

.device-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.device-name {
  font-weight: 500;
  color: var(--text-color);
}

.monitor-btn {
  padding: 0.5rem 1rem;
  background-color: var(--accent-color);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.monitor-btn:hover:not(:disabled) {
  background-color: var(--accent-color-hover);
}

.monitor-btn:disabled {
  background-color: var(--disabled-color);
  cursor: not-allowed;
}

.no-devices {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2rem;
  color: var(--text-color-secondary);
}

.no-devices i {
  font-size: 2rem;
}
</style> 