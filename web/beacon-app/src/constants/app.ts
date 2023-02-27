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
  PROJECTS_LIST: '/{:tenantID}/projects',
  GETTING_STARTED: '/onboarding/getting-started',
  ONBOARDING_SETUP: '/onboarding/setup',
  MEMBERS_LIST: 'tenant/{:tenantID}/members',
  MEMBERS: '/members',
  ORG_DETAIL: '/organization',
  ORGANIZATION: '/organization',
};

// quaterdeck api routes

export const QDK_API_ROUTE = {
  LOGIN: '/login',
  REGISTER: '/register',
  REFRESH_TOKEN: '/refresh',
};
