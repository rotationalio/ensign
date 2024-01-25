import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { StatusResponse } from '../types/StatusService';

export function statusRequest(request: Request): ApiAdapters['getStatus'] {
  return async () => {
    const response = (await request(`${APP_ROUTE.STATUS}`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<StatusResponse>(response);
  };
}
