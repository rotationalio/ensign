import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';
import type { ProjectDetailDTO, ProjectResponse } from '@/features/projects/types/projectService';

export function projectRequest(request: Request): ApiAdapters['projectDetail'] {
  return async (id: ProjectDetailDTO) => {
    const response = (await request(`${APP_ROUTE.PROJECTS}/${id}`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<ProjectResponse>(response);
  };
}
