import invariant from 'invariant';

import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';

export function deleteMemberRequest(request: Request): ApiAdapters['deleteMember'] {
  return async (memberId: string) => {
    invariant(memberId, 'member id is missing');
    const response = (await request(`/members/${memberId}`, {
      method: 'DELETE',
    })) as any;
    return getValidApiResponse<any>(response);
  };
}
