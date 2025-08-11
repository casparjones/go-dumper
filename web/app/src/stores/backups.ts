import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Backup } from '@/types'
import { targetsApi, backupsApi } from '@/services/api'
import { useToastStore } from './toasts'

export const useBackupsStore = defineStore('backups', () => {
  const backups = ref<Backup[]>([])
  const loading = ref(false)
  const toastStore = useToastStore()

  const fetchBackups = async (targetId: number) => {
    loading.value = true
    try {
      backups.value = await targetsApi.getBackups(targetId)
    } catch (error) {
      console.error('Failed to fetch backups:', error)
    } finally {
      loading.value = false
    }
  }

  const fetchAllBackups = async () => {
    loading.value = true
    try {
      backups.value = await backupsApi.getAllBackups()
    } catch (error) {
      console.error('Failed to fetch all backups:', error)
      toastStore.addToast('error', 'Error', 'Failed to load backups')
    } finally {
      loading.value = false
    }
  }

  const downloadBackup = async (id: number): Promise<boolean> => {
    try {
      await backupsApi.download(id)
      toastStore.addToast('success', 'Download Started', 'Backup download has started')
      return true
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || 'Failed to download backup'
      toastStore.addToast('error', 'Error', errorMessage)
      return false
    }
  }

  const restoreBackup = async (id: number): Promise<boolean> => {
    try {
      const result = await backupsApi.restore(id)
      toastStore.addToast('warning', 'Restore Started', result.message)
      return true
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || 'Failed to start restore'
      toastStore.addToast('error', 'Error', errorMessage)
      return false
    }
  }

  const deleteBackup = async (id: number): Promise<boolean> => {
    loading.value = true
    try {
      await backupsApi.delete(id)
      backups.value = backups.value.filter(b => b.id !== id)
      toastStore.addToast('success', 'Success', 'Backup deleted successfully')
      return true
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || 'Failed to delete backup'
      toastStore.addToast('error', 'Error', errorMessage)
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    backups,
    loading,
    fetchBackups,
    fetchAllBackups,
    downloadBackup,
    restoreBackup,
    deleteBackup
  }
})