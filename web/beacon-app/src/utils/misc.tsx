import { PATH_DASHBOARD } from '@/application';
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
  // alright, let's talk about this...
  // Normally with remix, you'd update the params via useSearchParams from react-router-dom
  // and updating the search params will trigger the search to update for you.
  // However, it also triggers a navigation to the new url, which will trigger
  // the loader to run which we do not want because all our data is already
  // on the client and we're just doing client-side filtering of data we
  // already have. So we manually call `window.history.pushState` to avoid
  // the router from triggering the loader.
  window.history.replaceState(null, '', newUrl);
}
