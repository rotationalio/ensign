import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { MembersResponse } from '../types/memberServices';

export function memberRequest(request: Request): ApiAdapters['getMemberDetail'] {
  return async (memberID: string) => {
    const response = (await request(`${APP_ROUTE.MEMBERS}/${memberID}`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<MembersResponse>(response);
  };
}
