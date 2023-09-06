import { Trans } from '@lingui/macro';

import { useUpdateMember } from '@/features/onboarding/hooks/useUpdateMember';
import { useOrgStore } from '@/store';

import StepCounter from '../StepCounter';
import OrganizationForm from './form';

const OrganizationStep = () => {
  const { updateMember } = useUpdateMember();
  const orgDataState = useOrgStore.getState() as any;
  const { user } = orgDataState;
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const handleSubmitOrganizationForm = (values: any) => {
    const payload = {
      memberID: user,
      onboardingPayload: {
        organization: values?.organization,
      },
    };
    updateMember(payload);
    increaseStep();
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
