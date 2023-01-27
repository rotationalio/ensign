import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { UserAuthResponse } from '../types/LoginService';
import { APP_ROUTE } from '@/constants';

export function loginRequest(request: Request): ApiAdapters['authenticateUser'] {
    return async (user) => {
        const response = await request(`${APP_ROUTE.LOGIN}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(user),
        }) as any;

        return getValidApiResponse<UserAuthResponse>(response);
    };
}

