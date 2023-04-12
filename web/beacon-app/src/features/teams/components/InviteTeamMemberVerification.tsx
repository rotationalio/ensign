import { useLoaderData } from 'react-router-dom';

import ExistingUserInvitationPage from './ExistingUserInvitationPage';
import NewUserInvitationPage from './NewUserInvitationPage';

export const InviteTeamMemberVerification = () => {
  const loaderData = useLoaderData() as any;

  if (loaderData && loaderData.has_account) {
    return <ExistingUserInvitationPage data={loaderData} />;
  }
  return <NewUserInvitationPage data={loaderData} />;
};
