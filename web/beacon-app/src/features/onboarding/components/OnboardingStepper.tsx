import OnboardingStepperStep from './OnboardingStepperStep';

const OnboardingStepper = () => {
  return (
    <>
      <ol className="relative border-l border-gray-200 text-white">
        <OnboardingStepperStep title="Step 1 of 4" description="Your Team Name" />
        <OnboardingStepperStep title="Step 2 of 4" description="Your Workspace URL" />
        <OnboardingStepperStep title="Step 3 of 4" description="Your Name" />
        <OnboardingStepperStep title="Step 4 of 4" description="Your Goals" />
      </ol>
    </>
  );
};

export default OnboardingStepper;
