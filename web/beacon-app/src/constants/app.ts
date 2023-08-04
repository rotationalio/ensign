export const APP_ROUTE = {
  ROOT: '/',
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  FORGOT_PASSWORD: '/forgot-password',
  RESET_PASSWORD: '/reset-password',
  DASHBOARD: '/app',
  TENANTS: '/tenant',
  APIKEYS: '/apikeys',
  PROJECTS: '/projects',
  TOPICS: '/topics',
  GETTING_STARTED: '/onboarding/getting-started',
  ONBOARDING_SETUP: '/onboarding/setup',
  MEMBERS_LIST: '/members',
  MEMBERS: '/members',
  ORG_DETAIL: '/organization',
  ORGANIZATION: '/organization',
  TOKEN: '/verify',
  INVITE: '/invites',
  SWITCH: '/switch',
};

// quaterdeck api routes

export const QDK_API_ROUTE = {
  LOGIN: '/login',
  REGISTER: '/register',
  REFRESH_TOKEN: '/refresh',
};

// mime type constants
export enum MIME_TYPE {
  JSON = 'application/json',
  XML = 'application/xml',
  TEXT_PLAIN = 'text/plain',
  TEXT_HTML = 'text/html',
  CSV = 'text/csv',
  BSON = 'application/bson',
  PDF = 'application/pdf',
  AVRO = 'application/avro',
  PROTOBUF = 'application/protobuf',
  LD_JSON = 'application/ld+json',
  JSONLINES = 'application/jsonlines',
  UB_JSON = 'application/ubjson',
  ATOM_XML = 'application/atom+xml',
  MSGPACK = 'application/msgpack',
  JAVA_ARCHIVE = 'application/java-archive',
  PYTHON_PICKLE = 'application/python-pickle',
  CALENDAR = 'text/calendar',
  OCTET_STREAM = 'application/octet-stream',
  PARQUET = 'application/parquet',
}
