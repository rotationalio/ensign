import React from 'react';
import ReactDOM from 'react-dom/client';
import './input.css';
import * as Sentry from '@sentry/react';
import initSentry from './config/sentry';
import App from './App';

initSentry();

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <Sentry.ErrorBoundary fallback={<p>An error has occured</p>}>
    <App />
    </Sentry.ErrorBoundary>
  </React.StrictMode>
);