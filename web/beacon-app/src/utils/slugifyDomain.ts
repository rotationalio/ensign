/* eslint-disable prettier/prettier */
/* Creates a URL with an organization domain */
/* Ex. Rotational Labs -> ensign.rotational.io/rotational-labs */
import { slugify as translugify } from 'transliteration';
export function slugify(domain: string, org?: string) {
  const site = 'https://rotational.app';
  if (!org) {
    return `${site}`;
  }
  return `${site}/${stringify_org(org)}/${stringify_org(domain)}`;
}

// sligify organization name to create a URL

export const stringify_org = (input: string) => {
  const string = input
    .normalize('NFKD')
    .toLowerCase()
    .trim()
    // handle russian characters (cyrillic) with transliteration
    .replace(/[\u0400-\u04FF]/g, (x) => translugify(x, { separator: '-' }))
    // handle chinese characters
    .replace(
      /[\u4E00-\u9FCC\u3400-\u4DB5\uFA0E\uFA0F\uFA11\uFA13\uFA14\uFA1F\uFA21\uFA23\uFA24\uFA27-\uFA29]|[\ud840-\ud868][\udc00-\udfff]|\ud869[\udc00-\uded6\udf00-\udfff]|[\ud86a-\ud86c][\udc00-\udfff]|\ud86d[\udc00-\udf34\udf40-\udfff]|\ud86e[\udc00-\udc1d]/g,
      // TODO: add a separator for chinese characters
      (x) => translugify(x, { separator: '-' })
    )
    .replace(/&amp;/g, '-')

    .replace(/-{2,}/g, '-')
    // replace all double or moore dashes with one dash
    .replace(/-+/g, '-')
    .replace(/'/g, '-')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .replace(/-$/, '')
    .replace(/_+/g, '-');

  return (
    string
      .replace(/-/g, '\\-')
      // remove all non-word chars
      .replace(/[^\w-]/g, '')
      .replace(/-+/g, '-')
      .replace(/-+$/, '')
      .replace(/^-/, '')
  );
};
