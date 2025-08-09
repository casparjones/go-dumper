export interface Target {
  id: number
  name: string
  host: string
  port: number
  db_name: string
  user: string
  comment: string
  schedule_time: string
  retention_days: number
  auto_compress: boolean
  created_at: string
  updated_at: string
}

export interface CreateTargetRequest {
  name: string
  host: string
  port: number
  db_name: string
  user: string
  password: string
  comment?: string
  schedule_time?: string
  retention_days?: number
  auto_compress?: boolean
}

export interface UpdateTargetRequest {
  name: string
  host: string
  port: number
  db_name: string
  user: string
  password?: string
  comment?: string
  schedule_time?: string
  retention_days?: number
  auto_compress?: boolean
}

export interface Backup {
  id: number
  target_id: number
  started_at: string
  finished_at?: string
  size_bytes: number
  status: 'running' | 'success' | 'failed'
  file_path: string
  notes: string
}

export interface ToastMessage {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  title: string
  message: string
}