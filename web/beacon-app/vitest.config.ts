import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./vitest.setup.ts'],
    coverage: {
      provider: 'istanbul',
    },
    // path resolution
    alias: {
      '@': './src',
    },
  },
});
