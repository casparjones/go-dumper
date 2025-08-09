<template>
  <div class="space-y-6">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-base-content mb-2">Manage Backups</h1>
      <p class="text-base-content/70">View, download, and manage all your backups</p>
    </div>

    <!-- Filters and Search -->
    <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <div class="grid gap-4 md:grid-cols-4">
        <div class="form-control">
          <label class="label">
            <span class="label-text">Filter by Configuration</span>
          </label>
          <select v-model="selectedFilter" class="select select-bordered select-sm">
            <option value="">All configurations</option>
            <option v-for="target in targetStore.targets" :key="target.id" :value="target.id">
              {{ target.name }}
            </option>
          </select>
        </div>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">Status Filter</span>
          </label>
          <select v-model="statusFilter" class="select select-bordered select-sm">
            <option value="">All statuses</option>
            <option value="success">Successful</option>
            <option value="failed">Failed</option>
            <option value="running">Running</option>
          </select>
        </div>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">Date Range</span>
          </label>
          <select v-model="dateFilter" class="select select-bordered select-sm">
            <option value="">All time</option>
            <option value="today">Today</option>
            <option value="week">This week</option>
            <option value="month">This month</option>
          </select>
        </div>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">Search</span>
          </label>
          <div class="relative">
            <input 
              v-model="searchQuery" 
              type="text" 
              placeholder="Search backups..." 
              class="input input-bordered input-sm w-full pr-8"
            />
            <svg xmlns="http://www.w3.org/2000/svg" class="absolute right-2 top-2 h-4 w-4 text-base-content/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Bulk Actions -->
    <div v-if="selectedBackups.length > 0" class="bg-base-100 p-4 rounded-lg shadow-sm border border-base-300">
      <div class="flex items-center justify-between">
        <span class="text-sm text-base-content">{{ selectedBackups.length }} backup(s) selected</span>
        <div class="flex gap-2">
          <button @click="downloadSelected" class="btn btn-primary btn-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4-4m0 0L8 8m4-4v12" />
            </svg>
            Download
          </button>
          <button @click="deleteSelected" class="btn btn-error btn-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Delete
          </button>
        </div>
      </div>
    </div>

    <!-- Backups Table -->
    <div class="bg-base-100 rounded-lg shadow-sm border border-base-300 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table table-zebra">
          <thead>
            <tr class="bg-base-200">
              <th>
                <label>
                  <input 
                    type="checkbox" 
                    class="checkbox checkbox-sm"
                    :checked="allSelected"
                    @change="toggleSelectAll"
                  />
                </label>
              </th>
              <th>Configuration</th>
              <th>Database</th>
              <th>Created</th>
              <th>Size</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="backup in filteredBackups" :key="backup.id">
              <td>
                <label>
                  <input 
                    type="checkbox" 
                    class="checkbox checkbox-sm"
                    :checked="selectedBackups.includes(backup.id)"
                    @change="toggleSelectBackup(backup.id)"
                  />
                </label>
              </td>
              <td>
                <div class="font-medium text-sm">{{ backup.targetName }}</div>
                <div class="text-xs text-base-content/70">{{ backup.host }}</div>
              </td>
              <td>
                <div class="font-medium text-sm">{{ backup.database }}</div>
              </td>
              <td>
                <div class="text-sm">{{ formatDate(backup.created_at) }}</div>
                <div class="text-xs text-base-content/70">{{ formatTime(backup.created_at) }}</div>
              </td>
              <td>
                <div class="text-sm font-medium">{{ backup.size }}</div>
              </td>
              <td>
                <div class="flex items-center gap-2">
                  <div 
                    class="w-2 h-2 rounded-full"
                    :class="{
                      'bg-success': backup.status === 'success',
                      'bg-error': backup.status === 'failed',
                      'bg-warning animate-pulse': backup.status === 'running'
                    }"
                  ></div>
                  <span class="text-sm capitalize">{{ backup.status }}</span>
                </div>
              </td>
              <td>
                <div class="flex gap-1">
                  <button 
                    v-if="backup.status === 'success'"
                    @click="downloadBackup(backup.id)"
                    class="btn btn-ghost btn-xs"
                    title="Download"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4-4m0 0L8 8m4-4v12" />
                    </svg>
                  </button>
                  <button 
                    v-if="backup.status === 'success'"
                    @click="restoreBackup(backup.id)"
                    class="btn btn-ghost btn-xs"
                    title="Restore"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                  </button>
                  <button 
                    @click="viewDetails(backup.id)"
                    class="btn btn-ghost btn-xs"
                    title="Details"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </button>
                  <button 
                    @click="deleteBackup(backup.id)"
                    class="btn btn-ghost btn-xs text-error hover:bg-error hover:text-error-content"
                    title="Delete"
                  >
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
      
      <!-- Empty state -->
      <div v-if="filteredBackups.length === 0" class="p-8 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 mx-auto text-base-content/20 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="text-lg font-medium text-base-content mb-2">No backups found</h3>
        <p class="text-base-content/70 mb-4">No backups match your current filter criteria</p>
        <router-link to="/backup" class="btn btn-primary">
          Create First Backup
        </router-link>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex justify-center">
      <div class="join">
        <button 
          @click="currentPage = Math.max(1, currentPage - 1)"
          class="join-item btn"
          :disabled="currentPage === 1"
        >
          «
        </button>
        <button 
          v-for="page in visiblePages" 
          :key="page"
          @click="currentPage = page"
          class="join-item btn"
          :class="{ 'btn-active': currentPage === page }"
        >
          {{ page }}
        </button>
        <button 
          @click="currentPage = Math.min(totalPages, currentPage + 1)"
          class="join-item btn"
          :disabled="currentPage === totalPages"
        >
          »
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTargetsStore } from '@/stores/targets'

