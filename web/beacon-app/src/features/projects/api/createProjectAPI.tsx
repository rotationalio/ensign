import { AxiosResponse } from 'axios';

import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { NewProjectDTO } from '../types/createProjectService';
import type { ProjectResponse } from '../types/projectService';

export function createProjectAPI(request: Request): ApiAdapters['createNewProject'] {
  return async ({ tenantID, name, description }: NewProjectDTO) => {
    const response = (await request(`${APP_ROUTE.TENANTS}/${tenantID}/projects`, {
      method: 'POST',
      data: JSON.stringify({
        name,
        description,
      }),
    })) as unknown as AxiosResponse;

    return getValidApiResponse<ProjectResponse>(response);
  };
}
