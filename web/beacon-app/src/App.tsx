import { i18n } from '@lingui/core';
import { I18nProvider } from '@lingui/react';
import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { useEffect } from 'react';
import { RouterProvider } from 'react-router-dom';

import { isDevEnv } from './application/config/appEnv';
import { queryClient } from './application/config/react-query';
import router from './application/routes/root';
import GoogleAnalyticsWrapper from './components/GaWrapper';
import useTracking from './hooks/useTracking';
import { defaultLocale, dynamicActivate } from './I18n';

function App() {
  useEffect(() => {
    dynamicActivate(defaultLocale);
  }, []);

  const { isInitialized } = useTracking();

  return (
    <I18nProvider i18n={i18n}>
      <QueryClientProvider client={queryClient}>
        <ReactQueryDevtools initialIsOpen={!!isDevEnv} />
        <GoogleAnalyticsWrapper isInitialized={isInitialized}>
          <RouterProvider router={router} />
        </GoogleAnalyticsWrapper>
      </QueryClientProvider>
    </I18nProvider>
  );
}

export default App;
