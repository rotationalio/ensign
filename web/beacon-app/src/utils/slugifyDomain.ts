/* eslint-disable prettier/prettier */
/* Creates a URL with an organization domain */
/* Ex. Rotational Labs -> ensign.rotational.io/rotational-labs */

export function slugify(domain: string, org?: string) {
  const site = 'https://rotational.app';
  if (!org) {
    return `${site}`;
  }
  return `${site}/${stringify_org(org)}/${stringify_org(domain)}`;
}

// sligify organization name to create a URL

export const stringify_org = (org: string) => {
  return (
    org &&
    org
      .normalize('NFKD')
      .replace(/[\u0300-\u036f]/g, '')
      .toLowerCase()
      .trim()
      .replace(/\s+/g, '-')
      // eslint-disable-next-line no-useless-escape
      .replace(/[^\w\-]+/g, '')
      // eslint-disable-next-line no-useless-escape
      .replace(/\-\-+/g, '-')
  );
};
