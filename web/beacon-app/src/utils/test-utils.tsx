/* eslint-disable import/export */
import { i18n } from '@lingui/core';
import { I18nProvider } from '@lingui/react';
import { render as rtlRender, RenderOptions } from '@testing-library/react';
import { ReactNode } from 'react';
import { BrowserRouter } from 'react-router-dom';

const AllTheProviders = ({ children }: { children: ReactNode }) => {
  return (
    <I18nProvider i18n={i18n}>
      <BrowserRouter>{children}</BrowserRouter>
    </I18nProvider>
  );
};

export const render = (ui: any, options: RenderOptions = {}) =>
  rtlRender(ui, { wrapper: AllTheProviders, ...options });

export * from '@testing-library/react';
