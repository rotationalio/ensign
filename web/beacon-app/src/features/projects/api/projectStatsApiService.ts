import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { ProjectQuickViewResponse } from '../types/projectService';

function projectStatsApiRequest(request: Request): ApiAdapters['getProjectStats'] {
  return async (tenantID: string) => {
    //console.log('[] tenantId', tenantID);
    const response = (await request(`${APP_ROUTE.TENANTS}/${tenantID}/projects/stats`, {
      method: 'GET',
    })) as any;
    return getValidApiResponse<ProjectQuickViewResponse>(response);
  };
}

export default projectStatsApiRequest;
