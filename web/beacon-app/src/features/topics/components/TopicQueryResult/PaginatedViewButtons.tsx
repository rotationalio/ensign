import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { ImNext2, ImPrevious2 } from 'react-icons/im';

const PaginatedViewButtons = () => {
  return (
    <div className="mt-2 flex justify-between">
      <Button type="button">
        <div className="flex justify-center">
          <ImPrevious2 fontSize={20} />
          <Trans>Previous</Trans>
        </div>
      </Button>
      <Button type="button">
        <div className="flex justify-center">
          <Trans>Next</Trans>
          <ImNext2 fontSize={20} />
        </div>
      </Button>
    </div>
  );
};

export default PaginatedViewButtons;
