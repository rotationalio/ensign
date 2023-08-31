/* eslint-disable prettier/prettier */
/* Creates a URL with an organization domain */
/* Ex. Rotational Labs -> ensign.rotational.io/rotational-labs */
import { slugify as Slugify } from 'transliteration';
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
    // replace all spaces with -
    .replace(/\s+/g, '-')
    // remplace all ' found in the string with -
    .replace(/'/g, '-')
    // handle chinese characters
    .replace(/[\u4E00-\u9FCC\u3400]/g, (a) => {
      console.log('a repeat: ', a);
      return Slugify(a);
    })

    // remove &amp; and replace with - and remove all other special characters
    .replace(/&amp;/g, '-')
    .replace(/[^A-Za-z0-9\s]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    //remove - at the end of the string
    .replace(/-$/, '');

  return string;
};
