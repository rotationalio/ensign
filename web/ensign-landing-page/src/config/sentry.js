import * as Sentry from '@sentry/react';
import { BrowserTracing } from '@sentry/tracing';

const initSentry = () => {
    const dsn = process.env.REACT_APP_SENTRY_DSN;
    const environment = process.env.REACT_APP_SENTRY_ENVIRONMENT ? process.env.REACT_APP_SENTRY_ENVIRONMENT : process.env.NODE_ENV;

    if (dsn) {
      Sentry.init({
          dsn: dsn,
          integrations: [new BrowserTracing()],
          environment: environment,
          tracesSampleRate: 1.0,
      });

      // eslint-disable-next-line no-console
      console.log("Sentry tracing initialized");
    } else {
      // eslint-disable-next-line no-console
      console.log("no Sentry configuration available");
    }
}

export default initSentry