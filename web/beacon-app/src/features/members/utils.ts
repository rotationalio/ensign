import { formatDate } from '@/utils/formatDate';

import { MemberStatusEnum } from '../teams/types/member';
export const formatMemberData = (m: any) => {
  //console.log('[formatMemberData] m: ', m);
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

export const isOnboardedMember = (status: string) => {
  return status === MemberStatusEnum.CONFIRMED;
};
