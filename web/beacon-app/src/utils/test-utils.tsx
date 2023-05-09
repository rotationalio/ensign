import { i18n } from '@lingui/core';
import { I18nProvider } from '@lingui/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, RenderOptions } from '@testing-library/react';
import { ReactElement, ReactNode, useState } from 'react';

import { LanguageContext } from '@/contexts/LanguageContext';

const queryClient = new QueryClient();

type CustomRenderOptions = { locale?: string } & Omit<RenderOptions, 'queries'>;

export const customRender = (
  ui: ReactElement,
  { locale, ...options }: CustomRenderOptions = {}
) => {
  function Wrapper({ children }: { children: ReactNode }) {
    const [language, setLanguage] = useState<string>(locale || 'en');
    return (
      <I18nProvider i18n={i18n as any}>
        <LanguageContext.Provider value={[language, setLanguage]}>
          <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
        </LanguageContext.Provider>
      </I18nProvider>
    );
  }

  return render(ui, { ...options, wrapper: Wrapper });
};
