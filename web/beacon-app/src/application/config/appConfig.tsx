const appConfig = {
  nodeENV: import.meta.env.MODE,

  // Backend connection information
  quaterDeckApiUrl: import.meta.env.REACT_APP_QUARTERDECK_BASE_URL || 'http://localhost:8088/v1',
  tenantApiUrl: import.meta.env.REACT_APP_TENANT_BASE_URL || 'http://localhost:8080/v1',

  // Sentry information
  sentryDSN: import.meta.env.REACT_APP_SENTRY_DSN,
  sentryENV: import.meta.env.REACT_APP_SENTRY_ENVIRONMENT,
  sentryEventID: import.meta.env.REACT_APP_SENTRY_EVENT_ID,

  // Google Analytics tag
  analyticsID: import.meta.env.REACT_APP_ANALYTICS_ID,

  // App version information from build workflow
  version: import.meta.env.REACT_APP_VERSION_NUMBER || '0.0.0',
  revision: import.meta.env.REACT_APP_GIT_REVISION || '0.0.0',
  // TODO: need to parse boolean from environment variable
  useDashLocale: import.meta.env.REACT_APP_USE_DASH_LOCALE,

  beaconVersion: import.meta.env.REACT_APP_BEACON_VERSION,
  apiVersion: import.meta.env.REACT_APP_API_VERSION,
};

export default appConfig;
