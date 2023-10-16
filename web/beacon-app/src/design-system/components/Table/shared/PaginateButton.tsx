import { mergeClassnames } from '../../../utils';
import { Button } from '../../Button';
interface PaginateButtonProps {
  className?: string;
  onClick?: () => void;
  disabled?: boolean;
  children?: React.ReactNode;
}

export const PaginateButton = ({ className, onClick, disabled, children }: PaginateButtonProps) => {
  return (
    <Button
      className={mergeClassnames('mx-2 flex h-8 w-24 items-center justify-center', className)}
      onClick={onClick}
      disabled={disabled}
      variant="ghost"
      size="custom"
    >
      {children}
    </Button>
  );
};
