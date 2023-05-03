import { TooltipProvider } from '@radix-ui/react-tooltip';
import * as Tooltip from '@radix-ui/react-tooltip';

import HintIcon from '@/components/icons/hint';

function APIKeyToolTip() {
  return (
    <TooltipProvider>
      <Tooltip.Root>
        <Tooltip.Trigger asChild>
          <button>
            <HintIcon />
          </button>
        </Tooltip.Trigger>
        <Tooltip.Portal>
          <Tooltip.Content
            className="w-full max-w-[275px] rounded-md bg-secondary-slate p-4 text-sm text-white"
            sideOffset={5}
            align="start"
          >
            <p>
              Each key consists of two parts - a ClientID and a ClientSecret. You'll need both to
              establish a client connection, create Ensign topics, publishers, and subscribers. Keep
              your API keys private -- if you misplace your keys, you can revoke them and generate
              new ones.
            </p>
          </Tooltip.Content>
        </Tooltip.Portal>
      </Tooltip.Root>
    </TooltipProvider>
  );
}

export default APIKeyToolTip;
