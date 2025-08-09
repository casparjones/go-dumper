import axios, { AxiosResponse } from 'axios'
import type { Target, CreateTargetRequest, UpdateTargetRequest, Backup } from '@/types'
import { useToastStore } from '@/stores/toasts'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

api.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error) => {
    const toastStore = useToastStore()
    
    if (error.response?.status === 401) {
      toastStore.addToast('error', 'Authentication Error', 'Please check your credentials')
    } else if (error.response?.status === 403) {
      toastStore.addToast('error', 'Access Denied', 'You do not have permission to perform this action')
    } else if (error.response?.status >= 500) {
      toastStore.addToast('error', 'Server Error', 'An internal server error occurred')
    } else if (error.code === 'ECONNABORTED') {
      toastStore.addToast('error', 'Timeout', 'Request timed out')
    } else if (!error.response) {
      toastStore.addToast('error', 'Network Error', 'Unable to connect to the server')
    }

    return Promise.reject(error)
  }
)

export const targetsApi = {
  async getAll(): Promise<Target[]> {
    const response = await api.get<Target[]>('/targets')
    return response.data
  },

  async getById(id: number): Promise<Target> {
    const response = await api.get<Target>(`/targets/${id}`)
    return response.data
  },

  async create(target: CreateTargetRequest): Promise<Target> {
    const response = await api.post<Target>('/targets', target)
    return response.data
  },

  async update(id: number, target: UpdateTargetRequest): Promise<Target> {
    const response = await api.put<Target>(`/targets/${id}`, target)
    return response.data
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/targets/${id}`)
  },

  async createBackup(id: number): Promise<{ message: string; backup_id: number; status: string }> {
    const response = await api.post(`/targets/${id}/backup`)
    return response.data
  },

  async getBackups(id: number): Promise<Backup[]> {
    const response = await api.get<Backup[]>(`/targets/${id}/backups`)
    return response.data
  }
}

export const backupsApi = {
  async download(id: number): Promise<void> {
    const response = await api.get(`/backups/${id}/download`, {
      responseType: 'blob'
    })
    
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    
    const contentDisposition = response.headers['content-disposition']
    let filename = 'backup.sql.gz'
    if (contentDisposition) {
      const matches = /filename=([^;]+)/.exec(contentDisposition)
      if (matches) {
        filename = matches[1].replace(/"/g, '')
      }
    }
    
    link.setAttribute('download', filename)
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)
  },

  async restore(id: number): Promise<{ message: string; backup_id: number }> {
    const response = await api.post(`/backups/${id}/restore`)
    return response.data
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/backups/${id}`)
  }
}

export const healthApi = {
  async check(): Promise<{ status: string; service: string }> {
    const response = await api.get('/healthz')
    return response.data
  }
}