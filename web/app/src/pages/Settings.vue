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
          <h2 class="card-title flex items-center gap-2 mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zM21 5a2 2 0 00-2-2h-4a2 2 0 00-2 2v12a4 4 0 004 4h4a2 2 0 002-2V5z" />
            </svg>
            Theme Selection
          </h2>
          
          <div class="space-y-4">
            <p class="text-sm text-base-content/70">
              Choose your preferred theme. The theme will be saved and remembered across sessions.
            </p>
            
            <!-- Current Theme Display -->
            <div class="flex items-center gap-3 p-3 bg-base-300 rounded-lg">
              <div class="w-4 h-4 rounded-full bg-primary"></div>
              <div class="w-4 h-4 rounded-full bg-secondary"></div>
              <div class="w-4 h-4 rounded-full bg-accent"></div>
              <span class="text-sm font-medium">Current: {{ appStore.currentTheme }}</span>
            </div>

            <!-- Theme Selection Grid -->
            <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              <div 
                v-for="theme in availableThemes" 
                :key="theme.name"
                @click="selectTheme(theme.name)"
                class="card bg-base-300 cursor-pointer border-2 transition-all hover:scale-105"
                :class="{ 'border-primary shadow-lg': appStore.currentTheme === theme.name }"
              >
                <div class="card-body p-4">
                  <div class="flex items-center justify-between mb-2">
                    <h4 class="text-sm font-medium">{{ theme.label }}</h4>
                    <input 
                      type="radio" 
                      :checked="appStore.currentTheme === theme.name"
                      class="radio radio-primary radio-xs"
                      readonly
                    />
                  </div>
                  <!-- Theme Preview -->
                  <div class="flex gap-1 mb-2">
                    <div 
                      v-for="color in theme.preview" 
                      :key="color"
                      class="w-3 h-3 rounded-full"
                      :style="{ backgroundColor: color }"
                    ></div>
                  </div>
                  <span class="text-xs text-base-content/60">{{ theme.description }}</span>
                </div>
              </div>
            </div>

            <!-- Apply Button -->
            <div class="flex justify-end">
              <button 
                @click="saveThemeSettings"
                class="btn btn-primary"
                :disabled="saving"
              >
                <span v-if="saving" class="loading loading-spinner loading-sm"></span>
                <span v-else class="flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  Save Theme Settings
                </span>
              </button>
            </div>
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
import { useAppStore } from '@/stores/app'

const toastStore = useToastStore()
const appStore = useAppStore()

const encryptionKey = ref('')
const checking = ref(false)
const saving = ref(false)
const healthStatus = ref<{ status: string; service: string } | null>(null)

