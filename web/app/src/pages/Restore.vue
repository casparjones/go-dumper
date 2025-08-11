<template>
  <div class="mx-auto max-w-3xl space-y-6">
    <!-- Header -->
    <div class="mb-2">
      <h1 class="text-3xl font-bold mb-1">Restore Database</h1>
      <p class="text-base-content/70">Restore your database from existing backups</p>
    </div>

    <!-- Configuration Selection -->
    <div class="card bg-base-100 border border-base-200 shadow-sm">
      <div class="card-body gap-4">
        <h3 class="card-title">Select Target Configuration</h3>

        <div class="grid gap-4 md:grid-cols-2">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Database Configuration</span>
            </label>
            <select v-model="selectedConfig" class="select select-bordered w-full">
              <option value="">Select a configuration...</option>
              <option v-for="target in targetStore.targets" :key="target.id" :value="target.id">
                {{ target.name }} ({{ target.host }}:{{ target.port }})
              </option>
            </select>
            <label class="label" v-if="!selectedConfig">
              <span class="label-text-alt">Choose a target to see available backups</span>
            </label>
          </div>

          <div v-if="selectedConfig" class="form-control">
            <label class="label">
              <span class="label-text">Available Backups</span>
            </label>
            
            <div v-if="loadingBackups" class="flex justify-center py-4">
              <span class="loading loading-spinner loading-sm"></span>
              <span class="ml-2">Loading backups...</span>
            </div>
            
            <div v-else-if="groupedBackups.length === 0" class="text-center py-8 text-base-content/70">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 mx-auto mb-3 text-base-content/30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <p>No backups found for this configuration</p>
              <router-link to="/backup" class="btn btn-primary btn-sm mt-2">Create Backup</router-link>
            </div>
            
            <div v-else class="space-y-4 max-h-96 overflow-y-auto">
              <div v-for="group in groupedBackups" :key="group.database" class="border border-base-200 rounded-lg overflow-hidden">
                <div class="bg-base-200 px-3 py-2 font-medium text-sm">
                  {{ group.database }}
                  <span class="text-xs text-base-content/70 ml-2">({{ group.backups.length }} backups)</span>
                </div>
                <div class="divide-y divide-base-200">
                  <label 
                    v-for="backup in group.backups" 
                    :key="backup.id"
                    class="flex items-center p-3 hover:bg-base-100 cursor-pointer"
                    :class="{ 'bg-primary/10': selectedBackup === backup.id }"
                  >
                    <input
                      v-model="selectedBackup"
                      type="radio"
                      :value="backup.id"
                      name="backup_selection"
                      class="radio radio-primary radio-sm"
                    />
                    <div class="ml-3 flex-1">
                      <div class="flex items-center justify-between">
                        <div class="font-medium text-sm">{{ formatDate(backup.started_at) }}</div>
                        <div class="text-xs text-base-content/70">{{ formatTime(backup.started_at) }}</div>
                      </div>
                      <div class="flex items-center justify-between mt-1">
                        <div class="text-xs text-base-content/70">{{ formatBytes(backup.size_bytes) }}</div>
                        <div class="flex items-center gap-2">
                          <div 
                            class="w-2 h-2 rounded-full"
                            :class="{
                              'bg-success': backup.status === 'success',
                              'bg-error': backup.status === 'failed',
                              'bg-warning animate-pulse': backup.status === 'running'
                            }"
                          ></div>
                          <span class="text-xs capitalize" :class="{
                            'text-success': backup.status === 'success',
                            'text-error': backup.status === 'failed',
                            'text-warning': backup.status === 'running'
                          }">{{ backup.status }}</span>
                        </div>
                      </div>
                    </div>
                  </label>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Backup Details -->
    <div v-if="selectedBackupInfo" class="card bg-base-100 border border-base-200 shadow-sm">
      <div class="card-body gap-4">
        <h3 class="card-title">Backup Details</h3>

        <div class="grid gap-4 md:grid-cols-4">
          <div>
            <div class="text-sm text-base-content/70">Database</div>
            <div class="font-medium">{{ selectedBackupInfo.database_name }}</div>
          </div>
          <div>
            <div class="text-sm text-base-content/70">Created</div>
            <div class="font-medium">{{ formatDate(selectedBackupInfo.started_at) }}</div>
          </div>
          <div>
            <div class="text-sm text-base-content/70">Size</div>
            <div class="font-medium">{{ formatBytes(selectedBackupInfo.size_bytes) }}</div>
          </div>
          <div>
            <div class="text-sm text-base-content/70">Status</div>
            <div class="flex items-center gap-2">
              <span
                  class="badge"
                  :class="{
                  'badge-success': selectedBackupInfo.status === 'success',
                  'badge-error': selectedBackupInfo.status === 'failed',
                  'badge-warning': selectedBackupInfo.status === 'running'
                }"
              >
                {{ selectedBackupInfo.status }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Restore Options -->
    <div v-if="selectedBackup" class="card bg-base-100 border border-base-200 shadow-sm">
      <div class="card-body gap-4">
        <h3 class="card-title">Restore Options</h3>

        <div class="alert alert-warning">
          <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 9v2m0 4h.01M4.93 20h14.14c1.54 0 2.5-1.67 1.73-2.5L13.73 4c-.77-.83-1.69-.83-2.46 0L3.2 17.5C2.43 18.33 3.39 20 4.93 20z"/>
          </svg>
          <div>
            <h3 class="font-bold">Warning</h3>
            <div class="text-xs">
              This replaces all data in the target database. Ensure you have a current backup.
            </div>
          </div>
        </div>

        <div class="space-y-3">
          <label class="label cursor-pointer justify-between">
            <span class="label-text">Drop existing tables before restore</span>
            <input type="checkbox" v-model="dropExistingTables" class="toggle toggle-primary" />
          </label>

          <label class="label cursor-pointer justify-between opacity-60">
            <span class="label-text">Temporarily disable foreign key checks</span>
            <input type="checkbox" v-model="disableForeignKeys" class="toggle toggle-primary" checked disabled />
          </label>

          <label class="label cursor-pointer justify-between">
            <span class="label-text">Create database if it doesn't exist</span>
            <input type="checkbox" v-model="createDatabase" class="toggle toggle-primary" />
          </label>
        </div>
      </div>
    </div>

    <!-- Confirmation -->
    <div v-if="selectedBackup" class="card bg-base-100 border border-base-200 shadow-sm">
      <div class="card-body gap-4">
        <h3 class="card-title">Confirmation</h3>

        <div class="form-control">
          <label class="label">
            <span class="label-text">Type "RESTORE" to confirm this dangerous operation</span>
          </label>
          <input
              v-model="confirmationText"
              type="text"
              placeholder="RESTORE"
              class="input input-bordered w-full"
              :class="{ 'input-error': confirmationText && confirmationText !== 'RESTORE' }"
          />
          <label v-if="!canRestore && confirmationText" class="label">
            <span class="label-text-alt text-error">Please type "RESTORE" exactly to confirm</span>
          </label>
        </div>

        <div class="flex flex-wrap gap-3">
          <button
              @click="performRestore"
              class="btn btn-error"
              :disabled="!canRestore || restoring"
          >
            <span v-if="restoring" class="loading loading-spinner loading-sm"></span>
            <span v-else class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none"
                   viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8-4-4m0 0L8 8m4-4v12"/>
              </svg>
              Restore Database
            </span>
          </button>

          <button @click="resetForm" class="btn btn-ghost">
            Cancel
          </button>
        </div>
      </div>
    </div>

    <!-- Restore Progress -->
    <div v-if="restoreProgress" class="card bg-base-100 border border-base-200 shadow-sm">
      <div class="card-body gap-4">
        <h3 class="card-title">Restore Progress</h3>
        <div class="space-y-3">
          <div class="flex justify-between text-sm">
            <span>{{ restoreProgress.current_step }}</span>
            <span>{{ restoreProgress.progress }}%</span>
          </div>
          <progress class="progress progress-primary w-full" :value="restoreProgress.progress" max="100"></progress>
          <div class="text-xs text-base-content/70">{{ restoreProgress.message }}</div>
        </div>
      </div>
    </div>

    <!-- Sticky Action Bar -->
    <div
        v-if="selectedBackup"
        class="btm-nav btm-nav-sm z-30 shadow-lg bg-base-100 border-t border-base-200 md:hidden"
    >
      <button class="text-error" :disabled="!canRestore || restoring" @click="performRestore">
        <span v-if="restoring" class="loading loading-spinner loading-xs"></span>
        <span v-else>Restore</span>
      </button>
      <button class="text-base-content/70" @click="resetForm">Cancel</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useTargetsStore } from '@/stores/targets'
