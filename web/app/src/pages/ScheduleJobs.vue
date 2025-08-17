<template>
  <div class="mx-auto max-w-6xl space-y-6">
    <!-- Header -->
    <div class="mb-2">
      <h1 class="text-3xl font-bold mb-1">Schedule Jobs</h1>
      <p class="text-base-content/70">Manage automated backup schedules for your databases</p>
    </div>

    <!-- Jobs Overview -->
    <div class="grid gap-4 md:grid-cols-4">
      <!-- Total Jobs -->
      <div class="stat bg-base-100 border border-base-200 shadow-sm rounded-lg">
        <div class="stat-title">Total Jobs</div>
        <div class="stat-value text-primary">{{ jobs.length }}</div>
        <div class="stat-desc">{{ activeJobs.length }} active</div>
      </div>

      <!-- Running Jobs -->
      <div class="stat bg-base-100 border border-base-200 shadow-sm rounded-lg">
        <div class="stat-title">Running</div>
        <div class="stat-value text-warning">{{ runningJobs.length }}</div>
        <div class="stat-desc">Currently executing</div>
      </div>

      <!-- Last 24h Success -->
      <div class="stat bg-base-100 border border-base-200 shadow-sm rounded-lg">
        <div class="stat-title">Success (24h)</div>
        <div class="stat-value text-success">{{ recentSuccessJobs.length }}</div>
        <div class="stat-desc">Completed successfully</div>
      </div>

      <!-- Failed Jobs -->
      <div class="stat bg-base-100 border border-base-200 shadow-sm rounded-lg">
        <div class="stat-title">Failed</div>
        <div class="stat-value text-error">{{ failedJobs.length }}</div>
        <div class="stat-desc">Need attention</div>
      </div>
    </div>

    <!-- Jobs Table -->
    <div class="card bg-base-100 border border-base-200 shadow-sm">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <h3 class="card-title">All Jobs</h3>
          <div class="flex gap-2">
            <select v-model="statusFilter" class="select select-bordered select-sm">
              <option value="">All Status</option>
              <option value="active">Active Only</option>
              <option value="inactive">Inactive Only</option>
            </select>
            <router-link to="/backup" class="btn btn-primary btn-sm">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Create Job
            </router-link>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="table w-full">
            <thead>
              <tr>
                <th>Name</th>
                <th>Target</th>
                <th>Schedule</th>
                <th>Last Run</th>
                <th>Next Run</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="filteredJobs.length === 0">
                <td colspan="7" class="text-center text-base-content/70 py-8">
                  <div class="flex flex-col items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 text-base-content/30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span>No scheduled jobs found</span>
                    <router-link to="/backup" class="btn btn-sm btn-primary">Create your first job</router-link>
                  </div>
                </td>
              </tr>
              <tr v-for="job in filteredJobs" :key="job.id" class="hover">
                <td>
                  <div>
                    <div class="font-medium">{{ job.name }}</div>
                    <div v-if="job.description" class="text-sm text-base-content/70">{{ job.description }}</div>
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-2">
                    <div class="w-2 h-2 rounded-full" :class="job.target ? 'bg-success' : 'bg-error'"></div>
                    <span>{{ job.target?.name || 'Unknown' }}</span>
                  </div>
                </td>
                <td>
                  <div class="text-sm">
                    <div class="font-medium">{{ formatSchedule(job.schedule_config) }}</div>
                  </div>
                </td>
                <td>
                  <div v-if="job.last_run_at" class="text-sm">
                    <div class="font-medium">{{ formatDateTime(job.last_run_at) }}</div>
                    <div class="flex items-center gap-1">
                      <div class="w-2 h-2 rounded-full" :class="getStatusColor(job.last_run_status)"></div>
                      <span class="text-xs">{{ job.last_run_status || 'pending' }}</span>
                    </div>
                  </div>
                  <span v-else class="text-base-content/50">Never</span>
                </td>
                <td>
                  <div v-if="job.next_run_at" class="text-sm">
                    <div class="font-medium">{{ formatDateTime(job.next_run_at) }}</div>
                    <div class="text-xs text-base-content/70">{{ getTimeUntil(job.next_run_at) }}</div>
                  </div>
                  <span v-else class="text-base-content/50">Not scheduled</span>
                </td>
                <td>
                  <div class="flex items-center gap-2">
                    <input type="checkbox" :checked="job.is_active" @change="toggleJobStatus(job)" class="toggle toggle-sm toggle-success" />
                    <span class="text-sm">{{ job.is_active ? 'Active' : 'Inactive' }}</span>
                  </div>
                </td>
                <td>
                  <div class="flex gap-1">
                    <button @click="runJobNow(job)" class="btn btn-ghost btn-xs" title="Run Now">
                      <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6 4h1m4 0h1m2-7a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                    </button>
                    <button @click="editJob(job)" class="btn btn-ghost btn-xs" title="Edit">
                      <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </button>
                    <button @click="deleteJob(job)" class="btn btn-ghost btn-xs text-error" title="Delete">
                      <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTargetsStore } from '@/stores/targets'
