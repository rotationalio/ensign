export const StatusIconMap = {};

export const isCurrentMenuPath = (href: string, pathname: string) => {
  const hrefArray = href.split('/');
  const pathnameArray = pathname.split('/');

  if (hrefArray.length > 2 && pathnameArray.length > 2) {
    return hrefArray[1] === pathnameArray[1] && hrefArray[2] === pathnameArray[2];
  }

  if (href === pathname) {
    return true;
  }
};
