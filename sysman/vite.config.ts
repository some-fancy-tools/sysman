import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    proxy: {
      "/api/": {
        target: "http://ginger.shark-scala.ts.net:8000",
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
