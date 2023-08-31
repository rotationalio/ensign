/* eslint-disable prettier/prettier */

import { slugify as Slugify } from 'transliteration';
export function slugify(domain: string, org?: string) {
  const site = 'https://rotational.app';
  if (!org) {
    return `${site}`;
  }
  return `${site}/${stringify_org(org)}/${stringify_org(domain)}`;
}

export const stringify_org = (input: string) => {
  const string = input
    .normalize('NFKD')
    .toLowerCase()
    .trim()
    .replace(/\s+/g, '-')

    .replace(/'/g, '-')
    // handle chinese characters
    .replace(/[\u4E00-\u9FCC\u3400]/g, (a) => {
      console.log('a repeat: ', a);
      return Slugify(a);
    })

    .replace(/&amp;/g, '-')
    .replace(/[^A-Za-z0-9\s]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')

    .replace(/-$/, '');

  return string;
};
