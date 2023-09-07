import { AxiosResponse } from 'axios';

import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { MemberResponse, UpdateMemberDTO } from '../types/memberServices';

export function updateMemberAPI(request: Request): ApiAdapters['updateMember'] {
  return async ({ memberID, payload }: UpdateMemberDTO) => {
    const response = (await request(`${APP_ROUTE.MEMBERS}/${memberID}`, {
      method: 'PUT',
      data: JSON.stringify({
        ...payload,
      }),
    })) as unknown as AxiosResponse;
    return getValidApiResponse<MemberResponse>(response);
  };
}
