import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
import { FC } from 'react';

type OnBoardingHeaderProps = {
  data: any;
};

const renderInvitedMessage = (profile: any) => {
  return (
    <>
      Thank you for accepting the invitation to join the workspace for{' '}
      <span className="font-bold"> {profile.organization}</span>. Please complete our brief
      onboarding survey to get started.
    </>
  );
};

const renderNewUserMessage = () => {
  return <>Welcome to Ensign! Please complete our brief onboarding survey to get started.</>;
};

const renderUserWithOrganizationMessage = (profile: any) => {
  return (
    <>
      Welcome to the workspace for <span className="font-bold"> {profile.organization} </span> on
      Ensign! Please complete our brief onboarding survey to get started.
    </>
  );
};

const OnBoardingHeader: FC<OnBoardingHeaderProps> = ({ data }) => {
  return (
    <Heading as="h1" className=" mt-20 px-4 text-xl xl:ml-12 xl:mt-20 xl:px-28">
      <Trans>
        {data?.invited
          ? renderInvitedMessage(data)
          : data?.organization && !data?.invited
          ? renderUserWithOrganizationMessage(data)
          : renderNewUserMessage()}
      </Trans>
    </Heading>
  );
};

export default OnBoardingHeader;
