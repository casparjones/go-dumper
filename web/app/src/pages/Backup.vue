<template>
  <div class="space-y-6">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-base-content mb-2">Create Backup</h1>
      <p class="text-base-content/70">Create manual backups for your configured databases</p>
    </div>

    <!-- Configuration Selection -->
    <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4">Select Configuration</h3>
      
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
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">Backup Options</span>
          </label>
          <div class="space-y-2">
            <label class="cursor-pointer label justify-start gap-2">
              <input type="checkbox" v-model="compressBackup" class="checkbox checkbox-primary checkbox-sm" />
              <span class="label-text">Compress backup</span>
            </label>
            <label class="cursor-pointer label justify-start gap-2">
              <input type="checkbox" v-model="includeStructure" class="checkbox checkbox-primary checkbox-sm" checked disabled />
              <span class="label-text">Include table structure</span>
            </label>
            <label class="cursor-pointer label justify-start gap-2">
              <input type="checkbox" v-model="includeData" class="checkbox checkbox-primary checkbox-sm" checked />
              <span class="label-text">Include data</span>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- Selected Configuration Details -->
    <div v-if="selectedTargetInfo" class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4">Configuration Details</h3>
      <div class="grid gap-4 md:grid-cols-3">
        <div>
          <div class="text-sm text-base-content/70">Host</div>
          <div class="font-medium">{{ selectedTargetInfo.host }}:{{ selectedTargetInfo.port }}</div>
        </div>
        <div>
          <div class="text-sm text-base-content/70">Database</div>
          <div class="font-medium">{{ selectedTargetInfo.db_name }}</div>
        </div>
        <div>
          <div class="text-sm text-base-content/70">User</div>
          <div class="font-medium">{{ selectedTargetInfo.user }}</div>
        </div>
      </div>
    </div>

    <!-- Create Backup Action -->
    <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4">Create Backup</h3>
      
      <div class="flex gap-4 items-center">
        <button 
          @click="createBackup" 
          class="btn btn-primary"
          :disabled="!selectedConfig || targetStore.loading"
        >
          <span v-if="targetStore.loading" class="loading loading-spinner loading-sm"></span>
          <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
          </svg>
          Create Backup Now
        </button>
        
        <div v-if="!selectedConfig" class="text-sm text-base-content/50">
          Please select a configuration first
        </div>
      </div>
    </div>

    <!-- Recent Backups for Selected Config -->
    <div v-if="selectedConfig" class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4">Recent Backups</h3>
      <div class="text-sm text-base-content/70">
        Last 5 backups for this configuration will be shown here...
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTargetsStore } from '@/stores/targets'

const targetStore = useTargetsStore()

const selectedConfig = ref<number | ''>('')
const compressBackup = ref(true)
const includeStructure = ref(true)
const includeData = ref(true)

const selectedTargetInfo = computed(() => {
  if (!selectedConfig.value) return null
  return targetStore.targets.find(t => t.id === selectedConfig.value)
})

const createBackup = async () => {
  if (!selectedConfig.value) return
  
  await targetStore.createBackup(selectedConfig.value as number)
}

onMounted(() => {
  targetStore.fetchTargets()
})
</script>