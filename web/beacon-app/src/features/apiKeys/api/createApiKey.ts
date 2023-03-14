import { AxiosResponse } from 'axios';

import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import type { APIKey } from '../types/apiKeyService';
import { APIKeyDTO } from '../types/createApiKeyService';

export function createProjectAPIKey(request: Request): ApiAdapters['createProjectAPIKey'] {
  return async ({ projectID, name, permissions }: APIKeyDTO) => {
    const response = (await request(`${APP_ROUTE.PROJECTS}/${projectID}/apikeys`, {
      method: 'POST',
      data: JSON.stringify({
        name,
        permissions,
      }),
    })) as unknown as AxiosResponse;

    return getValidApiResponse<APIKey>(response);
  };
}
