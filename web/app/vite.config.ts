import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), '')
    const API_TARGET = env.VITE_DEV_API_TARGET || 'http://localhost:8080'

    return {
        plugins: [vue()],
        resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
        server: {
            proxy: {
                '/api': {
                    target: API_TARGET,
                    changeOrigin: true,
                    secure: false,
                },
            },
        },
    }
})
