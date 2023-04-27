import { i18n } from '@lingui/core';
import { en, fr } from 'make-plural/plurals';

export const locales = {
  en: 'English',
  fr: 'French',
};
export const DEFAULT_LOCALE = 'en';

i18n.loadLocaleData({
  en: { plurals: en },
  fr: { plurals: fr },
});

/**
 * We do a dynamic import of just the catalog that we need
 * @param locale any locale string
 */
export async function dynamicActivate(locale = 'en') {
  const { messages } = await import(`./locales/${locale}/messages.po`);

  i18n.load(locale, messages);
  i18n.activate(locale);
}
