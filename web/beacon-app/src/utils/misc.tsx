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
