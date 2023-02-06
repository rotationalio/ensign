const appConfig = {
  apiUrl: import.meta.env.REACT_APP_QUARTERDECK_BASE_URL || 'http://localhost:8088',
  apiVersion: 'v1',
  sentryDSN: import.meta.env.REACT_APP_SENTRY_DSN,
  sentryENV: import.meta.env.REACT_APP_SENTRY_ENVIRONMENT,
  nodeENV: import.meta.env.NODE_ENV,
};

export default appConfig;
