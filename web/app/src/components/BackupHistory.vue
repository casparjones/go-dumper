<template>
  <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
    <h3 class="text-lg font-semibold text-base-content mb-4 flex items-center gap-2">
      <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      Recent Backup History
    </h3>
    
    <div class="space-y-3">
      <div 
        v-for="backup in recentBackups" 
        :key="backup.id"
        class="flex items-center justify-between p-3 bg-base-200 rounded-lg hover:bg-base-300 transition-colors"
      >
        <div class="flex items-center gap-3">
          <div class="flex-shrink-0">
            <div 
              class="w-2 h-2 rounded-full"
              :class="{
                'bg-success': backup.status === 'success',
                'bg-error': backup.status === 'failed',
                'bg-warning animate-pulse': backup.status === 'running'
              }"
            ></div>
          </div>
          
          <div>
            <div class="font-medium text-sm text-base-content">{{ backup.targetName }}</div>
            <div class="text-xs text-base-content/70">{{ backup.database }}</div>
          </div>
        </div>
        
        <div class="text-right">
          <div class="text-xs text-base-content/80">{{ formatDate(backup.createdAt) }}</div>
          <div class="text-xs text-base-content/60">{{ backup.size }}</div>
        </div>
      </div>
      
      <div v-if="recentBackups.length === 0" class="text-center py-8">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 mx-auto text-base-content/20 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <p class="text-sm text-base-content/50">No backup history available</p>
      </div>
    </div>
    
    <div class="mt-4 pt-3 border-t border-base-300">
      <router-link 
        to="/manage" 
        class="text-sm text-primary hover:text-primary-focus font-medium flex items-center gap-1"
      >
        View all backups
        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface BackupHistoryItem {
  id: number
  targetName: string
  database: string
  status: 'success' | 'failed' | 'running'
  createdAt: string
  size: string
}

const recentBackups = ref<BackupHistoryItem[]>([
  {
    id: 1,
    targetName: 'Production Server',
    database: 'ecommerce_prod',
    status: 'success',
    createdAt: '2024-01-15T14:30:00Z',
    size: '2.4 GB'
  },
  {
    id: 2,
    targetName: 'Staging Environment',
    database: 'staging_db',
    status: 'success',
    createdAt: '2024-01-15T02:00:00Z',
    size: '856 MB'
  },
  {
    id: 3,
    targetName: 'Analytics DB',
    database: 'analytics_main',
    status: 'failed',
    createdAt: '2024-01-14T20:15:00Z',
    size: '0 B'
  },
  {
    id: 4,
    targetName: 'Development DB',
    database: 'dev_database',
    status: 'success',
    createdAt: '2024-01-14T16:45:00Z',
    size: '124 MB'
  },
  {
    id: 5,
    targetName: 'User Management',
    database: 'user_data',
    status: 'success',
    createdAt: '2024-01-14T08:30:00Z',
    size: '1.2 GB'
  }
])

const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInHours = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60))
  
  if (diffInHours < 1) {
    return 'Just now'
  } else if (diffInHours < 24) {
    return `${diffInHours}h ago`
  } else if (diffInHours < 48) {
    return 'Yesterday'
  } else {
    return date.toLocaleDateString()
  }
}

onMounted(async () => {
  // In a real implementation, you would fetch recent backup history from the API
  try {
    // const response = await backupsApi.getRecent(5)
    // recentBackups.value = response.data
  } catch (error) {
    console.error('Failed to fetch backup history:', error)
  }
})
</script>