import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import React from 'react';

type StepButtonsProps = {
  isSubmitting?: boolean;
  isDisabled?: boolean;
};
const StepButtons = ({ isSubmitting, isDisabled }: StepButtonsProps) => {
  const handlePreviousClick = () => {
    console.log('previous clicked');
  };
  return (
    <div className="flex flex-row items-stretch gap-3 pt-10">
      <Button type="submit" isLoading={isSubmitting} disabled={isDisabled || isSubmitting}>
        <Trans>Next</Trans>
      </Button>
      <Button
        onClick={handlePreviousClick}
        isLoading={isSubmitting}
        disabled={isSubmitting}
        variant="ghost"
        className="hover:border-black-600 hover:text-black-600"
      >
        <Trans>Back</Trans>
      </Button>
    </div>
  );
};

export default StepButtons;
