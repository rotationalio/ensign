import * as Sentry from '@sentry/react';
import { BrowserTracing } from '@sentry/tracing';

const initSentry = () => {
    const environment = process.env.REACT_APP_SENTRY_ENVIRONMENT ? process.env.REACT_APP_SENTRY_ENVIRONMENT : process.env.NODE_ENV;

    Sentry.init({
        dsn: process.env.REACT_APP_SENTRY_DSN,
        integrations: [new BrowserTracing()],
        environment: environment,
        tracesSampleRate: 1.0,
    });
}

export default initSentry