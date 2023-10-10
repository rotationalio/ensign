import { PATH_DASHBOARD } from '@/application';

import { removeCookie } from './cookies';
export const StatusIconMap = {};

export const isCurrentMenuPath = (href: string, pathname: string, href_linked?: string) => {
  const hrefArray = href.split('/');
  const hrefLinkedArray = href_linked?.split('/') || [];
  const pathnameArray = pathname.split('/');

  if (href_linked === PATH_DASHBOARD.TOPICS && pathnameArray[2] === 'topics') {
    return hrefLinkedArray[1] === pathnameArray[1] && hrefLinkedArray[2] === pathnameArray[2];
  }

  if (hrefArray.length > 2 && pathnameArray.length > 2) {
    return hrefArray[1] === pathnameArray[1] && hrefArray[2] === pathnameArray[2];
  }

  return href === pathname;
};

export function updateQueryStringValueWithoutNavigation(
  queryKey: string,
  queryValue: string | null
) {
  const currentSearchParams = new URLSearchParams(window.location.search);
  const oldQuery = currentSearchParams.get(queryKey) ?? '';
  if (queryValue === oldQuery) return;

  if (queryValue) {
    currentSearchParams.set(queryKey, queryValue);
  } else {
    currentSearchParams.delete(queryKey);
  }
  const newUrl = [window.location.pathname, currentSearchParams.toString()]
    .filter(Boolean)
    .join('?');

  window.history.replaceState(null, '', newUrl);
}

export const cleanCookiesOnDashboard = () => {
  // ensure we don't have any of this cookies on dashboard
  removeCookie('invitee_token');
  removeCookie('isInvitedUser');
  removeCookie('esg.new.user');
};
