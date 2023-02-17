import './index.css';

import React from 'react';
import ReactDOM from 'react-dom/client';

import App from './App';
import { appConfig } from './application/config';
import initSentry from './application/config/sentry';

// eslint-disable-next-line no-console
console.info('initializing beacon ui', appConfig.nodeENV, appConfig.version, appConfig.revision);
initSentry();

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
