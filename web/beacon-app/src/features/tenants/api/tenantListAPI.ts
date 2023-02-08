import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import type { UserTenantResponse } from '../types/tenantServices';

export function tenantRequest(request: Request): ApiAdapters['getTenantList'] {
  return async () => {
    const response = (await request(`${APP_ROUTE.TENANTS}`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<UserTenantResponse>(response);
  };
}
