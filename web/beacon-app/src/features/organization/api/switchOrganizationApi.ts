import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';
import type { UserAuthResponse } from '@/features/auth/types/LoginService';

export function switchOrganizationRequest(request: Request): ApiAdapters['switchOrganization'] {
  return async (orgId) => {
    const response = (await request(`${APP_ROUTE.SWITCH}`, {
      method: 'POST',
      data: {
        org_id: orgId,
      },
    })) as any;

    return getValidApiResponse<UserAuthResponse>(response);
  };
}
