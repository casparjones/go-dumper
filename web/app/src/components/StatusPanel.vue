<template>
  <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
    <h3 class="text-lg font-semibold text-base-content mb-4 flex items-center gap-2">
      <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-info" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      System Information
    </h3>
    
    <div class="space-y-6">
      <!-- GoDumper Information -->
      <div>
        <h4 class="font-medium text-base-content mb-2">GoDumper</h4>
        <div class="space-y-1 text-sm text-base-content/80">
          <div class="flex justify-between">
            <span>Version:</span>
            <span class="font-mono">{{ systemInfo.version }}</span>
          </div>
          <div class="flex justify-between">
            <span>Installation Path:</span>
            <span class="font-mono text-xs">{{ systemInfo.installPath }}</span>
          </div>
          <div class="flex justify-between">
            <span>Backup Directory:</span>
            <span class="font-mono text-xs">{{ systemInfo.backupDir }}</span>
          </div>
          <div class="flex justify-between">
            <span>Config Database:</span>
            <span class="font-mono text-xs">{{ systemInfo.configDb }}</span>
          </div>
        </div>
      </div>

      <!-- Operating System -->
      <div>
        <h4 class="font-medium text-base-content mb-2">Operating System</h4>
        <div class="space-y-1 text-sm text-base-content/80">
          <div class="flex justify-between">
            <span>OS:</span>
            <span class="font-mono">{{ systemInfo.os.name }}</span>
          </div>
          <div class="flex justify-between">
            <span>Architecture:</span>
            <span class="font-mono">{{ systemInfo.os.arch }}</span>
          </div>
          <div class="flex justify-between">
            <span>Kernel:</span>
            <span class="font-mono">{{ systemInfo.os.kernel }}</span>
          </div>
        </div>
      </div>

      <!-- MySQL Information -->
      <div>
        <h4 class="font-medium text-base-content mb-2">Database Server</h4>
        <div class="space-y-1 text-sm text-base-content/80">
          <div class="flex justify-between">
            <span>Server:</span>
            <span class="font-mono">{{ systemInfo.mysql.server }}</span>
          </div>
          <div class="flex justify-between">
            <span>Version:</span>
            <span class="font-mono">{{ systemInfo.mysql.version }}</span>
          </div>
          <div class="flex justify-between">
            <span>Status:</span>
            <span class="flex items-center gap-1">
              <div class="w-2 h-2 rounded-full" :class="systemInfo.mysql.status === 'Connected' ? 'bg-success' : 'bg-error'"></div>
              {{ systemInfo.mysql.status }}
            </span>
          </div>
        </div>
      </div>

      <!-- Statistics -->
      <div>
        <h4 class="font-medium text-base-content mb-2">Statistics</h4>
        <div class="grid grid-cols-2 gap-4">
          <div class="text-center p-3 bg-base-200 rounded-lg">
            <div class="text-lg font-bold text-primary">{{ systemInfo.stats.totalConfigs }}</div>
            <div class="text-xs text-base-content/70">Configurations</div>
          </div>
          <div class="text-center p-3 bg-base-200 rounded-lg">
            <div class="text-lg font-bold text-success">{{ systemInfo.stats.totalBackups }}</div>
            <div class="text-xs text-base-content/70">Total Backups</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface SystemInfo {
  version: string
  installPath: string
  backupDir: string
  configDb: string
  os: {
    name: string
    arch: string
    kernel: string
  }
  mysql: {
    server: string
    version: string
    status: string
  }
  stats: {
    totalConfigs: number
    totalBackups: number
  }
}

const systemInfo = ref<SystemInfo>({
  version: 'v1.0.0',
  installPath: '/usr/local/bin/go-dumper',
  backupDir: '/data/backups',
  configDb: '/data/app/app.db',
  os: {
    name: 'Linux Ubuntu 22.04',
    arch: 'amd64',
    kernel: '5.15.0-72-generic'
  },
  mysql: {
    server: 'MySQL',
    version: '8.0.33',
    status: 'Connected'
  },
  stats: {
    totalConfigs: 3,
    totalBackups: 127
  }
})

onMounted(async () => {
  // In a real implementation, you would fetch this data from the API
  try {
    // const response = await systemApi.getInfo()
    // systemInfo.value = response.data
  } catch (error) {
    console.error('Failed to fetch system info:', error)
  }
})
</script>