<template>
  <div class="space-y-6">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-base-content mb-2">System Logs</h1>
      <p class="text-base-content/70">View application logs and system activity</p>
    </div>

    <!-- Log Filters -->
    <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <div class="grid gap-4 md:grid-cols-5">
        <div class="form-control">
          <label class="label">
            <span class="label-text">Log Level</span>
          </label>
          <select v-model="selectedLevel" class="select select-bordered select-sm">
            <option value="">All levels</option>
            <option value="DEBUG">Debug</option>
            <option value="INFO">Info</option>
            <option value="WARN">Warning</option>
            <option value="ERROR">Error</option>
            <option value="FATAL">Fatal</option>
          </select>
        </div>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">Component</span>
          </label>
          <select v-model="selectedComponent" class="select select-bordered select-sm">
            <option value="">All components</option>
            <option value="backup">Backup</option>
            <option value="restore">Restore</option>
            <option value="scheduler">Scheduler</option>
            <option value="api">API</option>
            <option value="database">Database</option>
            <option value="auth">Authentication</option>
          </select>
        </div>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">Time Range</span>
          </label>
          <select v-model="timeRange" class="select select-bordered select-sm">
            <option value="1h">Last hour</option>
            <option value="24h">Last 24 hours</option>
            <option value="7d">Last 7 days</option>
            <option value="30d">Last 30 days</option>
            <option value="all">All time</option>
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
              placeholder="Search logs..." 
              class="input input-bordered input-sm w-full pr-8"
            />
            <svg xmlns="http://www.w3.org/2000/svg" class="absolute right-2 top-2 h-4 w-4 text-base-content/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </div>
        
        <div class="form-control">
          <label class="label">
            <span class="label-text">Actions</span>
          </label>
          <div class="flex gap-2">
            <button @click="refreshLogs" class="btn btn-primary btn-sm">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Refresh
            </button>
            <button @click="clearLogs" class="btn btn-ghost btn-sm">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              Clear
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Auto-refresh Toggle -->
    <div class="bg-base-100 p-4 rounded-lg shadow-sm border border-base-300">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <input 
            id="auto-refresh" 
            type="checkbox" 
            v-model="autoRefresh" 
            class="toggle toggle-primary"
          />
          <label for="auto-refresh" class="text-sm font-medium">Auto-refresh logs</label>
          <span v-if="autoRefresh" class="text-xs text-base-content/70">(every {{ refreshInterval }}s)</span>
        </div>
        
        <div class="text-sm text-base-content/70">
          Showing {{ filteredLogs.length }} of {{ logs.length }} entries
        </div>
      </div>
    </div>

    <!-- Log Entries -->
    <div class="bg-base-100 rounded-lg shadow-sm border border-base-300">
      <div class="p-4 border-b border-base-300">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-semibold">Log Entries</h3>
          <div class="flex items-center gap-2">
            <span class="text-xs text-base-content/70">Auto-scroll</span>
            <input 
              type="checkbox" 
              v-model="autoScroll" 
              class="toggle toggle-xs toggle-primary"
            />
          </div>
        </div>
      </div>
      
      <div 
        ref="logContainer"
        class="h-96 overflow-y-auto font-mono text-sm"
        :class="{ 'scroll-smooth': autoScroll }"
      >
        <div v-for="(log, index) in filteredLogs" :key="index" class="border-b border-base-200 last:border-b-0">
          <div 
            class="flex items-start gap-3 p-3 hover:bg-base-50 transition-colors"
            :class="{
              'bg-error/10': log.level === 'ERROR' || log.level === 'FATAL',
              'bg-warning/10': log.level === 'WARN',
              'bg-info/10': log.level === 'DEBUG'
            }"
          >
            <!-- Timestamp -->
            <div class="flex-shrink-0 w-20 text-xs text-base-content/60">
              {{ formatTime(log.timestamp) }}
            </div>
            
            <!-- Level Badge -->
            <div class="flex-shrink-0">
              <span 
                class="badge badge-xs font-mono"
                :class="{
                  'badge-error': log.level === 'ERROR' || log.level === 'FATAL',
                  'badge-warning': log.level === 'WARN',
                  'badge-info': log.level === 'INFO',
                  'badge-secondary': log.level === 'DEBUG'
                }"
              >
                {{ log.level }}
              </span>
            </div>
            
            <!-- Component -->
            <div class="flex-shrink-0 w-16 text-xs text-base-content/60">
              {{ log.component }}
            </div>
            
            <!-- Message -->
            <div class="flex-1 text-sm leading-relaxed">
              <span v-html="highlightSearch(log.message)"></span>
              <div v-if="log.details" class="mt-1 text-xs text-base-content/70 bg-base-200 p-2 rounded">
                {{ log.details }}
              </div>
            </div>
            
            <!-- Actions -->
            <div class="flex-shrink-0">
              <button 
                @click="copyLogEntry(log)"
                class="btn btn-ghost btn-xs"
                title="Copy log entry"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
              </button>
            </div>
          </div>
        </div>
        
        <!-- Empty state -->
        <div v-if="filteredLogs.length === 0" class="p-8 text-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 mx-auto text-base-content/20 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <h3 class="text-lg font-medium text-base-content mb-2">No logs found</h3>
          <p class="text-base-content/70">No log entries match your current filter criteria</p>
        </div>
        
        <!-- Loading indicator -->
        <div v-if="loading" class="p-4 text-center">
          <span class="loading loading-spinner loading-md"></span>
          <div class="text-sm text-base-content/70 mt-2">Loading logs...</div>
        </div>
      </div>
    </div>

    <!-- Download Logs -->
    <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4">Export Logs</h3>
      <div class="flex items-center gap-4">
        <button @click="downloadLogs('txt')" class="btn btn-outline">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Download as TXT
        </button>
        <button @click="downloadLogs('json')" class="btn btn-outline">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Download as JSON
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'

