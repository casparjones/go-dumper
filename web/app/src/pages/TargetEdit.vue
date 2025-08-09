<template>
  <div>
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-base-content">
        {{ isEditing ? 'Edit Target' : 'Add Target' }}
      </h1>
      <p class="text-base-content/70">
        {{ isEditing ? 'Update your backup target configuration' : 'Configure a new MySQL/MariaDB backup target' }}
      </p>
    </div>

    <form @submit.prevent="handleSubmit" class="max-w-2xl">
      <div class="space-y-6">
        <!-- Basic Information -->
        <div class="card bg-base-200 shadow">
          <div class="card-body">
            <h2 class="card-title mb-4">Basic Information</h2>
            
            <div class="form-control">
              <label class="label">
                <span class="label-text">Name *</span>
              </label>
              <input
                v-model="form.name"
                type="text"
                placeholder="My Database"
                class="input input-bordered"
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
                class="textarea textarea-bordered"
                placeholder="Optional description"
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Connection Settings -->
        <div class="card bg-base-200 shadow">
          <div class="card-body">
            <h2 class="card-title mb-4">Connection Settings</h2>
            
            <div class="grid gap-4 md:grid-cols-2">
              <div class="form-control">
                <label class="label">
                  <span class="label-text">Host *</span>
                </label>
                <input
                  v-model="form.host"
                  type="text"
                  placeholder="localhost"
                  class="input input-bordered"
                  :class="{ 'input-error': errors.host }"
                  required
                />
                <label v-if="errors.host" class="label">
                  <span class="label-text-alt text-error">{{ errors.host }}</span>
                </label>
              </div>

              <div class="form-control">
                <label class="label">
                  <span class="label-text">Port *</span>
                </label>
                <input
                  v-model.number="form.port"
                  type="number"
                  placeholder="3306"
                  class="input input-bordered"
                  :class="{ 'input-error': errors.port }"
                  min="1"
                  max="65535"
                  required
                />
                <label v-if="errors.port" class="label">
                  <span class="label-text-alt text-error">{{ errors.port }}</span>
                </label>
              </div>
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text">Database Name *</span>
              </label>
              <input
                v-model="form.db_name"
                type="text"
                placeholder="mydb"
                class="input input-bordered"
                :class="{ 'input-error': errors.db_name }"
                required
              />
              <label v-if="errors.db_name" class="label">
                <span class="label-text-alt text-error">{{ errors.db_name }}</span>
              </label>
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text">Username *</span>
              </label>
              <input
                v-model="form.user"
                type="text"
                placeholder="root"
                class="input input-bordered"
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
                placeholder="password"
                class="input input-bordered"
                :class="{ 'input-error': errors.password }"
                :required="!isEditing"
              />
              <label v-if="errors.password" class="label">
                <span class="label-text-alt text-error">{{ errors.password }}</span>
              </label>
            </div>
          </div>
        </div>

        <!-- Backup Settings -->
        <div class="card bg-base-200 shadow">
          <div class="card-body">
            <h2 class="card-title mb-4">Backup Settings</h2>
            
            <div class="form-control">
              <label class="label">
                <span class="label-text">Schedule (UTC)</span>
              </label>
              <input
                v-model="form.schedule_time"
                type="time"
                class="input input-bordered"
                placeholder="HH:MM"
              />
              <label class="label">
                <span class="label-text-alt">Leave empty to disable automatic backups</span>
              </label>
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text">Retention Days</span>
              </label>
              <input
                v-model.number="form.retention_days"
                type="number"
                placeholder="30"
                class="input input-bordered"
                min="1"
                max="365"
              />
              <label class="label">
                <span class="label-text-alt">How many days to keep backups (default: 30)</span>
              </label>
            </div>

            <div class="form-control">
              <label class="cursor-pointer label">
                <span class="label-text">Enable Compression</span>
                <input
                  v-model="form.auto_compress"
                  type="checkbox"
                  class="checkbox checkbox-primary"
                />
              </label>
              <label class="label">
                <span class="label-text-alt">Compress backup files with gzip</span>
              </label>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex gap-4">
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
  db_name: '',
  user: '',
  password: '',
  comment: '',
  schedule_time: '',
  retention_days: 30,
  auto_compress: true
})

const errors = ref<Record<string, string>>({})

const validateForm = (): boolean => {
  errors.value = {}

  if (!form.value.name.trim()) {
    errors.value.name = 'Name is required'
  }

  if (!form.value.host.trim()) {
    errors.value.host = 'Host is required'
  }

  if (!form.value.port || form.value.port < 1 || form.value.port > 65535) {
    errors.value.port = 'Port must be between 1 and 65535'
  }

  if (!form.value.db_name.trim()) {
    errors.value.db_name = 'Database name is required'
  }

  if (!form.value.user.trim()) {
    errors.value.user = 'Username is required'
  }

  if (!isEditing.value && !form.value.password) {
    errors.value.password = 'Password is required'
  }

  if (form.value.schedule_time && !/^([01]?[0-9]|2[0-3]):[0-5][0-9]$/.test(form.value.schedule_time)) {
    errors.value.schedule_time = 'Schedule time must be in HH:MM format'
  }

  return Object.keys(errors.value).length === 0
}

const handleSubmit = async () => {
  if (!validateForm()) {
    return
  }

  let success = false
  
  if (isEditing.value) {
    const id = Number(route.params.id)
    success = await targetStore.updateTarget(id, form.value as UpdateTargetRequest)
  } else {
    success = await targetStore.createTarget(form.value as CreateTargetRequest)
  }

  if (success) {
    router.push('/targets')
  }
}

const testConnection = async () => {
  if (!validateForm()) {
    return
  }

  testing.value = true
  try {
    // In a real implementation, you'd call an API endpoint to test the connection
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
        db_name: target.db_name,
        user: target.user,
        password: '',
        comment: target.comment,
        schedule_time: target.schedule_time,
        retention_days: target.retention_days,
        auto_compress: target.auto_compress
      }
    } else {
      toastStore.addToast('error', 'Error', 'Target not found')
      router.push('/targets')
    }
  }
})
</script>