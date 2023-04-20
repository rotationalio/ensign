import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { OrgListResponse } from '../types/organizationService';

export function organizationRequest(request: Request): ApiAdapters['getOrganizationList'] {
  return async () => {
    const response = (await request(`${APP_ROUTE.ORGANIZATION}`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<OrgListResponse>(response);
  };
}
