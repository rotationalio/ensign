import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import React from 'react';

import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
// import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { useOrgStore } from '@/store';

import useHandlePreviousBtn from '../hooks/useHandlePreviousBtn';
import { ONBOARDING_STEPS } from '../shared/constants';
import { isInvitedUser } from '../shared/utils';
type StepButtonsProps = {
  isSubmitting?: boolean;
  isDisabled?: boolean;
  formValues: any;
  hasErrored?: boolean;
};
const StepButtons = ({ isSubmitting, isDisabled, formValues }: StepButtonsProps) => {
  const state = useOrgStore((state: any) => state) as any;
  const { currentStep } = state.onboarding as any;
  const { profile } = useFetchProfile();
  const isInvited = isInvitedUser(profile);
  const shouldDisplayBackButton = currentStep !== ONBOARDING_STEPS.ORGANIZATION;
  const { handlePrevious } = useHandlePreviousBtn();

  const handlePreviousClick = () => {
    handlePrevious({ currentStep, values: formValues, isInvited });
  };
  return (
    <div className="flex flex-row items-stretch gap-3 pt-10">
      <Button
        type="submit"
        isLoading={isSubmitting}
        disabled={isDisabled || isSubmitting}
        data-cy="next-bttn"
      >
        <Trans>Next</Trans>
      </Button>
      {shouldDisplayBackButton && (
        <Button
          onClick={handlePreviousClick}
          isLoading={isSubmitting}
          disabled={isSubmitting}
          variant="ghost"
          data-cy="back-bttn"
          className="hover:border-black-600 hover:text-black-600"
        >
          <Trans>Back</Trans>
        </Button>
      )}
    </div>
  );
};

export default StepButtons;
