<template>
  <div class="main-container">
    <div class="header">
      <h1>Packet Analyzer</h1>
      <button class="pause-btn" @click="togglePause" :class="{ 'paused': isPaused }">
        <i :class="isPaused ? 'fas fa-play' : 'fas fa-pause'"></i>
        {{ isPaused ? 'Resume' : 'Pause' }}
      </button>
    </div>
    <div class="content">
      <div v-if="!isConnected" class="disconnected-message">
        <i class="fas fa-wifi-slash"></i>
        <h2>Disconnected</h2>
        <p>Please ensure that both devices and monitor are active in the settings.</p>
      </div>
      <div v-else class="packets-container">
        <div v-for="(packet, index) in packets" :key="index" class="packet-card">
          <div class="packet-header">
            <span class="timestamp">{{ formatTimestamp(packet.timestamp) }}</span>
            <span class="packet-type">{{ packet.ethernet?.eth_type || 'Unknown' }}</span>
          </div>
          
          <div class="packet-details">
            <!-- Ethernet Info -->
            <div v-if="packet.ethernet" class="detail-section">
              <h3>Ethernet</h3>
              <div class="detail-grid">
                <div class="detail-item">
                  <span class="label">Source MAC:</span>
                  <span class="value">{{ packet.ethernet.src_mac }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">Destination MAC:</span>
                  <span class="value">{{ packet.ethernet.dst_mac }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">Type:</span>
                  <span class="value">{{ packet.ethernet.eth_type }}</span>
                </div>
              </div>
            </div>

            <!-- Radiotap Info -->
            <div v-if="packet.radiotap" class="detail-section">
              <h3>RadioTap</h3>
              <div class="detail-grid">
                <div class="detail-item">
                  <span class="label">Channel:</span>
                  <span class="value">{{ packet.radiotap.channel_freq }} MHz</span>
                </div>
                <div class="detail-item">
                  <span class="label">Signal:</span>
                  <span class="value">{{ packet.radiotap.dbm_antenna_sig }} dBm</span>
                </div>
                <div class="detail-item">
                  <span class="label">Data Rate:</span>
                  <span class="value">{{ packet.radiotap.data_rate }} Mbps</span>
                </div>
              </div>
            </div>

            <!-- 802.11 Info -->
            <div v-if="packet.dot11" class="detail-section">
              <h3>802.11</h3>
              <div class="detail-grid">
                <div class="detail-item">
                  <span class="label">Type:</span>
                  <span class="value">{{ packet.dot11.type }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">Subtype:</span>
                  <span class="value">{{ packet.dot11.subtype }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">Address 1:</span>
                  <span class="value">{{ packet.dot11.addr1 }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">Address 2:</span>
                  <span class="value">{{ packet.dot11.addr2 }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">Address 3:</span>
                  <span class="value">{{ packet.dot11.addr3 }}</span>
                </div>
                <div v-if="packet.dot11.addr4" class="detail-item">
                  <span class="label">Address 4:</span>
                  <span class="value">{{ packet.dot11.addr4 }}</span>
                </div>
              </div>
            </div>

            <!-- LLC Info -->
            <div v-if="packet.llc" class="detail-section">
              <h3>LLC</h3>
              <div class="detail-grid">
                <div class="detail-item">
                  <span class="label">DSAP:</span>
                  <span class="value">0x{{ packet.llc.dsap.toString(16) }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">SSAP:</span>
                  <span class="value">0x{{ packet.llc.ssap.toString(16) }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">Control:</span>
                  <span class="value">0x{{ packet.llc.control.toString(16) }}</span>
                </div>
              </div>
            </div>

            <!-- Payload Length -->
            <div class="detail-section">
              <h3>Packet Info</h3>
              <div class="detail-grid">
                <div class="detail-item">
                  <span class="label">Payload Length:</span>
                  <span class="value">{{ packet.payload_length }} bytes</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const packets = ref([])
const ws = ref(null)
const isConnected = ref(false)
const isPaused = ref(false)

const checkStatus = async () => {
  try {
    const response = await fetch('http://localhost:3000/status')
    const status = await response.json()
    isConnected.value = status.devices && status.monitor
    if (isConnected.value && !ws.value) {
      connectWebSocket()
    }
  } catch (error) {
    console.error('Error checking status:', error)
    isConnected.value = false
  }
}

const connectWebSocket = () => {
  if (ws.value) {
    return
  }

  ws.value = new WebSocket('ws://localhost:3000/ws/analizer/raw')

  ws.value.onmessage = (event) => {
    if (!isPaused.value) {
      const packet = JSON.parse(event.data)
      packets.value.unshift(packet)
      
      // Mantener solo los últimos 10 paquetes
      if (packets.value.length > 10) {
        packets.value.pop()
      }
    }
  }

  ws.value.onerror = (error) => {
    console.error('WebSocket error:', error)
    isConnected.value = false
  }

  ws.value.onclose = () => {
    console.log('WebSocket connection closed')
    ws.value = null
    if (!isPaused.value) {
      isConnected.value = false
      setTimeout(checkStatus, 5000)
    }
  }
}

const togglePause = () => {
  isPaused.value = !isPaused.value
}

const formatTimestamp = (timestamp) => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString()
}

onMounted(() => {
  checkStatus()
  setInterval(checkStatus, 5000)
})

onUnmounted(() => {
  if (ws.value) {
    ws.value.close()
  }
})
</script>

<style scoped>
.main-container {
  padding: 2rem;
  background-color: var(--primary-color);
  color: var(--text-color);
  height: 100vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

h1 {
  margin-bottom: 1rem;
  color: var(--accent-color);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.pause-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background-color: var(--accent-color);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.pause-btn:hover {
  background-color: var(--accent-color-hover);
}

.pause-btn.paused {
  background-color: var(--success-color);
}

.pause-btn i {
  font-size: 1rem;
}

.content {
  background-color: var(--secondary-color);
  padding: 1rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
  flex: 1;
  overflow-y: auto;
}

.packets-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.packet-card {
  background-color: var(--primary-color);
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.packet-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--border-color);
}

.timestamp {
  color: var(--text-color-secondary);
  font-size: 0.875rem;
}

.packet-type {
  background-color: var(--accent-color);
  color: white;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
}

.detail-section {
  margin-bottom: 1rem;
}

.detail-section h3 {
  color: var(--accent-color);
  font-size: 1rem;
  margin-bottom: 0.5rem;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 0.5rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.label {
  color: var(--text-color-secondary);
  font-size: 0.875rem;
}

.value {
  font-family: monospace;
  word-break: break-all;
}

.disconnected-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  color: var(--text-color-secondary);
}

.disconnected-message i {
  font-size: 4rem;
  margin-bottom: 1rem;
  color: var(--error-color);
}

.disconnected-message h2 {
  font-size: 2rem;
  margin-bottom: 1rem;
  color: var(--error-color);
}

.disconnected-message p {
  font-size: 1.1rem;
  max-width: 400px;
  line-height: 1.5;
}
</style> 