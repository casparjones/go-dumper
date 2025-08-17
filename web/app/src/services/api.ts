import axios from 'axios'
import type { AxiosResponse, AxiosError, AxiosInstance } from 'axios'
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
  },

  async discoverDatabases(host: string, port: number, user: string, password: string): Promise<{databases: {name: string}[]}> {
    const response = await api.post('/targets/discover', {
      host,
      port,
      user,
      password
    })
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
  },

  async getAllBackups(): Promise<any[]> {
    const response = await api.get('/backups')
    return response.data
  }
}

export interface ScheduleJob {
  id: number
  target_id: number
  name: string
  description: string
  is_active: boolean
  schedule_config: string
  backup_options: string
  meta_config: string
  last_run_at?: string
  last_run_status?: string
  last_run_notes?: string
  next_run_at?: string
  created_at: string
  updated_at: string
  target?: Target
}

export interface CreateJobRequest {
  target_id: number
  name: string
  description?: string
  schedule_config: {
    frequency: string
    minutes?: number[]
    hours?: number[]
    weekdays?: number[]
    days_of_month?: number[]
    months?: number[]
  }
  backup_options: {
    compress: boolean
    include_structure: boolean
    include_data: boolean
    databases?: string[]
  }
  meta_config?: Record<string, any>
}

export interface UpdateJobRequest {
  name: string
  description: string
  is_active: boolean
  schedule_config: {
    frequency: string
    minutes?: number[]
    hours?: number[]
    weekdays?: number[]
    days_of_month?: number[]
    months?: number[]
  }
  backup_options: {
    compress: boolean
    include_structure: boolean
    include_data: boolean
    databases?: string[]
  }
  meta_config?: Record<string, any>
}

export const jobsApi = {
  async getAll(): Promise<ScheduleJob[]> {
    const response = await api.get<ScheduleJob[]>('/jobs')
    return response.data
  },

  async getById(id: number): Promise<ScheduleJob> {
    const response = await api.get<ScheduleJob>(`/jobs/${id}`)
    return response.data
  },

  async create(job: CreateJobRequest): Promise<ScheduleJob> {
    const response = await api.post<ScheduleJob>('/jobs', job)
    return response.data
  },

  async update(id: number, job: UpdateJobRequest): Promise<ScheduleJob> {
    const response = await api.put<ScheduleJob>(`/jobs/${id}`, job)
    return response.data
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/jobs/${id}`)
  },

  async runNow(id: number): Promise<{ message: string; job_id: number }> {
    const response = await api.post(`/jobs/${id}/run`)
    return response.data
  }
}

export const healthApi = {
  async check(): Promise<{ status: string; service: string }> {
    const response = await api.get('/healthz')
    return response.data
  }
}

export interface AppConfig {
  id: number
  key: string
  value: string
  created_at: string
  updated_at: string
}

export const configApi = {
  async getTheme(): Promise<{ theme: string }> {
    const response = await api.get('/config/theme')
    return response.data
  },

  async setTheme(theme: string): Promise<{ message: string; theme: string }> {
    const response = await api.post('/config/theme', { theme })
    return response.data
  },

  async getConfig(key: string): Promise<AppConfig> {
    const response = await api.get(`/config/${key}`)
    return response.data
  },

  async setConfig(key: string, value: string): Promise<{ message: string }> {
    const response = await api.post('/config', { key, value })
    return response.data
  },

  async getAllConfigs(): Promise<AppConfig[]> {
    const response = await api.get('/config')
    return response.data
  }
}