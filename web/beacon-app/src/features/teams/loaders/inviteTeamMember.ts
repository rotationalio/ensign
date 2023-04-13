import { LoaderFunctionArgs } from 'react-router-dom';

import axiosInstance from '@/application/api/ApiService';
import { setCookie } from '@/utils/cookies';

import { getInviteTeamMemberRequest } from '../api/getInviteTeamMemberRequest';

export const inviteTeamMemberLoader = async ({ request }: LoaderFunctionArgs) => {
  const url = new URL(request.url);
  const token = url.searchParams.get('token') || '';

  if (token) {
    setCookie('invite_token', token);
  }

  const response = await getInviteTeamMemberRequest(axiosInstance)(token);
  return response;
};
