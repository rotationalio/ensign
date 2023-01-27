import * as Sentry from '@sentry/react';
import { BrowserTracing } from '@sentry/tracing';

import appConfig from './appConfig';

const initSentry = () => {
  const dsn = appConfig.sentryDSN;
  const environment = appConfig.sentryENV ? appConfig.sentryENV : appConfig.nodeENV;
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
