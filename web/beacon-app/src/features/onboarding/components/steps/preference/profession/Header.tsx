import { Trans } from '@lingui/macro';
import { memo } from 'react';
const ProfessionSegmentHeader = () => {
  return (
    <>
      <p className="text-base font-bold">
        <Trans>What will you use Ensign for?</Trans>
      </p>
      <p className="pt-3 text-base">
        <Trans>This will help us plan new features and improvements.</Trans>
      </p>
    </>
  );
};

export default memo(ProfessionSegmentHeader);
