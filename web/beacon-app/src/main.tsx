import './index.css';

import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';

import { queryClient, QueryClientProvider } from '@/application/config/react-query';
import router from '@/application/routes/root';

import { appConfig } from './application/config';
import initSentry from './application/config/sentry';

// eslint-disable-next-line no-console
console.info('initializing beacon ui', appConfig.nodeENV, appConfig.version, appConfig.revision);
initSentry();

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <ReactQueryDevtools initialIsOpen={false} />
      <RouterProvider router={router} />
    </QueryClientProvider>
  </React.StrictMode>
);