import { useToastStore } from '@/stores/toasts'
import { jobsApi } from '@/services/api'
import type { ScheduleJob } from '@/services/api'

const targetStore = useTargetsStore()
const toastStore = useToastStore()

// State
const jobs = ref<ScheduleJob[]>([])
const statusFilter = ref('')
const loading = ref(false)

// Computed
const activeJobs = computed(() => jobs.value.filter(job => job.is_active))
const runningJobs = computed(() => jobs.value.filter(job => job.last_run_status === 'running'))
const failedJobs = computed(() => jobs.value.filter(job => job.last_run_status === 'failed'))
const recentSuccessJobs = computed(() => {
  const oneDayAgo = new Date(Date.now() - 24 * 60 * 60 * 1000)
  return jobs.value.filter(job => 
    job.last_run_status === 'success' && 
    job.last_run_at && 
    new Date(job.last_run_at) > oneDayAgo
  )
})

const filteredJobs = computed(() => {
  if (!statusFilter.value) return jobs.value
  if (statusFilter.value === 'active') return jobs.value.filter(job => job.is_active)
  if (statusFilter.value === 'inactive') return jobs.value.filter(job => !job.is_active)
  return jobs.value
})

// Methods
const loadJobs = async () => {
  loading.value = true
  try {
    jobs.value = await jobsApi.getAll()
  } catch (error) {
    console.error('Failed to load jobs:', error)
    toastStore.addToast('error', 'Error', 'Failed to load scheduled jobs')
  } finally {
    loading.value = false
  }
}

const formatSchedule = (scheduleConfigStr: string) => {
  try {
    const config = JSON.parse(scheduleConfigStr)
    let result = config.frequency
    
    if (config.hours && config.hours.length > 0) {
      result += ` at ${config.hours[0]}:${(config.minutes?.[0] || 0).toString().padStart(2, '0')}`
    } else if (config.minutes && config.minutes.length > 0) {
      result += ` at :${config.minutes[0].toString().padStart(2, '0')}`
    }
    
    return result
  } catch {
    return 'Invalid schedule'
  }
}

const formatDateTime = (dateStr: string) => {
  return new Date(dateStr).toLocaleString()
}

const getTimeUntil = (dateStr: string) => {
  const now = new Date()
  const target = new Date(dateStr)
  const diff = target.getTime() - now.getTime()
  
  if (diff < 0) return 'Overdue'
  
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  
  if (hours > 24) {
    const days = Math.floor(hours / 24)
    return `in ${days}d ${hours % 24}h`
  }
  return `in ${hours}h ${minutes}m`
}

const getStatusColor = (status?: string) => {
  switch (status) {
    case 'success': return 'bg-success'
    case 'failed': return 'bg-error'
    case 'running': return 'bg-warning'
    default: return 'bg-base-300'
  }
}

const toggleJobStatus = async (job: ScheduleJob) => {
  try {
    const updatedJob = await jobsApi.update(job.id, {
      name: job.name,
      description: job.description,
      is_active: !job.is_active,
      schedule_config: JSON.parse(job.schedule_config),
      backup_options: JSON.parse(job.backup_options),
      meta_config: JSON.parse(job.meta_config || '{}')
    })
    
    // Update local job
    const index = jobs.value.findIndex(j => j.id === job.id)
    if (index !== -1) {
      jobs.value[index] = updatedJob
    }
    
    toastStore.addToast('success', 'Success', 
      `Job "${job.name}" ${updatedJob.is_active ? 'activated' : 'deactivated'}`)
  } catch (error: any) {
    const errorMessage = error.response?.data?.error || 'Failed to toggle job status'
    toastStore.addToast('error', 'Error', errorMessage)
  }
}

const runJobNow = async (job: ScheduleJob) => {
  try {
    const result = await jobsApi.runNow(job.id)
    toastStore.addToast('success', 'Job Started', result.message)
    
    // Reload jobs to get updated status
    await loadJobs()
  } catch (error: any) {
    const errorMessage = error.response?.data?.error || 'Failed to run job'
    toastStore.addToast('error', 'Error', errorMessage)
  }
}

const editJob = (job: ScheduleJob) => {
  // TODO: Navigate to edit page or open modal
  console.log(`Editing job ${job.id}`)
  toastStore.addToast('info', 'Info', 'Edit functionality will be implemented soon')
}

const deleteJob = async (job: ScheduleJob) => {
  if (confirm(`Are you sure you want to delete the job "${job.name}"?`)) {
    try {
      await jobsApi.delete(job.id)
      jobs.value = jobs.value.filter(j => j.id !== job.id)
      toastStore.addToast('success', 'Success', `Job "${job.name}" deleted successfully`)
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || 'Failed to delete job'
      toastStore.addToast('error', 'Error', errorMessage)
    }
  }
}

onMounted(() => {
  loadJobs()
  targetStore.fetchTargets()
})
</script>