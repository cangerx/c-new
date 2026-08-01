import { fileURLToPath, URL } from 'node:url'

import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

// The homepage plugin is served by the Go backend, not by new-api's own
// Rsbuild pipeline. Assets go under /homepage-assets/ so they never collide
// with new-api's own /assets/ prefix, which its static middleware owns.
export default defineConfig({
  base: '/homepage-assets/',
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    outDir: '../dist',
    emptyOutDir: true,
    assetsDir: '.',
  },
})
