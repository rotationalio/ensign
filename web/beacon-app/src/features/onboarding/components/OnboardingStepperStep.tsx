import { Trans } from '@lingui/macro';

interface OnboardingStepperStepProps {
  title: string;
  description: string;
}

const OnboardingStepperStep = ({ title, description }: OnboardingStepperStepProps) => {
  return (
    <>
      <li className="mb-10 ml-6">
        <span className="absolute -left-[4px] mt-1 flex h-2 w-2 items-center justify-center rounded-full bg-gray-100 ring-4 ring-white"></span>
        <h3 className="font-medium leading-tight">
          <Trans>{title}</Trans>
        </h3>
        <button onClick={() => console.log('Step')} className="text-sm">
          <Trans>{description}</Trans>
        </button>
      </li>
    </>
  );
};

export default OnboardingStepperStep;
