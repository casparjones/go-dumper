<template>
  <div class="mx-auto max-w-3xl space-y-6">
    <!-- Header -->
    <div class="mb-2">
      <h1 class="text-3xl font-bold mb-1">Create Backup</h1>
      <p class="text-base-content/70">Create manual backups for your configured databases</p>
    </div>

    <!-- No Configuration Selected Warning -->
    <div v-if="!selectedTargetInfo" class="card bg-warning/10 border border-warning/20 shadow-sm">
      <div class="card-body">
        <div class="flex items-center gap-3">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-warning" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
          <div>
            <h3 class="font-semibold text-warning">No Configuration Selected</h3>
            <p class="text-sm text-base-content/70">Please select a configuration in the sidebar to create a backup.</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Selected Configuration Info & Backup Options -->
    <div v-if="selectedTargetInfo" class="card bg-base-100 border border-base-200 shadow-sm">
      <div class="card-body gap-4">
        <h3 class="card-title">Selected Configuration</h3>

        <div class="grid gap-4 md:grid-cols-2">
          <!-- Configuration Info -->
          <div class="space-y-3">
            <div>
              <div class="text-sm text-base-content/70">Configuration</div>
              <div class="font-medium">{{ selectedTargetInfo.name }}</div>
            </div>
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

          <!-- Backup Options -->
          <div class="form-control">
            <label class="label">
              <span class="label-text">Backup Options</span>
            </label>

            <div class="space-y-3 flex flex-col">
              <label class="label cursor-pointer flex-row items-start gap-3">
                <input type="checkbox" v-model="compressBackup" class="toggle toggle-primary" />
                <span class="label-text">Compress Backup</span>
              </label>

              <label class="label cursor-pointer flex-row items-start gap-3 opacity-60">
                <input type="checkbox" v-model="includeStructure" class="toggle toggle-primary" checked disabled />
                <span class="label-text">Include table structure</span>
              </label>

              <label class="label cursor-pointer flex-row items-start gap-3">
                <input type="checkbox" v-model="includeData" class="toggle toggle-primary" />
                <span class="label-text">Include data</span>
              </label>

              <label class="label cursor-pointer flex-row items-start gap-3">
                <input type="checkbox" v-model="selectSpecificDatabases" class="toggle toggle-primary" />
                <span class="label-text">Select specific databases</span>
              </label>

              <label class="label cursor-pointer flex-row items-start gap-3">
                <input type="checkbox" v-model="createScheduledBackup" class="toggle toggle-primary" />
                <span class="label-text">Create a scheduled backup</span>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Database Selection Override -->
    <div v-if="selectedTargetInfo && selectSpecificDatabases" class="card bg-base-100 border border-base-200 shadow-sm">
      <div class="card-body gap-4">
        <h3 class="card-title">Select Databases</h3>
        <p class="text-sm text-base-content/70">Override the default configuration to backup only specific databases.</p>
        
        <div class="grid gap-2">
          <div v-if="availableDatabases.length === 0" class="text-sm text-base-content/70">
            Loading databases...
          </div>
          <div v-else class="space-y-2">
            <label v-for="db in availableDatabases" :key="db" class="label cursor-pointer justify-start gap-3">
              <input type="checkbox" v-model="selectedDatabases" :value="db" class="checkbox checkbox-primary" />
              <span class="label-text">{{ db }}</span>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- Schedule Options -->
    <div v-if="selectedTargetInfo && createScheduledBackup" class="card bg-base-100 border border-base-200 shadow-sm">
      <div class="card-body gap-4">
        <h3 class="card-title">Schedule Settings</h3>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">Frequency</span>
          </label>
          <select v-model="scheduleFrequency" class="select select-bordered w-full">
            <option value="hourly">Hourly</option>
            <option value="daily">Daily</option>
            <option value="weekly">Weekly</option>
            <option value="monthly">Monthly</option>
            <option value="yearly">Yearly</option>
          </select>
        </div>

        <!-- Minutes (Available for all frequencies) -->
        <div class="form-control">
          <label class="label">
            <span class="label-text">Minutes</span>
          </label>
          <div class="grid grid-cols-4 gap-2">
            <label v-for="minute in hourlyMinutes" :key="minute" class="label cursor-pointer justify-start gap-2">
              <input type="checkbox" v-model="selectedMinutes" :value="minute" class="checkbox checkbox-sm checkbox-primary" />
              <span class="label-text text-sm">:{{ minute.toString().padStart(2, '0') }}</span>
            </label>
          </div>
        </div>

        <!-- Hours (Available for daily, weekly, monthly, yearly) -->
        <div v-if="scheduleFrequency !== 'hourly'" class="form-control">
          <label class="label">
            <span class="label-text">Hours (0-23)</span>
          </label>
          <div class="grid grid-cols-6 gap-2">
            <label v-for="hour in 24" :key="hour-1" class="label cursor-pointer justify-start gap-2">
              <input type="checkbox" v-model="selectedHours" :value="hour-1" class="checkbox checkbox-sm checkbox-primary" />
              <span class="label-text text-sm">{{ (hour-1).toString().padStart(2, '0') }}</span>
            </label>
          </div>
        </div>

        <!-- Days of Week (Available for weekly, monthly, yearly) -->
        <div v-if="['weekly', 'monthly', 'yearly'].includes(scheduleFrequency)" class="form-control">
          <label class="label">
            <span class="label-text">Days of Week</span>
          </label>
          <div class="grid grid-cols-7 gap-2">
            <label v-for="day in weekdays" :key="day.value" class="label cursor-pointer justify-start gap-2">
              <input type="checkbox" v-model="selectedWeekdays" :value="day.value" class="checkbox checkbox-sm checkbox-primary" />
              <span class="label-text text-sm">{{ day.label.substring(0, 3) }}</span>
            </label>
          </div>
        </div>

        <!-- Days of Month (Available for monthly, yearly) -->
        <div v-if="['monthly', 'yearly'].includes(scheduleFrequency)" class="form-control">
          <label class="label">
            <span class="label-text">Days of Month</span>
          </label>
          <div class="grid grid-cols-8 gap-1">
            <label v-for="day in 31" :key="day" class="label cursor-pointer justify-start gap-1">
              <input type="checkbox" v-model="selectedDaysOfMonth" :value="day" class="checkbox checkbox-xs checkbox-primary" />
              <span class="label-text text-xs">{{ day }}</span>
            </label>
          </div>
        </div>

        <!-- Months (Available for yearly) -->
        <div v-if="scheduleFrequency === 'yearly'" class="form-control">
          <label class="label">
            <span class="label-text">Months</span>
          </label>
          <div class="grid grid-cols-4 gap-2">
            <label v-for="month in months" :key="month.value" class="label cursor-pointer justify-start gap-2">
              <input type="checkbox" v-model="selectedMonths" :value="month.value" class="checkbox checkbox-sm checkbox-primary" />
              <span class="label-text text-sm">{{ month.label.substring(0, 3) }}</span>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Backup or Job -->
    <div v-if="selectedTargetInfo" class="card bg-base-100 border border-base-200 shadow-sm">
      <div class="card-body gap-4">
        <h3 class="card-title">{{ createScheduledBackup ? 'Create Scheduled Job' : 'Create Backup' }}</h3>

        <!-- Job Name and Description (only for scheduled backups) -->
        <div v-if="createScheduledBackup" class="grid gap-4 md:grid-cols-2">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Job Name *</span>
            </label>
            <input v-model="jobName" type="text" placeholder="e.g. Daily Production Backup" 
                   class="input input-bordered w-full" required />
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Description</span>
            </label>
            <input v-model="jobDescription" type="text" placeholder="Optional description" 
                   class="input input-bordered w-full" />
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-3">
          <button
              @click="createBackupOrJob"
              class="btn btn-primary"
              :disabled="targetStore.loading || (createScheduledBackup && !jobName.trim())"
          >
            <span v-if="targetStore.loading" class="loading loading-spinner loading-sm"></span>
            <span v-else class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none"
                   viewBox="0 0 24 24" stroke="currentColor">
                <path v-if="!createScheduledBackup" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              {{ createScheduledBackup ? 'Save Job' : 'Create Backup Now' }}
            </span>
          </button>

          <div class="text-sm text-base-content/70">
            <div>Compression: <span class="font-medium">{{ compressBackup ? 'Enabled' : 'Disabled' }}</span></div>
            <div>Include Data: <span class="font-medium">{{ includeData ? 'Yes' : 'No' }}</span></div>
            <div v-if="selectSpecificDatabases && selectedDatabases.length">
              Databases: <span class="font-medium">{{ selectedDatabases.length }} selected</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Backups for Selected Config -->
    <div
        v-if="selectedTargetInfo"
        class="card bg-base-100 border border-base-200 shadow-sm"
    >
      <div class="card-body gap-2">
        <h3 class="card-title">Recent Backups</h3>
        <div class="text-sm text-base-content/70">
          Last 5 backups for "{{ selectedTargetInfo.name }}" will be shown here...
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useTargetsStore } from '@/stores/targets'
import { useAppStore } from '@/stores/app'
import { useToastStore } from '@/stores/toasts'
import { jobsApi } from '@/services/api'
import type { CreateJobRequest } from '@/services/api'

