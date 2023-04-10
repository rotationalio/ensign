function path(root: string, sublink: string) {
  return `${root}${sublink}`;
}
const ROOT = '/app';

export const ROUTES = {
  HOME: '/',
  DASHBOARD: '/dashboard',
  DOCS: 'https://ensign.rotational.dev/getting-started/',
  SUPPORT: 'support@rotational.io',
  PROFILE: '/profile',
  WELCOME: '/welcome',
  SETUP: '/onboarding/setup',
  COMPLETE: '/onboarding/complete',
  VERIFY_PAGE: '/verify-account',
  VERIFY_EMAIL: '/verify',
  NEW_INVITATION: '/new-user-invitation',
  EXISTING_INVITATION: '/existing-user-invitation',
};

export const PATH_DASHBOARD = {
  ROOT: ROOT,
  HOME: path(ROOT, ''),
  PROJECTS: path(ROOT, '/projects'),
  PROFILE: path(ROOT, '/profile'),
  ORGANIZATION: path(ROOT, '/organization'),
  TEAMS: path(ROOT, '/team'),
};

export const EXTRENAL_LINKS = {
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
};
