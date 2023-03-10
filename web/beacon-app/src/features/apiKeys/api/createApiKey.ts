import { AxiosResponse } from 'axios';

import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import type { APIKey } from '../types/apiKeyService';

export function createProjectAPIKey(request: Request): ApiAdapters['createProjectAPIKey'] {
  return async (projectID: string) => {
    const response = (await request(`${APP_ROUTE.PROJECTS}/${projectID}/apikeys`, {
      method: 'POST',
    })) as unknown as AxiosResponse;

    return getValidApiResponse<APIKey>(response);
  };
}
