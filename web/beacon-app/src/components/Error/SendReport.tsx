import * as Sentry from '@sentry/react';

import { appConfig } from '@/application/config';
export const SendReport = () => Sentry.showReportDialog({ eventId: appConfig.sentryEventID });
