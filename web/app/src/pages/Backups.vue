<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-3xl font-bold text-base-content">
          Backups{{ target ? ` - ${target.name}` : '' }}
        </h1>
        <p class="text-base-content/70">View and manage backups for this target</p>
      </div>
      <div class="flex gap-2">
        <button
          v-if="target"
          @click="createBackup"
          class="btn btn-primary"
          :disabled="targetStore.loading"
        >
          <span v-if="targetStore.loading" class="loading loading-spinner loading-xs"></span>
          Create Backup
        </button>
        <router-link to="/targets" class="btn btn-outline">Back to Targets</router-link>
      </div>
    </div>

    <div v-if="target" class="card bg-base-200 shadow mb-6">
      <div class="card-body">
        <h3 class="card-title text-lg">Target Information</h3>
        <div class="grid gap-2 text-sm">
          <div><strong>Host:</strong> {{ target.host }}:{{ target.port }}</div>
          <div><strong>Database:</strong> {{ target.db_name }}</div>
          <div><strong>User:</strong> {{ target.user }}</div>
          <div v-if="target.schedule_time"><strong>Schedule:</strong> Daily at {{ target.schedule_time }} UTC</div>
          <div><strong>Retention:</strong> {{ target.retention_days }} days</div>
        </div>
      </div>
    </div>

    <div v-if="backupStore.loading" class="flex justify-center py-8">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="backupStore.backups.length === 0" class="text-center py-12">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-16 h-16 mx-auto text-base-content/30 mb-4">
        <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z" />
      </svg>
      <h3 class="text-lg font-semibold text-base-content mb-2">No backups found</h3>
      <p class="text-base-content/70 mb-4">Create your first backup to get started</p>
      <button
        @click="createBackup"
        class="btn btn-primary"
        :disabled="targetStore.loading"
      >
        Create Backup
      </button>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="backup in backupStore.backups"
        :key="backup.id"
        class="card bg-base-200 shadow"
      >
        <div class="card-body">
          <div class="flex justify-between items-start">
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <div class="badge" :class="{
                  'badge-success': backup.status === 'success',
                  'badge-error': backup.status === 'failed',
                  'badge-warning': backup.status === 'running'
                }">
                  {{ backup.status }}
                </div>
                <div class="text-sm text-base-content/70">
                  Started: {{ formatDate(backup.started_at) }}
                </div>
                <div v-if="backup.finished_at" class="text-sm text-base-content/70">
                  Finished: {{ formatDate(backup.finished_at) }}
                </div>
              </div>
              
              <div class="text-sm text-base-content/70 space-y-1">
                <div v-if="backup.size_bytes > 0">
                  Size: {{ formatBytes(backup.size_bytes) }}
                </div>
                <div v-if="backup.notes" class="text-error">
                  {{ backup.notes }}
                </div>
              </div>

              <div v-if="backup.status === 'running'" class="mt-2">
                <progress class="progress progress-primary w-full"></progress>
                <div class="text-xs text-base-content/50 mt-1">Backup in progress...</div>
              </div>
            </div>

            <div v-if="backup.status === 'success'" class="dropdown dropdown-end">
              <div tabindex="0" role="button" class="btn btn-ghost btn-sm">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 12.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 18.75a.75.75 0 110-1.5.75.75 0 010 1.5z" />
                </svg>
              </div>
              <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
                <li>
                  <button @click="downloadBackup(backup.id)">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" />
                    </svg>
                    Download
                  </button>
                </li>
                <li>
                  <button @click="confirmRestore(backup)" class="text-warning">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
                    </svg>
                    Restore
                  </button>
                </li>
                <li>
                  <button @click="confirmDelete(backup)" class="text-error">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                    </svg>
                    Delete
                  </button>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Restore Confirmation Modal -->
    <div v-if="backupToRestore" class="modal modal-open">
      <div class="modal-box">
        <h3 class="font-bold text-lg text-warning">⚠️ Confirm Restore</h3>
        <div class="py-4">
          <p class="mb-2">
            Are you sure you want to restore this backup?
          </p>
          <div class="text-sm text-base-content/70 space-y-1">
            <div><strong>Backup Date:</strong> {{ formatDate(backupToRestore.started_at) }}</div>
            <div><strong>Size:</strong> {{ formatBytes(backupToRestore.size_bytes) }}</div>
          </div>
          <div class="alert alert-warning mt-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4.5c-.77-.833-2.694-.833-3.464 0L3.34 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
            <span>This will overwrite the current database content. This action cannot be undone!</span>
          </div>
        </div>
        <div class="modal-action">
          <button @click="backupToRestore = null" class="btn">Cancel</button>
          <button
            @click="restoreBackup"
            class="btn btn-warning"
            :disabled="backupStore.loading"
          >
            <span v-if="backupStore.loading" class="loading loading-spinner loading-xs"></span>
            Restore Database
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="backupToDelete" class="modal modal-open">
      <div class="modal-box">
        <h3 class="font-bold text-lg">Confirm Deletion</h3>
        <p class="py-4">
          Are you sure you want to delete this backup? This action cannot be undone.
        </p>
        <div class="modal-action">
          <button @click="backupToDelete = null" class="btn">Cancel</button>
          <button
            @click="deleteBackup"
            class="btn btn-error"
            :disabled="backupStore.loading"
          >
            <span v-if="backupStore.loading" class="loading loading-spinner loading-xs"></span>
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTargetsStore } from '@/stores/targets'
import { useBackupsStore } from '@/stores/backups'
import type { Backup } from '@/types'

const route = useRoute()
const router = useRouter()
const targetStore = useTargetsStore()
const backupStore = useBackupsStore()

const backupToRestore = ref<Backup | null>(null)
const backupToDelete = ref<Backup | null>(null)

const targetId = computed(() => Number(route.params.id))
const target = computed(() => targetStore.getTargetById(targetId.value))

const createBackup = async () => {
  const success = await targetStore.createBackup(targetId.value)
  if (success) {
    // Refresh backups list after a short delay
    setTimeout(() => {
      backupStore.fetchBackups(targetId.value)
    }, 1000)
  }
}

const downloadBackup = async (backupId: number) => {
  await backupStore.downloadBackup(backupId)
}

const confirmRestore = (backup: Backup) => {
  backupToRestore.value = backup
}

const restoreBackup = async () => {
  if (backupToRestore.value) {
    const success = await backupStore.restoreBackup(backupToRestore.value.id)
    if (success) {
      backupToRestore.value = null
    }
  }
}

const confirmDelete = (backup: Backup) => {
  backupToDelete.value = backup
}

const deleteBackup = async () => {
  if (backupToDelete.value) {
    const success = await backupStore.deleteBackup(backupToDelete.value.id)
    if (success) {
      backupToDelete.value = null
    }
  }
}

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleString()
}

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

onMounted(async () => {
  await targetStore.fetchTargets()
  
  if (!target.value) {
    router.push('/targets')
    return
  }

  await backupStore.fetchBackups(targetId.value)

  // Auto-refresh backups every 10 seconds
  const interval = setInterval(() => {
    if (route.name === 'TargetBackups') {
      backupStore.fetchBackups(targetId.value)
    } else {
      clearInterval(interval)
    }
  }, 10000)
})
</script>