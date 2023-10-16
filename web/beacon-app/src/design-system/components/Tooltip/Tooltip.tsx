import * as RadixTooltip from '@radix-ui/react-tooltip';

import { StyledArrow, StyledContent } from './Tooltip.styles';
import { TooltipProps } from './Tooltip.types';

function Tooltip({ children, title }: TooltipProps) {
  return (
    <RadixTooltip.Provider>
      <RadixTooltip.Root>
        <RadixTooltip.Trigger asChild>
          <span tabIndex={0}>
            <button disabled style={{ pointerEvents: 'none' }} className="p-0">
              {children}
            </button>
          </span>
        </RadixTooltip.Trigger>
        <RadixTooltip.Portal>
          {title && (
            <StyledContent sideOffset={5}>
              {title}
              <StyledArrow />
            </StyledContent>
          )}
        </RadixTooltip.Portal>
      </RadixTooltip.Root>
    </RadixTooltip.Provider>
  );
}

export default Tooltip;
