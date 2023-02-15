import { i18n } from '@lingui/core';
import { I18nProvider } from '@lingui/react';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { useEffect } from 'react';
import { RouterProvider } from 'react-router-dom';

import { queryClient, QueryClientProvider } from '@/application/config/react-query';
import router from '@/application/routes/root';

import { defaultLocale, dynamicActivate } from './application/config/i18n';

function App() {
  useEffect(() => {
    dynamicActivate(defaultLocale);
  }, []);

  return (
    <I18nProvider i18n={i18n}>
      <QueryClientProvider client={queryClient}>
        <ReactQueryDevtools initialIsOpen={false} />
        <RouterProvider router={router} />
      </QueryClientProvider>
    </I18nProvider>
  );
}

export default App;
