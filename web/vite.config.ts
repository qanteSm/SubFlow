import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
  },
  define: {
    __ARCHITECT__: JSON.stringify('Muhammet-Ali-Buyuk'),
  },
});
