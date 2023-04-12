import React from 'react';
import ReactGA from 'react-ga4';

interface IProps {
  children: React.ReactNode;
  isInitialized: boolean;
}

const GoogleAnalyticsWrapper: React.FC<IProps> = ({ children, isInitialized }) => {
  const pathname = window.location.href;
  const search = window.location.search;

  React.useEffect(() => {
    if (isInitialized) {
      // ReactGA.set({ page: location.pathname });
      ReactGA.send({ hitType: 'pageview', page: pathname + search });
    }
  }, [isInitialized, pathname, search]);

  return <>{children}</>;
};

export default GoogleAnalyticsWrapper;
