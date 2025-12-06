import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
    extensions: ['.mjs', '.js', '.mts', '.ts', '.jsx', '.tsx', '.json'],
  },
  
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
  
  server: {
    host: '0.0.0.0',
    port: 5173,
  },
  
  // إعدادات Vitest
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
    include: ['src/**/*.{test,spec}.{js,jsx,ts,tsx}'],
    
    // إعدادات خاصة لحل مشاكل الاستيراد
    deps: {
      inline: [
        /@mui/,
        /@testing-library/,
      ],
    },
    
    // Mock للملفات غير الموجودة
    mockReset: true,
    
    // إعدادات التغطية
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/test/',
        '**/*.d.ts',
        '**/*.config.*',
        '**/types.ts',
      ],
    },
  },
})