const targetStore = useTargetsStore()
const appStore = useAppStore()
const toastStore = useToastStore()
const router = useRouter()

// Basic backup options
const compressBackup = ref(true)
const includeStructure = ref(true)
const includeData = ref(true)

// Advanced options
const selectSpecificDatabases = ref(false)
const createScheduledBackup = ref(false)

// Job details
const jobName = ref('')
const jobDescription = ref('')

// Database selection
const availableDatabases = ref<string[]>([])
const selectedDatabases = ref<string[]>([])

// Schedule options
const scheduleFrequency = ref('daily')
const selectedMinutes = ref<number[]>([])
const selectedHours = ref<number[]>([])
const selectedWeekdays = ref<number[]>([])
const selectedDaysOfMonth = ref<number[]>([])
const selectedMonths = ref<number[]>([])

// Static data for UI
const hourlyMinutes = [0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55]

const weekdays = [
  { value: 1, label: 'Monday' },
  { value: 2, label: 'Tuesday' },
  { value: 3, label: 'Wednesday' },
  { value: 4, label: 'Thursday' },
  { value: 5, label: 'Friday' },
  { value: 6, label: 'Saturday' },
  { value: 0, label: 'Sunday' }
]

const months = [
  { value: 1, label: 'January' },
  { value: 2, label: 'February' },
  { value: 3, label: 'March' },
  { value: 4, label: 'April' },
  { value: 5, label: 'May' },
  { value: 6, label: 'June' },
  { value: 7, label: 'July' },
  { value: 8, label: 'August' },
  { value: 9, label: 'September' },
  { value: 10, label: 'October' },
  { value: 11, label: 'November' },
  { value: 12, label: 'December' }
]

