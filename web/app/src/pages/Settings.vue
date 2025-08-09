<template>
  <div>
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-base-content">Settings</h1>
      <p class="text-base-content/70">Application settings and system information</p>
    </div>

    <div class="space-y-6">
      <!-- System Information -->
      <div class="card bg-base-200 shadow">
        <div class="card-body">
          <h2 class="card-title mb-4">System Information</h2>
          <div class="grid gap-4 md:grid-cols-2">
            <div>
              <div class="text-sm font-semibold text-base-content/70">Application</div>
              <div>Go Dumper v1.0.0</div>
            </div>
            <div>
              <div class="text-sm font-semibold text-base-content/70">Status</div>
              <div class="flex items-center gap-2">
                <div class="w-2 h-2 bg-success rounded-full"></div>
                <span>Running</span>
              </div>
            </div>
            <div>
              <div class="text-sm font-semibold text-base-content/70">Backend</div>
              <div>Go with Gin Framework</div>
            </div>
            <div>
              <div class="text-sm font-semibold text-base-content/70">Frontend</div>
              <div>Vue 3 with TypeScript</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Theme Settings -->
      <div class="card bg-base-200 shadow">
        <div class="card-body">
          <h2 class="card-title mb-4">Theme Settings</h2>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Theme</span>
            </label>
            <select v-model="selectedTheme" @change="changeTheme" class="select select-bordered">
              <option value="light">Light</option>
              <option value="dark">Dark</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Configuration Help -->
      <div class="card bg-base-200 shadow">
        <div class="card-body">
          <h2 class="card-title mb-4">Configuration Help</h2>
          <div class="space-y-4">
            <div>
              <h3 class="font-semibold mb-2">Environment Variables</h3>
              <div class="text-sm text-base-content/70 space-y-1">
                <div><code class="bg-base-300 px-2 py-1 rounded">APP_PORT</code> - Server port (default: 8080)</div>
                <div><code class="bg-base-300 px-2 py-1 rounded">SQLITE_PATH</code> - SQLite database path (default: /data/app/app.db)</div>
                <div><code class="bg-base-300 px-2 py-1 rounded">BACKUP_DIR</code> - Backup storage directory (default: /data/backups)</div>
                <div><code class="bg-base-300 px-2 py-1 rounded">APP_ENC_KEY</code> - 32-byte base64 encryption key (required)</div>
                <div><code class="bg-base-300 px-2 py-1 rounded">ADMIN_USER</code> - Basic auth username (optional)</div>
                <div><code class="bg-base-300 px-2 py-1 rounded">ADMIN_PASS</code> - Basic auth password (optional)</div>
              </div>
            </div>
            
            <div>
              <h3 class="font-semibold mb-2">Docker Usage</h3>
              <div class="mockup-code text-xs">
                <pre><code>docker run -p 8080:8080 \
  -e BACKUP_DIR=/data/backups \
  -e SQLITE_PATH=/data/app/app.db \
  -e APP_ENC_KEY=BASE64_32_BYTES_KEY \
  -v backups:/data/backups \
  -v app:/data/app \
  ghcr.io/user/go-dumper:latest</code></pre>
              </div>
            </div>

            <div>
              <h3 class="font-semibold mb-2">Generate Encryption Key</h3>
              <div class="flex gap-2 items-center">
                <code class="bg-base-300 px-2 py-1 rounded text-xs flex-1 overflow-hidden">{{ encryptionKey }}</code>
                <button @click="generateKey" class="btn btn-xs btn-outline">Generate</button>
                <button @click="copyKey" class="btn btn-xs btn-outline">Copy</button>
              </div>
              <div class="text-xs text-base-content/50 mt-1">
                Use this as your APP_ENC_KEY environment variable
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Health Check -->
      <div class="card bg-base-200 shadow">
        <div class="card-body">
          <h2 class="card-title mb-4">Health Check</h2>
          <div class="flex gap-2">
            <button @click="checkHealth" class="btn btn-outline btn-sm" :disabled="checking">
              <span v-if="checking" class="loading loading-spinner loading-xs"></span>
              Check System Health
            </button>
          </div>
          <div v-if="healthStatus" class="mt-4">
            <div class="alert" :class="{
              'alert-success': healthStatus.status === 'ok',
              'alert-error': healthStatus.status !== 'ok'
            }">
              <div>
                <div class="font-bold">{{ healthStatus.status === 'ok' ? 'System Healthy' : 'System Error' }}</div>
                <div class="text-sm">{{ healthStatus.service }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { healthApi } from '@/services/api'
import { useToastStore } from '@/stores/toasts'

const toastStore = useToastStore()

const selectedTheme = ref('light')
const encryptionKey = ref('')
const checking = ref(false)
const healthStatus = ref<{ status: string; service: string } | null>(null)

const changeTheme = () => {
  document.documentElement.setAttribute('data-theme', selectedTheme.value)
  localStorage.setItem('theme', selectedTheme.value)
}

const generateKey = () => {
  const array = new Uint8Array(32)
  crypto.getRandomValues(array)
  encryptionKey.value = btoa(String.fromCharCode(...array))
}

const copyKey = async () => {
  try {
    await navigator.clipboard.writeText(encryptionKey.value)
    toastStore.addToast('success', 'Copied', 'Encryption key copied to clipboard')
  } catch (error) {
    toastStore.addToast('error', 'Error', 'Failed to copy to clipboard')
  }
}

const checkHealth = async () => {
  checking.value = true
  try {
    healthStatus.value = await healthApi.check()
  } catch (error) {
    healthStatus.value = { status: 'error', service: 'go-dumper' }
  } finally {
    checking.value = false
  }
}

onMounted(() => {
  // Load saved theme
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    selectedTheme.value = savedTheme
    document.documentElement.setAttribute('data-theme', savedTheme)
  } else {
    // Default to system preference
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      selectedTheme.value = 'dark'
      document.documentElement.setAttribute('data-theme', 'dark')
    }
  }

  // Generate initial encryption key
  generateKey()

  // Check health on load
  checkHealth()
})
</script>