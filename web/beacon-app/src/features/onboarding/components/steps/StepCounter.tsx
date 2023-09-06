import { Heading } from '@rotational/beacon-core';

import { useOrgStore } from '@/store';
// TODO: dynamic step counter based on the current step
interface StepCounterProps {
  totalSteps?: number;
}
const StepCounter = ({ totalSteps = 4 }: StepCounterProps) => {
  const orgDataState = useOrgStore.getState() as any;

  const { currentStep } = orgDataState.onboarding;

  console.log('currentStep', currentStep);

  return (
    <div>
      {currentStep && (
        <Heading as="h1" className="mb-2 space-y-3 text-xl font-bold">
          Step {currentStep} of {totalSteps}
        </Heading>
      )}
    </div>
  );
};

export default StepCounter;