const router = useRouter()
const targetStore = useTargetsStore()

// Filters
const selectedFilter = ref<number | ''>('')
const statusFilter = ref('')
const dateFilter = ref('')
const searchQuery = ref('')

// Selection
const selectedBackups = ref<number[]>([])

// Pagination
const currentPage = ref(1)
const itemsPerPage = 10

// Mock backup data - in real app this would come from API
const allBackups = ref([
  {
    id: 1,
    target_id: 1,
    targetName: 'Production Server',
    host: 'prod.example.com',
    database: 'ecommerce_prod',
    created_at: '2024-01-15T14:30:00Z',
    size: '2.4 GB',
    status: 'success'
  },
  {
    id: 2,
    target_id: 2,
    targetName: 'Staging Environment',
    host: 'staging.example.com',
    database: 'staging_db',
    created_at: '2024-01-15T02:00:00Z',
    size: '856 MB',
    status: 'success'
  },
  {
    id: 3,
    target_id: 3,
    targetName: 'Analytics DB',
    host: 'analytics.example.com',
    database: 'analytics_main',
    created_at: '2024-01-14T20:15:00Z',
    size: '0 B',
    status: 'failed'
  },
  {
    id: 4,
    target_id: 1,
    targetName: 'Production Server',
    host: 'prod.example.com',
    database: 'ecommerce_prod',
    created_at: '2024-01-14T14:30:00Z',
    size: '2.3 GB',
    status: 'success'
  },
  {
    id: 5,
    target_id: 4,
    targetName: 'Development DB',
    host: 'dev.example.com',
    database: 'dev_database',
    created_at: '2024-01-14T16:45:00Z',
    size: '124 MB',
    status: 'success'
  }
])

