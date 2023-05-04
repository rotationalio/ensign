import * as TooltipPrimitive from '@radix-ui/react-tooltip';
import { FC, ReactNode } from 'react';

import HintIcon from '@/components/icons/hint';
export interface HelpTooltipProps {
  children: ReactNode;
  open?: boolean;
  defaultOpen?: boolean;
  onOpenChange?: (open: boolean) => void;
}

const HelpTooltip: FC<HelpTooltipProps> = ({
  children,
  open,
  defaultOpen,
  onOpenChange,
  ...props
}) => {
  return (
    <>
      <TooltipPrimitive.Provider>
        <TooltipPrimitive.Root open={open} defaultOpen={defaultOpen} onOpenChange={onOpenChange}>
          <TooltipPrimitive.Trigger asChild>
            <button>
              <HintIcon />
            </button>
          </TooltipPrimitive.Trigger>
          <TooltipPrimitive.Content
            className="w-full max-w-[275px] rounded-md bg-secondary-slate p-4 text-sm text-white"
            sideOffset={5}
            align="start"
            {...props}
          >
            {children}
          </TooltipPrimitive.Content>
        </TooltipPrimitive.Root>
      </TooltipPrimitive.Provider>
    </>
  );
};

export { HelpTooltip };
