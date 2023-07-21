import { AxiosResponse } from 'axios';

import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

export function deleteAPIKeyRequest(request: Request): ApiAdapters['deleteAPIKey'] {
  return async (apiKey: string) => {
    const response = (await request(`${APP_ROUTE.APIKEYS}/${apiKey}`, {
      method: 'DELETE',
    })) as unknown as AxiosResponse;

    return getValidApiResponse<any>(response);
  };
}
