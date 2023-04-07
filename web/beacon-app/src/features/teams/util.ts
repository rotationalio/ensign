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

export const getMembers = (members: any) => {
  if (!members?.members || members?.members?.length === 0) return [];
  return Object.keys(members?.members).map((m) => {
    const { name, email, role, status, last_activity, date_added } = members.members[m];
    return {
      name,
      email,
      role,
      status,
      last_activity,
      date_added,
      actions: [
        {
          label: 'Change Role',
          onClick: () => alert('not yet implemented'),
        },
        {
          label: 'Remove',
          onClick: () => alert('not yet implemented'),
        },
      ],
    };
  }) as any;
};
