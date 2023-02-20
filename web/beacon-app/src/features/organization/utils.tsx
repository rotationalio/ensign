// import { Project } from '../types';

export const getOrgData = (org: any | undefined) => {
  if (org && org !== null) {
    const { id, name, domain, created } = org;
    return [
      {
        label: 'Name',
        value: name,
      },
      {
        label: 'URL',
        value: domain,
      },
      {
        label: 'Org ID',
        value: id,
      },
      {
        label: 'Owner',
        value: 'owner',
      },
      {
        label: 'Date Created',
        value: Intl.DateTimeFormat('en-US', { dateStyle: 'full' }).format(new Date(created)),
      },
    ];
  }
  return [];
};
