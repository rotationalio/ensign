import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

export function getInviteTeamMemberRequest(request: Request): ApiAdapters['getInviteTeamMember'] {
  return async (token: string) => {
    const response = (await request(`${APP_ROUTE.INVITE}/${token}`)) as any;

    return getValidApiResponse<any>(response);
  };
}
