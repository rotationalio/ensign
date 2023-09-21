import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
import { FC } from 'react';

type OnBoardingHeaderProps = {
  data: any;
};

const renderInvitedMessage = (profile: any) => {
  return (
    <>
      Thank you for accepting `${profile.name}`'s invitation to join the workspace for `$
      {profile.organization}`. Please complete our brief onboarding survey to get started.
    </>
  );
};

const renderNewUserMessage = () => {
  return <>Welcome to Beacon! Please complete our brief onboarding survey to get started.</>;
};

const renderUserWithOrganizationMessage = (profile: any) => {
  return (
    <>
      Welcome to the workspace for `${profile.organization}` on Ensign!. Please complete our brief
      onboarding survey to get started.
    </>
  );
};

const OnBoardingHeader: FC<OnBoardingHeaderProps> = ({ data }) => {
  return (
    <Heading as="h1" className=" m-10 mt-20 px-4  text-xl font-bold xl:mt-20 xl:px-28">
      <Trans>
        {data?.profile?.invited
          ? renderInvitedMessage(data)
          : data?.profile?.organization
          ? renderUserWithOrganizationMessage(data)
          : renderNewUserMessage()}
      </Trans>
    </Heading>
  );
};

export default OnBoardingHeader;
