import { formatDate } from '@/utils/formatDate';

export const formatMemberData = (data: any) => {
  if (!data) return [];
  return [
    {
      label: 'Name',
      value: data.name,
    },
    {
      label: 'Email',
      value: data.email,
    },
    {
      label: 'Role',
      value: data.role,
    },
    {
      label: 'Status',
      value: data.status,
    },
    {
      label: 'Last Activity',
      value: formatDate(new Date(data.last_activity)),
    },
    {
      label: 'Joined Date',
      value: formatDate(new Date(data.date_added)),
    },
  ];
};

export const getMembers = (members: any) => {
  if (!members?.members || members?.members.length === 0) return [];
  return Object.keys(members?.members).map((m) => {
    const { name, email, role, status, last_activity, date_added } = members.members[m];
    return { name, email, role, status, last_activity, date_added };
  }) as any;
};
