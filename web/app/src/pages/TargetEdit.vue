<template>
  <div class="mx-auto max-w-3xl">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-3xl font-bold">
        {{ isEditing ? 'Edit Target' : 'Add Target' }}
      </h1>
      <p class="text-base-content/70">
        {{ isEditing ? 'Update your backup target configuration' : 'Configure a new MySQL/MariaDB backup target' }}
      </p>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Basic Information -->
      <div class="card bg-base-100 border border-base-200 shadow-sm">
        <div class="card-body gap-4">
          <h2 class="card-title">Basic Information</h2>

          <div class="form-control">
            <label class="label">
              <span class="label-text">Name <span class="text-error">*</span></span>
            </label>
            <input
                v-model="form.name"
                type="text"
                placeholder="My Database"
                class="input input-bordered w-full"
                :class="{ 'input-error': errors.name }"
                required
            />
            <label v-if="errors.name" class="label">
              <span class="label-text-alt text-error">{{ errors.name }}</span>
            </label>
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">Comment</span>
            </label>
            <textarea
                v-model="form.comment"
                class="textarea textarea-bordered w-full min-h-24"
                placeholder="Optional description"
            ></textarea>
          </div>
        </div>
      </div>

      <!-- Connection Settings -->
      <div class="card bg-base-100 border border-base-200 shadow-sm">
        <div class="card-body gap-4">
          <h2 class="card-title">Connection Settings</h2>

          <!-- Host + Port as a joined input -->
          <div class="form-control">
            <label class="label">
              <span class="label-text">Host &amp; Port <span class="text-error">*</span></span>
            </label>
            <div class="join w-full">
              <input
                  v-model="form.host"
                  type="text"
                  placeholder="localhost"
                  class="input input-bordered join-item w-full"
                  :class="{ 'input-error': errors.host }"
                  required
              />
              <input
                  v-model.number="form.port"
                  type="number"
                  placeholder="3306"
                  class="input input-bordered join-item max-w-28 text-center"
                  :class="{ 'input-error': errors.port }"
                  min="1"
                  max="65535"
                  required
              />
            </div>
            <div class="mt-1 flex gap-4">
              <span v-if="errors.host" class="text-xs text-error">{{ errors.host }}</span>
              <span v-if="errors.port" class="text-xs text-error">{{ errors.port }}</span>
            </div>
          </div>

          <div class="grid gap-4 md:grid-cols-2">
            <div class="form-control">
              <label class="label">
                <span class="label-text">Username <span class="text-error">*</span></span>
              </label>
              <input
                  v-model="form.user"
                  type="text"
                  placeholder="root"
                  class="input input-bordered w-full"
                  :class="{ 'input-error': errors.user }"
                  required
              />
              <label v-if="errors.user" class="label">
                <span class="label-text-alt text-error">{{ errors.user }}</span>
              </label>
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text">Password {{ isEditing ? '(leave empty to keep current)' : '*' }}</span>
              </label>
              <input
                  v-model="form.password"
                  type="password"
                  placeholder="••••••••"
                  class="input input-bordered w-full"
                  :class="{ 'input-error': errors.password }"
                  :required="!isEditing"
              />
              <label v-if="errors.password" class="label">
                <span class="label-text-alt text-error">{{ errors.password }}</span>
              </label>
            </div>
          </div>

          <!-- Database Configuration -->
          <div class="form-control">
            <label class="label">
              <span class="label-text">Database Configuration <span class="text-error">*</span></span>
            </label>
            <div class="flex flex-col space-y-3">
              <label class="label cursor-pointer justify-start gap-3">
                <input
                    v-model="form.database_mode"
                    type="radio"
                    value="all"
                    name="database_mode"
                    class="radio radio-primary"
                    @change="onDatabaseModeChange"
                />
                <span class="label-text">Backup all databases (recommended)</span>
              </label>
              <label class="label cursor-pointer justify-start gap-3">
                <input
                    v-model="form.database_mode"
                    type="radio"
                    value="selected"
                    name="database_mode"
                    class="radio radio-primary"
                    @change="onDatabaseModeChange"
                />
                <span class="label-text">Select specific databases</span>
              </label>
            </div>
          </div>

          <!-- Database Discovery and Selection -->
          <div v-if="form.database_mode === 'selected'" class="form-control">
            <label class="label">
              <span class="label-text">Available Databases</span>
            </label>
            <div class="space-y-3">
              <button
                  type="button"
                  @click="discoverDatabases"
                  class="btn btn-secondary btn-sm"
                  :disabled="discovering || !canDiscover"
              >
                <span v-if="discovering" class="loading loading-spinner loading-xs"></span>
                {{ discovering ? 'Discovering...' : 'Discover Databases' }}
              </button>
              <div v-if="!canDiscover" class="text-sm text-base-content/70">
                Please fill in host, port, username and password first
              </div>
            </div>

            <div v-if="availableDatabases.length > 0" class="mt-4 space-y-2">
              <div class="text-sm font-medium">Select databases to backup:</div>
              <div class="grid gap-2 md:grid-cols-2 lg:grid-cols-3">
                <label
                    v-for="database in availableDatabases"
                    :key="database.name"
                    class="label cursor-pointer justify-start gap-2"
                >
                  <input
                      v-model="form.selected_databases"
                      type="checkbox"
                      :value="database.name"
                      class="checkbox checkbox-primary checkbox-sm"
                  />
                  <span class="label-text">{{ database.name }}</span>
                </label>
              </div>
              <div v-if="errors.selected_databases" class="text-sm text-error mt-1">
                {{ errors.selected_databases }}
              </div>
            </div>
          </div>

        </div>
      </div>

      <!-- Actions -->
      <div class="flex flex-wrap gap-3">
        <button
            type="submit"
            class="btn btn-primary"
            :disabled="targetStore.loading"
        >
          <span v-if="targetStore.loading" class="loading loading-spinner loading-xs"></span>
          {{ isEditing ? 'Update Target' : 'Create Target' }}
        </button>

        <router-link to="/targets" class="btn btn-outline">
          Cancel
        </router-link>

        <button
            v-if="isEditing"
            type="button"
            @click="testConnection"
            class="btn btn-secondary"
            :disabled="testing"
        >
          <span v-if="testing" class="loading loading-spinner loading-xs"></span>
          Test Connection
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTargetsStore } from '@/stores/targets'
import { useToastStore } from '@/stores/toasts'
import type { CreateTargetRequest, UpdateTargetRequest } from '@/types'

