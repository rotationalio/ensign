import './index.css';

import React from 'react';
import ReactDOM from 'react-dom/client';

import App from './App';
import Footer from './components/Footer/Footer';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <App />
    <Footer />
  </React.StrictMode>
);