const availableThemes = ref([
  { name: 'light', label: 'Light', description: 'Clean light theme', preview: ['#ffffff', '#f3f4f6', '#3b82f6'] },
  { name: 'dark', label: 'Dark', description: 'Classic dark theme', preview: ['#1f2937', '#374151', '#60a5fa'] },
  { name: 'dim', label: 'Dim', description: 'Dimmed dark theme', preview: ['#2a2e37', '#383e47', '#66cc8a'] },
  { name: 'winter', label: 'Winter', description: 'Cool light theme', preview: ['#f8fafc', '#e2e8f0', '#0ea5e9'] },
  { name: 'cupcake', label: 'Cupcake', description: 'Sweet pink theme', preview: ['#fef7f0', '#f7c3d7', '#ec4899'] },
  { name: 'bumblebee', label: 'Bumblebee', description: 'Bright yellow theme', preview: ['#fef9c3', '#fde047', '#eab308'] },
  { name: 'emerald', label: 'Emerald', description: 'Fresh green theme', preview: ['#ecfdf5', '#10b981', '#059669'] },
  { name: 'corporate', label: 'Corporate', description: 'Professional blue', preview: ['#f8fafc', '#3b82f6', '#1e40af'] },
  { name: 'synthwave', label: 'Synthwave', description: 'Retro neon theme', preview: ['#2d1b69', '#ff007c', '#00d9ff'] },
  { name: 'retro', label: 'Retro', description: 'Vintage style', preview: ['#ede6d3', '#d2691e', '#8b4513'] },
  { name: 'cyberpunk', label: 'Cyberpunk', description: 'Futuristic neon', preview: ['#0a0a0a', '#ff007c', '#00ff9f'] },
  { name: 'valentine', label: 'Valentine', description: 'Romantic pink', preview: ['#ffe0e6', '#ff69b4', '#dc143c'] },
  { name: 'halloween', label: 'Halloween', description: 'Spooky orange', preview: ['#1a1a1a', '#ff6600', '#990099'] },
  { name: 'garden', label: 'Garden', description: 'Natural green', preview: ['#f0fdf4', '#22c55e', '#16a34a'] },
  { name: 'forest', label: 'Forest', description: 'Deep green', preview: ['#052e16', '#166534', '#22c55e'] },
  { name: 'aqua', label: 'Aqua', description: 'Ocean blue', preview: ['#f0fdff', '#06b6d4', '#0891b2'] },
  { name: 'lofi', label: 'Lofi', description: 'Calm and minimal', preview: ['#f7f3f3', '#9ca3af', '#6b7280'] },
  { name: 'pastel', label: 'Pastel', description: 'Soft colors', preview: ['#fef3f2', '#fbbf24', '#a78bfa'] },
  { name: 'fantasy', label: 'Fantasy', description: 'Magical purple', preview: ['#f3e8ff', '#a855f7', '#7c3aed'] },
  { name: 'wireframe', label: 'Wireframe', description: 'Minimal black/white', preview: ['#ffffff', '#000000', '#6b7280'] },
  { name: 'black', label: 'Black', description: 'Pure black theme', preview: ['#000000', '#333333', '#ffffff'] },
  { name: 'luxury', label: 'Luxury', description: 'Elegant gold', preview: ['#1a1a1a', '#d4af37', '#ffd700'] },
  { name: 'dracula', label: 'Dracula', description: 'Popular dark theme', preview: ['#282a36', '#ff79c6', '#50fa7b'] },
  { name: 'cmyk', label: 'CMYK', description: 'Print colors', preview: ['#ffffff', '#00ffff', '#ff00ff'] },
  { name: 'autumn', label: 'Autumn', description: 'Warm fall colors', preview: ['#fffbeb', '#d97706', '#dc2626'] },
  { name: 'business', label: 'Business', description: 'Professional gray', preview: ['#f9fafb', '#4b5563', '#1f2937'] },
  { name: 'acid', label: 'Acid', description: 'Bright lime', preview: ['#f7fee7', '#84cc16', '#65a30d'] },
  { name: 'lemonade', label: 'Lemonade', description: 'Fresh yellow', preview: ['#fffbeb', '#fbbf24', '#f59e0b'] },
  { name: 'night', label: 'Night', description: 'Deep blue night', preview: ['#0f172a', '#1e293b', '#0ea5e9'] },
  { name: 'coffee', label: 'Coffee', description: 'Warm brown', preview: ['#44403c', '#78716c', '#a8a29e'] },
  { name: 'nord', label: 'Nord', description: 'Arctic inspired', preview: ['#2e3440', '#5e81ac', '#88c0d0'] },
  { name: 'sunset', label: 'Sunset', description: 'Warm orange/pink', preview: ['#fff7ed', '#fb923c', '#f97316'] }
])

const selectTheme = (themeName: string) => {
  appStore.setTheme(themeName)
}

const saveThemeSettings = async () => {
  saving.value = true
  try {
    localStorage.setItem('theme', appStore.currentTheme)
    await appStore.saveThemeToBackend(appStore.currentTheme)
    toastStore.addToast('success', 'Settings Saved', `Theme "${appStore.currentTheme}" has been saved`)
  } catch (error: any) {
    const errorMessage = error.response?.data?.error || 'Failed to save theme settings'
    toastStore.addToast('error', 'Error', errorMessage)
  } finally {
    saving.value = false
  }
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
  // Generate initial encryption key
  generateKey()

  // Check health on load
  checkHealth()
})
</script>