export const formatMemberData = (m: any) => {
  console.log('[formatMemberData] m: ', m);
  if (m) {
    const { id, name, role, created } = m;
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
