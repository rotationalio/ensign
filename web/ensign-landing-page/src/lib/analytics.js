import React, { useEffect } from 'react';
import GA4React from 'ga-4-react';

const gaCode = process.env.ENSIGN_UI_ANALYTICS_ID;

const withTracker = (Component) => (props) => {
    useEffect(() => {
      switch(GA4React.isInitialized()) {
        case true:
          const ga4 = GA4React.getGA4React();
          if (ga4) {
            ga4.pageview(window.location.pathname);
          }
          break
        default:
        case false:
          const ga4react = new GA4React(gaCode);
          ga4react.initialize().then((ga4) => {
            console.log(`tracking initialized with ${gaCode}`);
            ga4.pageview(window.location.pathname);
          }, (err) => {
            console.error(err);
          });
          break
      }
    });
    return <Component {...props} />;
  };
  
  export default withTracker;