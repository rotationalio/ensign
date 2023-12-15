function path(root: string, sublink: string) {
  return `${root}${sublink}`;
}
const ROOT = '/app';

/* ROUTES CONSTANTS 

 * - ROUTES are used in the following files to define the routes of the application
 * - PATH_DASHBOARD is used in the following files to define the routes of the dashboard
 * - EXTERNAL_LINKS are used in the following files to define the external links of the application and came be found in external links file in the constants folder

 */

// TODO: Move the external links to a separate file in the constants folder, only routes should be here

export const ROUTES = {
  HOME: '/',
  DASHBOARD: '/dashboard',
  PROFILE: '/profile',
  WELCOME: '/welcome',
  SETUP: '/onboarding/setup',
  COMPLETE: '/onboarding/complete',
  VERIFY_PAGE: '/verify-account',
  VERIFY_EMAIL: '/verify',
  NEW_INVITATION: '/new-user-invitation',
  EXISTING_INVITATION: '/existing-user-invitation',
  LOGIN: '/',
  REGISTER: '/register',
  FORGOT_PASSWORD: '/forgot-password',
  RESET_PASSWORD: '/reset-password',
  RESET_VERIFICATION: '/reset-verification',
};

export const PATH_DASHBOARD = {
  ROOT: ROOT,
  HOME: path(ROOT, ''),
  PROJECTS: path(ROOT, '/projects'),
  PROJECT: path(ROOT, '/projects/:id'),
  PROFILE: path(ROOT, '/profile'),
  ORGANIZATION: path(ROOT, '/organization'),
  TEAMS: path(ROOT, '/team'),
  TOPICS: path(ROOT, '/topics'),
  TEMPLATES: path(ROOT, '/templates'), // can be updated when is ready
  ONBOARDING: path(ROOT, '/onboarding'),
};

export const EXTERNAL_LINKS = {
  ABOUT: 'https://rotational.io/about',
  BLOG: 'https://rotational.io/blog',
  CONTACT: 'https://rotational.io/contact',
  DOCUMENTATION: 'https://ensign.rotational.dev/getting-started/',
  EMAIL_US: 'mailto:info@rotational.io',
  GITHUB: 'https://github.com/rotationalio',
  LINKEDIN: 'https://www.linkedin.com/company/rotational',
  OPEN_SOURCE: 'https://rotational.io/opensource',
  OTHERS: 'https://twitter.com/In_Otter_News2',
  PRIVACY: 'https://rotational.io/privacy',
  ROTATIONAL: 'https://rotational.io',
  SDKs: 'https://ensign.rotational.dev/sdk/',
  SERVER: 'https://status.rotational.dev/',
  SERVICES: 'https://rotational.io/services',
  TERMS: 'https://rotational.io/terms',
  TUTORIAL: 'https://youtube.com/@rotationalio',
  TWITTER: 'https://twitter.com/rotationalio',
  ENSQL: 'https://ensign.rotational.dev/ensql/',
  OFFICE_HOURS_SCHEDULE: 'https://calendar.app.google/1r7PuDPzKp2jjHPX8',
  DATA_FLOW_OVERVIEW: 'https://www.youtube.com/watch?v=XUnEHGZXxmM&t=1600s',
  NAMING_TOPICS_GUIDE: 'https://ensign.rotational.dev/getting-started/topics/',
  DATA_PLAYGROUND: 'https://rotational.io/data-playground/',
  SDK_DOCUMENTATION: 'https://ensign.rotational.dev/sdk',
  ENSIGN_UNIVERSITY: 'https://rotational.io/ensign-u/',
  USE_CASES: 'https://ensign.rotational.dev/eventing/use_cases/',
  DOCS: 'https://ensign.rotational.dev/getting-started/',
  SUPPORT: 'support@rotational.io',
  PROTECT_API_KEYS_VIDEO: 'https://youtu.be/EEpIDkKJopY',
  ENSIGN_PRICING: 'https://rotational.io/ensign-pricing/',
  RESOURCES: 'https://rotational.io/resources/',
};
