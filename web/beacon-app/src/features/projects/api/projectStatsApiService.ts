import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { ProjectsResponse } from '../types/projectService';

export function projectStatsApiRequest(request: Request): ApiAdapters['getProjectList'] {
  return async (tenantID: string) => {
    console.log('[] tenantId', tenantID);
    const response = (await request(`${APP_ROUTE.TENANTS}/${tenantID}/projects/stats`, {
      method: 'GET',
    })) as any;
    return getValidApiResponse<ProjectsResponse>(response);
  };
}
