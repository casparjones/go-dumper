<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-3xl font-bold text-base-content">Backup Targets</h1>
        <p class="text-base-content/70">Manage your MySQL/MariaDB backup targets</p>
      </div>
      <router-link to="/targets/new" class="btn btn-primary">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
        Add Target
      </router-link>
    </div>

    <div v-if="targetStore.loading" class="flex justify-center py-8">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="targetStore.targets.length === 0" class="text-center py-12">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-16 h-16 mx-auto text-base-content/30 mb-4">
        <path stroke-linecap="round" stroke-linejoin="round" d="M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375m16.5 0v3.75m-16.5-3.75v3.75m16.5 0v3.75C20.25 16.153 16.556 18 12 18s-8.25-1.847-8.25-4.125v-3.75m16.5 0c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125" />
      </svg>
      <h3 class="text-lg font-semibold text-base-content mb-2">No targets configured</h3>
      <p class="text-base-content/70 mb-4">Get started by adding your first backup target</p>
      <router-link to="/targets/new" class="btn btn-primary">Add Target</router-link>
    </div>

    <div v-else class="grid gap-6">
      <div
        v-for="target in targetStore.targets"
        :key="target.id"
        class="card bg-base-200 shadow"
      >
        <div class="card-body">
          <div class="flex justify-between items-start">
            <div class="flex-1">
              <h3 class="card-title text-lg">{{ target.name }}</h3>
              <div class="text-sm text-base-content/70 space-y-1">
                <div>{{ target.host }}:{{ target.port }} / {{ target.db_name }}</div>
                <div>User: {{ target.user }}</div>
                <div v-if="target.comment">{{ target.comment }}</div>
                <div v-if="target.schedule_time">
                  Scheduled: Daily at {{ target.schedule_time }} UTC
                </div>
              </div>
              <div class="flex gap-2 mt-2">
                <div class="badge badge-outline">{{ target.retention_days }} days retention</div>
                <div v-if="target.auto_compress" class="badge badge-outline">Compressed</div>
              </div>
            </div>
            <div class="dropdown dropdown-end">
              <div tabindex="0" role="button" class="btn btn-ghost btn-sm">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 12.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 18.75a.75.75 0 110-1.5.75.75 0 010 1.5z" />
                </svg>
              </div>
              <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
                <li>
                  <button
                    @click="createBackup(target.id)"
                    :disabled="targetStore.loading"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m.75 12l3 3m0 0l3-3m-3 3v-6m-1.06-9.75L5.5 8.25M3 14.25V17.25A2.25 2.25 0 005.25 19.5h13.5A2.25 2.25 0 0021 17.25v-3A2.25 2.25 0 0018.75 12H5.25A2.25 2.25 0 003 14.25z" />
                    </svg>
                    Create Backup
                  </button>
                </li>
                <li>
                  <router-link :to="`/targets/${target.id}/backups`">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12.75V12A2.25 2.25 0 014.5 9.75h15A2.25 2.25 0 0121.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z" />
                    </svg>
                    View Backups
                  </router-link>
                </li>
                <li>
                  <router-link :to="`/targets/${target.id}/edit`">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                    </svg>
                    Edit
                  </router-link>
                </li>
                <li>
                  <button
                    @click="confirmDelete(target)"
                    class="text-error"
                    :disabled="targetStore.loading"
                  >
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

    <!-- Delete Confirmation Modal -->
    <div v-if="targetToDelete" class="modal modal-open">
      <div class="modal-box">
        <h3 class="font-bold text-lg">Confirm Deletion</h3>
        <p class="py-4">
          Are you sure you want to delete the target <strong>{{ targetToDelete.name }}</strong>?
          This will also delete all associated backups and cannot be undone.
        </p>
        <div class="modal-action">
          <button @click="targetToDelete = null" class="btn">Cancel</button>
          <button
            @click="deleteTarget"
            class="btn btn-error"
            :disabled="targetStore.loading"
          >
            <span v-if="targetStore.loading" class="loading loading-spinner loading-xs"></span>
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useTargetsStore } from '@/stores/targets'
import type { Target } from '@/types'

const targetStore = useTargetsStore()
const targetToDelete = ref<Target | null>(null)

const createBackup = async (targetId: number) => {
  await targetStore.createBackup(targetId)
}

const confirmDelete = (target: Target) => {
  targetToDelete.value = target
}

const deleteTarget = async () => {
  if (targetToDelete.value) {
    const success = await targetStore.deleteTarget(targetToDelete.value.id)
    if (success) {
      targetToDelete.value = null
    }
  }
}

onMounted(() => {
  targetStore.fetchTargets()
})
</script>