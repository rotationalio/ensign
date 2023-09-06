import { AxiosResponse } from 'axios';

import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';
import { MemberResponse } from '@/features/members/types/memberServices';

import { UpdateMemberDTO } from '../types/onboardingServices';

export function updateMemberAPI(request: Request): ApiAdapters['updateOnboardingMember'] {
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
