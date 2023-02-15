import { i18n } from '@lingui/core';
import { en, fr } from 'make-plural/plurals';

export const locales = {
  en: 'English',
  fr: 'French',
};

type Locale = keyof typeof locales;

export const defaultLocale = 'en';

i18n.loadLocaleData({
  en: { plurals: en },
  fr: { plurals: fr },
});

/**
 * We do a dynamic import of just the catalog that we need
 * @param locale any locale string
 */
export async function dynamicActivate(locale: Locale) {
  const importedMessages = await import(`../../locales/${locale}/messages.js`);
  const messages = importedMessages.default.messages;

  i18n.load(locale, messages);
  i18n.activate(locale);
}
