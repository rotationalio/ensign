import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import type { APIKey } from '../types/ApiKeyServices';

export function apiKeysRequest(request: Request): ApiAdapters['getApiKeys'] {
  return async () => {
    const response = (await request(`${APP_ROUTE.APIKEYS}`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<APIKey>(response);
  };
}
