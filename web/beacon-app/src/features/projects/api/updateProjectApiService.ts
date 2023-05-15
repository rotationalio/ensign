import { AxiosResponse } from 'axios';

import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import type { ProjectResponse } from '../types/projectService';
import { UpdateProjectDTO } from '../types/updateProjectService';

export function updateProjectAPI(request: Request): ApiAdapters['updateProject'] {
  return async ({ projectID, projectPayload }: UpdateProjectDTO) => {
    const response = (await request(`${APP_ROUTE.PROJECTS}/${projectID}`, {
      method: 'PATCH',
      data: JSON.stringify({
        ...projectPayload,
      }),
    })) as unknown as AxiosResponse;

    return getValidApiResponse<ProjectResponse>(response);
  };
}
