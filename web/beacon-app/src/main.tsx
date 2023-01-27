import './index.css';

import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';

import router from '@/application/routes/root';

import App from './App';
import TenantSetup from './components/ui/SetupTenant/TenantSetup';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <App />
    <TenantSetup />
  </React.StrictMode>
);
