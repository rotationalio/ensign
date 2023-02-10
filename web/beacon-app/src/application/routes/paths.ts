function path(root: string, sublink: string) {
  return `${root}${sublink}`;
}

export const routes = {
  home: '/',
  docs: '/docs',
  support: '/support',
  profile: '/profile',
  welcome: '/welcome',
  setup: '/onboarding/setup',
  complete: '/onboarding/complete',
};

const ROOTS_DASHBOARD = '/app';

export const PATH_DASHBOARD = {
  root: ROOTS_DASHBOARD,
  home: path(ROOTS_DASHBOARD, ''),
  project: path(ROOTS_DASHBOARD, '/projects'),
  profile: path(ROOTS_DASHBOARD, '/profile'),
};
