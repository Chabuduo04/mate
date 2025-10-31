import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      // 代理后端不同路由分组，开发时将这些路径转发到后端服务
      '/chat': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/user': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/role': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/voice': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
