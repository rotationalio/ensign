const appConfig = {
  apiUrl: 'http://localhost:8088',
  apiVersion: 'v1',
  sentryDSN: process.env.REACT_APP_SENTRY_DSN,
  sentryENV: process.env.REACT_APP_SENTRY_ENVIRONMENT,
  nodeENV: process.env.NODE_ENV,
};

export default appConfig;
