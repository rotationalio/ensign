import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';
import { UserAuthResponse } from '@/features/auth';

export function getInviteAuthenticationRequest(
  request: Request
): ApiAdapters['getInviteAuthenticationRequest'] {
  return async (token: string) => {
    const response = (await request(`${APP_ROUTE.INVITE}/accept`, {
      method: 'POST',
      data: {
        token,
      },
    })) as any;

    return getValidApiResponse<UserAuthResponse>(response);
  };
}
