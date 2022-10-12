import * as Sentry from '@sentry/react';
import { BrowserTracing } from '@sentry/tracing';

const initSentry = () => {
    const environment = process.env.ENSIGN_UI_SENTRY_ENVIRONMENT ? process.ENSIGN_UI_SENTRY_ENVIRONMENT : process.env.NODE_ENV;

    Sentry.init({
        dsn: process.env.ENSIGN_UI_SENTRY_DSN,
        integrations: [new BrowserTracing({tracingOrigins})],
        environment: environment,
        tracesSampleRate: 1.0,
    });
}

export default initSentry