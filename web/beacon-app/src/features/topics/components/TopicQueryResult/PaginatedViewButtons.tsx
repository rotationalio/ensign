import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { ImNext2, ImPrevious2 } from 'react-icons/im';

interface PaginatedViewButtonsProps {
  onClickNext: () => void;
  onClickPrevious: () => void;
  isNextDisabled?: boolean;
  isPreviousDisabled?: boolean;
}

const PaginatedViewButtons: React.FC<PaginatedViewButtonsProps> = ({
  onClickNext,
  onClickPrevious,
  isNextDisabled = true,
  isPreviousDisabled = true,
}) => {
  return (
    <div className="mt-2 flex justify-between">
      <Button
        type="button"
        onClick={onClickPrevious}
        disabled={isPreviousDisabled}
        data-testid="prev-query-btn"
        data-cy="prev-query-btn"
      >
        <div className="flex justify-center">
          <ImPrevious2 fontSize={20} />
          <Trans>Previous</Trans>
        </div>
      </Button>
      <Button
        type="button"
        onClick={onClickNext}
        disabled={isNextDisabled}
        data-testid="next-query-btn"
        data-cy="next-query-btn"
      >
        <div className="flex justify-center">
          <Trans>Next</Trans>
          <ImNext2 fontSize={20} />
        </div>
      </Button>
    </div>
  );
};

export default PaginatedViewButtons;
