import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { useTargetsStore } from './targets'
import { configApi } from '@/services/api'

export const useAppStore = defineStore('app', () => {
  const selectedConfigId = ref<number | null>(null)
  const selectedDatabase = ref<string>('')
  const currentTheme = ref('dim')

  // Load saved configuration from localStorage
  const loadSavedConfig = () => {
    const saved = localStorage.getItem('go-dumper-selected-config')
    if (saved) {
      try {
        const configId = parseInt(saved, 10)
        if (!isNaN(configId)) {
          selectedConfigId.value = configId
        }
      } catch (error) {
        console.warn('Failed to load saved configuration:', error)
        localStorage.removeItem('go-dumper-selected-config')
      }
    }
  }

  // Validate saved configuration exists
  const validateSavedConfig = () => {
    if (selectedConfigId.value) {
      const targetStore = useTargetsStore()
      const config = targetStore.getTargetById(selectedConfigId.value)
      if (!config) {
        // Configuration no longer exists, clear selection
        selectedConfigId.value = null
        localStorage.removeItem('go-dumper-selected-config')
      }
    }
  }

  const setSelectedConfig = (configId: number | null) => {
    selectedConfigId.value = configId
    selectedDatabase.value = '' // Reset database selection when config changes
    
    // Save to localStorage
    if (configId !== null) {
      localStorage.setItem('go-dumper-selected-config', configId.toString())
    } else {
      localStorage.removeItem('go-dumper-selected-config')
    }
  }

  const setSelectedDatabase = (database: string) => {
    selectedDatabase.value = database
  }

  // Detect system theme and set default
  const detectSystemTheme = () => {
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      currentTheme.value = 'dim'
    } else {
      currentTheme.value = 'winter'
    }
    document.documentElement.setAttribute('data-theme', currentTheme.value)
  }

  const setTheme = (theme: string) => {
    currentTheme.value = theme
    document.documentElement.setAttribute('data-theme', theme)
  }

  // Load theme from backend
  const loadThemeFromBackend = async () => {
    try {
      const response = await configApi.getTheme()
      currentTheme.value = response.theme
      document.documentElement.setAttribute('data-theme', response.theme)
    } catch (error) {
      console.warn('Failed to load theme from backend, using system default:', error)
      detectSystemTheme()
    }
  }

  // Save theme to backend
  const saveThemeToBackend = async (theme: string) => {
    try {
      await configApi.setTheme(theme)
    } catch (error) {
      console.error('Failed to save theme to backend:', error)
      throw error
    }
  }

  const selectedConfig = computed(() => {
    if (!selectedConfigId.value) return null
    const targetStore = useTargetsStore()
    return targetStore.getTargetById(selectedConfigId.value) || null
  })

  // Watch for system theme changes
  if (window.matchMedia) {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', detectSystemTheme)
  }

  return {
    selectedConfigId,
    selectedDatabase,
    currentTheme,
    selectedConfig,
    setSelectedConfig,
    setSelectedDatabase,
    setTheme,
    loadSavedConfig,
    validateSavedConfig,
    detectSystemTheme,
    loadThemeFromBackend,
    saveThemeToBackend
  }
})