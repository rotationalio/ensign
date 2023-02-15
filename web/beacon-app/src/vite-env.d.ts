/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly REACT_APP_QUARTERDECK_BASE_URL: string;
  readonly REACT_APP_TENANT_BASE_URL: string;

  readonly REACT_APP_SENTRY_DSN: string;
  readonly REACT_APP_SENTRY_ENVIRONMENT: string;
  readonly REACT_APP_SENTRY_EVENT_ID: string;

  readonly REACT_APP_ANALYTICS_ID: string;

  readonly REACT_APP_VERSION_NUMBER: string;
  readonly REACT_APP_GIT_REVISION: string;

  readonly REACT_APP_USE_DASH_LOCALE: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
