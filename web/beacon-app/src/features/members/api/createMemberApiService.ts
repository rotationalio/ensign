import invariant from 'invariant';

import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';

import { hasMemberRequiredFields, MemberResponse, NewMemberDTO } from '../types/memberServices';

export function createMemberRequest(request: Request): ApiAdapters['createMember'] {
  return async (member: NewMemberDTO) => {
    // check if account has all the required fields defined
    const hasRequiredFields = hasMemberRequiredFields(member);
    invariant(hasRequiredFields, 'Member is missing required fields');

    const response = (await request(`/members`, {
      method: 'POST',
      data: JSON.stringify(member),
    })) as any;
    return getValidApiResponse<MemberResponse>(response);
  };
}
