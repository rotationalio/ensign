import { defineConfig } from 'vite';
import path from 'path';
import react from '@vitejs/plugin-react';
import eslint from 'vite-plugin-eslint';
import svgrPlugin from 'vite-plugin-svgr';
import tsConfigPaths from 'vite-tsconfig-paths';
import { lingui } from '@lingui/vite-plugin';
import UnpluginFonts from 'unplugin-fonts/vite';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react({
      babel: {
        plugins: ['macros'],
      },
    }),
    lingui(),
    eslint(),
    svgrPlugin(),
    tsConfigPaths(),
    UnpluginFonts({
      google: {
        families: ['Quattrocento', 'PT Mono'],
      },
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    outDir: 'build',
    rollupOptions: {
      output: {
        entryFileNames: 'assets/[name].js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name].[ext]',
      },
    },
  },
  envPrefix: ['VITE_', 'REACT_APP_'],
  appType: 'spa',
  server: {
    port: 3000,
    strictPort: true,
    host: true,
  },
});
