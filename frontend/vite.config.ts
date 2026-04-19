import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  optimizeDeps: {
    include: ["@syncfusion/ej2-react-querybuilder" ]
  },
  resolve: {
    conditions: ["module", "browser", "import"]
  }
})