import { MembersResponse } from './types/memberServices';
export const formatMemberData = (m: MembersResponse) => {
  if (m && m.member.length > 0) {
    const { id, name, role, created } = m.member[0];
    return [
      {
        label: 'id',
        value: id,
      },
      {
        label: 'name',
        value: name,
      },
      {
        label: 'roles',
        value: role,
      },
      {
        label: 'Date Created',
        value: Intl.DateTimeFormat('en-US', { dateStyle: 'full' }).format(new Date(created)),
      },
    ];
  }
  return [];
};
