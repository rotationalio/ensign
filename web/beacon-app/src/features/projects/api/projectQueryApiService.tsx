import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { ProjectQueryDTO } from '../types/projectQueryService';
// TODO: add projectQueryresponse type

export function projectQueryAPI(request: Request): ApiAdapters['projectQuery'] {
  return async ({ projectID, query }: ProjectQueryDTO) => {
    const response = (await request(`${APP_ROUTE.PROJECTS}/${projectID}/query`, {
      method: 'POST',
      data: JSON.stringify({
        query,
      }),
    })) as unknown as any;

    return getValidApiResponse<any>(response);
  };
}
