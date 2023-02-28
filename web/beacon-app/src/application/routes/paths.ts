function path(root: string, sublink: string) {
  return `${root}${sublink}`;
}
const ROOT = '/app';

export const ROUTES = {
  HOME: '/',
  DASHBOARD: '/dashboard',
  DOCS: 'https://ensign.rotational.dev/getting-started/',
  SUPPORT: '',
  PROFILE: '/profile',
  WELCOME: '/welcome',
  SETUP: '/onboarding/setup',
  COMPLETE: '/onboarding/complete',
  VERIFY_PAGE: '/verify-account',
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
