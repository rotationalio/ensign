import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

export function createTenantRequest(request: Request): ApiAdapters['createTenant'] {
  return async () => {
    const response = (await request(`${APP_ROUTE.TENANTS}`, {
      method: 'POST',
    })) as any;

    return getValidApiResponse<any>(response);
  };
}
