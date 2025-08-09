<template>
  <div class="space-y-6">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-base-content mb-2">Restore Database</h1>
      <p class="text-base-content/70">Restore your database from existing backups</p>
    </div>

    <!-- Configuration Selection -->
    <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4">Select Target Configuration</h3>
      
      <div class="grid gap-4 md:grid-cols-2">
        <div class="form-control">
          <label class="label">
            <span class="label-text">Database Configuration</span>
          </label>
          <select v-model="selectedConfig" class="select select-bordered">
            <option value="">Select a configuration...</option>
            <option v-for="target in targetStore.targets" :key="target.id" :value="target.id">
              {{ target.name }} ({{ target.host }}:{{ target.port }})
            </option>
          </select>
        </div>
        
        <div v-if="selectedConfig" class="form-control">
          <label class="label">
            <span class="label-text">Available Backups</span>
          </label>
          <select v-model="selectedBackup" class="select select-bordered">
            <option value="">Select a backup...</option>
            <option v-for="backup in availableBackups" :key="backup.id" :value="backup.id">
              {{ formatDate(backup.created_at) }} - {{ backup.size }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <!-- Backup Details -->
    <div v-if="selectedBackupInfo" class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4">Backup Details</h3>
      <div class="grid gap-4 md:grid-cols-3">
        <div>
          <div class="text-sm text-base-content/70">Created</div>
          <div class="font-medium">{{ formatDate(selectedBackupInfo.created_at) }}</div>
        </div>
        <div>
          <div class="text-sm text-base-content/70">Size</div>
          <div class="font-medium">{{ selectedBackupInfo.size }}</div>
        </div>
        <div>
          <div class="text-sm text-base-content/70">Status</div>
          <div class="flex items-center gap-2">
            <div 
              class="w-2 h-2 rounded-full"
              :class="{
                'bg-success': selectedBackupInfo.status === 'success',
                'bg-error': selectedBackupInfo.status === 'failed',
                'bg-warning': selectedBackupInfo.status === 'running'
              }"
            ></div>
            <span class="font-medium text-sm capitalize">{{ selectedBackupInfo.status }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Restore Options -->
    <div v-if="selectedBackup" class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4">Restore Options</h3>
      
      <div class="space-y-4">
        <div class="alert alert-warning">
          <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
          <div>
            <h3 class="font-bold">Warning!</h3>
            <div class="text-xs">This operation will replace all data in the target database. Make sure you have a current backup before proceeding.</div>
          </div>
        </div>

        <div class="space-y-3">
          <label class="cursor-pointer label justify-start gap-2">
            <input type="checkbox" v-model="dropExistingTables" class="checkbox checkbox-primary checkbox-sm" />
            <span class="label-text">Drop existing tables before restore</span>
          </label>
          <label class="cursor-pointer label justify-start gap-2">
            <input type="checkbox" v-model="disableForeignKeys" class="checkbox checkbox-primary checkbox-sm" checked disabled />
            <span class="label-text">Temporarily disable foreign key checks</span>
          </label>
          <label class="cursor-pointer label justify-start gap-2">
            <input type="checkbox" v-model="createDatabase" class="checkbox checkbox-primary checkbox-sm" />
            <span class="label-text">Create database if it doesn't exist</span>
          </label>
        </div>
      </div>
    </div>

    <!-- Confirmation -->
    <div v-if="selectedBackup" class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4">Confirmation</h3>
      
      <div class="space-y-4">
        <div class="form-control">
          <label class="label">
            <span class="label-text">Type "RESTORE" to confirm this dangerous operation</span>
          </label>
          <input 
            v-model="confirmationText" 
            type="text" 
            placeholder="RESTORE" 
            class="input input-bordered"
            :class="{ 'input-error': confirmationText && confirmationText !== 'RESTORE' }"
          />
        </div>
        
        <div class="flex gap-4">
          <button 
            @click="performRestore" 
            class="btn btn-error"
            :disabled="!canRestore || restoring"
          >
            <span v-if="restoring" class="loading loading-spinner loading-sm"></span>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
            </svg>
            {{ restoring ? 'Restoring...' : 'Restore Database' }}
          </button>
          
          <button @click="resetForm" class="btn btn-ghost">
            Cancel
          </button>
        </div>
        
        <div v-if="!canRestore && confirmationText" class="text-sm text-error">
          Please type "RESTORE" exactly to confirm
        </div>
      </div>
    </div>

    <!-- Restore Progress -->
    <div v-if="restoreProgress" class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4">Restore Progress</h3>
      <div class="space-y-3">
        <div class="flex justify-between text-sm">
          <span>{{ restoreProgress.current_step }}</span>
          <span>{{ restoreProgress.progress }}%</span>
        </div>
        <progress class="progress progress-primary w-full" :value="restoreProgress.progress" max="100"></progress>
        <div class="text-xs text-base-content/70">
          {{ restoreProgress.message }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useTargetsStore } from '@/stores/targets'

const targetStore = useTargetsStore()

const selectedConfig = ref<number | ''>('')
const selectedBackup = ref<number | ''>('')
const confirmationText = ref('')
const restoring = ref(false)

// Restore options
const dropExistingTables = ref(true)
const disableForeignKeys = ref(true)
const createDatabase = ref(false)

// Mock backup data - in real app this would come from API
const availableBackups = ref([
  {
    id: 1,
    target_id: 1,
    created_at: '2024-01-15T14:30:00Z',
    size: '2.4 GB',
    status: 'success',
    file_path: '/data/backups/prod_2024-01-15_14-30-00.sql.gz'
  },
  {
    id: 2,
    target_id: 1,
    created_at: '2024-01-14T14:30:00Z',
    size: '2.3 GB',
    status: 'success',
    file_path: '/data/backups/prod_2024-01-14_14-30-00.sql.gz'
  }
])

const restoreProgress = ref<{
  progress: number
  current_step: string
  message: string
} | null>(null)

const selectedTargetInfo = computed(() => {
  if (!selectedConfig.value) return null
  return targetStore.targets.find(t => t.id === selectedConfig.value)
})

const selectedBackupInfo = computed(() => {
  if (!selectedBackup.value) return null
  return availableBackups.value.find(b => b.id === selectedBackup.value)
})

const canRestore = computed(() => {
  return selectedConfig.value && 
         selectedBackup.value && 
         confirmationText.value === 'RESTORE' &&
         !restoring.value
})

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleString()
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
    // Simulate restore progress
    const steps = [
      { step: 'Preparing restore...', message: 'Validating backup file and target database', progress: 10 },
      { step: 'Connecting to database...', message: 'Establishing connection to target database', progress: 20 },
      { step: 'Dropping existing tables...', message: 'Removing existing database structure', progress: 30 },
      { step: 'Decompressing backup...', message: 'Extracting SQL from compressed backup file', progress: 40 },
      { step: 'Creating database structure...', message: 'Restoring tables, indexes, and constraints', progress: 60 },
      { step: 'Importing data...', message: 'Restoring table data from backup', progress: 85 },
      { step: 'Finalizing restore...', message: 'Re-enabling foreign keys and cleaning up', progress: 95 },
      { step: 'Restore completed!', message: 'Database has been successfully restored', progress: 100 }
    ]
    
    for (const stepInfo of steps) {
      restoreProgress.value = stepInfo
      await new Promise(resolve => setTimeout(resolve, 1500))
    }
    
    // Reset form after successful restore
    setTimeout(() => {
      resetForm()
      restoreProgress.value = null
    }, 2000)
    
  } catch (error) {
    console.error('Restore failed:', error)
    restoreProgress.value = {
      progress: 0,
      current_step: 'Restore failed!',
      message: 'An error occurred during the restore process'
    }
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

// Watch for config changes to load available backups
watch(selectedConfig, async (newConfig) => {
  if (newConfig) {
    // In real app: load backups for this target
    // const backups = await backupsApi.getForTarget(newConfig)
    // availableBackups.value = backups
  }
  selectedBackup.value = ''
})

onMounted(() => {
  targetStore.fetchTargets()
})
</script>