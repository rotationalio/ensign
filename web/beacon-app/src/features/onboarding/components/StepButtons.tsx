import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import React from 'react';

import { useOrgStore } from '@/store';

import { ONBOARDING_STEPS } from '../shared/constants';
type StepButtonsProps = {
  isSubmitting?: boolean;
  isDisabled?: boolean;
};
const StepButtons = ({ isSubmitting, isDisabled }: StepButtonsProps) => {
  const state = useOrgStore((state: any) => state) as any;
  const { currentStep } = state.onboarding as any;
  const handlePreviousClick = () => {
    if (!currentStep || currentStep === ONBOARDING_STEPS.ORGANIZATION) return;
    state.decrementStep();
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
      {currentStep !== ONBOARDING_STEPS.ORGANIZATION && (
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
