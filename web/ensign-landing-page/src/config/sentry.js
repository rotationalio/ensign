import * as Sentry from '@sentry/react';
import { BrowserTracing } from '@sentry/tracing';

const initSentry = () => {
    Sentry.init({
        dsn: process.env.ENSIGN_UI_SENTRY_DSN,
        integrations: [new BrowserTracing()],
    
        tracesSampleRate: 1.0,
    });
}

export default initSentry