import { useBackupsStore } from '@/stores/backups'

const targetStore = useTargetsStore()
const backupsStore = useBackupsStore()

const selectedConfig = ref<number | ''>('')
const selectedBackup = ref<number | ''>('')
const confirmationText = ref('')
const restoring = ref(false)

// Restore options
const dropExistingTables = ref(true)
const disableForeignKeys = ref(true)
const createDatabase = ref(false)

// Loading state and real backup data
const loadingBackups = ref(false)
const availableBackups = ref<any[]>([])

const restoreProgress = ref<{
  progress: number
  current_step: string
  message: string
} | null>(null)

const selectedTargetInfo = computed(() => {
  if (!selectedConfig.value) return null
  return targetStore.targets.find(t => t.id === selectedConfig.value)
})

const groupedBackups = computed(() => {
  if (!availableBackups.value.length) return []
  
  // Group backups by database name
  const groups = availableBackups.value.reduce((acc, backup) => {
    const dbName = backup.database_name || 'Unknown Database'
    if (!acc[dbName]) {
      acc[dbName] = []
    }
    acc[dbName].push(backup)
    return acc
  }, {})
  
  // Convert to array and sort backups within each group
  return Object.entries(groups).map(([database, backups]: [string, any]) => ({
    database,
    backups: backups.sort((a: any, b: any) => 
      new Date(b.started_at).getTime() - new Date(a.started_at).getTime()
    )
  })).sort((a, b) => a.database.localeCompare(b.database))
})

