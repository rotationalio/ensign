/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly REACT_APP_SENTRY_ENVIRONMENT: string;
  readonly REACT_APP_SENTRY_DSN: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
