import { LoaderFunctionArgs } from 'react-router-dom';

import axiosInstance from '@/application/api/ApiService';

import { getInviteTeamMemberRequest } from '../api/getInviteTeamMemberRequest';

export const inviteTeamMemberLoader = async ({ request }: LoaderFunctionArgs) => {
  const url = new URL(request.url);
  const token = url.searchParams.get('token') || '';

  const response = await getInviteTeamMemberRequest(axiosInstance)(token);
  return response;
};
