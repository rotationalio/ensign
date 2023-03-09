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
};

export const PATH_DASHBOARD = {
  ROOT: ROOT,
  HOME: path(ROOT, ''),
  PROJECTS: path(ROOT, '/projects'),
  PROFILE: path(ROOT, '/profile'),
  ORGANIZATION: path(ROOT, '/organization'),
};

export const FOOTER = {
  ABOUT: 'https://rotational.io/about',
  CONTACT: 'https://rotational.io/contact',
  SERVER: '',
};

export const EXTRENAL_LINKS = {
  DOCUMENTATION: 'https://ensign.rotational.dev/getting-started/',
  TUTORIAL: 'https://youtube.com/@rotationalio',
  OTHERS: 'https://twitter.com/In_Otter_News2',
};
