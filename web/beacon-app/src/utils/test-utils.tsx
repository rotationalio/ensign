import { i18n } from '@lingui/core';
import { I18nProvider } from '@lingui/react';

export const I18nWrapper = ({ children }: any) => (
  <I18nProvider i18n={i18n}>{children}</I18nProvider>
);
