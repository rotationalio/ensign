export const APP_ROUTE = {
  ROOT: '/',
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  FORGOT_PASSWORD: '/forgot-password',
  RESET_PASSWORD: '/reset-password',
  DASHBOARD: '/dashboard',
  TENANTS: '/tenant',
  APIKEYS: '/apikeys',
  PROJECTS: '/projects',
  PROJECTS_LIST: 'tenant/{:tenantID}/projects',
  GETTING_STARTED: '/onboarding/getting-started',
  ONBOARDING_SETUP: '/onboarding/setup',
  MEMBERS_LIST: 'tenant/{:tenantID}/members',
  ORG_DETAIL: '/organization/{:orgID}'
};

// quaterdeck api routes

export const QDK_API_ROUTE = {
  LOGIN: '/login',
  REGISTER: '/register',
  REFRESH_TOKEN: '/refresh',
};
