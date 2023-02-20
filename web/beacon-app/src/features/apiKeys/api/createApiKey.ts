import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import type { APIKey } from '../types/apiKeyService';

export function createAPIKey(request: Request): ApiAdapters['createAPIKey'] {
  return async () => {
    const response = (await request(`${APP_ROUTE.APIKEYS}`, {
      method: 'POST',
    })) as any;

    return getValidApiResponse<APIKey>(response);
  };
}
