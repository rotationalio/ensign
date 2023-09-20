import { Member } from './types/member';

export const formatMemberData = (data: any) => {
  if (!data) return [];
  return [
    {
      label: 'Name',
      value: data?.name,
    },
    {
      label: 'Email',
      value: data?.email,
    },
    {
      label: 'Role',
      value: data?.role,
    },
    {
      label: 'Status',
      value: data?.status,
    },
    {
      label: 'Last Activity',
      value: data?.last_activity,
    },
    {
      label: 'Joined Date',
      value: data?.date_added,
    },
  ];
};

type Actions = {
  handleOpenChangeRoleModal: (member: Member) => void;
  handleOpenDeleteMemberModal: (member: Member) => void;
};

export const getMembers = (members: any, actions?: Actions) => {
  if (!members?.members || members?.members?.length === 0) return [];
  return Object.keys(members?.members).map((m) => {
    const { name, email, role, onboarding_status, last_activity, date_added } = members.members[m];
    return {
      name: name ? name : '-',
      email,
      role,
      status: onboarding_status,
      last_activity,
      date_added,
      actions: [
        {
          label: 'Change Role',
          onClick: () => actions?.handleOpenChangeRoleModal(members.members[m]),
        },
        // {
        //   label: 'Remove',
        //   onClick: () => actions?.handleOpenDeleteMemberModal(members.members[m]),
        // },
      ],
    };
  }) as any;
};
