import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { OrgDetailDTO, OrgResponse } from '../types/organizationService';

export function orgRequest(request: Request): ApiAdapters['orgDetail'] {
  return async (id: string) => {
    const response = (await request(`${APP_ROUTE.ORG_DETAIL}/${id}`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<OrgResponse>(response);
  };
}
