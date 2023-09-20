import * as RadixTooltip from '@radix-ui/react-tooltip';
import { ReactNode } from 'react';

type PasswordTooltipProps = {
  isFocused: boolean;
  isMobile: boolean;
  triggerNode: ReactNode;
  contentNode?: ReactNode;
};

const PasswordTooltip = ({
  isFocused,
  isMobile,
  triggerNode,
  contentNode,
}: PasswordTooltipProps) => {
  return (
    <RadixTooltip.Provider>
      <RadixTooltip.Root open={isFocused}>
        <RadixTooltip.Trigger asChild>
          <div>{triggerNode}</div>
        </RadixTooltip.Trigger>
        <RadixTooltip.Portal>
          <RadixTooltip.Content
            className="select-none rounded-[4px] bg-white px-[15px] py-[10px] text-xs text-[15px] leading-none shadow-[hsl(206_22%_7%_/_35%)_0px_10px_38px_-10px,_hsl(206_22%_7%_/_20%)_0px_10px_20px_-15px] will-change-[transform,opacity] data-[state=delayed-open]:data-[side=top]:animate-slideDownAndFade data-[state=delayed-open]:data-[side=right]:animate-slideLeftAndFade data-[state=delayed-open]:data-[side=left]:animate-slideRightAndFade data-[state=delayed-open]:data-[side=bottom]:animate-slideUpAndFade"
            sideOffset={2}
            side={isMobile ? 'bottom' : 'right'}
          >
            {contentNode}
            <RadixTooltip.Arrow className="fill-white" />
          </RadixTooltip.Content>
        </RadixTooltip.Portal>
      </RadixTooltip.Root>
    </RadixTooltip.Provider>
  );
};

export default PasswordTooltip;
