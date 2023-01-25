import './index.css';

import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';

import router from '@/application/routes/root';

import App from './App';
import OnboardingHeader from './components/Ensign Welcome/OnboardingHeader';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <App />
    <OnboardingHeader />
  </React.StrictMode>
);
