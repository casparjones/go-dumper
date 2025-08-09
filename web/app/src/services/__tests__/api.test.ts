import { describe, it, expect, beforeEach, vi } from 'vitest'
import axios from 'axios'
import { targetsApi, backupsApi } from '../api'
import type { Target, CreateTargetRequest } from '@/types'

// Mock axios
vi.mock('axios')
const mockedAxios = vi.mocked(axios, true)

// Mock pinia store
vi.mock('@/stores/toasts', () => ({
  useToastStore: () => ({
    addToast: vi.fn()
  })
}))

describe('API Service', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockedAxios.create.mockReturnValue(mockedAxios)
  })

  describe('targetsApi', () => {
    it('should get all targets', async () => {
      const mockTargets: Target[] = [
        {
          id: 1,
          name: 'Test Target',
          host: 'localhost',
          port: 3306,
          db_name: 'testdb',
          user: 'testuser',
          comment: '',
          schedule_time: '',
          retention_days: 30,
          auto_compress: true,
          created_at: '2023-01-01T00:00:00Z',
          updated_at: '2023-01-01T00:00:00Z'
        }
      ]

      mockedAxios.get.mockResolvedValue({ data: mockTargets })

      const result = await targetsApi.getAll()

      expect(mockedAxios.get).toHaveBeenCalledWith('/targets')
      expect(result).toEqual(mockTargets)
    })

    it('should get target by id', async () => {
      const mockTarget: Target = {
        id: 1,
        name: 'Test Target',
        host: 'localhost',
        port: 3306,
        db_name: 'testdb',
        user: 'testuser',
        comment: '',
        schedule_time: '',
        retention_days: 30,
        auto_compress: true,
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-01-01T00:00:00Z'
      }

      mockedAxios.get.mockResolvedValue({ data: mockTarget })

      const result = await targetsApi.getById(1)

      expect(mockedAxios.get).toHaveBeenCalledWith('/targets/1')
      expect(result).toEqual(mockTarget)
    })

    it('should create target', async () => {
      const createRequest: CreateTargetRequest = {
        name: 'New Target',
        host: 'localhost',
        port: 3306,
        db_name: 'newdb',
        user: 'newuser',
        password: 'password123'
      }

      const createdTarget: Target = {
        id: 2,
        ...createRequest,
        comment: '',
        schedule_time: '',
        retention_days: 30,
        auto_compress: false,
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-01-01T00:00:00Z'
      }

      mockedAxios.post.mockResolvedValue({ data: createdTarget })

      const result = await targetsApi.create(createRequest)

      expect(mockedAxios.post).toHaveBeenCalledWith('/targets', createRequest)
      expect(result).toEqual(createdTarget)
    })

    it('should update target', async () => {
      const updateRequest = {
        name: 'Updated Target',
        host: 'localhost',
        port: 3306,
        db_name: 'updateddb',
        user: 'updateduser'
      }

      const updatedTarget: Target = {
        id: 1,
        ...updateRequest,
        comment: '',
        schedule_time: '',
        retention_days: 30,
        auto_compress: false,
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-01-01T01:00:00Z'
      }

      mockedAxios.put.mockResolvedValue({ data: updatedTarget })

      const result = await targetsApi.update(1, updateRequest)

      expect(mockedAxios.put).toHaveBeenCalledWith('/targets/1', updateRequest)
      expect(result).toEqual(updatedTarget)
    })

    it('should delete target', async () => {
      mockedAxios.delete.mockResolvedValue({})

      await targetsApi.delete(1)

      expect(mockedAxios.delete).toHaveBeenCalledWith('/targets/1')
    })

    it('should create backup', async () => {
      const mockResponse = {
        message: 'Backup started',
        backup_id: 123,
        status: 'running'
      }

      mockedAxios.post.mockResolvedValue({ data: mockResponse })

      const result = await targetsApi.createBackup(1)

      expect(mockedAxios.post).toHaveBeenCalledWith('/targets/1/backup')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('backupsApi', () => {
    it('should download backup', async () => {
      const mockBlob = new Blob(['backup data'], { type: 'application/gzip' })
      
      mockedAxios.get.mockResolvedValue({
        data: mockBlob,
        headers: {
          'content-disposition': 'attachment; filename="backup.sql.gz"'
        }
      })

      // Mock URL.createObjectURL and related functions
      const createObjectURL = vi.fn().mockReturnValue('blob:url')
      const revokeObjectURL = vi.fn()
      global.URL.createObjectURL = createObjectURL
      global.URL.revokeObjectURL = revokeObjectURL

      // Mock document.createElement and related functions
      const mockLink = {
        href: '',
        setAttribute: vi.fn(),
        click: vi.fn(),
        remove: vi.fn()
      }
      const createElement = vi.fn().mockReturnValue(mockLink)
      const appendChild = vi.fn()
      
      Object.defineProperty(document, 'createElement', {
        value: createElement,
        configurable: true
      })
      Object.defineProperty(document.body, 'appendChild', {
        value: appendChild,
        configurable: true
      })

      await backupsApi.download(1)

      expect(mockedAxios.get).toHaveBeenCalledWith('/backups/1/download', {
        responseType: 'blob'
      })
      expect(createElement).toHaveBeenCalledWith('a')
      expect(mockLink.setAttribute).toHaveBeenCalledWith('download', 'backup.sql.gz')
      expect(mockLink.click).toHaveBeenCalled()
    })

    it('should restore backup', async () => {
      const mockResponse = {
        message: 'Restore started',
        backup_id: 1
      }

      mockedAxios.post.mockResolvedValue({ data: mockResponse })

      const result = await backupsApi.restore(1)

      expect(mockedAxios.post).toHaveBeenCalledWith('/backups/1/restore')
      expect(result).toEqual(mockResponse)
    })

    it('should delete backup', async () => {
      mockedAxios.delete.mockResolvedValue({})

      await backupsApi.delete(1)

      expect(mockedAxios.delete).toHaveBeenCalledWith('/backups/1')
    })
  })
})