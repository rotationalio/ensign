import { Trans } from '@lingui/macro';

import StepCounter from '../StepCounter';
import NameForm from './form';

const NameStep = () => {
  const handleSubmitNameForm = (values: any) => {
    console.log(values);
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
