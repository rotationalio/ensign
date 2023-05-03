import * as TooltipPrimitive from '@radix-ui/react-tooltip';
import { ReactNode } from 'react';

export interface HelpTooltipProps {
  children?: ReactNode;
  content: string;
  open?: boolean;
  defaultOpen?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export function HelpTooltip({
  children,
  content,
  open,
  defaultOpen,
  onOpenChange,
}: HelpTooltipProps) {
  return (
    <>
      <TooltipPrimitive.Provider>
        <TooltipPrimitive.Root open={open} defaultOpen={defaultOpen} onOpenChange={onOpenChange}>
          <TooltipPrimitive.Trigger asChild>{children}</TooltipPrimitive.Trigger>
          <TooltipPrimitive.Content
            className="w-full max-w-[275px] rounded-md bg-secondary-slate p-4 text-sm text-white"
            sideOffset={5}
            align="start"
          >
            {content}
          </TooltipPrimitive.Content>
        </TooltipPrimitive.Root>
      </TooltipPrimitive.Provider>
    </>
  );
}
