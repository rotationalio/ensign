import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { ProjectsResponse } from '../types/projectService';

export function projectsRequest(request: Request): ApiAdapters['getProjectList'] {
  return async (tenantID: string) => {
    const response = (await request(`${APP_ROUTE.TENANTS}/${tenantID}/projects`, {
      method: 'GET',
    })) as any;
    return getValidApiResponse<ProjectsResponse>(response);
  };
}
