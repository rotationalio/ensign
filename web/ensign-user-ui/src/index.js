import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import * as Sentry from '@sentry/react';
import initSentry from './config/sentry';
import App from './App';
import {store} from './app/store'
import { Provider } from 'react-redux'

initSentry()

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <Sentry.ErrorBoundary fallback={<p>An error has occurred.</p>}>
      <Provider store={store}>
        <App />
      </Provider>
    </Sentry.ErrorBoundary>
  </React.StrictMode>
);
