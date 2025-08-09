<template>
  <div class="drawer lg:drawer-open">
    <input id="drawer-toggle" type="checkbox" class="drawer-toggle" />
    
    <!-- Sidebar -->
    <Sidebar />
    
    <!-- Main Content -->
    <div class="drawer-content flex flex-col min-h-screen">
      <!-- Top Navigation -->
      <div class="navbar bg-base-300 lg:hidden">
        <div class="navbar-start">
          <label for="drawer-toggle" class="btn btn-square btn-ghost">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </label>
        </div>
        <div class="navbar-center">
          <span class="text-lg font-semibold">GoDumper</span>
        </div>
        <div class="navbar-end">
          <!-- Theme toggle or other controls -->
        </div>
      </div>

      <!-- Page Content -->
      <main class="flex-1 p-6 bg-base-100">
        <router-view />
      </main>
    </div>

    <Toasts />
  </div>
</template>

<script setup lang="ts">
import Sidebar from '@/components/Sidebar.vue'
import Toasts from '@/components/Toasts.vue'
</script>