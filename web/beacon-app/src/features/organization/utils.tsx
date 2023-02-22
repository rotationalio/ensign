// import { Project } from '../types';

import { t } from '@lingui/macro';

export const getOrgData = (org: any | undefined) => {
  if (org && org !== null) {
    const { id, name, domain, created } = org;
    return [
      {
        label: t`Name`,
        value: name,
      },
      {
        label: t`URL`,
        value: domain,
      },
      {
        label: t`Org ID`,
        value: id,
      },
      {
        label: t`Owner`,
        value: 'owner',
      },
      {
        label: t`Date Created`,
        value: Intl.DateTimeFormat('en-US', { dateStyle: 'full' }).format(new Date(created)),
      },
    ];
  }
  return [];
};
