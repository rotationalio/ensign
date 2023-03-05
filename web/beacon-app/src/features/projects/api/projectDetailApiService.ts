import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';
import type { ProjectResponse } from '@/features/projects/types/projectService';

export function projectRequest(request: Request): ApiAdapters['projectDetail'] {
  return async (projectID: string) => {
    console.log('typeof projectID', typeof projectID);
    const response = (await request(`${APP_ROUTE.PROJECTS}/${projectID}`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<ProjectResponse>(response);
  };
}