interface LogEntry {
  timestamp: string
  level: 'DEBUG' | 'INFO' | 'WARN' | 'ERROR' | 'FATAL'
  component: string
  message: string
  details?: string
}

// Filters
const selectedLevel = ref('')
const selectedComponent = ref('')
const timeRange = ref('24h')
const searchQuery = ref('')

// UI state
const autoRefresh = ref(false)
const autoScroll = ref(true)
const refreshInterval = ref(5)
const loading = ref(false)

// Refs
const logContainer = ref<HTMLElement>()

// Auto-refresh timer
let refreshTimer: NodeJS.Timeout | null = null

// Mock log data - in real app this would come from API
const logs = ref<LogEntry[]>([
  {
    timestamp: '2024-01-15T14:30:25Z',
    level: 'INFO',
    component: 'backup',
    message: 'Starting backup for target: Production Server',
    details: 'Target ID: 1, Database: ecommerce_prod'
  },
  {
    timestamp: '2024-01-15T14:30:26Z',
    level: 'DEBUG',
    component: 'backup',
    message: 'Connecting to database server prod.example.com:3306'
  },
  {
    timestamp: '2024-01-15T14:30:27Z',
    level: 'INFO',
    component: 'backup',
    message: 'Database connection established successfully'
  },
  {
    timestamp: '2024-01-15T14:31:45Z',
    level: 'INFO',
    component: 'backup',
    message: 'Backup completed successfully',
    details: 'File: /data/backups/prod_2024-01-15_14-30-25.sql.gz, Size: 2.4 GB'
  },
  {
    timestamp: '2024-01-15T14:32:00Z',
    level: 'INFO',
    component: 'scheduler',
    message: 'Next scheduled backup for Production Server: 2024-01-16T14:30:00Z'
  },
  {
    timestamp: '2024-01-15T15:00:00Z',
    level: 'WARN',
    component: 'api',
    message: 'Rate limit nearly exceeded for client 192.168.1.100'
  },
  {
    timestamp: '2024-01-15T15:15:30Z',
    level: 'ERROR',
    component: 'backup',
    message: 'Failed to connect to Analytics DB server',
    details: 'Error: dial tcp 192.168.1.50:3306: connection refused'
  },
  {
    timestamp: '2024-01-15T15:15:31Z',
    level: 'ERROR',
    component: 'backup',
    message: 'Backup job failed for target: Analytics DB',
    details: 'Reason: Database connection failed'
  },
  {
    timestamp: '2024-01-15T16:00:00Z',
    level: 'INFO',
    component: 'auth',
    message: 'User authenticated successfully via Basic Auth'
  },
  {
    timestamp: '2024-01-15T16:05:12Z',
    level: 'DEBUG',
    component: 'database',
    message: 'SQLite database maintenance completed'
  }
])

