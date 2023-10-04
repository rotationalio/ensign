import { AxiosResponse } from 'axios';

import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { ResetPasswordDTO } from '../types/ResetPasswordService';

export function resetPasswordRequest(request: Request): ApiAdapters['resetPassword'] {
  return async (payload: ResetPasswordDTO) => {
    const response = (await request(`${APP_ROUTE.RESET_PASSWORD}`, {
      method: 'POST',
      data: JSON.stringify(payload),
    })) as unknown as AxiosResponse;

    return getValidApiResponse<any>(response);
  };
}
