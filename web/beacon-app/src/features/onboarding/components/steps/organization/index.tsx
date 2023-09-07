import { Trans } from '@lingui/macro';

import useUserLoader from '@/features/members/loaders/userLoader';
import { useOrgStore } from '@/store';

import StepCounter from '../StepCounter';
import OrganizationForm from './form';

const OrganizationStep = () => {
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const { member } = useUserLoader();

  const submitFormHandler = (values: any) => {
    console.log(values);
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
      <OrganizationForm
        onSubmit={submitFormHandler}
        initialValues={{
          organization: member?.organization,
        }}
      />
    </>
  );
};

export default OrganizationStep;