const filteredLogs = computed(() => {
  let filtered = [...logs.value]
  
  // Apply filters
  if (selectedLevel.value) {
    filtered = filtered.filter(log => log.level === selectedLevel.value)
  }
  
  if (selectedComponent.value) {
    filtered = filtered.filter(log => log.component === selectedComponent.value)
  }
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(log => 
      log.message.toLowerCase().includes(query) ||
      log.component.toLowerCase().includes(query) ||
      (log.details && log.details.toLowerCase().includes(query))
    )
  }
  
  // Apply time range filter
  if (timeRange.value !== 'all') {
    const now = new Date()
    const cutoff = new Date()
    
    switch (timeRange.value) {
      case '1h':
        cutoff.setHours(now.getHours() - 1)
        break
      case '24h':
        cutoff.setDate(now.getDate() - 1)
        break
      case '7d':
        cutoff.setDate(now.getDate() - 7)
        break
      case '30d':
        cutoff.setDate(now.getDate() - 30)
        break
    }
    
    filtered = filtered.filter(log => new Date(log.timestamp) >= cutoff)
  }
  
  // Sort by timestamp (newest first)
  filtered.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
  
  return filtered
})

const formatTime = (timestamp: string): string => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString()
}

const highlightSearch = (text: string): string => {
  if (!searchQuery.value) return text
  
  const regex = new RegExp(`(${searchQuery.value})`, 'gi')
  return text.replace(regex, '<mark class="bg-yellow-200 dark:bg-yellow-800">$1</mark>')
}

const refreshLogs = async () => {
  loading.value = true
  try {
    // In real app: fetch logs from API
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Add some new mock entries
    const newLogs = [
      {
        timestamp: new Date().toISOString(),
        level: 'INFO' as const,
        component: 'api',
        message: 'Log refresh requested by user'
      }
    ]
    
    logs.value = [...newLogs, ...logs.value]
  } finally {
    loading.value = false
  }
}

const clearLogs = () => {
  if (confirm('Are you sure you want to clear all log entries? This action cannot be undone.')) {
    logs.value = []
  }
}

const copyLogEntry = async (log: LogEntry) => {
  const logText = `[${log.timestamp}] ${log.level} ${log.component}: ${log.message}${log.details ? '\n' + log.details : ''}`
  
  try {
    await navigator.clipboard.writeText(logText)
    // Show success feedback (would need toast component)
    console.log('Log entry copied to clipboard')
  } catch (error) {
    console.error('Failed to copy log entry:', error)
  }
}

const downloadLogs = (format: 'txt' | 'json') => {
  const logsToDownload = filteredLogs.value
  let content: string
  let filename: string
  let mimeType: string
  
  if (format === 'json') {
    content = JSON.stringify(logsToDownload, null, 2)
    filename = `godumper-logs-${new Date().toISOString().split('T')[0]}.json`
    mimeType = 'application/json'
  } else {
    content = logsToDownload.map(log => 
      `[${log.timestamp}] ${log.level} ${log.component}: ${log.message}${log.details ? '\n' + log.details : ''}`
    ).join('\n')
    filename = `godumper-logs-${new Date().toISOString().split('T')[0]}.txt`
    mimeType = 'text/plain'
  }
  
  const blob = new Blob([content], { type: mimeType })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const scrollToBottom = () => {
  if (autoScroll.value && logContainer.value) {
    nextTick(() => {
      logContainer.value!.scrollTop = logContainer.value!.scrollHeight
    })
  }
}

// Watch for new logs to auto-scroll
watch(() => logs.value.length, scrollToBottom)

// Auto-refresh functionality
watch(autoRefresh, (enabled) => {
  if (enabled) {
    refreshTimer = setInterval(() => {
      refreshLogs()
    }, refreshInterval.value * 1000)
  } else if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
})

onMounted(() => {
  // Initial scroll to bottom
  scrollToBottom()
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>