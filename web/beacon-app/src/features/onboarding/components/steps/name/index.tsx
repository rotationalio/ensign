import { Trans } from '@lingui/macro';

import { useOrgStore } from '@/store';

import StepCounter from '../StepCounter';
import NameForm from './form';
const NameStep = () => {
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const handleSubmitNameForm = (values: any) => {
    console.log(values);
    increaseStep();
  };
  return (
    <>
      <StepCounter />
      <p className="mt-4 font-bold">
        <Trans>What's your name?</Trans>
      </p>
      <p className="my-4">
        <Trans>
          Adding your name will make it easier for your teammates to collaborate with you.
        </Trans>
      </p>
      <NameForm onSubmit={handleSubmitNameForm} />
    </>
  );
};

export default NameStep;
