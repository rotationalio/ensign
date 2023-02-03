import { defineConfig } from 'vite';
import path from 'path';
import react from '@vitejs/plugin-react-swc';
import eslint from 'vite-plugin-eslint';
import svgrPlugin from 'vite-plugin-svgr';
import tsConfigPaths from 'vite-tsconfig-paths';
import { swcReactRefresh } from "vite-plugin-swc-react-refresh";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    swcReactRefresh(),
    eslint(),
    svgrPlugin(),
    tsConfigPaths(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    outDir: 'build',
  },
});
