import { i18n } from '@lingui/core';
import { I18nProvider } from '@lingui/react';
import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { useEffect } from 'react';
import { DefaultToastOptions, Toaster } from 'react-hot-toast';
import { RouterProvider } from 'react-router-dom';

import { isDevEnv } from './application/config/appEnv';
import { queryClient } from './application/config/react-query';
import router from './application/routes/root';
import GoogleAnalyticsWrapper from './components/GaWrapper';
import useTracking from './hooks/useTracking';
import { defaultLocale, dynamicActivate } from './I18n';

const TOAST_DURATION = 5 * 1000;

const toasterOptions: DefaultToastOptions = {
  duration: TOAST_DURATION,
};

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
      <Toaster toastOptions={toasterOptions} />
    </I18nProvider>
  );
}

export default App;
