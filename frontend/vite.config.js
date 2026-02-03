import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,  
    host: '0.0.0.0',  
    proxy: {
      '/auth': {
        target: 'http://backend:8080', 
        changeOrigin: true,  
      },
      '/bookmarks': {
        target: 'http://backend:8080',
        changeOrigin: true,
      },
    },
  },
})
