import { i18n } from '@lingui/core';
import { detect, fromStorage } from '@lingui/detect-locale';
import { I18nProvider } from '@lingui/react';
import {
  createContext,
  Dispatch,
  ReactNode,
  SetStateAction,
  useContext,
  useEffect,
  useState,
} from 'react';

import { LANG_KEY } from '@/constants/lang-key';
import { DEFAULT_LOCALE, dynamicActivate } from '@/I18n';

type State = [string, Dispatch<SetStateAction<string>>];

export const LanguageContext = createContext<State | null>(null);

type LanguageProviderProps = {
  children: ReactNode;
};

const DEFAULT_FALLBACK = () => DEFAULT_LOCALE;
const detectedLanguage = detect(fromStorage(LANG_KEY), DEFAULT_FALLBACK);

const LanguageProvider = ({ children }: LanguageProviderProps) => {
  const [language, setLanguage] = useState<string>(detectedLanguage || DEFAULT_FALLBACK);

  useEffect(() => {
    dynamicActivate(language);
  }, [language]);

  return (
    <I18nProvider i18n={i18n as any}>
      <LanguageContext.Provider value={[language, setLanguage]}>
        {children}
      </LanguageContext.Provider>
    </I18nProvider>
  );
};

const useLanguageProvider = () => {
  const context = useContext(LanguageContext);

  if (!context) {
    throw new Error(`useLanguageProvider should be used within a LanguageProvider`);
  }

  return context;
};

export { LanguageProvider, useLanguageProvider };