const filteredBackups = computed(() => {
  let filtered = [...allBackups.value]
  
  // Apply filters
  if (selectedFilter.value) {
    filtered = filtered.filter(b => b.target_id === selectedFilter.value)
  }
  
  if (statusFilter.value) {
    filtered = filtered.filter(b => b.status === statusFilter.value)
  }
  
  if (dateFilter.value) {
    const now = new Date()
    const filterDate = new Date()
    
    switch (dateFilter.value) {
      case 'today':
        filterDate.setHours(0, 0, 0, 0)
        break
      case 'week':
        filterDate.setDate(now.getDate() - 7)
        break
      case 'month':
        filterDate.setMonth(now.getMonth() - 1)
        break
    }
    
    filtered = filtered.filter(b => new Date(b.created_at) >= filterDate)
  }
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(b => 
      b.targetName.toLowerCase().includes(query) ||
      b.database.toLowerCase().includes(query) ||
      b.host.toLowerCase().includes(query)
    )
  }
  
  // Sort by creation date (newest first)
  filtered.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  
  // Apply pagination
  const startIndex = (currentPage.value - 1) * itemsPerPage
  const endIndex = startIndex + itemsPerPage
  
  return filtered.slice(startIndex, endIndex)
})

const totalPages = computed(() => {
  const totalItems = allBackups.value.filter(backup => {
    // Apply same filters for total count
    let include = true
    
    if (selectedFilter.value && backup.target_id !== selectedFilter.value) include = false
    if (statusFilter.value && backup.status !== statusFilter.value) include = false
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      include = backup.targetName.toLowerCase().includes(query) ||
                backup.database.toLowerCase().includes(query) ||
                backup.host.toLowerCase().includes(query)
    }
    
    return include
  }).length
  
  return Math.ceil(totalItems / itemsPerPage)
})

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, start + 4)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  return pages
})

const allSelected = computed(() => {
  return filteredBackups.value.length > 0 && 
         filteredBackups.value.every(b => selectedBackups.value.includes(b.id))
})

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString()
}

const formatTime = (dateString: string): string => {
  return new Date(dateString).toLocaleTimeString()
}

const toggleSelectAll = () => {
  if (allSelected.value) {
    selectedBackups.value = selectedBackups.value.filter(
      id => !filteredBackups.value.some(b => b.id === id)
    )
  } else {
    const currentPageIds = filteredBackups.value.map(b => b.id)
    selectedBackups.value = [...new Set([...selectedBackups.value, ...currentPageIds])]
  }
}

const toggleSelectBackup = (backupId: number) => {
  const index = selectedBackups.value.indexOf(backupId)
  if (index > -1) {
    selectedBackups.value.splice(index, 1)
  } else {
    selectedBackups.value.push(backupId)
  }
}

const downloadBackup = async (backupId: number) => {
  // In real app: call API to download backup
  console.log('Downloading backup:', backupId)
}

const downloadSelected = async () => {
  // In real app: call API to download multiple backups
  console.log('Downloading selected backups:', selectedBackups.value)
}

const restoreBackup = (backupId: number) => {
  // Navigate to restore page with pre-selected backup
  router.push({ name: 'Restore', query: { backup: backupId } })
}

const viewDetails = (backupId: number) => {
  // In real app: show backup details modal or navigate to details page
  console.log('Viewing backup details:', backupId)
}

const deleteBackup = async (backupId: number) => {
  if (confirm('Are you sure you want to delete this backup? This action cannot be undone.')) {
    // In real app: call API to delete backup
    const index = allBackups.value.findIndex(b => b.id === backupId)
    if (index > -1) {
      allBackups.value.splice(index, 1)
    }
    
    // Remove from selection if selected
    const selIndex = selectedBackups.value.indexOf(backupId)
    if (selIndex > -1) {
      selectedBackups.value.splice(selIndex, 1)
    }
  }
}

const deleteSelected = async () => {
  if (confirm(`Are you sure you want to delete ${selectedBackups.value.length} selected backup(s)? This action cannot be undone.`)) {
    // In real app: call API to delete multiple backups
    allBackups.value = allBackups.value.filter(b => !selectedBackups.value.includes(b.id))
    selectedBackups.value = []
  }
}

onMounted(() => {
  targetStore.fetchTargets()
})
</script>