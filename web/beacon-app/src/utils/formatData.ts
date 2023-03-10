import { UserTenantResponse } from '@/features/tenants/types/tenantServices';
export const getRecentTenant = (tenants: UserTenantResponse) => {
  if (!tenants || tenants.tenants.length === 0) {
    return null;
  }
  // get the recent tenant by modified date
  const recentTenant = tenants?.tenants.reduce((prev: any, current: any) => {
    return prev.modified > current.modified ? prev : current;
  });
  return recentTenant;
};
