import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(({ mode }) => ({
  plugins: [react()],

  server: {
    host: true,
    port: 5173,
    strictPort: true,
    hmr: {
      host: 'localhost',
      port: 5173,
    },

    proxy: {
      '/api': {
        target: 'http://backend:8080',
        changeOrigin: true,
      },
    },
  },

  preview: {
    host: true,
    port: 4173,
    strictPort: true,
  },

  build: {
    outDir: 'dist',
    sourcemap: mode !== 'production',
    emptyOutDir: true,
  },

  define: {
    __APP_ENV__: JSON.stringify(mode),
  },
}))