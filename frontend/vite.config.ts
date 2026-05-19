import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig, loadEnv } from 'vite'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const target = env.VITE_API_URL || 'https://backend.pharma-hub.ru'

  return {
    server: {
      port: 3000,
      proxy: {
        '/api': {
          target,
          changeOrigin: true,
          secure: false,
        },
        '^/auth/(google|register|refresh|logout|send-code|verify-code)': {
          target,
          changeOrigin: true,
          secure: false,
        }
      }
    },
    plugins: [react(), tailwindcss()],
    resolve: {
      tsconfigPaths: true,
    },
  }
})
