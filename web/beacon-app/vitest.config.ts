import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: './setupTests.js',
    coverage: {
      provider: 'istanbul',
    },
    // path resolution
    alias: {
      '@': './src',
    },
  },
});
