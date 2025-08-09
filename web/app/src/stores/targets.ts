import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Target, CreateTargetRequest, UpdateTargetRequest } from '@/types'
import { targetsApi } from '@/services/api'
import { useToastStore } from './toasts'

export const useTargetsStore = defineStore('targets', () => {
  const targets = ref<Target[]>([])
  const loading = ref(false)
  const toastStore = useToastStore()

  const fetchTargets = async () => {
    loading.value = true
    try {
      targets.value = await targetsApi.getAll()
    } catch (error) {
      console.error('Failed to fetch targets:', error)
    } finally {
      loading.value = false
    }
  }

  const createTarget = async (target: CreateTargetRequest): Promise<boolean> => {
    loading.value = true
    try {
      const newTarget = await targetsApi.create(target)
      targets.value.push(newTarget)
      toastStore.addToast('success', 'Success', `Target "${target.name}" created successfully`)
      return true
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || 'Failed to create target'
      toastStore.addToast('error', 'Error', errorMessage)
      return false
    } finally {
      loading.value = false
    }
  }

  const updateTarget = async (id: number, target: UpdateTargetRequest): Promise<boolean> => {
    loading.value = true
    try {
      const updatedTarget = await targetsApi.update(id, target)
      const index = targets.value.findIndex(t => t.id === id)
      if (index !== -1) {
        targets.value[index] = updatedTarget
      }
      toastStore.addToast('success', 'Success', `Target "${target.name}" updated successfully`)
      return true
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || 'Failed to update target'
      toastStore.addToast('error', 'Error', errorMessage)
      return false
    } finally {
      loading.value = false
    }
  }

  const deleteTarget = async (id: number): Promise<boolean> => {
    loading.value = true
    try {
      await targetsApi.delete(id)
      targets.value = targets.value.filter(t => t.id !== id)
      toastStore.addToast('success', 'Success', 'Target deleted successfully')
      return true
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || 'Failed to delete target'
      toastStore.addToast('error', 'Error', errorMessage)
      return false
    } finally {
      loading.value = false
    }
  }

  const createBackup = async (id: number): Promise<boolean> => {
    try {
      const result = await targetsApi.createBackup(id)
      toastStore.addToast('success', 'Backup Started', result.message)
      return true
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || 'Failed to start backup'
      toastStore.addToast('error', 'Error', errorMessage)
      return false
    }
  }

  const getTargetById = (id: number): Target | undefined => {
    return targets.value.find(t => t.id === id)
  }

  return {
    targets,
    loading,
    fetchTargets,
    createTarget,
    updateTarget,
    deleteTarget,
    createBackup,
    getTargetById
  }
})