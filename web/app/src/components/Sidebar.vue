<template>
  <div class="drawer-side">
    <label for="drawer-toggle" class="drawer-overlay"></label>
    <aside class="w-80 min-h-full bg-base-200 flex flex-col">
      <!-- Logo Section -->
      <div class="p-6 border-b border-base-300">
        <!-- Text Logo -->
        <div class="text-center mb-4">
          <h1 class="text-2xl font-bold text-primary">GoDumper</h1>
          <p class="text-xs text-base-content/70">MySQL/MariaDB Backup Tool</p>
        </div>
        
        <!-- Logo Placeholder -->
        <div class="flex justify-center mb-2">
          <div class="w-16 h-16 bg-primary/20 rounded-lg flex items-center justify-center border-2 border-primary/30">
            <svg xmlns="http://www.w3.org/2000/svg"
                 class="w-8 h-8 text-primary"
                 viewBox="0 0 24 24"
                 fill="none"
                 stroke="currentColor"
                 stroke-width="2"
                 stroke-linecap="round"
                 stroke-linejoin="round">
              <ellipse cx="12" cy="5" rx="8" ry="3"></ellipse>
              <path d="M4 5v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5"></path>
              <path d="M4 12c0 1.66 3.58 3 8 3s8-1.34 8-3"></path>
              <path d="M4 18c0 1.66 3.58 3 8 3s8-1.34 8-3"></path>
            </svg>
          </div>
        </div>
      </div>

      <!-- Navigation Menu -->
      <nav class="flex-1 p-4">
        <ul class="menu space-y-1">
          <li>
            <router-link to="/configure" class="flex items-center gap-3" :class="{ 'active': $route.path === '/configure' }">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>

              Database Config
            </router-link>
          </li>
          
          <li>
            <router-link to="/backup" class="flex items-center gap-3" :class="{ 'active': $route.path === '/backup' }">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
              </svg>
              Backup
            </router-link>
          </li>

          <li>
            <router-link to="/schedule-jobs" class="flex items-center gap-3" :class="{ 'active': $route.path === '/schedule-jobs' }">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Schedule Jobs
            </router-link>
          </li>
          
          <li>
            <router-link to="/restore" class="flex items-center gap-3" :class="{ 'active': $route.path === '/restore' }">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
              </svg>
              Restore
            </router-link>
          </li>
          
          <li>
            <router-link to="/manage" class="flex items-center gap-3" :class="{ 'active': $route.path === '/manage' }">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
              Manage
            </router-link>
          </li>
          
          <li>
            <router-link to="/log" class="flex items-center gap-3" :class="{ 'active': $route.path === '/log' }">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              Log
            </router-link>
          </li>

          <li>
            <router-link to="/settings" class="flex items-center gap-3" :class="{ 'active': $route.path === '/settings' }">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>

              Settings
            </router-link>
          </li>

          <li>
            <router-link to="/help" class="flex items-center gap-3" :class="{ 'active': $route.path === '/help' }">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Credits/Help
            </router-link>
          </li>
        </ul>
      </nav>

      <!-- Configuration Selectors -->
      <div class="p-4 border-t border-base-300 space-y-4">
        <div>
          <label class="block text-sm font-medium text-base-content mb-2">
            Configuration:
          </label>
          <select 
            v-model="selectedConfig" 
            @change="onConfigChange"
            class="select select-bordered select-sm w-full"
          >
            <option value="">Select Configuration...</option>
            <option 
              v-for="config in configurations" 
              :key="config.id" 
              :value="config.id"
            >
              {{ config.name }}
            </option>
          </select>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-base-content mb-2">
            Select a Database:
          </label>
          <select 
            v-model="selectedDatabase" 
            @change="onDatabaseChange"
            class="select select-bordered select-sm w-full"
            :disabled="!selectedConfig"
          >
            <option value="">Select Database...</option>
            <option 
              v-for="database in databases" 
              :key="database.name" 
              :value="database.name"
            >
              {{ database.name }} ({{ database.size }})
            </option>
          </select>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTargetsStore } from '@/stores/targets'
import { useAppStore } from '@/stores/app'

interface Configuration {
  id: number
  name: string
  host: string
  port: number
  user: string
}

interface Database {
  name: string
  size: string
  tables: number
}

const targetStore = useTargetsStore()
const appStore = useAppStore()

const configurations = computed<Configuration[]>(() => {
  return targetStore.targets.map(target => ({
    id: target.id,
    name: target.name,
    host: target.host,
    port: target.port,
    user: target.user
  }))
})

const databases = ref<Database[]>([
  { name: 'production_db', size: '2.4 GB', tables: 47 },
  { name: 'staging_db', size: '856 MB', tables: 42 },
  { name: 'test_db', size: '124 MB', tables: 35 },
  { name: 'analytics_db', size: '5.7 GB', tables: 89 }
])

const selectedConfig = computed({
  get: () => appStore.selectedConfigId || '',
  set: (value: number | '') => {
    appStore.setSelectedConfig(value === '' ? null : value)
  }
})

const selectedDatabase = computed({
  get: () => appStore.selectedDatabase,
  set: (value: string) => {
    appStore.setSelectedDatabase(value)
  }
})

const onConfigChange = () => {
  // Database selection is automatically reset in the store
}

const onDatabaseChange = () => {
  // Handle database selection
  console.log('Selected database:', selectedDatabase.value)
}

onMounted(async () => {
  // Load saved configuration from localStorage
  appStore.loadSavedConfig()
  
  // Load targets and then validate saved configuration
  await targetStore.fetchTargets()
  appStore.validateSavedConfig()
  
  // Load theme from backend, fallback to system detection
  await appStore.loadThemeFromBackend()
})
</script>
