import { t } from '@lingui/macro';

import userLoader from '@/features/members/loaders/userLoader';

import { WORKSPACE_DOMAIN_BASE } from '../../shared/constants';
import { stepperContents } from '../../shared/utils';
import StepperStep from './StepperStep';

const Stepper = () => {
  const { member } = userLoader();
  const isInvitedUser = member?.invited;

  return (
    <>
      <ol className="stepper-items relative border-l border-gray-200 text-white">
        {stepperContents.map((step, idx) => (
          <StepperStep
            title={step.title}
            description={
              isInvitedUser && idx === 0
                ? t`Organization`
                : isInvitedUser && idx === 1
                ? t`Workspace URL`
                : step.description
            }
            value={
              isInvitedUser && idx === 0
                ? member?.organization
                : isInvitedUser && idx === 1
                ? `${WORKSPACE_DOMAIN_BASE}${member?.workspace}`
                : undefined
            }
            key={idx}
          />
        ))}
      </ol>
    </>
  );
};

export default Stepper;
