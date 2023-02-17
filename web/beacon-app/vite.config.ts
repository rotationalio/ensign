import { defineConfig } from 'vite';
import path from 'path';
import react from '@vitejs/plugin-react';
import eslint from 'vite-plugin-eslint';
import svgrPlugin from 'vite-plugin-svgr';
import tsConfigPaths from 'vite-tsconfig-paths';
import lingui from '@lingui/vite-plugin';
import macrosPlugin from 'vite-plugin-babel-macros';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), eslint(), svgrPlugin(), tsConfigPaths(), macrosPlugin(), lingui()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    outDir: 'build',
  },
  envPrefix: ['VITE_', 'REACT_APP_'],
  appType: 'spa',
});
