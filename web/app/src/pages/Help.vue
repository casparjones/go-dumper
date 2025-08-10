<template>
  <div class="space-y-6">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-base-content mb-2">Help & Credits</h1>
      <p class="text-base-content/70">Documentation, support information, and project credits</p>
    </div>

    <!-- Quick Links -->
    <div class="grid gap-4 md:grid-cols-3">
      <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 mx-auto text-primary mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
        </svg>
        <h3 class="font-semibold mb-2">Documentation</h3>
        <p class="text-sm text-base-content/70 mb-4">Complete user guide and API documentation</p>
        <a href="#" class="btn btn-primary btn-sm">View Docs</a>
      </div>
      
      <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 mx-auto text-secondary mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192L5.636 18.364M12 12h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="font-semibold mb-2">Support</h3>
        <p class="text-sm text-base-content/70 mb-4">Get help with issues and questions</p>
        <a href="#" class="btn btn-secondary btn-sm">Get Support</a>
      </div>
      
      <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 mx-auto text-accent mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
        </svg>
        <h3 class="font-semibold mb-2">GitHub</h3>
        <p class="text-sm text-base-content/70 mb-4">Source code and issue tracking</p>
        <a href="#" class="btn btn-accent btn-sm">View on GitHub</a>
      </div>
    </div>

    <!-- System Information -->
    <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        System Information
      </h3>
      
      <div class="grid gap-4 md:grid-cols-2">
        <div class="space-y-3">
          <div class="flex justify-between">
            <span class="text-sm text-base-content/70">Application</span>
            <span class="text-sm font-medium">GoDumper v{{ systemInfo.version }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-sm text-base-content/70">Build</span>
            <span class="text-sm font-mono">{{ systemInfo.build }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-sm text-base-content/70">Go Version</span>
            <span class="text-sm font-mono">{{ systemInfo.goVersion }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-sm text-base-content/70">Platform</span>
            <span class="text-sm font-mono">{{ systemInfo.platform }}</span>
          </div>
        </div>
        
        <div class="space-y-3">
          <div class="flex justify-between">
            <span class="text-sm text-base-content/70">Uptime</span>
            <span class="text-sm font-medium">{{ systemInfo.uptime }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-sm text-base-content/70">Backup Directory</span>
            <span class="text-sm font-mono">{{ systemInfo.backupDir }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-sm text-base-content/70">Database Path</span>
            <span class="text-sm font-mono">{{ systemInfo.databasePath }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-sm text-base-content/70">Config Files</span>
            <span class="text-sm font-medium">{{ systemInfo.configFiles.join(', ') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Start Guide -->
    <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
        </svg>
        Quick Start Guide
      </h3>
      
      <div class="steps steps-vertical lg:steps-horizontal">
        <div class="step step-primary">
          <div class="step-content text-left">
            <h4 class="font-medium">Configure Database</h4>
            <p class="text-sm text-base-content/70">Add your MySQL/MariaDB connection details</p>
          </div>
        </div>
        <div class="step step-primary">
          <div class="step-content text-left">
            <h4 class="font-medium">Create Backup</h4>
            <p class="text-sm text-base-content/70">Run manual backup or schedule automatic ones</p>
          </div>
        </div>
        <div class="step step-primary">
          <div class="step-content text-left">
            <h4 class="font-medium">Manage Files</h4>
            <p class="text-sm text-base-content/70">Download, restore, or delete your backups</p>
          </div>
        </div>
      </div>
      
      <div class="mt-6 flex gap-3">
        <router-link to="/configure" class="btn btn-primary">
          Get Started
        </router-link>
        <button @click="showTutorial = true" class="btn btn-outline">
          Interactive Tutorial
        </button>
      </div>
    </div>

    <!-- FAQ Section -->
    <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
      <h3 class="text-lg font-semibold text-base-content mb-4 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-info" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Frequently Asked Questions
      </h3>
      
      <div class="space-y-4">
        <div v-for="(faq, index) in faqs" :key="index" class="collapse collapse-arrow bg-base-200">
          <input type="radio" name="faq-accordion" :id="`faq-${index}`" />
          <div class="collapse-title text-sm font-medium">
            <label :for="`faq-${index}`" class="cursor-pointer">{{ faq.question }}</label>
          </div>
          <div class="collapse-content">
            <div class="text-sm text-base-content/80" v-html="faq.answer"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Credits and License -->
    <div class="grid gap-4 md:grid-cols-2">
      <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
        <h3 class="text-lg font-semibold text-base-content mb-4 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-warning" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
          </svg>
          Credits
        </h3>
        
        <div class="space-y-3">
          <div>
            <div class="text-sm font-medium">GoDumper</div>
            <div class="text-xs text-base-content/70">Open-source MySQL/MariaDB backup tool</div>
          </div>
          <div>
            <div class="text-sm font-medium">Built with</div>
            <div class="text-xs text-base-content/70">Go, Vue.js, TailwindCSS, DaisyUI</div>
          </div>
          <div>
            <div class="text-sm font-medium">Inspired by</div>
            <div class="text-xs text-base-content/70">MySQLDumper project</div>
          </div>
        </div>
      </div>
      
      <div class="bg-base-100 p-6 rounded-lg shadow-sm border border-base-300">
        <h3 class="text-lg font-semibold text-base-content mb-4 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
          </svg>
          License & Security
        </h3>
        
        <div class="space-y-3">
          <div>
            <div class="text-sm font-medium">MIT License</div>
            <div class="text-xs text-base-content/70">Free and open source software</div>
          </div>
          <div>
            <div class="text-sm font-medium">Security</div>
            <div class="text-xs text-base-content/70">AES-GCM password encryption</div>
          </div>
          <div>
            <div class="text-sm font-medium">Privacy</div>
            <div class="text-xs text-base-content/70">Self-hosted, no external data sharing</div>
          </div>
        </div>
        
        <div class="mt-4 pt-4 border-t border-base-300">
          <div class="text-xs text-base-content/60">
            © 2024 GoDumper. All rights reserved.
          </div>
        </div>
      </div>
    </div>

    <!-- Tutorial Modal -->
    <div v-if="showTutorial" class="modal modal-open">
      <div class="modal-box max-w-4xl">
        <h3 class="font-bold text-lg mb-4">Interactive Tutorial</h3>
        <p class="text-base-content/70 mb-4">
          This feature would provide a step-by-step interactive tutorial to help new users get started with GoDumper.
        </p>
        <div class="text-center py-8">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-24 h-24 mx-auto text-base-content/20 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1.01M15 10h1.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="text-base-content/50">Tutorial coming soon!</p>
        </div>
        <div class="modal-action">
          <button @click="showTutorial = false" class="btn">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const showTutorial = ref(false)

const systemInfo = ref({
  version: '1.0.0',
  build: 'a1b2c3d',
  goVersion: 'go1.22.0',
  platform: 'linux/amd64',
  uptime: '2 days, 14 hours',
  backupDir: '/data/backups',
  databasePath: '/data/app/app.db',
  configFiles: ['.env', '.env']
})

const faqs = ref([
  {
    question: 'How do I configure my first database connection?',
    answer: 'Go to the <strong>Configure</strong> section and click "Add New Configuration". Enter your database connection details including host, port, database name, username, and password. Test the connection before saving.'
  },
  {
    question: 'How are passwords stored securely?',
    answer: 'Passwords are encrypted using AES-GCM encryption with a 256-bit key. The encryption key is stored as an environment variable and passwords are never stored in plain text.'
  },
  {
    question: 'Can I schedule automatic backups?',
    answer: 'Yes! When configuring a database target, you can set a daily backup time in UTC. The system will automatically create backups at the scheduled time and manage retention according to your settings.'
  },
  {
    question: 'What backup formats are supported?',
    answer: 'GoDumper creates standard SQL dump files (.sql) which can be optionally compressed to .sql.gz format. These files are compatible with MySQL and MariaDB and can be restored using standard tools.'
  },
  {
    question: 'How do I restore a backup to a different database?',
    answer: 'Use the <strong>Restore</strong> section to select any backup file and any target database configuration. You can restore to the same database or a different one. Always ensure you have a current backup before restoring.'
  },
  {
    question: 'Can I run GoDumper in Docker?',
    answer: 'Yes! GoDumper is designed to run in Docker containers. Mount volumes for <code>/data/backups</code> (backup storage) and <code>/data/app</code> (application database). See the documentation for complete Docker setup instructions.'
  },
  {
    question: 'What happens if a backup fails?',
    answer: 'Failed backups are logged with detailed error information. You can view these in the <strong>Logs</strong> section. The system will retry scheduled backups automatically, and you can also run manual backups to troubleshoot issues.'
  },
  {
    question: 'Is there an API for automation?',
    answer: 'Yes! GoDumper provides a REST API at <code>/api/*</code> endpoints. You can manage configurations, trigger backups, and download backup files programmatically. API documentation is available in the full documentation.'
  }
])

onMounted(async () => {
  // In real app: fetch system info from API
  try {
    // const response = await api.getSystemInfo()
    // systemInfo.value = response.data
  } catch (error) {
    console.error('Failed to fetch system info:', error)
  }
})
</script>