import { formatDate } from '@/utils/formatDate';
export const formatMemberData = (m: any) => {
  if (m) {
    const { id, name, role, created } = m;
    return [
      {
        label: 'User ID',
        value: id,
      },
      {
        label: 'Name',
        value: name,
      },
      {
        label: 'Role',
        value: role,
      },
      {
        label: 'Date Created',
        value: formatDate(new Date(created)),
      },
    ];
  }
  return [];
};
