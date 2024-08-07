import * as Sentry from '@sentry/react';

import appConfig from './appConfig';

const initSentry = () => {
  const dsn = appConfig.sentryDSN;
  const environment = appConfig.sentryENV ? appConfig.sentryENV : appConfig.nodeENV;
  if (dsn) {
    Sentry.init({
      dsn: dsn,
      integrations: [Sentry.browserTracingIntegration()],
      environment: environment,
      tracesSampleRate: 1.0,
    });

    // eslint-disable-next-line no-console
    console.info('sentry tracing initialized');
  } else {
    // eslint-disable-next-line no-console
    console.warn('no sentry configuration available');
  }
};

export default initSentry;
