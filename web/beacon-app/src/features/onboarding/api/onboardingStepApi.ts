import { AxiosResponse } from 'axios';

import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';
import { MemberResponse } from '@/features/members/types/memberServices';

import { UpdateMemberOnboardingDTO } from '../types/onboardingServices';

export function onboardingStepAPI(request: Request): ApiAdapters['updateOnboardingMember'] {
  return async ({ memberID, onboardingPayload }: UpdateMemberOnboardingDTO) => {
    const response = (await request(`${APP_ROUTE.MEMBERS}/${memberID}`, {
      method: 'PUT',
      data: JSON.stringify({
        ...onboardingPayload,
      }),
    })) as unknown as AxiosResponse;
    return getValidApiResponse<MemberResponse>(response);
  };
}
