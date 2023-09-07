import { t } from '@lingui/macro';

import { userLoader } from '@/features/members/loaders';

import Indicator from './Indicator';
interface StepperProps {
  title: string;
  description: string;
}

const stepperContents = [
  {
    title: t`Step 1 of 4`,
    description: t`Your Team Name`,
  },
  {
    title: t`Step 2 of 4`,
    description: t`Your Workspace URL`,
  },
  {
    title: t`Step 3 of 4`,
    description: t`Your Name`,
  },
  {
    title: t`Step 4 of 4`,
    description: t`Your Goals`,
  },
];

const Step = ({ title, description }: StepperProps) => {
  return (
    <>
      <li className="mb-10 ml-6">
        <Indicator />
        <h3 className="font-medium leading-tight">{title}</h3>
        <button onClick={() => console.log('Step')} className="text-sm">
          {description}
        </button>
      </li>
    </>
  );
};

const Stepper = () => {
  const { member } = userLoader();
  const isInvitedUser = member?.invited;

  return (
    <>
      <ol className="stepper-items relative border-l border-gray-200 text-white">
        {stepperContents.map((step, idx) => (
          <Step
            title={step.title}
            description={
              isInvitedUser && idx === 0
                ? member?.organization
                : isInvitedUser && idx === 1
                ? member?.workspace
                : step.description
            }
            key={idx}
          />
        ))}
      </ol>
    </>
  );
};

export default Stepper;
