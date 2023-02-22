import { defineConfig } from 'cypress';
import cucumber from 'cypress-cucumber-preprocessor';
import * as browserify from '@cypress/browserify-preprocessor';

export default defineConfig({
    numTestsKeptInMemory: 0,
    //defaultCommandTimeout: 10000,
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
        setupNodeEvents(on) {
            const options = {
               ...browserify.defaultOptions,
                typescript: require.resolve('typescript'),
            };

            on('file:preprocessor', cucumber(options));
        },
        specPattern: 'cypress/e2e/**/*.{feature,features}',
    },
    "reporter": "cypress-multi-reporters",
    "reporterOptions": {
      "reporterEnabled": ["mochawesome"],
      "mochawesomeReporterOptions": {
        "reportDir": "cypress/results",
        "overwrite": false,
        "html": false,
        "json": true
      }

    }
});
