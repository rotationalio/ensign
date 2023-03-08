import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

export function checkVerifyTokenRequest(request: Request): ApiAdapters['checkToken'] {
  return async (token: string) => {
    const response = (await request(`${APP_ROUTE.TOKEN}`, {
      method: 'POST',
      data: JSON.stringify(token),
    })) as any;
    return getValidApiResponse<any>(response);
  };
}
