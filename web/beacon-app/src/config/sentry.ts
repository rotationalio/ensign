import * as Sentry from '@sentry/react';
import { BrowserTracing } from '@sentry/tracing';

const initSentry = () => {
  const dsn = import.meta.env.REACT_APP_SENTRY_DSN;
  const environment = import.meta.env.REACT_APP_SENTRY_ENVIRONMENT
    ? import.meta.env.REACT_APP_SENTRY_ENVIRONMENT
    : import.meta.env.NODE_ENV;

  if (dsn) {
    Sentry.init({
      dsn: dsn,
      integrations: [new BrowserTracing()],
      environment: environment,
      tracesSampleRate: 1.0,
    });

    // eslint-disable-next-line no-console
    console.log('Sentry tracing initialized');
  } else {
    // eslint-disable-next-line no-console
    console.log('no Sentry configuration available');
  }
};

export default initSentry;
