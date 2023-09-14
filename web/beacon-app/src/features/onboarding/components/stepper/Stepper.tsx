import { t } from '@lingui/macro';

import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';

import { WORKSPACE_DOMAIN_BASE } from '../../shared/constants';
import { stepperContents } from '../../shared/utils';
import StepperStep from './StepperStep';

const Stepper = () => {
  const { profile: userProfile } = useFetchProfile();
  const isInvitedUser = userProfile?.invited;

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
                ? userProfile?.organization
                : isInvitedUser && idx === 1
                ? `${WORKSPACE_DOMAIN_BASE}${userProfile?.workspace}`
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
