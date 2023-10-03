import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

export function forgotPasswordRequest(request: Request): ApiAdapters['forgotPassword'] {
  return async (email) => {
    const response = (await request(`${APP_ROUTE.FORGOT_PASSWORD}`, {
      method: 'POST',
      data: JSON.stringify(email),
    })) as any;

    return getValidApiResponse<any>(response);
  };
}