const route = useRoute()
const router = useRouter()
const targetStore = useTargetsStore()
const toastStore = useToastStore()

const isEditing = computed(() => !!route.params.id)
const testing = ref(false)

const form = ref<CreateTargetRequest | UpdateTargetRequest>({
  name: '',
  host: 'localhost',
  port: 3306,
  user: '',
  password: '',
  comment: '',
  schedule_time: '',
  retention_days: 30,
  auto_compress: true,
  database_mode: 'all',
  selected_databases: []
})

const errors = ref<Record<string, string>>({})
const discovering = ref(false)
const availableDatabases = ref<{name: string}[]>([])

const canDiscover = computed(() => {
  return form.value.host.trim() && 
         form.value.port > 0 && 
         form.value.user.trim() && 
         form.value.password.trim()
})

const validateForm = (): boolean => {
  errors.value = {}

  if (!form.value.name.trim()) errors.value.name = 'Name is required'
  if (!form.value.host.trim()) errors.value.host = 'Host is required'
  if (!form.value.port || form.value.port < 1 || form.value.port > 65535) errors.value.port = 'Port must be between 1 and 65535'
  if (!form.value.user.trim()) errors.value.user = 'Username is required'
  if (!isEditing.value && !form.value.password) errors.value.password = 'Password is required'
  if (form.value.database_mode === 'selected' && form.value.selected_databases?.length === 0) {
    errors.value.selected_databases = 'Please select at least one database'
  }
  if (form.value.schedule_time && !/^([01]?[0-9]|2[0-3]):[0-5][0-9]$/.test(form.value.schedule_time)) {
    errors.value.schedule_time = 'Schedule time must be in HH:MM format'
  }

  return Object.keys(errors.value).length === 0
}

const handleSubmit = async () => {
  if (!validateForm()) return

  let success = false
  if (isEditing.value) {
    const id = Number(route.params.id)
    success = await targetStore.updateTarget(id, form.value as UpdateTargetRequest)
  } else {
    success = await targetStore.createTarget(form.value as CreateTargetRequest)
  }

  if (success) router.push('/targets')
}

const onDatabaseModeChange = () => {
  if (form.value.database_mode === 'all') {
    form.value.selected_databases = []
    availableDatabases.value = []
  }
}

const discoverDatabases = async () => {
  if (!canDiscover.value) return

  discovering.value = true
  try {
    const response = await fetch('/api/targets/discover', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        host: form.value.host,
        port: form.value.port,
        user: form.value.user,
        password: form.value.password
      })
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to discover databases')
    }

    const data = await response.json()
    availableDatabases.value = data.databases || []
    
    toastStore.addToast('success', 'Success', `Found ${availableDatabases.value.length} databases`)
  } catch (error) {
    console.error('Database discovery failed:', error)
    toastStore.addToast('error', 'Discovery Failed', error.message || 'Unable to discover databases')
  } finally {
    discovering.value = false
  }
}

const testConnection = async () => {
  if (!validateForm()) return

  testing.value = true
  try {
    toastStore.addToast('info', 'Test Connection', 'This feature would test the database connection')
  } catch (error) {
    toastStore.addToast('error', 'Connection Failed', 'Unable to connect to the database')
  } finally {
    testing.value = false
  }
}

onMounted(async () => {
  if (isEditing.value) {
    const id = Number(route.params.id)
    await targetStore.fetchTargets()
    const target = targetStore.getTargetById(id)

    if (target) {
      form.value = {
        name: target.name,
        host: target.host,
        port: target.port,
        user: target.user,
        password: '',
        comment: target.comment,
        schedule_time: target.schedule_time,
        retention_days: target.retention_days,
        auto_compress: target.auto_compress,
        database_mode: target.database_mode || 'all',
        selected_databases: target.selected_databases || []
      }
    } else {
      toastStore.addToast('error', 'Error', 'Target not found')
      router.push('/targets')
    }
  }
})
</script>
