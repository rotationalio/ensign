import { useMutation } from '@tanstack/react-query';

import { ApiAdapters } from '@/application/api/ApiAdapters';
import axiosInstance, { getValidApiResponse, Request } from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { APP_ROUTE } from '@/constants';

import { QUERY_KEY } from '../constants/query-key';
import { ChangeRoleFormDto } from '../types/changeRoleFormDto';

export function updateMemberRequest(request: Request): ApiAdapters['updateMemberRole'] {
  return async (memberId: string, role: ChangeRoleFormDto['role']) => {
    const response = (await request(`${APP_ROUTE.MEMBERS}/${memberId}`, {
      method: 'POST',
      data: JSON.stringify({ role }),
    })) as any;

    return getValidApiResponse<any>(response);
  };
}

export function useUpdateMemberRole() {
  const { error, isLoading, isError, data, mutate } = useMutation(
    ({ memberID = '', role }: { memberID?: string; role: ChangeRoleFormDto['role'] }) =>
      updateMemberRequest(axiosInstance)(memberID, role),
    {
      onSuccess() {
        queryClient.invalidateQueries([QUERY_KEY.MEMBERS_LIST]);
      },
    }
  );

  return {
    data,
    isUpdatingMemberRole: isLoading,
    hasError: isError,
    error,
    updateMemberRole: mutate,
  };
}