const selectedTargetInfo = computed(() => {
  return appStore.selectedConfig
})


// Load available databases when target changes or specific database selection is enabled
watch([() => appStore.selectedConfigId, selectSpecificDatabases], async ([targetId, enabled]) => {
  if (targetId && enabled) {
    await loadAvailableDatabases(targetId)
  }
}, { immediate: true })

const loadAvailableDatabases = async (targetId: number) => {
  try {
    // This would need to be implemented in the API
    // For now, using dummy data
    availableDatabases.value = ['database1', 'database2', 'database3', 'test_db']
  } catch (error) {
    console.error('Failed to load databases:', error)
  }
}

const createBackupOrJob = async () => {
  if (!appStore.selectedConfigId) return
  
  const backupOptions = {
    compress: compressBackup.value,
    include_structure: includeStructure.value,
    include_data: includeData.value,
    databases: selectSpecificDatabases.value && selectedDatabases.value.length > 0 
      ? selectedDatabases.value 
      : undefined
  }
  
  if (createScheduledBackup.value) {
    // Create a scheduled job
    const jobRequest: CreateJobRequest = {
      target_id: appStore.selectedConfigId,
      name: jobName.value.trim(),
      description: jobDescription.value.trim(),
      schedule_config: {
        frequency: scheduleFrequency.value,
        minutes: selectedMinutes.value.length > 0 ? selectedMinutes.value : undefined,
        hours: selectedHours.value.length > 0 ? selectedHours.value : undefined,
        weekdays: selectedWeekdays.value.length > 0 ? selectedWeekdays.value : undefined,
        days_of_month: selectedDaysOfMonth.value.length > 0 ? selectedDaysOfMonth.value : undefined,
        months: selectedMonths.value.length > 0 ? selectedMonths.value : undefined
      },
      backup_options: backupOptions,
      meta_config: {} // For future extensions
    }
    
    try {
      const createdJob = await jobsApi.create(jobRequest)
      toastStore.addToast('success', 'Job Created', `Scheduled job "${createdJob.name}" created successfully`)
      
      // Reset form
      jobName.value = ''
      jobDescription.value = ''
      createScheduledBackup.value = false
      
      // Navigate to jobs page
      router.push('/schedule-jobs')
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || 'Failed to create scheduled job'
      toastStore.addToast('error', 'Error', errorMessage)
    }
  } else {
    // Create immediate backup
    try {
      await targetStore.createBackup(appStore.selectedConfigId)
    } catch (error) {
      console.error('Failed to create backup:', error)
    }
  }
}

onMounted(() => {
  targetStore.fetchTargets()
})
</script>
