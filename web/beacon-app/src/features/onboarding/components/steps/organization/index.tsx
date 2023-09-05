import { Trans } from '@lingui/macro';

import { useUpdateMember } from '@/features/onboarding/hooks/useUpdateMember';
import { useOrgStore } from '@/store';

import StepCounter from '../StepCounter';
import OrganizationForm from './form';

const OrganizationStep = () => {
  const { updateMember } = useUpdateMember();
  // Get the member ID from the store
  const { user } = useOrgStore.getState() as any;
  const handleSubmitOrganizationForm = (values: any) => {
    const payload = {
      memberID: user,
      onboardingPayload: {
        organization: values?.organization,
      },
    };
    // console.log('payload', payload);
    updateMember(payload);
  };
  return (
    <>
      <StepCounter />
      <p className="mt-4 font-bold">
        <Trans>What's the name of your team or organization?</Trans>
      </p>
      <p className="my-4">
        <Trans>
          This will be the name of your workspace where you create projects and collaborate, so
          choose something you and your teammates will recognize.
        </Trans>
      </p>
      <OrganizationForm onSubmit={handleSubmitOrganizationForm} />
    </>
  );
};

export default OrganizationStep;
