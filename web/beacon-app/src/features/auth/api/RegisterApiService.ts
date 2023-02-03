import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';
import invariant from 'invariant';
import type { NewUserResponseData, NewUserAccount } from '../types/RegisterService';
import { hasUserRequiredFields } from '../types/RegisterService';

export function createAccountRequest(request: Request): ApiAdapters['createNewAccount'] {
  return async (account: NewUserAccount) => {
    // check if account has all the required fields defined
    const hasRequiredFields = hasUserRequiredFields(account);
    invariant(hasRequiredFields, 'Account is missing required fields');

    const response = (await request(`${APP_ROUTE.REGISTER}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(account),
    })) as any;

    return getValidApiResponse<NewUserResponseData>(response);
  };
}
