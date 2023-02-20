import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { useEffect } from 'react';
import { RouterProvider } from 'react-router-dom';

import { queryClient } from './application/config/react-query';
import router from './application/routes/root';
//import { defaultLocale, dynamicActivate } from './I18n';

function App() {
  useEffect(() => {
    //dynamicActivate(defaultLocale);
  }, []);

  return (
    //<I18nProvider i18n={i18n}>
    <QueryClientProvider client={queryClient}>
      <ReactQueryDevtools initialIsOpen={false} />
      <RouterProvider router={router} />
    </QueryClientProvider>
    //</I18nProvider>
  );
}

export default App;