const selectedBackupInfo = computed(() => {
  if (!selectedBackup.value) return null
  return availableBackups.value.find(b => b.id === selectedBackup.value)
})

const canRestore = computed(() => {
  return Boolean(selectedConfig.value && selectedBackup.value && confirmationText.value === 'RESTORE' && !restoring.value)
})

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString()
}

const formatTime = (dateString: string): string => {
  return new Date(dateString).toLocaleTimeString()
}

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const performRestore = async () => {
  if (!canRestore.value) return

  restoring.value = true
  restoreProgress.value = {
    progress: 0,
    current_step: 'Preparing restore...',
    message: 'Validating backup file and target database'
  }

  try {
    // Call the real restore API
    const response = await fetch(`/api/backups/${selectedBackup.value}/restore`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        create_database: createDatabase.value
      })
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to start restore')
    }

    // Show progress simulation (since restore runs in background)
    const steps = [
      { step: 'Restore started...', message: 'Backup restore has been initiated', progress: 20 },
      { step: 'Decompressing backup...', message: 'Extracting SQL from compressed backup file', progress: 40 },
      { step: 'Creating database structure...', message: 'Restoring tables, indexes, and constraints', progress: 60 },
      { step: 'Importing data...', message: 'Restoring table data from backup', progress: 80 },
      { step: 'Restore in progress...', message: 'Restore is running in the background', progress: 100 }
    ]

    for (const s of steps) {
      restoreProgress.value = {
        progress: s.progress,
        current_step: s.step,
        message: s.message
      }
      await new Promise(r => setTimeout(r, 1000))
    }

    // Final message
    restoreProgress.value = {
      progress: 100,
      current_step: 'Restore initiated!',
      message: 'The restore process is running in the background. Check your database.'
    }

    // Reset after delay
    setTimeout(() => {
      resetForm()
      restoreProgress.value = null
    }, 3000)

  } catch (e) {
    console.error('Restore failed:', e)
    restoreProgress.value = {
      progress: 0,
      current_step: 'Restore failed!',
      message: e.message || 'An error occurred during the restore process'
    }
    setTimeout(() => {
      restoreProgress.value = null
    }, 5000)
  } finally {
    restoring.value = false
  }
}

const resetForm = () => {
  selectedConfig.value = ''
  selectedBackup.value = ''
  confirmationText.value = ''
  dropExistingTables.value = true
  createDatabase.value = false
}

const loadBackupsForTarget = async (targetId: number) => {
  if (!targetId) return
  
  loadingBackups.value = true
  try {
    // Load backups for the selected target
    const response = await fetch(`/api/targets/${targetId}/backups`)
    if (!response.ok) {
      throw new Error('Failed to load backups')
    }
    const backups = await response.json()
    availableBackups.value = backups.filter(b => b.status === 'success') // Only show successful backups for restore
  } catch (error) {
    console.error('Failed to load backups:', error)
    availableBackups.value = []
  } finally {
    loadingBackups.value = false
  }
}

watch(selectedConfig, async (newConfig) => {
  selectedBackup.value = ''
  availableBackups.value = []
  
  if (newConfig) {
    await loadBackupsForTarget(newConfig as number)
  }
})

onMounted(() => {
  targetStore.fetchTargets()
})
</script>
