import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

const apiProxyTarget = process.env.VITE_API_PROXY_TARGET || process.env.BACKEND_URL || ''
const port = Number(process.env.PORT) || undefined

export default defineConfig({
  base: './',
  plugins: [react(), tailwindcss()],
  server: {
    port,
    ...(apiProxyTarget
      ? {
          proxy: {
            '/api': {
              target: apiProxyTarget,
              changeOrigin: true,
            },
          },
        }
      : {}),
  },
})
