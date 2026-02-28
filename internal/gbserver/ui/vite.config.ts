import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // 匹配所有 API 路径
      '^/(device|channel|control|play|index|swagger|subscribe)': {
        target: 'http://localhost:18080',
        changeOrigin: true
      }
    }
  }
})
