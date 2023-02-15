import { ITenant } from '../types/tenantServices';

export const getRecentTenant = (tenants: ITenant[]) => {
  // for now, just return the first tenant
  if (tenants && tenants.length > 0) {
    const recent = tenants[0];
    const { name, id } = recent;
    return [
      {
        label: 'Tenant Name',
        value: name,
      },
      {
        label: 'Tenant ID',
        value: id,
      },
    ];
  }
  return [];
};
