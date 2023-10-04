import * as browserify from '@cypress/browserify-preprocessor';
import { defineConfig } from 'cypress';
import cucumber from 'cypress-cucumber-preprocessor';

export default defineConfig({
  numTestsKeptInMemory: 0,
  env: {
    uncaughtCypressException: false,
    hideXhr: true,
  },
  chromeWebSecurity: false,
  retries: {
    runMode: 1,
    openMode: 0,
  },
  e2e: {
    baseUrl: 'http://localhost:3000',
    setupNodeEvents(on, config) {
      const options = {
        ...browserify.defaultOptions,
        typescript: require.resolve('typescript'),
      };

      on('file:preprocessor', cucumber(options));
      // eslint-disable-next-line @typescript-eslint/no-var-requires
      require('@cypress/code-coverage/task')(on, config);

      return config;
    },
    specPattern: 'cypress/e2e/**/*.{feature,features}',
    env: {
      API_URL: 'http://localhost:8080/v1',
    },
  },
  reporter: 'cypress-multi-reporters',
  reporterOptions: {
    reporterEnabled: ['mochawesome'],
    mochawesomeReporterOptions: {
      reportDir: 'cypress/results',
      overwrite: false,
      html: false,
      json: true,
    },
  },
});
