import * as RadixTooltip from '@radix-ui/react-tooltip';
import { FC, ReactNode } from 'react';

export interface IconTooltipProps {
  icon: ReactNode;
  content: ReactNode;
}

const IconTooltip: FC<IconTooltipProps> = ({ icon, content, ...props }) => {
  return (
    <>
      <RadixTooltip.Provider>
        <RadixTooltip.Root>
          <RadixTooltip.Trigger asChild>
            <button>{icon}</button>
          </RadixTooltip.Trigger>
          <RadixTooltip.Portal>
            <RadixTooltip.Content
              className="rounded-md bg-secondary-slate p-2 text-sm text-white"
              {...props}
            >
              {content}
              <RadixTooltip.Arrow />
            </RadixTooltip.Content>
          </RadixTooltip.Portal>
        </RadixTooltip.Root>
      </RadixTooltip.Provider>
    </>
  );
};

export { IconTooltip };
