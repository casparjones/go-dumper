<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-base-content mb-2">Dashboard</h1>
      <p class="text-base-content/70">Overview of your backup system and recent activity</p>
    </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <div class="stat bg-base-200 rounded-lg shadow-sm">
        <div class="stat-figure text-primary">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        </div>
        <div class="stat-title">Configurations</div>
        <div class="stat-value text-primary">{{ targetStore.targets.length }}</div>
        <div class="stat-desc">Active targets</div>
      </div>

      <div class="stat bg-base-200 rounded-lg shadow-sm">
        <div class="stat-figure text-success">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="stat-title">Successful</div>
        <div class="stat-value text-success">124</div>
        <div class="stat-desc">This month</div>
      </div>

      <div class="stat bg-base-200 rounded-lg shadow-sm">
        <div class="stat-figure text-warning">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="stat-title">Scheduled</div>
        <div class="stat-value text-warning">{{ scheduledTargets }}</div>
        <div class="stat-desc">Auto backups</div>
      </div>

      <div class="stat bg-base-200 rounded-lg shadow-sm">
        <div class="stat-figure text-info">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 1.79 4 4 4h8c2.21 0 4-1.79 4-4V7M4 7l8-4 8 4M4 7l8 4 8-4" />
          </svg>
        </div>
        <div class="stat-title">Storage Used</div>
        <div class="stat-value text-info">{{ formatBytes(totalStorageUsed) }}</div>
        <div class="stat-desc">Total backups</div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 xl:grid-cols-3 gap-6">
      <!-- Left Column - Quick Actions & Active Targets -->
      <div class="xl:col-span-1 space-y-6">
        <!-- Quick Actions -->
        <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
          <h3 class="text-lg font-semibold text-base-content mb-4 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
            Quick Actions
          </h3>
          <div class="space-y-3">
            <router-link to="/configure" class="btn btn-primary w-full justify-start">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Add New Configuration
            </router-link>
            <router-link to="/backup" class="btn btn-outline w-full justify-start">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
              </svg>
              Create Backup
            </router-link>
            <router-link to="/restore" class="btn btn-outline w-full justify-start">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
              </svg>
              Restore Database
            </router-link>
            <button @click="refreshData" class="btn btn-ghost w-full justify-start" :disabled="loading">
              <span v-if="loading" class="loading loading-spinner loading-xs"></span>
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Refresh Status
            </button>
          </div>
        </div>

        <!-- Active Configurations -->
        <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
          <h3 class="text-lg font-semibold text-base-content mb-4 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            Active Configurations
          </h3>
          <div class="space-y-3">
            <div
              v-for="target in recentTargets"
              :key="target.id"
              class="flex items-center justify-between p-3 bg-base-200 rounded-lg"
            >
              <div>
                <div class="font-medium text-sm text-base-content">{{ target.name }}</div>
                <div class="text-xs text-base-content/70">{{ target.host }}:{{ target.port }}</div>
              </div>
              <div class="flex gap-1">
                <button
                  @click="createBackup(target.id)"
                  class="btn btn-xs btn-primary"
                  :disabled="targetStore.loading"
                >
                  Backup
                </button>
              </div>
            </div>
            <div v-if="recentTargets.length === 0" class="text-center text-base-content/50 py-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 mx-auto text-base-content/20 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
              <p class="text-sm">No configurations yet</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column - Status & History -->
      <div class="xl:col-span-2 space-y-6">
        <StatusPanel />
        <BackupHistory />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useTargetsStore } from '@/stores/targets'
import StatusPanel from '@/components/StatusPanel.vue'
import BackupHistory from '@/components/BackupHistory.vue'

const targetStore = useTargetsStore()
const loading = ref(false)

const scheduledTargets = computed(() => {
  return targetStore.targets.filter(t => t.schedule_time).length
})

const recentTargets = computed(() => {
  return [...targetStore.targets]
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .slice(0, 5)
})

const totalStorageUsed = computed(() => {
  // This would need to be calculated from actual backup data
  return 8547483648 // 8.5 GB as example
})

const refreshData = async () => {
  loading.value = true
  try {
    await targetStore.fetchTargets()
  } finally {
    loading.value = false
  }
}

const createBackup = async (targetId: number) => {
  await targetStore.createBackup(targetId)
}

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

onMounted(() => {
  targetStore.fetchTargets()
})
</script>