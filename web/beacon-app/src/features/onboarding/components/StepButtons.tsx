import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import React from 'react';

import { userLoader } from '@/features/members/loaders';
import { useOrgStore } from '@/store';

import { ONBOARDING_STEPS } from '../shared/constants';
import { isInvitedUser } from '../shared/utils';
type StepButtonsProps = {
  isSubmitting?: boolean;
  isDisabled?: boolean;
};
const StepButtons = ({ isSubmitting, isDisabled }: StepButtonsProps) => {
  const state = useOrgStore((state: any) => state) as any;
  const { currentStep } = state.onboarding as any;
  const { member } = userLoader();
  const isInvited = isInvitedUser(member);
  const shouldDisplayBackButton =
    currentStep !== ONBOARDING_STEPS.ORGANIZATION ||
    (currentStep > 2 && isInvited && currentStep !== ONBOARDING_STEPS.NAME); // don't show back button on name step for invited users
  console.log(isInvited);
  const handlePreviousClick = () => {
    if (!currentStep || currentStep === ONBOARDING_STEPS.ORGANIZATION) return;
    state.decrementStep();
  };
  return (
    <div className="flex flex-row items-stretch gap-3 pt-10">
      <Button type="submit" isLoading={isSubmitting} disabled={isDisabled || isSubmitting}>
        <Trans>Next</Trans>
      </Button>
      {shouldDisplayBackButton && (
        <Button
          onClick={handlePreviousClick}
          isLoading={isSubmitting}
          disabled={isSubmitting}
          variant="ghost"
          className="hover:border-black-600 hover:text-black-600"
        >
          <Trans>Back</Trans>
        </Button>
      )}
    </div>
  );
};

export default StepButtons;
