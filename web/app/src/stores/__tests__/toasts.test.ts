import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useToastStore } from '../toasts'

describe('Toast Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
  })

  it('should add a toast', () => {
    const store = useToastStore()
    
    store.addToast('success', 'Test Title', 'Test message')
    
    expect(store.toasts).toHaveLength(1)
    expect(store.toasts[0]).toMatchObject({
      type: 'success',
      title: 'Test Title',
      message: 'Test message'
    })
    expect(store.toasts[0].id).toBeDefined()
  })

  it('should add multiple toasts', () => {
    const store = useToastStore()
    
    store.addToast('success', 'Success', 'Success message')
    store.addToast('error', 'Error', 'Error message')
    
    expect(store.toasts).toHaveLength(2)
    expect(store.toasts[0].type).toBe('success')
    expect(store.toasts[1].type).toBe('error')
  })

  it('should remove a toast by id', () => {
    const store = useToastStore()
    
    store.addToast('info', 'Info', 'Info message')
    const toastId = store.toasts[0].id
    
    store.removeToast(toastId)
    
    expect(store.toasts).toHaveLength(0)
  })

  it('should not remove toast with invalid id', () => {
    const store = useToastStore()
    
    store.addToast('warning', 'Warning', 'Warning message')
    store.removeToast('invalid-id')
    
    expect(store.toasts).toHaveLength(1)
  })

  it('should auto-remove toast after duration', () => {
    const store = useToastStore()
    
    store.addToast('success', 'Auto Remove', 'This will disappear', 1000)
    
    expect(store.toasts).toHaveLength(1)
    
    vi.advanceTimersByTime(1000)
    
    expect(store.toasts).toHaveLength(0)
  })

  it('should not auto-remove toast with 0 duration', () => {
    const store = useToastStore()
    
    store.addToast('error', 'Persistent', 'This stays', 0)
    
    expect(store.toasts).toHaveLength(1)
    
    vi.advanceTimersByTime(10000)
    
    expect(store.toasts).toHaveLength(1)
  })

  it('should generate unique ids', () => {
    const store = useToastStore()
    
    store.addToast('info', 'First', 'First message')
    store.addToast('info', 'Second', 'Second message')
    
    expect(store.toasts[0].id).not.toBe(store.toasts[1].id)
  })
})