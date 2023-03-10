// import { Project } from '../types';
import { formatDate } from '@/utils/formatDate';
export const getOrgData = (org: any | undefined) => {
  if (org && org !== null) {
    const { id, name, created, owner } = org;
    return [
      {
        label: 'Name',
        value: name,
      },
      {
        label: 'Org ID',
        value: id,
      },
      {
        label: 'Owner',
        value: owner,
      },
      {
        label: 'Date Created',
        value: formatDate(new Date(created)),
      },
    ];
  }
  return [];
};
