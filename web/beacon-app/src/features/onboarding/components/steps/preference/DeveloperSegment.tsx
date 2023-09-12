import { Trans } from '@lingui/macro';
import { memo } from 'react';
const DeveloperSegment = () => {
  return (
    <div className="mt-5">
      <p className="text-base font-bold">
        <Trans>What kind of work do you do?</Trans>
      </p>
      <p className="pt-3 text-base">
        <Trans>Select up to 3. This will help us personalize your experience.</Trans>
      </p>
    </div>
  );
};

export default memo(DeveloperSegment);
