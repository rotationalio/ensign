import './index.css';

import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';

import router from '@/application/routes/root';

import App from './App';
import MaintenaceMode from './components/MaintenanceMode/MaintenaceMode';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <App />
    <MaintenaceMode />
  </React.StrictMode>
);
