import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import type { User } from '../types/RegisterService';

export function createAccountRequest(request: Request): ApiAdapters['createNewAccount'] {
  return async (account) => {
    const response = (await request(`${APP_ROUTE.REGISTER}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(account),
    })) as any;

    return getValidApiResponse<User>(response);
  };
}
