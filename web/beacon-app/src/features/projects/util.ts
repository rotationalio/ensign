export const formatProjectData = (data: any) => {
  if (!data) return [];
  return [
    {
      label: 'Name',
      value: data?.name,
    },
    {
      label: 'Permissions',
      value: data?.permissions,
    },
    {
      label: 'Permissions',
      value: data?.permissions,
    },
    {
      label: 'Status',
      value: data?.status,
    },
    {
      label: 'Date Created',
      value: data?.created,
    },
  ];
};

export const getApiKeys = (apiKeys: any) => {
  if (!apiKeys?.api_keys || apiKeys?.api_keys.length === 0) return [];
  return Object.keys(apiKeys?.api_keys).map((key) => {
    const { id, name, client_id, permissions, status, last_used, created } = apiKeys.api_keys[key];
    return { id, name, client_id, permissions, status, last_used, created };
  }) as any;
};
