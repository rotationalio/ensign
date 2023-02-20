function path(root: string, sublink: string) {
  return `${root}${sublink}`;
}
const ROOT = '/app';

export const ROUTES = {
  HOME: '/',
  DASHBOARD: '/dashboard',
  DOCS: '/docs',
  SUPPORT: '/support',
  PROFILE: '/profile',
  WELCOME: '/welcome',
  SETUP: '/onboarding/setup',
  COMPLETE: '/onboarding/complete',
};

export const PATH_DASHBOARD = {
  ROOT: ROOT,
  HOME: path(ROOT, ''),
  PROJECTS: path(ROOT, '/projects'),
  PROFILE: path(ROOT, '/profile'),
  ORGANIZATION: path(ROOT, '/organization'),
};
