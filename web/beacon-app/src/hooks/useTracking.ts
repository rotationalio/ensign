import React, { useEffect } from 'react';
import ReactGA from 'react-ga4';

import { isProdEnv } from '@/application/config/appEnv';

import appConfig from '../application/config/appConfig';

const useTracking = () => {
  const [isInitialized, setIsInitialized] = React.useState<boolean>(false);
  const trackingID: any = appConfig.analyticsID;
  useEffect(() => {
    //  initialize google analytics only in production environment
    if (isProdEnv && trackingID) {
      // eslint-disable-next-line no-console
      console.log('initializing google analytics');
      ReactGA.initialize(trackingID, {
        gaOptions: {
          siteSpeedSampleRate: 100,
        },
      });
    }
    setIsInitialized(true);
  }, [trackingID]);
  return {
    isInitialized,
  };
};

export default useTracking;
