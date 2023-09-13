import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { MembersResponse } from '../types/memberServices';

export function profileRequest(request: Request): ApiAdapters['getProfile'] {
  return async () => {
    const response = (await request(`${APP_ROUTE.PROFILE}`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<MembersResponse>(response);
  };
}
